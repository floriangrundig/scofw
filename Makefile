CURRENT_DIRECTORY="$(shell pwd)"

run:
	go run scofw.go

run-linux:
	docker run -it --rm -v $(CURRENT_DIRECTORY):/usr/src/myapp -w /usr/src/myapp scofw-golang:1.6 go-wrapper run scofw.go

build-go: build-go-linux build-go-darwin
	#"Finished..."

build-go-linux:
	# To install scofw-golang:1.6 docker image (execute buildDockerImage.sh in scripts folder)
	docker run -e "GOOS=linux" -e "GOARCH=amd64" -it --rm -v $(CURRENT_DIRECTORY):/usr/src/myapp -w /usr/src/myapp scofw-golang:1.6 go-wrapper build -o bin/scofw_linux scofw.go

build-go-darwin:
	# Install go with brew install go --with-cc-all
	# to get all compilers
	# env GOOS="darwin" go build -o bin/scofw scofw.go
	# env GOOS="linux" GOARCH=amd64 go build -o bin/scofw_linux scofw.go
	# go build -o bin/scofw
	env GOOS="darwin" go build -o bin/scofw_darwin scofw.go

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
