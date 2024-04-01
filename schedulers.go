package main

import (
	"fmt"
	"io"
)

type (
	Process struct {
		ProcessID     string
		ArrivalTime   int64
		BurstDuration int64
		Priority      int64
	}
	TimeSlice struct {
		PID   string
		Start int64
		Stop  int64
	}
)

//region Schedulers

// FCFSSchedule outputs a schedule of processes in a GANTT chart and a table of timing given:
// • an output writer
// • a title for the chart
// • a slice of processes
func FCFSSchedule(w io.Writer, title string, processes []Process) {
	var (
		serviceTime     int64
		totalWait       float64
		totalTurnaround float64
		lastCompletion  float64
		waitingTime     int64
		schedule        = make([][]string, len(processes))
		gantt           = make([]TimeSlice, 0)
	)
	for i := range processes {
		if processes[i].ArrivalTime > 0 {
			waitingTime = serviceTime - processes[i].ArrivalTime
		}
		totalWait += float64(waitingTime)

		start := waitingTime + processes[i].ArrivalTime

		turnaround := processes[i].BurstDuration + waitingTime
		totalTurnaround += float64(turnaround)

		completion := processes[i].BurstDuration + processes[i].ArrivalTime + waitingTime
		lastCompletion = float64(completion)

		schedule[i] = []string{
			fmt.Sprint(processes[i].ProcessID),
			fmt.Sprint(processes[i].Priority),
			fmt.Sprint(processes[i].BurstDuration),
			fmt.Sprint(processes[i].ArrivalTime),
			fmt.Sprint(waitingTime),
			fmt.Sprint(turnaround),
			fmt.Sprint(completion),
		}
		serviceTime += processes[i].BurstDuration

		gantt = append(gantt, TimeSlice{
			PID:   processes[i].ProcessID,
			Start: start,
			Stop:  serviceTime,
		})
	}

	count := float64(len(processes))
	aveWait := totalWait / count
	aveTurnaround := totalTurnaround / count
	aveThroughput := count / lastCompletion

	outputTitle(w, title)
	outputGantt(w, gantt)
	outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)
}

func SJFSchedule(w io.Writer, title string, processes []Process) {
    // Sort the processes based on Burst Duration
    sort.Slice(processes, func(i, j int) bool {
        return processes[i].BurstDuration < processes[j].BurstDuration
    })

    FCFSSchedule(w, title, processes)

}
func SJFPrioritySchedule(w io.Writer, title string, processes []Process) {
    // Sort the processes based on Priority first, then by Burst Duration
    sort.Slice(processes, func(i, j int) bool {
        if processes[i].Priority == processes[j].Priority {
            return processes[i].BurstDuration < processes[j].BurstDuration
        }
        return processes[i].Priority < processes[j].Priority
    })

    FCFSSchedule(w, title, processes)
}
func RRSchedule(w io.Writer, title string, processes []Process, quantum int64) {
    var (
        queue           []Process
        currentTime     int64
        schedule        [][]string
        gantt           []TimeSlice
        totalWait       float64
        totalTurnaround float64
    )

    for len(processes) > 0 || len(queue) > 0 {
        if len(queue) == 0 {
            currentTime = processes[0].ArrivalTime
        }

        // Add arrived processes to the queue
        for i := 0; i < len(processes); i++ {
            if processes[i].ArrivalTime <= currentTime {
                queue = append(queue, processes[i])
                processes = append(processes[:i], processes[i+1:]...)
                i--
            }
        }

        // Process the first process in the queue
        p := queue[0]
        queue = queue[1:]
        
        waitingTime := currentTime - p.ArrivalTime
        totalWait += float64(waitingTime)
        
        burst := p.BurstDuration
        if burst > quantum {
            burst = quantum
        }

        currentTime += burst
        p.BurstDuration -= burst
        
        turnaround := currentTime - p.ArrivalTime
        totalTurnaround += float64(turnaround)
        
        gantt = append(gantt, TimeSlice{PID: p.ProcessID, Start: currentTime - burst, Stop: currentTime})
        
        // If the process is not finished, add it back to the queue
        if p.BurstDuration > 0 {
            queue = append(queue, p)
        }
        
        schedule = append(schedule, []string{
            fmt.Sprint(p.ProcessID),
            fmt.Sprint(p.Priority),
            fmt.Sprint(p.BurstDuration + burst),
            fmt.Sprint(p.ArrivalTime),
            fmt.Sprint(waitingTime),
            fmt.Sprint(turnaround),
            fmt.Sprint(currentTime),
        })
    }

    count := float64(len(schedule))
    aveWait := totalWait / count
    aveTurnaround := totalTurnaround / count
    aveThroughput := count / float64(currentTime)

    outputTitle(w, title)
    outputGantt(w, gantt)
    outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)

//endregion
