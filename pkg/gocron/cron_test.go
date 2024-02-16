package gocron

import (
	"fmt"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestInitAndRun(t *testing.T) {
	count := 0
	task1 := func() {
		fmt.Println("running task list:", GetRunningTasks())
	}
	task2 := func() {
		time.Sleep(time.Second)
	}
	task3 := func() {
		count++
		if count%3 == 0 {
			panic("trigger panic")
		}
	}

	tasks := []*Task{
		{
			Name:     "task1",
			TimeSpec: "@every 1s",
			Fn:       task1,
		},
		{
			Name:     "task2",
			TimeSpec: "@every 2s",
			Fn:       task2,
		},
		{
			Name:     "task3",
			TimeSpec: "@every 3s",
			Fn:       task3,
		},
	}

	err := Init()
	if err != nil {
		t.Fatal(err)
	}
	err = Run(tasks...)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 7)
}

func TestRunOnce(t *testing.T) {
	myTask := func() {
		taskName := "myTask"
		fmt.Println("running task list:", GetRunningTasks())
		fmt.Printf("the task '%s' is executed only once\n", taskName)
		DeleteTask(taskName)
		fmt.Println("running task list:", GetRunningTasks())
	}

	tasks := []*Task{
		{
			Name:     "myTask",
			Fn:       myTask,
			TimeSpec: "@every 2s",
		},
	}

	err := Init(WithLog(zap.NewNop()))
	if err != nil {
		t.Fatal(err)
	}
	err = Run(tasks...)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 5)
	Stop()
}

func TestEvery(t *testing.T) {
	task_1 := func() {
		fmt.Println("this is task_1")
		fmt.Println("running task list:", GetRunningTasks())
	}
	task_2 := func() {
		fmt.Println("this is task_2")
	}
	task_3 := func() {
		fmt.Println("this is task_3")
	}
	task_4 := func() {
		fmt.Println("this is task_4")
	}

	tasks := []*Task{
		{
			TimeSpec: EverySecond(5),
			Name:     "task_1",
			Fn:       task_1,
		},
		{
			TimeSpec: EveryMinute(1),
			Name:     "task_2",
			Fn:       task_2,
		},
		{
			TimeSpec: EveryHour(1),
			Name:     "task_3",
			Fn:       task_3,
		},
		{
			TimeSpec: Everyday(1),
			Name:     "task_4",
			Fn:       task_4,
		},
	}

	err := Init()
	if err != nil {
		t.Fatal(err)
	}
	err = Run(tasks...)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 7)
}
