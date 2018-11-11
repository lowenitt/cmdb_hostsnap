
# define build version
VERSION ?= 1.0.0
BUILD_COMMIT := `git rev-parse HEAD`
BUILD_TIME := `date "+%Y-%m-%d %H:%M:%S"`
GO_VERSION := `go version`

# go build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X 'main.BuildCommit=$(BUILD_COMMIT)' \
-X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GoVersion=$(GO_VERSION)'"

# binary name
BINARY := cmdb_hostsnap
PLATFORMS ?= linux windows darwin
BUILDARC := amd64

# paths
BUILDROOT=$(shell readlink -f .)
BUILDPATH=$(BUILDROOT)/release/build/$(VERSION)
PACKAGEPATH=$(BUILDROOT)/release/package


GO_FILES=$(shell find . -type f -name "*.go" -not -path "./vendor/*" -not -path "./release/*")

all: $(PLATFORMS)

# support multi platform
$(PLATFORMS):
	@mkdir -p $(BUILDPATH)
	@echo "building \033[34m$(BUILDPATH)/$@/plugins/bin/$(BINARY)\033[0m"
	GOOS="$@" GOARC=$(BUILDARC) go build $(LDFLAGS) -o $(BUILDPATH)/$@/plugins/bin/$(BINARY)
	@mkdir -p $(BUILDPATH)/$@/plugins/etc/
	@cp $(BINARY).json $(BUILDPATH)/$@/plugins/etc/$(BINARY).json
	@echo "finish build \033[34m$(BUILDPATH)/$@/plugins/bin/$(BINARY)\033[0m\n"

package: 
	@mkdir -p $(PACKAGEPATH) 
	@for PLATFORM in $(PLATFORMS); \
	do \
		PACKAGEFILE=$(PACKAGEPATH)/$(BINARY)-$$PLATFORM-$(VERSION)-x86_64.tgz ; \
		if [ -d $(BUILDPATH)/$$PLATFORM ]; then \
		echo "packaging \033[34m$$PACKAGEFILE\033[0m"; \
		rm -rf $$PACKAGEFILE ;\
		tar --group=root --owner=root -zcf $$PACKAGEFILE -C $(BUILDPATH)/$$PLATFORM plugins ; \
		echo "finish package \033[34m$$PACKAGEFILE\033[0m\n" ; \
		fi \
	done

.PHONY:tool
tool:
	rm -rf release/tool
	mkdir -p release/tool

gentest:
	gotests -all -excl main -w $(GO_FILES)

.PHONY:test
test:
	rm -rf cover.pprof
	go test -v -cover=true --covermode=count -coverprofile=cover.pprof ./... | tee /dev/stderr

cover:test
	go tool cover -func=cover.pprof

fmt:
	gofmt -l -w $(GO_FILES)

check:
	@test -z $(shell gofmt -l $(GO_FILES) | tee /dev/stderr) || echo "[WARN] Fix formatting issues with 'make fmt'"
	@for d in $$(go list ./... | grep -v /vendor/); do golint $${d}; done
	@go tool vet ${GO_FILES}

clean:
	$(GOCLEAN)
	rm -rf release/
	rm -rf *.pprof

