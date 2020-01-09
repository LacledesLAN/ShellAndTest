package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Handler will handle the incoming from the Cobra module.
func Handler(path string, dir bool, showOutput bool) {

	// If not directory, just run.
	if !dir {
		runner(path, showOutput)
		return
	}

	var files []string

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".json" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		runner(file, showOutput)
	}
}

func runner(file string, showOutput bool) {
	criteria, err := readCriteriaFile(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Pre Tasks
	if len(criteria.Target.PreTasks) >= 1 {
		log.Print("Starting Pre-Tasks...")
		tasksRunner(criteria.Target.PreTasks, showOutput)
	}

	// Tasks
	log.Print("Starting Test Tasks...")
	output := criteria.testRunner(showOutput)

	// Parse Output
	log.Print("Starting Parse Output...")
	criteria.parseOutput(output)

	// Post Tasks
	if len(criteria.Target.PostTasks) >= 1 {
		log.Print("Starting Post-Tasks...")
		tasksRunner(criteria.Target.PostTasks, showOutput)
	}

}

// readCriteriaFile attempts to extract testing criteria from a file
func readCriteriaFile(filePath string) (testCriteria, error) {
	filePath = strings.TrimSpace(filePath)

	if len(filePath) == 0 {
		return testCriteria{}, errors.New("Invalid path for a critera file: path cannot be empty or whitespace")
	}

	absolutePath, err := filepath.Abs(filePath)
	if err != nil {
		return testCriteria{}, fmt.Errorf("Could not resolve absolute path of file %q: %w", filePath, err)
	}

	info, err := os.Stat(absolutePath)
	if os.IsNotExist(err) {
		return testCriteria{}, fmt.Errorf("Path %q could not be found - double check it exist and you have read access: %w", filePath, err)

	}

	if info.IsDir() {
		return testCriteria{}, errors.New("%q is a directory and not a file")
	}

	encodedJSON, err := ioutil.ReadFile(absolutePath)
	if err != nil {
		return testCriteria{}, fmt.Errorf("Couldn't read content from file %q: %w", absolutePath, err)
	}

	var target testCriteria
	err = json.Unmarshal([]byte(encodedJSON), &target)
	if err != nil {
		return testCriteria{}, fmt.Errorf("Error decoding JSON from file %q: %w", absolutePath, err)
	}

	return target, nil
}
