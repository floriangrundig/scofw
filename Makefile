CURRENT_DIRECTORY="$(shell pwd)"

run:
	go run scofw.go -c ".sco.config.example.json" -v

run-humio:
	go run scofw.go -c ".sco.config.humio.json" -v

run-global:
	go run scofw.go -c "/Users/flg/.sco.config" -v

linux-bash:
	docker run -it --rm -v $(CURRENT_DIRECTORY):/usr/src/myapp -w /usr/src/myapp scofw-golang:1.6 bash

run-linux:
	docker run -it --rm --name scolinux  -v $(CURRENT_DIRECTORY):/usr/src/myapp -w /usr/src/myapp scofw-golang:1.6 go-wrapper run scofw.go -c ".sco.config.example.linux.json" -v

build-go: build-go-linux build-go-darwin
	#"Finished...."

build-go-linux:
	# TODO seems to be broken since a vendor dependency can not be found -> use run-linux target and use go-wrapper download to fetch dependencies and then run manually from within the container
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

