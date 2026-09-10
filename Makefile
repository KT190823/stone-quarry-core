.PHONY: all run seed run-seed build test clean help

# Default target
all: run

## Chạy API server (không nạp lại seed, khởi động nhanh)
run:
	@echo "🚀 Khởi động Mo_da Backend API..."
	go run .

## Nạp lại dữ liệu mẫu vào PostgreSQL (Database Seed)
seed:
	@echo "🌱 Đang nạp dữ liệu mẫu vào cơ sở dữ liệu..."
	go run . -seed

## Nạp seed xong rồi tự động khởi động API server
run-seed:
	@echo "🌱 Nạp dữ liệu mẫu và khởi động Mo_da Backend API..."
	go run . -with-seed

## Biên dịch binary cho môi trường production
build:
	@echo "📦 Đang biên dịch Mo_da backend binary..."
	@mkdir -p bin
	go build -ldflags="-s -w" -o bin/mo-da-backend main.go
	@echo "✅ Biên dịch hoàn tất tại bin/mo-da-backend"

## Chạy kiểm thử tự động
test:
	@echo "🧪 Đang chạy unit tests..."
	go test -v ./...

## Dọn dẹp thư mục build
clean:
	@echo "🧹 Dọn dẹp bản build..."
	rm -rf bin/

## Hiển thị danh sách các lệnh hỗ trợ
help:
	@echo "=========================================================="
	@echo "  MO_DA BACKEND MAKEFILE WORKFLOW"
	@echo "=========================================================="
	@echo "  make run        : Khởi động API server ngay (MẶC ĐỊNH, không seed)"
	@echo "  make seed       : Chỉ chạy nạp dữ liệu mẫu (Seed Database) rồi thoát"
	@echo "  make run-seed   : Chạy seed dữ liệu trước rồi khởi động API server"
	@echo "  make build      : Biên dịch binary ra thư mục bin/mo-da-backend"
	@echo "  make test       : Chạy toàn bộ unit tests"
	@echo "  make clean      : Xóa thư mục bin/"
	@echo "=========================================================="
