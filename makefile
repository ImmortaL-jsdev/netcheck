.PHONY: build run test lint clean fmt

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

# Удалить собранный бинарник
clean:
	rm -rf bin/