package shelltest

import "errors"

// define variables, mainly static errors
var (
	ErrCouldNotReadContent     = errors.New("couldn't read content from file")
	ErrCouldNotResolvePath     = errors.New("could not resolve absolute path of file")
	ErrCriteriaNotMet          = errors.New("a should-have or should-lack criteria was not met")
	ErrDecodingJSON            = errors.New("error decoding JSON from file")
	ErrFailedToCreateFile      = errors.New("failed to create file")
	ErrFailedToCreateStdinPipe = errors.New("failed to create stdin pipe")
	ErrFailedToParseOutput     = errors.New("failed to parse output of tasks")
	ErrFailedToReadCriteria    = errors.New("failed to read criteria file")
	ErrFailedToRunTasks        = errors.New("failed to run test tasks")
	ErrInvalidCriteriaPath     = errors.New("invalid path for a criteria file: path cannot be empty or whitespace")
	ErrNoFileSpecified         = errors.New("no file specified")
	ErrPathIsDirectory         = errors.New("path is a directory and not a file")
	ErrPathNotFound            = errors.New("path could not be found - double check it exists and you have read access")
	ErrTaskFailed              = errors.New("one or more tasks failed")
)
