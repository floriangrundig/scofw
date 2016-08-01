#SCO-Filewatcher client

Unless otherwise noted, the scofw source files are distributed under the GPL-3.0 license found in the LICENSE file.


## Configuration
Use `scofw --help` to see the list of flags and commands.

## Problems

1. I get an error: "Too many open files"
- MAC: `ulimit -a` you need to increase "Max open file descriptors" -> when executing with bash it should run out of the box -> this error occurs when starting via fish or csh



## Gource
Installing Gource on a Mac can be done with homebrew: `brew install gource`


## Development

### Run

Use `make run` to run the application with default settings

Use `go run scofw.go --help` to list of flags...


### Add new vendor packages as dependency
We use https://github.com/kardianos/govendor to store our dependencies within the project

Use e.g. `govendor fetch github.com/satori/go.uuid`.
