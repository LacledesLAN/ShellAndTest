package shelltest

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

func Runner(file string, showOutput bool, errCount chan int32) error {
	var err error
	criteria, err := readCriteriaFile(file, errCount)
	if err != nil {
		log.Error().Err(err).Msg(ErrFailedToReadCriteria.Error())
		return ErrFailedToReadCriteria
	}

	// Pre Tasks
	if len(criteria.Target.PreTasks) >= 1 {
		log.Info().Msg("starting pre-tasks")
		if err := tasksRunner(criteria.Target.PreTasks, showOutput, errCount); err != nil {
			log.Error().
				Err(err).
				Msg("failed to run pre-tasks")
			errCount <- 1
		}
	} else {
		log.Warn().Msg("no pre-tasks found")
	}

	// Tasks
	log.Info().Msg("starting test tasks")
	output, err := criteria.testRunner(showOutput, errCount)
	if err != nil {
		log.Error().Err(err).Msg(ErrFailedToRunTasks.Error())
		errCount <- 1
		return ErrFailedToRunTasks
	}

	// Parse Output
	log.Info().Msg("attempting to parse output")
	err = criteria.parseOutput(output, errCount)
	if err != nil {
		log.Error().Err(err).Msg(ErrFailedToParseOutput.Error())
		errCount <- 1
		return ErrFailedToParseOutput
	}

	// Post Tasks
	if len(criteria.Target.PostTasks) >= 1 {
		log.Print("Starting Post-Tasks...")
		if err := tasksRunner(criteria.Target.PostTasks, showOutput, errCount); err != nil {
			log.Error().
				Err(err).
				Msg("failed to run post-tasks")
			errCount <- 1
		}
	}

	return nil
}

// readCriteriaFile attempts to extract testing criteria from a file
func readCriteriaFile(filePath string, errCount chan int32) (testCriteria, error) {
	filePath = strings.TrimSpace(filePath)

	if len(filePath) == 0 {
		errCount <- 1
		return testCriteria{}, ErrInvalidCriteriaPath
	}

	absolutePath, err := filepath.Abs(filePath)
	if err != nil {
		errCount <- 1
		return testCriteria{}, ErrCouldNotResolvePath
	}

	info, err := os.Stat(absolutePath)
	if os.IsNotExist(err) {
		errCount <- 1
		return testCriteria{}, ErrPathNotFound

	}

	if info.IsDir() {
		errCount <- 1
		return testCriteria{}, ErrPathIsDirectory
	}

	encodedYAML, err := os.ReadFile(absolutePath)
	if err != nil {
		errCount <- 1
		return testCriteria{}, ErrCouldNotReadContent
	}

	var target testCriteria
	err = yaml.Unmarshal(encodedYAML, &target)
	if err != nil {
		errCount <- 1
		return testCriteria{}, ErrDecodingJSON
	}

	return target, nil
}
