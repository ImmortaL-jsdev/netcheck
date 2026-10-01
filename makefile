.PHONY: build run test lint fmt install test-zapret build-all clean

# Сборка бинарника в bin/
build:
	go build -o bin/netcheck cmd/netcheck/main.go

# Запуск без сборки (для разработки)
run:
	go run cmd/netcheck/main.go

# Прогон всех тестов
test:
	go test -v ./...

# Линтер (если установлен golangci-lint)
lint:
	golangci-lint run ./...

# Форматирование кода
fmt:
	go fmt ./...

# Установка в $GOPATH/bin
install:
	go install ./cmd/netcheck

# Тесты только для zapret-модуля
test-zapret:
	go test -v ./internal/zapret/...

# Сборка под все платформы
build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/netcheck-linux-amd64 ./cmd/netcheck
	GOOS=windows GOARCH=amd64 go build -o bin/netcheck-windows-amd64.exe ./cmd/netcheck
	GOOS=darwin GOARCH=amd64 go build -o bin/netcheck-darwin-amd64 ./cmd/netcheck

# Удалить собранный бинарник
clean:
	rm -rf bin/