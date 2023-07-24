VERSION  = $(shell git describe --tags --abbrev=0 --dirty --match v[0-9]* 2> /dev/null)
GOFLAGS  = -tags osusergo,netgo,timetzdata
LDFLAGS  = -ldflags="-X 'main.version=$(VERSION)' -extldflags=-static"

GOSRC := $(shell find . -maxdepth 1 -type f -name '*.go')
DATA  := $(shell find data          -type f -name '*')
POSTS := $(shell find posts         -type f -name '*')

website: go.mod go.sum $(GOSRC) $(DATA) $(POSTS)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) $(GOFLAGS) -o website

.PHONY: image
image: website
	docker build -t "n0thub/website:dev" .
