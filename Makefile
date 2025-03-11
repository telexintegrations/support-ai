GO=go
GOTEST=${GO} test  ./tests/... -v
COLORIZE ?= | sed 's/PASS/✅ PASS/g' | sed 's/FAIL/❌ FAIL/g' | sed 's/SKIP/🔕 SKIP/g'

.PHONY: test
test:
	bash -c "set -e; set -o pipefail; $(GOTEST) . $(COLORIZE)"

.PHONY: run
run:
	$(GO) run main.go

.PHONY: test
test:
	$(GOTEST)
# Start Docker Compose
.PHONY: docker-up
docker-up:
	docker compose -f docker-compose.dev.yml up  -d

# Stop and remove containers
.PHONY: docker-down
docker-down:
	docker compose -f docker-compose.dev.yml down

# Build and run MongoDB manually (without Compose)
.PHONY: docker-container
docker-container:
	docker build -t my-mongodb .
	docker run -d -p 27017:27017 --name mongodb my-mongodb


.PHONY: clean-container
clean-container:
	docker stop mongodb
	docker rm mongodb
