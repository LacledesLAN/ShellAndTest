package shelltest

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/creack/pty"
	"github.com/rs/zerolog/log"
)

// testCriteria contains all the criteria and information needed from the JSON file.
type testCriteria struct {
	ShouldEcho []struct {
		Command    string `yaml:"command"`
		ShouldHave string `yaml:"should_have"`
	} `yaml:"should_echo"`
	ShouldHave []string `yaml:"should_have"`
	ShouldLack []string `yaml:"should_lack"`
	Target     struct {
		Execute         string   `yaml:"execute"`
		PostTasks       []string `yaml:"post_tasks"`
		PreTasks        []string `yaml:"pre_tasks"`
		ShouldEchoDelay int      `yaml:"should_echo_delay"`
		Timeout         int      `yaml:"timeout"`
	} `yaml:"target"`
}

// tasksRunner is utilized when performing the Pre/Post Tasks specified by the json file.
func tasksRunner(tasks []string, showOutput bool, errCount chan int32) error {
	var failedTask = 0
	for task := range tasks {
		log.Info().Str("task", tasks[task]).Msg("executing task")
		name, args := splitArgs(tasks[task])

		cmd := exec.Command(name, args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Error().Err(err).Str("task", tasks[task]).Msg("task failed")
			failedTask++
			continue
		}
		if output != nil {
			if showOutput {
				fmt.Fprintln(os.Stdout, bytes.NewBuffer(output).String())
			}
		}
	}

	if failedTask > 0 {
		if errCount != nil {
			errCount <- int32(failedTask)
		}
		return ErrTaskFailed
	}

	return nil
}

// testRunner is the main runner of the executable. This will perform the execution
// as well as the 'should-echo' requests (if applicable) and then returns the
// stdout/stderr buffer.

// NOTE(mattburchett): There is a lot of bad that's about to happen here. At the time of this
// writing, Go doesn't have a good module for generating a PTY interface so that we can interact
// with Docker images using the '-t' parameter. Therefore, there is a lot of weird stuff going on
// to make the one module I found do what I want it to do.
func (criteria *testCriteria) testRunner(showOutput bool, errCount chan int32) (string, error) {

	// create a log file since the pty won't let me direct it to a buffer...
	record, err := os.Create("shellandtest-" + time.Now().Format("2006-01-02_150405") + ".log")
	if err != nil {
		errCount <- 1
		return "", ErrFailedToCreateFile
	}
	defer record.Close()

	// set command configuration
	name, args := splitArgs(criteria.Target.Execute)
	cmd := exec.Command(name, args...)

	// Start the command with a pty.
	tty, err := pty.Start(cmd)
	if err != nil {
		errCount <- 1
		return "", err
	}
	// Make sure to close the pty at the end.
	defer tty.Close()

	// spin up a channel so we can wait for the command to exit successfully
	done := make(chan error)
	go func() { done <- cmd.Wait() }()

	// if there is should_echo criteria, then go run these after the wait in a goroutine
	if len(criteria.ShouldEcho) >= 1 {
		go func() {

			// sleep for the duration specified in the configuration file
			time.Sleep(time.Duration(criteria.Target.ShouldEchoDelay) * time.Second)

			// loop over all ShouldEcho criteria and execute the commands
			for echo := range criteria.ShouldEcho {
				log.Info().
					Str("Command", criteria.ShouldEcho[echo].Command).
					Msg("echoing into executable")
				_, err := tty.Write([]byte(criteria.ShouldEcho[echo].Command + "\n"))
				if err != nil {
					log.Error().
						Err(err).
						Str("Command", criteria.ShouldEcho[echo].Command).
						Msg("failed to echo command")
					if errCount != nil {
						errCount <- 1
					}
				}
			}

			// give it a small bit of time to respond and then call a SIGHUP to try
			// to cleanly exit the process
			time.Sleep(5 * time.Second)
			if err := cmd.Process.Signal(syscall.SIGHUP); err != nil {
				log.Error().Err(err).Msg("failed to send SIGHUP to process")
			}

		}()
	}

	// Start a timer just in case the process doesn't stop
	timeout := time.After(time.Duration(criteria.Target.Timeout) * time.Second)

	// The select statement allows us to execute based on which channel
	// we get a message from first.
	select {
	case <-timeout:
		// copy the output to a log file
		if _, err := io.Copy(record, tty); err != nil {
			log.Warn().Err(err).Msg("failed to copy output to log")
		}

		// kill the process since the timeout has been reached
		if err := cmd.Process.Signal(syscall.SIGKILL); err != nil {
			log.Error().
				Err(err).
				Msg("failed to kill process")
			if errCount != nil {
				errCount <- 1
			}
		}
		log.Warn().
			Int("timeout", criteria.Target.Timeout).
			Msg("command timeout exceeded.")
	case <-done:
		// copy the output to a log file
		if _, err := io.Copy(record, tty); err != nil {
			log.Warn().Err(err).Msg("failed to copy output to log")
		}
		log.Info().Msg("command completed before timeout.")
	}

	// return the output of the commands if showOutput is true
	if showOutput {
		file, err := os.ReadFile(record.Name())
		if err != nil {
			errCount <- 1
		}
		log.Info().Msg("returning output")
		fmt.Fprintln(os.Stdout, bytes.NewBuffer(file).String())
	}

	return record.Name(), err
}

// parseOutput will take in the stdout/stderr buffer and process it for the
// should-have and should-lack criteria.
func (criteria *testCriteria) parseOutput(filename string, errCount chan int32) error {
	var (
		missingShouldHave = 0
		missingShouldLack = 0
	)

	if filename == "" {
		errCount <- 1
		return ErrNoFileSpecified
	}

	file, err := os.ReadFile(filename)
	if err != nil {
		errCount <- 1
		return ErrCouldNotReadContent
	}

	outputStr := bytes.NewBuffer(file).String()

	for echo := range criteria.ShouldEcho {
		log.Info().
			Str("Command", criteria.ShouldEcho[echo].Command).
			Str("ShouldHave", criteria.ShouldEcho[echo].ShouldHave).
			Msg("checking output for should-have")

		if strings.Contains(outputStr, criteria.ShouldEcho[echo].ShouldHave) {
			log.Info().
				Str("Command", criteria.ShouldEcho[echo].Command).
				Str("ShouldHave", criteria.ShouldEcho[echo].ShouldHave).
				Msg("should-have found")
		} else {
			log.Error().
				Str("Command", criteria.ShouldEcho[echo].Command).
				Str("ShouldHave", criteria.ShouldEcho[echo].ShouldHave).
				Msg("should-have not found")
			if errCount != nil {
				errCount <- 1
			}
		}
	}

	log.Info().Msg("checking for each should-have")
	for have := range criteria.ShouldHave {
		log.Info().
			Str("ShouldHave", criteria.ShouldHave[have]).
			Msg("checking output for should-have")
		if strings.Contains(outputStr, criteria.ShouldHave[have]) {
			log.Info().
				Str("ShouldHave", criteria.ShouldHave[have]).
				Msg("found should-have")
		} else {
			log.Error().
				Str("ShouldHave", criteria.ShouldHave[have]).
				Msg("should-have not found")
			if errCount != nil {
				errCount <- 1
			}
			missingShouldHave++
		}
	}

	log.Info().Msg("checking for each should-lack")
	for lack := range criteria.ShouldLack {
		log.Info().
			Str("ShouldLack", criteria.ShouldLack[lack]).
			Msg("checking output for should-lack")
		if !strings.Contains(outputStr, criteria.ShouldLack[lack]) {
			log.Info().
				Str("ShouldLack", criteria.ShouldLack[lack]).
				Msg("should-lack not found")
		} else {
			log.Error().
				Str("ShouldLack", criteria.ShouldLack[lack]).
				Msg("should-lack found")
			if errCount != nil {
				errCount <- 1
			}
			missingShouldLack++
		}
	}

	if missingShouldHave > 0 || missingShouldLack > 0 {
		log.Error().
			Int("missingShouldHave", missingShouldHave).
			Int("missingShouldLack", missingShouldLack).
			Msg(ErrCriteriaNotMet.Error())
		if errCount != nil {
			errCount <- 1
		}
		return ErrCriteriaNotMet
	} else {
		log.Info().
			Int("missingShouldHave", missingShouldHave).
			Int("missingShouldLack", missingShouldLack).
			Msg("all should-have and should-lack criteria was met")
	}

	return nil
}

// splitArgs splits the executable, pre tasks, and post tasks into a string and []string
// to allow feeding to the exec.Command() function since the arguments are separate.
func splitArgs(input string) (string, []string) {
	split := strings.Split(input, " ")
	var cmd string
	var args []string
	if split[1:] != nil {
		cmd, args = split[0], split[1:]
	} else {
		cmd, args = split[0], nil
	}

	return cmd, args
}
