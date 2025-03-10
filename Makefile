GO=go
GOTEST=${GO} test -v
COLORIZE ?= | $(SED) 's/PASS/✅ PASS/g' | $(SED) 's/FAIL/❌ FAIL/g' | $(SED) 's/SKIP/🔕 SKIP/g'


.PHONY: test
test:
	bash -c "set -e; set -o pipefail; $(GOTEST) . $(COLORIZE)"