package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/LacledesLAN/ShellAndTest/internal/cli"
	"github.com/LacledesLAN/ShellAndTest/internal/shelltest"
	"github.com/alecthomas/kong"
)

type runCmd struct {
	ShowOutput bool   `kong:"name='output',short='o',help='Show the output of the commands.'"`
	TestDir    string `kong:"name='testdir',short='d',help='Provide the path to a directory: /path/to/dir'"`
	TestFile   string `kong:"name='testfile',short='f',help='Provide the path to a YAML file: /path/to/file.yml'"`
	LogLevel   string `kong:"name='log-level',default='info',enum='panic,fatal,error,warn,info,debug,trace',help='Set log level.'"`
}

var (
	ErrFoundErrors = errors.New("one or more steps failed, check output")
	errCounter     = atomic.Int32{}
)

func (cmd runCmd) Run() error {
	// init zerolog
	cli.InitZeroLog(cmd.LogLevel)

	errCount := make(chan int32)
	defer close(errCount)
	go func() {
		current := int32(0)
		for {
			current = <-errCount
			errCounter.Add(current)
		}
	}()

	var files []string
	// check if stdin, a directory, or testfile is specified
	switch {
	case cmd.TestFile != "":
		files = append(files, cmd.TestFile)
	case cmd.TestDir != "":
		if err := filepath.Walk(cmd.TestDir, func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(path) == ".json" {
				files = append(files, path)
			}
			return nil
		}); err != nil {
			errCount <- 1
			return err
		}
	}

	for file := range files {
		err := shelltest.Runner(files[file], cmd.ShowOutput, errCount)
		if err != nil {
			return err
		}
	}

	if errCounter.Load() > 0 {
		return fmt.Errorf("%d errors found. %w", errCounter.Load(), ErrFoundErrors)
	}

	return nil
}

func main() {
	// configure Kong CLI
	kongContext := kong.Parse(
		&runCmd{},
		kong.Name("ShellAndTest"),
		kong.Description("CLI utility for writing automated tests for 3rd-party CLI binaries."),
		kong.UsageOnError(),
	)

	err := kongContext.Run()
	kongContext.FatalIfErrorf(err)
}
