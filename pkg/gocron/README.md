## gocron

Scheduled task library encapsulated on [cron v3](github.com/robfig/cron).

<br>

### Example of use

```go
package main

import (
    "fmt"
    "time"

    "github.com/zhufuyi/sponge/pkg/gocron"
)

var task1 = func() {
     fmt.Println("this is task1")
     fmt.Println("running task list:", gocron.GetRunningTasks())
}

var taskOnce = func() {
    fmt.Println("this is task2, only run once")
    fmt.Println("running task list:", gocron.GetRunningTasks())
}

func main() {
    err := gocron.Init()
    if err != nil {
        panic(err)
    }
	
    gocron.Run([]*gocron.Task{
        {
            Name:     "task1",
            TimeSpec: "@every 2s",
            Fn:       task1,
        },
        {
            Name:     "taskOnce",
            TimeSpec: "@every 3s",
            Fn:       taskOnce,
            IsRunOnce: true,  // run only once
        },
    }...)

    time.Sleep(time.Second * 10)

    // stop task1
    gocron.DeleteTask("task1")

    // view running tasks
    fmt.Println("running task list:", gocron.GetRunningTasks())
}
```
