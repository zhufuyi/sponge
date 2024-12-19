## gocron

Scheduled task library encapsulated on [cron v3](github.com/robfig/cron).

<br>

### Example of use

```go
package main

import (
	"fmt"
	"time"

	"github.com/go-dev-frame/sponge/pkg/gocron"
	"github.com/go-dev-frame/sponge/pkg/logger"
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
	err := gocron.Init(
			gocron.WithLogger(logger.Get()),
			// gocron.WithLogger(logger.Get(), true), // only print error logs, ignore info logs
		)
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
			Name:      "taskOnce",
			TimeSpec:  "@every 3s",
			Fn:        taskOnce,
			IsRunOnce: true, // run only once
		},
	}...)

	time.Sleep(time.Second * 10)

	// stop task1
	gocron.DeleteTask("task1")

	// view running tasks
	fmt.Println("running task list:", gocron.GetRunningTasks())
}
```
