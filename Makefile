run:
	go run scofw.go

build-go:
		# Install go with brew install go --with-cc-all
		# to get all compilers
		# env GOOS="darwin" go build -o bin/scofw scofw.go
		# env GOOS="linux" GOARCH=amd64 go build -o bin/scofw_linux scofw.go
		go build -o bin/scofw

filebeat:
	filebeat -c filebeat.yml

local-gource:
	tail -F .sco/logs/*.gource.log | gource --highlight-dirs --realtime --log-format custom -

	# ==> Caveats
	# As of go 1.2, a valid GOPATH is required to use the `go get` command:
	#   https://golang.org/doc/code.html#GOPATH
	#
	# You may wish to add the GOROOT-based install location to your PATH:
	#   export PATH=$PATH:/usr/local/opt/go/libexec/bin
	# ==> Summary
