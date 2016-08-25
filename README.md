#SCO-Filewatcher

SCO-Filewatcher is still pre alpha alpha alpha alpha. So there's no complete documentation yet :)

Unless otherwise noted, the scofw source files are distributed under the GPL-3.0 license found in the LICENSE file.

TODO documentation

## Configuration
Use `scofw --help` to see the list of flags and commands.


### Filebeat
Works currently only with filebeat v5.(alpha). Have a look at the filebeat.yml located here in the project root folder.


## Problems

### I get an error: "Too many open files"
MAC:
`ulimit -a` you need to increase "Max open file descriptors" -> when executing with bash it should run out of the box or at leat configurable via `ulimit -n 4000` -> this error occurs when starting via fish or csh.
Example how to increase: `sysctl -w kern.maxfilesperproc=20000` 

Use `lsof | wc -l` to find out the current number of open file descriptors.


## Gource
Installing Gource on a Mac can be done with homebrew: `brew install gource`


## Development

### Run

Use `make run` to run the application with default settings

Use `go run scofw.go --help` to list of flags...


### Building binaries

#### Ubuntu
apt-get install -y pkg-config cmake
go get -d github.com/libgit2/git2go
cd /go/src/github.com/libgit2/git2go/
git checkout next
git submodule update --init # get libgit2
make install

cd /usr/src/myapp
go build -o bin/scofw_linux scofw.go

### Add new vendor packages as dependency
We use https://github.com/kardianos/govendor to store our dependencies within the project

Use e.g. `govendor fetch github.com/satori/go.uuid`.
