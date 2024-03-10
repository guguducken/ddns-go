default: build
WORKDIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: build
build:
	@cd $(WORKDIR) && go mod tidy && CGO_ENABLED=0 go build -o ddns-go main.go

.PHONY: clean
clean:
	@cd $(WORKDIR) && rm -rf ddns-go dist

.PHONY: image
image:
	@$(WORKDIR)/optools/image/build.sh
