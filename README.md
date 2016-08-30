#SCO-Filewatcher

SCO-Filewatcher is still pre alpha alpha alpha alpha. So there's no complete documentation yet :)

Unless otherwise noted, the scofw source files are distributed under the GPL-3.0 license found in the LICENSE file.

TODO documentation

## Configuration
Scofw needs to watch all your projects and all theirs files which are not excluded by some ignore patterns.
In general we recommend to exclude everything what's listed in your .gitignore (just add the correspondig path to your projects 
.gitignore file in the scofw-project configuration).

Use `scofw --help` to see the list of flags and commands.

### Configuration File 
Use the `-c <path_to_config_json>` flag to point to a configuration file. 
This flag is currently not mandatory but scofw will not start if you don't specify it (there is an Github issue for that).

**Important**
In the configuration file always use absolute paths! Relative paths don't work at the moment.

#### Example .sco-config.json
The following JSON can be used as a template:
```
{
    "projects" : [
        {
            "name" : "humio",
            "path" : "/Users/flg/code/humio",
            "ignoreFiles":[
                "/Users/flg/code/humio/.gitignore"
            ],
	         "scoDir": "/Users/flg/code/humio/.sco"
        },
        {
            "name" : "scofw",
            "path" : "/Users/flg/go/scofw",
            "ignoreFiles":[
                "/Users/flg/go/scofw/.gitignore"
            ]
        }
    ],

    "ignorePatterns" : [
        ".idea",
        ".DS_Store",
        "ui/node_modules", 
        "/bin"
        "
    ]      
}
```
To make incremental diffs scofw needs to store copies of your files which have changed. 
The location where scofw will store thes copies and some log output can be 
specified by the `scoDir` property per project.   
* If you don't specify a `scoDir` path for a project then the scofw internal files will be stored in <userHome>/.sco/<project>.  

Some Editors don't save file changes directly to a file, they use a temporary file which will be renamed to original file name.
If you use one of these editor we recommend to exclude the temporary files by listing them in the `ignorePatterns` section
which is used by all projects... 

Here're some examples - if you encounter other editors with similar behaviour please send me the pattern of the temporary files...
##### emacs
```
.#*
*.*~
```

##### vim
```
.*.*.swp
```

### Filebeat
Works currently only with filebeat v5.(alpha). Have a look at the filebeat.yml located here in the project root folder.


## Problems

### I get an error: "Too many open files"
 
#### Linux
If you only want to increase the number of allowed open files for the current shell sesssion:
```
ulimit -n 200000
```

For a system wide setting most Linux distros support
```
sysctl fs.file-max
fs.file-max = 200000
```


#### Mac (El Capitan)
Use `launchctl limit maxfiles` to get the current limit.
Use `lsof | wc -l` to get the current number of open files 

We need to create two files which have are owned by root:wheel and permissions -rw-r--r-- (644) (so I created these files after changing to root `sudo su`)).

**/Library/LaunchDaemons/limit.maxfiles.plist**
```
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
        "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>Label</key>
    <string>limit.maxfiles</string>
    <key>ProgramArguments</key>
    <array>
      <string>launchctl</string>
      <string>limit</string>
      <string>maxfiles</string>
      <string>200000</string>
      <string>200000</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>ServiceIPC</key>
    <false/>
  </dict>
</plist>
``` 

**/Library/LaunchDaemons/limit.maxproc.plist**
```
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple/DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
  <plist version="1.0">
    <dict>
      <key>Label</key>
        <string>limit.maxproc</string>
      <key>ProgramArguments</key>
        <array>
          <string>launchctl</string>
          <string>limit</string>
          <string>maxproc</string>
          <string>2048</string>
          <string>2048</string>
        </array>
      <key>RunAtLoad</key>
        <true />
      <key>ServiceIPC</key>
        <false />
    </dict>
  </plist>
```

In you `.bashrc` or `.bash_profile` or equivalent add
```
ulimit -n 200000
ulimit -u 2048
```

After restarting you can check the changes via `ulimit -a` in a shell or `launchctl limit maxfiles`

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
