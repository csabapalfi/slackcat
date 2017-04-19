SHELL = /bin/bash
repo = $(shell basename $(CURDIR))

release:
ifndef version
	@echo "Please provide a version"
	exit 1
endif
ifndef GITHUB_TOKEN
	@echo "Please set GITHUB_TOKEN in the environment"
	exit 1
endif
	git tag $(version)
	git push origin --tags
	mkdir -p releases/$(version)
	GOOS=linux GOARCH=amd64 go build -o releases/$(version)/$(repo)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o releases/$(version)/$(repo)-darwin-amd64 .
	GOOS=windows GOARCH=amd64 go build -o releases/$(version)/$(repo)-windows-amd64 .
ifndef RELEASE
	go get -u github.com/aktau/github-release
endif
	github-release release --user csabapalfi --repo $(repo) --tag $(version) || true
	github-release upload --user csabapalfi --repo $(repo) --tag $(version) --name $(repo)-linux-amd64 --file releases/$(version)/$(repo)-linux-amd64 || true
	github-release upload --user csabapalfi --repo $(repo) --tag $(version) --name $(repo)-darwin-amd64 --file releases/$(version)/$(repo)-darwin-amd64 || true
	github-release upload --user csabapalfi --repo $(repo) --tag $(version) --name $(repo)-windows-amd64 --file releases/$(version)/$(repo)-windows-amd64 || true
