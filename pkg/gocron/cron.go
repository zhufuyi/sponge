// Package gocron is scheduled task library.
package gocron

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/robfig/cron/v3"
)

var (
	c *cron.Cron
	// task name and id mapping, used to add, delete, modify and query tasks
	nameID = sync.Map{}
	// id and task name mapping, used in log printing
	idName = sync.Map{}
)

// Task scheduled task
type Task struct {
	// seconds (0-59) minutes (0- 59) hours (0-23) days (1-31) months (1-12) weeks (0-6)
	// "*/5 * * * * *"  means every five seconds.
	// "0 15,45 9-12 * * * "  indicates execution at the 15th and 45th minutes from 9 a.m. to 12 a.m. each day
	TimeSpec string

	Name      string // task name
	Fn        func() // task function
	IsRunOnce bool   // if the task is only run once
}

// Init initialize and start timed tasks
func Init(opts ...Option) error {
	o := defaultOptions()
	o.apply(opts...)

	log := &zapLog{zapLog: o.zapLog}
	cronOpts := []cron.Option{
		cron.WithSeconds(), // second-level granularity, default is minute-level granularity
		cron.WithLogger(log),
		cron.WithChain(
			cron.Recover(log),
		),
	}

	c = cron.New(cronOpts...)
	c.Start()

	return nil
}

// Run the tasks
func Run(tasks ...*Task) error {
	if c == nil {
		return errors.New("cron is not initialized")
	}

	var errs []string
	for _, task := range tasks {
		if IsRunningTask(task.Name) {
			errs = append(errs, fmt.Sprintf("task '%s' is already exists", task.Name))
			continue
		}

		if err := checkRunOnce(task); err != nil {
			errs = append(errs, err.Error())
			continue
		}

		id, err := c.AddFunc(task.TimeSpec, task.Fn)
		if err != nil {
			errs = append(errs, fmt.Sprintf("run task '%s' error: %v", task.Name, err))
			continue
		}
		idName.Store(id, task.Name)
		nameID.Store(task.Name, id)
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, " || "))
	}

	return nil
}

func checkRunOnce(task *Task) error {
	if task.Fn == nil {
		return fmt.Errorf("task '%s' is nil", task.Name)
	}
	if task.IsRunOnce {
		job := task.Fn
		task.Fn = func() {
			job()
			DeleteTask(task.Name)
		}
	}
	return nil
}

// IsRunningTask determine if the task is running
func IsRunningTask(name string) bool {
	_, ok := nameID.Load(name)
	return ok
}

// GetRunningTasks gets a list of running task names
func GetRunningTasks() []string {
	var names []string
	nameID.Range(func(key, value interface{}) bool {
		names = append(names, key.(string))
		return true
	})
	return names
}

// DeleteTask stop and delete the specified task
func DeleteTask(name string) {
	if id, ok := nameID.Load(name); ok {
		entryID, isOk := id.(cron.EntryID)
		if !isOk {
			return
		}
		c.Remove(entryID)
		nameID.Delete(name)
		idName.Delete(entryID)
	}
}

// Stop all scheduled tasks
func Stop() {
	if c != nil {
		c.Stop()
	}
}

// EverySecond every second size (1~59)
func EverySecond(size int) string {
	return fmt.Sprintf("@every %ds", size)
}

// EveryMinute every minute size (1~59)
func EveryMinute(size int) string {
	return fmt.Sprintf("@every %dm", size)
}

// EveryHour every hour size (1~23)
func EveryHour(size int) string {
	return fmt.Sprintf("@every %dh", size)
}

// Everyday size (1~31)
func Everyday(size int) string {
	return fmt.Sprintf("@every %dh", size*24)
}
