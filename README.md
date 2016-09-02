# gototcov
Generate total coverage from a golang cover profile and exit non 0 on limit

Can be easy used in a Makefile like this. Dependant on gocovmerge for generating the profile
```
cover-profile:
	@echo Running coverage
	go get github.com/wadey/gocovmerge
	$(eval PKGS := $(shell go list ./... | grep -v /vendor/))
	$(eval PKGS_DELIM := $(shell echo $(PKGS) | sed -e 's/ /,/g'))
	go list -f '{{if or (len .TestGoFiles) (len .XTestGoFiles)}}go test -test.v -test.timeout=120s -covermode=count -coverprofile={{.Name}}_{{len .Imports}}_{{len .Deps}}.coverprofile -coverpkg $(PKGS_DELIM) {{.ImportPath}}{{end}}' $(PKGS) | xargs -I {} bash -c {}
	gocovmerge `ls *.coverprofile` > cover.out
	rm *.coverprofile

cover-test: cover-profile
	gototcov -f cover.out -limit 80 -ignore-zero

```


usage:
```
Usage of gototcov:
  -f string
    	Filename of the cover profile
  -ignore-zero
    	Ignore files with 0%. Example main.go
  -limit float
    	% threshold to throw error
```
