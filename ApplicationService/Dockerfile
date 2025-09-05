# Первый этап: сборка бинаря
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем модули и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Копируем весь исходный код
COPY . .

# Сборка только API Gateway
RUN go build -o application-service ./cmd/server/main.go

# Второй этап: минимальный образ для запуска
FROM alpine:latest

WORKDIR /app

# Добавляем сертификаты для https-запросов
RUN apk add --no-cache ca-certificates

# Копируем только бинарь и конфиг
COPY --from=builder /app/application-service /app/application-service
COPY --from=builder /app/config-yaml /app/config-yaml

# Открываем порт 
EXPOSE 8082

# Запускаем бинарь
CMD ["./application-service"]
