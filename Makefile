.PHONY: all install-tools setup-components generate templ tailwind server build dev clean

# Detect binaries in PATH or fallback to user's ~/go/bin
TEMPL := $(shell which templ 2>/dev/null || echo $(HOME)/go/bin/templ)
AIR := $(shell which air 2>/dev/null || echo $(HOME)/go/bin/air)
TEMPLUI := $(shell which templui 2>/dev/null || echo $(HOME)/go/bin/templui)
TAILWIND := $(shell which tailwindcss 2>/dev/null || echo tailwindcss)

all: build

install-tools:
	@echo "Installing required dev tools..."
	@which templ >/dev/null 2>&1 || go install github.com/a-h/templ/cmd/templ@latest
	@which air >/dev/null 2>&1 || go install github.com/air-verse/air@latest
	@which templui >/dev/null 2>&1 || go install github.com/templui/templui/cmd/templui@latest

setup-components: install-tools
	@echo "Initializing and adding templui components..."
	$(TEMPLUI) --force init
	$(TEMPLUI) -force add "*"

generate: setup-components
	@echo "Generating initial templates and styles..."
	$(TEMPL) generate
	$(TAILWIND) -i ./assets/css/input.css -o ./assets/css/output.css

templ:
	@echo "Watching templates..."
	$(TEMPL) generate --watch --proxy="http://localhost:8080" --open-browser=false

tailwind:
	@echo "Watching tailwind styles..."
	$(TAILWIND) -i ./assets/css/input.css -o ./assets/css/output.css --watch

server:
	@echo "Starting hot-reloading server..."
	$(AIR)

build: generate
	@echo "Building production bundles..."
	@mkdir -p bin
	go build -o bin/server ./cmd/main.go
	@echo "Production build completed. Binary is in bin/server"

dev: generate
	@echo "Starting development environment..."
	make -j3 templ tailwind server

clean:
	@echo "Cleaning build artifacts..."
	rm -rf tmp/ bin/
	rm -rf assets/css/output.css
	rm -rf assets/js/templui/
	rm -rf views/components/
	rm -rf views/utils/templui.go
	find . -name "*_templ.go" -type f -delete
	@echo "Clean completed."
