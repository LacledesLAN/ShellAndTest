package shelltest

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTasksRunnerWithError(t *testing.T) {
	var errCounter = atomic.Int32{}
	errCount := make(chan int32)
	defer close(errCount)
	go func() {
		// NOTE(mattburchett):
		// I don't love having to put this sleep here, but in testing, it seems to get data
		// to quickly and causes the output to return zero. Maybe there's a better way...
		time.Sleep(1 * time.Second)
		current := int32(0)
		fmt.Println("foobar")
		for {
			current = <-errCount
			errCounter.Add(current)
		}
	}()

	tasks := []string{"eecho foo", "eecho bar"}
	err := tasksRunner(tasks, false, errCount)
	assert.Error(t, err)
	assert.Equal(t, errCounter.Load(), int32(2))
}

func TestTasksRunnerWithNoError(t *testing.T) {
	var errCounter = atomic.Int32{}
	errCount := make(chan int32)
	defer close(errCount)
	go func() {
		// NOTE(mattburchett):
		// I don't love having to put this sleep here, but in testing, it seems to get data
		// too quickly and causes the output to return zero. Maybe there's a better way...
		time.Sleep(1 * time.Second)
		current := int32(0)
		fmt.Println("foobar")
		for {
			current = <-errCount
			errCounter.Add(current)
		}
	}()

	tasks := []string{"echo foo", "echo bar"}
	err := tasksRunner(tasks, false, errCount)
	assert.NoError(t, err)
	assert.Equal(t, errCounter.Load(), int32(0))
}
