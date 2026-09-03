.PHONY: build test lint run clean build-quartz

BINARY_NAME=bot.exe
VAULT_REPO?=../vault
QUARTZ_DIR?=../quartz
KNOWLEDGE_DIST=internal/presentation/web/knowledge_dist

build: build-frontend
	go build -o $(BINARY_NAME) ./cmd/bot
	@if [ -d "$(KNOWLEDGE_DIST)" ]; then echo "Knowledge map included."; else echo "Tip: run 'make build-quartz' to include the knowledge map."; fi

build-frontend:
	cd frontend && bun install && bun run build

build-quartz:
	@echo "Building knowledge map from $(VAULT_REPO)..."
	cd $(QUARTZ_DIR) && npx quartz build
	rm -rf $(KNOWLEDGE_DIST)
	cp -r $(QUARTZ_DIR)/public $(KNOWLEDGE_DIST)
	@echo "Knowledge map built to $(KNOWLEDGE_DIST)"

test:
	go test -v ./...

lint:
	golangci-lint run

run: build
	./$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)
	rm -rf bin/
