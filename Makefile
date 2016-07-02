# Go commands
GOCMD = go
GOGET = $(GOCMD) get -v -t
GOTEST = $(GOCMD) test -v

# Coverage
COVERAGE_REPORT = coverage.txt
COVERAGE_PROFILE = profile.out
COVERAGE_MODE = atomic

dependencies:
	$(GOGET) ./...

test-coverage:
	cd $(TRAVIS_BUILD_DIR); \
	echo "" > $(COVERAGE_REPORT); \
	for dir in `find . -name "*.go" | grep -o '.*/' | sort -u`; do \
		$(GOTEST) $$dir -coverprofile=$(COVERAGE_PROFILE) -covermode=$(COVERAGE_MODE); \
		if [ $$? != 0 ]; then \
			exit 2; \
		fi; \
		if [ -f $(COVERAGE_PROFILE) ]; then \
			cat $(COVERAGE_PROFILE) >> $(COVERAGE_REPORT); \
			rm $(COVERAGE_PROFILE); \
		fi; \
	done; \
