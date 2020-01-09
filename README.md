# Shell And Test

Shell And Test is a CLI application that will capture `standard out` and `standard error` streams of a application as well as allow for sending commands via it's `standard in` streams. It examines the resulted captured in `standard out` and `standard error` and tests for criteria specified in a JSON file.

## Usage

Shell and Test accepts the following arguments:

```
  -h, --help              help for ShellAndTest
  -o, --output            Show command outputs
  -d, --testdir string    Provide a path to directory: /path/to/dir/
  -f, --testfile string   Provide a path to json: /path/to/file.json
```

Command line syntax is the following:
```bash
./shellandtest -f test.json
./shellandtest -d /my/app/here
cat test.json | ./shellandtest
```
