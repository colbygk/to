TIME_BUILT := $(shell date -u +%x/%X%z)
GIT_HASH := $(shell git rev-list -1 HEAD)
GIT_FETCH_URL := $(shell git remote show $(shell git config branch.$(shell git name-rev --name-only HEAD).remote) -n | grep Fetch | awk '{print $$3}')

to: to.go
	go build -ldflags "-X main.GitHash=$(GIT_HASH) -X main.TimeBuilt=$(TIME_BUILT) -X main.GitFetchURL=$(GIT_FETCH_URL)"
