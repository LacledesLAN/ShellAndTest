package test

import (
	"encoding/json"
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
	criteria := readJSON(file)

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

// readJSON will read the json files from the specified path(s).
func readJSON(input string) testCriteria {
	var criteria testCriteria

	if !strings.Contains(input, ".json") {
		err := json.Unmarshal([]byte(input), &criteria)
		if err != nil {
			log.Fatal("Failed to unmarshal json.")
		}
		return criteria
	}

	content, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatal("Failure to read file.")
	}
	err = json.Unmarshal([]byte(content), &criteria)
	if err != nil {
		fmt.Println("Failed to unmarshal json.")
		os.Exit(1)
	}

	return criteria

}
