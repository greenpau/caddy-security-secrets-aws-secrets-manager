.PHONY: test ctest covdir coverage linter gtest qtest clean dep release license build_info
APP_VERSION:=$(shell cat VERSION | head -1)
GIT_COMMIT:=$(shell git describe --dirty --always)
GIT_BRANCH:=$(shell git rev-parse --abbrev-ref HEAD -- | head -1)
BUILD_USER:=$(shell whoami)
BUILD_DATE:=$(shell date +"%Y-%m-%d")
VERBOSE:=-v
ifdef TEST
	TEST:="-run ${TEST}"
endif

all: build_info
	@echo "$@: complete"

build_info:
	@echo "Version: $(APP_VERSION), Branch: $(GIT_BRANCH), Revision: $(GIT_COMMIT)"
	@echo "Build on $(BUILD_DATE) by $(BUILD_USER)"

linter:
	@echo "Running lint checks"
	@golint -set_exit_status ./...
	@echo "$@: complete"

gtest:
	@go test $(VERBOSE) -coverprofile=.coverage/coverage.out ./...
	@echo "$@: complete"

test: build_info covdir linter gtest coverage
	@echo "$@: complete"

ctest: covdir linter
	@richgo version || go install github.com/kyoh86/richgo@latest
	@time richgo test $(VERBOSE) $(TEST) -coverprofile=.coverage/coverage.out ./...
	@echo "$@: complete"

covdir:
	@echo "Creating .coverage/ directory"
	@mkdir -p .coverage
	@echo "$@: complete"

coverage:
	@#go tool cover -help
	@go tool cover -html=.coverage/coverage.out -o .coverage/coverage.html
	@go test -covermode=count -coverprofile=.coverage/coverage.out ./...
	@go tool cover -func=.coverage/coverage.out | grep -v "100.0"
	@echo "$@: complete"

clean:
	@rm -rf .doc
	@rm -rf .coverage
	@echo "$@: complete"

qtest: covdir
	@echo "Perform quick tests ..."
	@time richgo test $(VERBOSE) $(TEST) -coverprofile=.coverage/coverage.out -run TestGetMetadata ./*.go
	@go tool cover -html=.coverage/coverage.out -o .coverage/coverage.html
	@#go tool cover -func=.coverage/coverage.out | grep -v "100.0"
	@go tool cover -func=.coverage/coverage.out
	@echo "$@: complete"

dep:
	@echo "Making dependencies check ..."
	@golint || go install golang.org/x/lint/golint@latest
	@go install github.com/kyoh86/richgo@latest
	@versioned || go install github.com/greenpau/versioned/cmd/versioned@latest
	@echo "$@: complete"

license:
	@versioned || go install github.com/greenpau/versioned/cmd/versioned@latest
	@for f in `find ./ -type f -name '*.go'`; do versioned -addlicense -copyright="Paul Greenberg greenpau@outlook.com" -year=2022 -filepath=$$f; done
	@#for f in `find ./ -type f -name '*.go'`; do versioned -striplicense -filepath=$$f; done
	@echo "$@: complete"

release:
	@echo "Making release"
	@go mod tidy
	@go mod verify
	@if [ $(GIT_BRANCH) != "main" ]; then echo "cannot release to non-main branch $(GIT_BRANCH)" && false; fi
	@git diff-index --quiet HEAD -- || ( echo "git directory is dirty, commit changes first" && false )
	@versioned -patch
	@echo "Patched version"
	@git add VERSION
	@git commit -m "released v`cat VERSION | head -1`"
	@git tag -a v`cat VERSION | head -1` -m "v`cat VERSION | head -1`"
	@git push
	@git push --tags
