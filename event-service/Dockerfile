FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Копируем исходный код приложения в контейнер
COPY . .

# Сборка приложения
RUN go build -o server ./cmd/server/main.go

# Второй этап: используем минимальный образ для запуска приложения
FROM golang:1.21

# Устанавливаем необходимые зависимости
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Рабочая директория в контейнере
WORKDIR /app

# Копируем скомпилированное приложение из первого этапа
COPY --from=builder /app/server /app/server
COPY --from=builder /app/config-yaml /app/config-yaml

# Открываем порт приложения
EXPOSE 8080

# Команда запуска приложения
CMD ["./server"]