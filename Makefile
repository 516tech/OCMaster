VERSION := $(shell cat VERSION | tr -d '[:space:]')

.PHONY: help build test run clean docker-up docker-down

help:
	@echo "超频大师 (OCMaster) v$(VERSION) 构建系统"
	@echo "  make run        - 本地启动 S端后端 (SQLite, 零Docker)"
	@echo "  make build      - 构建全部 (版本=$(VERSION))"
	@echo "  make build-c-end - 构建 C端 + 打包 Windows"
	@echo "  make test       - 运行所有测试"
	@echo "  make clean      - 清理构建产物"

# ======== 构建 ========
build: sync-version build-backend build-frontend

sync-version:
	@bash scripts/sync-version.sh

build-backend:
	cd s-end/backend && go build -ldflags="-s -w -X main.Version=$(VERSION)" -o bin/server ./cmd/server

build-frontend:
	cd s-end/frontend && npm run build

build-c-end:
	cd c-end/hardware-scanner && CGO_ENABLED=1 go build -ldflags="-s -w" -o ../bin/ocmaster ./cmd/gui

build-c-end-cli:
	cd c-end/hardware-scanner && CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/ocmaster-cli ./cmd/cli

# ======== 测试 ========
test: test-backend

test-backend:
	cd s-end/backend && go test ./... -count=1 -cover

test-frontend:
	cd s-end/frontend && npm run test

# ======== 运行 ========
run:
	cd s-end/backend && go run -ldflags="-X main.Version=$(VERSION)" ./cmd/server

run-frontend:
	cd s-end/frontend && npm run dev

# ======== 清理 ========
clean:
	rm -rf s-end/backend/bin s-end/backend/data
	rm -rf s-end/frontend/dist
	rm -rf c-end/bin c-end/hardware-scanner/bin
	rm -f .version.env
