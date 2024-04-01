package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "sort"
    "math/rand"
    "time"
)

type Process struct {
    ProcessID      int
    ArrivalTime    int64
    BurstDuration  int64
    Priority       int
}

type TimeSlice struct {
    PID   int
    Start int64
    Stop  int64
}

func FCFSSchedule(w io.Writer, title string, processes []Process) {
    // Implement your FCFS scheduling logic here
}

func SJFSchedule(w io.Writer, title string, processes []Process) {
    // Sort the processes based on Burst Duration
    sort.Slice(processes, func(i, j int) bool {
        return processes[i].BurstDuration < processes[j].BurstDuration
    })

    FCFSSchedule(w, title, processes)
}

func SJFPrioritySchedule(w io.Writer, title string, processes []Process) {
    // Sort the processes based on Priority first, then by Burst Duration
    sort.Slice(processes, func(i, j int) bool {
        if processes[i].Priority == processes[j].Priority {
            return processes[i].BurstDuration < processes[j].BurstDuration
        }
        return processes[i].Priority < processes[j].Priority
    })

    FCFSSchedule(w, title, processes)
}

func RRSchedule(w io.Writer, title string, processes []Process, quantum int64) {
    // Implement your Round-Robin scheduling logic here
}

func openProcessingFile(args ...string) (*os.File, func(), error) {
    // Implement your file opening logic here
    return nil, nil, nil
}

func loadProcesses(f *os.File) ([]Process, error) {
    // Implement your process loading logic here
    return nil, nil
}

func outputTitle(w io.Writer, title string) {
    // Implement your title output logic here
}

func outputGantt(w io.Writer, gantt []TimeSlice) {
    // Implement your Gantt chart output logic here
}

func outputSchedule(w io.Writer, schedule [][]string, aveWait, aveTurnaround, aveThroughput float64) {
    // Implement your schedule output logic here
}

func main() {
    // Seed the random number generator
    rand.Seed(time.Now().UnixNano())

    // Generate random processes for testing
    numProcesses := 5 // Adjust the number of processes as needed
    processes := make([]Process, numProcesses)

    for i := 0; i < numProcesses; i++ {
        processes[i] = Process{
            ProcessID:     i + 1,
            ArrivalTime:   int64(rand.Intn(10)), // Random arrival time between 0 and 9
            BurstDuration: int64(rand.Intn(10)), // Random burst duration between 0 and 9
            Priority:      rand.Intn(5),         // Random priority between 0 and 4
        }
    }

    // First-come, first-serve scheduling
    FCFSSchedule(os.Stdout, "First-come, first-serve", processes)
    
    // Shortest-job-first scheduling
    SJFSchedule(os.Stdout, "Shortest-job-first", processes)
    
    // SJF Priority scheduling
    SJFPrioritySchedule(os.Stdout, "Priority", processes)
    
    // Round-robin scheduling
    quantum := int64(2) // set your quantum value
    RRSchedule(os.Stdout, "Round-robin", processes, quantum)
}
