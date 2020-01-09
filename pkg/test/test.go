package test

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
)

// tasksRunner is utilized when performing the Pre/Post Tasks specified by the json file.
func tasksRunner(tasks []string, showOutput bool) {
	defer color.Unset()
	red, _ := colors()
	for _, pt := range tasks {
		log.Printf("Executing Task: %v", pt)
		name, args := splitArgs(pt)

		cmd := exec.Command(name, args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			red.Set()
			log.Fatalf("Task Failed: %v", err)
		}
		if output != nil {
			if showOutput {
				log.Print(string(output))
			}
		}
	}
}

// testRunner is the main runner of the executable. This will perform the execution
// as well as the 'should-echo' requests (if applicable) and then returns the
// stdout/stderr buffer.
func (c *testCriteria) testRunner(showOutput bool) bytes.Buffer {

	name, args := splitArgs(c.Target.Execute)
	cmd := exec.Command(name, args...)
	stdin, err := cmd.StdinPipe()
	red, green := colors()
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	defer out.Reset()

	if len(c.ShouldEcho) >= 1 {
		// This function uses the stdin buffer for entering commands into the executing
		// program.
		go func() {
			defer stdin.Close()
			time.Sleep(time.Duration(c.Target.ShouldEchoDelay) * time.Second)
			for _, echo := range c.ShouldEcho {
				log.Printf("Echoing '%v' into executable.", echo.Command)
				stdin.Write([]byte(echo.Command + "\n"))
				time.Sleep(1 * time.Second)
				log.Printf("Checking output for: %v", echo.ShouldHave)
				inputStr := out.String()

				if strings.Contains(inputStr, echo.ShouldHave) {
					green.Set()
					log.Printf("%v found!", echo.ShouldHave)
					color.Unset()
				} else {
					red.Set()
					log.Printf("%v not found!", echo.ShouldHave)
					color.Unset()
				}

			}

		}()
	}

	cmd.Start()

	done := make(chan error)
	go func() { done <- cmd.Wait() }()

	// Start a timer
	timeout := time.After(time.Duration(c.Target.Timeout) * time.Second)

	// The select statement allows us to execute based on which channel
	// we get a message from first.

	select {
	case <-timeout:
		// Timeout happened first, kill the process and gather the output.
		cmd.Process.Kill()
		log.Print("Command timeout exceeded.")
		if showOutput {
			log.Print("Returning output.")
			log.Print(out.String())
		}
	case <-done:
		// Command completed before timeout. Print output and error if it exists.
		log.Print("Command completed before timeout.")
		if showOutput {
			log.Print("Returning output.")
			log.Print(out.String())
		}
	}

	return out
}

// parseOutput will take in the stdout/stderr buffer and process it for the
// should-have and should-lack criteria.
func (c *testCriteria) parseOutput(input bytes.Buffer) {
	defer color.Unset()

	red, green := colors()

	inputStr := input.String()
	count := 0

	log.Print("Checking for each should-have:")

	for _, sH := range c.ShouldHave {
		log.Printf("Checking output for %v", sH)
		if strings.Contains(inputStr, sH) {
			green.Set()
			log.Printf("%v found!", sH)
			color.Unset()
		} else {
			red.Set()
			log.Printf("%v not found!", sH)
			color.Unset()
			count++
		}
	}

	log.Print("Checking for each should-lack:")
	for _, sL := range c.ShouldLack {
		log.Printf("Checking output for %v", sL)
		if !strings.Contains(inputStr, sL) {
			green.Set()
			log.Printf("%v not found!", sL)
			color.Unset()
		} else {
			red.Set()
			log.Printf("%v found!", sL)
			color.Unset()
			count++
		}
	}

	if count > 0 {
		red.Set()
		log.Fatal("A should-have or should-lack criteria was not met.")
	} else {
		green.Set()
		log.Print("All should-have and should-lack criteria is met.")
	}
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

// colors provides shell-friendly colors.
func colors() (*color.Color, *color.Color) {

	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)

	return red, green

}
