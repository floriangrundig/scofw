#SCO-Filewatcher client

TODO

## Configuration
Use `scofw --help` to see the list of flags and commands.

## Gource
Installing Gource on a Mac can be done with homebrew: `brew install gource`

## Development

### Run

Use `make run` to run the application with default settings

Use `go run scofw.go --help` to list of flags...


### Add new vendor packages as dependency
We use https://github.com/kardianos/govendor to store our dependencies within the project

Use e.g. `govendor fetch github.com/satori/go.uuid`.
