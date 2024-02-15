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
	taskName := "taskOnce"
    fmt.Println("this is taskOnce")
    gocron.DeleteTask(taskName)
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
            TimeSpec: "@every 5s",
            Fn:       taskOnce,
        },
    }...)

    time.Sleep(time.Minute)

    // delete task1
    gocron.DeleteTask("task1")

    // view running tasks
    fmt.Println("running task list:", gocron.GetRunningTasks())
}
```
