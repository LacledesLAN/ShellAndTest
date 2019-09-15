# Design Doc

## The Idea

ShellAndWait should launch a CLI application, capturing its `standard out` and `standard error` streams and sending it commands via its `standard in` streams.  It should examine the results captured `standard out` and `standard error` and test for the following:

- `Should have` - These strings should be found in either out stream; not having one ore more of these strings should result in a failure.
- `Should lack` - These strings should not exist in either out stream; having one or more of these strings should result in a failure.

After a delay `target.should-echo-delay` commands should be run, sent to the app via it's `standard in`, and then check for expected output to be matched (similar to `should have`).

## Requirements

- CLI application that can be compiled for cross-platform use (Linux x64 and Windows x64).
- App needs to work when Dockerized (both Linux and Windows).
- Test specifications to be stored in `.json` file format; see bellow for possible format.
- App needs to be provided with at least one 'test specification' json file to do anything:
  - Either receive one, and only one, from the OS via a pipe. This is who it will be used in Cloud CI/CD builds.            Example: `cat test.json | shellandtest`
  - Or receive one, and only one via an argument. Use the [`spf13/cobra`](https://github.com/spf13/cobra) golang package for this.
        Example: `shellandtest --testfile /path/to/file.json`
- Or receive a directory via an argument. Use the [`spf13/cobra`](https://github.com/spf13/cobra) golang package for this. If this option is used it should run through all tests in the directory.
        Example: `shellandtest --testdir /path/to/dir/`

## Possible Test Specification Format

This is just a suggestion -- feel free to modify as you see fit.

```javascript
{
	"target": {
        "pre-tasks": ["these commands", "should be run in order", "before executing the shell app", "and should not affect timeout"],
		"execute": "command use to execute the application being tested",
		"should-echo-delay": 30,
		"timeout": 90,  // if the app being test runs over this kill the process
        "expectedExitCode": 0,
        "post-tasks": ["test commands should be run in order", "after the application testing is complete"]
	},
	"should-have": [
		"these string should be found",
		"in the app's output (either stream)"
	],
	"should-lack": [
		"these strings should not be found",
		"in the app's output"
	],
	"should echo": [{
		"command": "command to send to application",
		"should-have": "this string should be captured in the output after the associated command has been sent"
	}]
}
```
