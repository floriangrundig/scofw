# ScoFW go container (for cross compilation of our executables)
# based on https://github.com/docker-library/golang/blob/9f666dc2f4f51df564613f787d28b3a2353243e0/1.6/Dockerfile
# modifications:
# - install pkg-config and cmake packages
# - install git2go and prepare for static linking (www.petethompson.net/blog/golang/2015/10/04/getting-going-with-git2go/)
FROM buildpack-deps:jessie-scm

# gcc for cgo
RUN apt-get update && apt-get install -y --no-install-recommends \
		g++ \
		gcc \
		libc6-dev \
		make \
    pkg-config \
    cmake \
	&& rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.6.3
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 cdde5e08530c0579255d6153b08fdb3b8e47caabbe717bc7bcd7561275a87aeb

RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
	&& echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
	&& tar -C /usr/local -xzf golang.tar.gz \
	&& rm golang.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

# prepare for static linking
RUN go get -d github.com/libgit2/git2go
RUN go get -u github.com/kardianos/govendor
RUN go get github.com/mitchellh/go-homedir

WORKDIR $GOPATH/src/github.com/libgit2/git2go
RUN git checkout next
RUN git submodule update --init # get libgit2
RUN make install


WORKDIR $GOPATH
COPY go-wrapper /usr/local/bin/
