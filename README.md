# Shell And Test

Shell And Test is a CLI application that will capture `standard out` and `standard error` streams of a application as well as allow for sending commands via it's `standard in` streams. It examines the resulted captured in `standard out` and `standard error` and tests for criteria specified in a YAML file.

## Usage

Shell and Test accepts the following arguments:

```
Flags:
  -h, --help                Show context-sensitive help.
  -o, --output              Show the output of the commands.
  -d, --testdir=STRING      Provide the path to a directory: /path/to/dir
  -f, --testfile=STRING     Provide the path to yml: /path/to/file.yml
      --log-level="info"    Set log level.
```

Command line syntax is the following:
```bash
./shellandtest -f test.yml
./shellandtest -d /my/app/here
```
