
#FROM golang:1.21.6-alpine3.18 AS BuildStage
#RUN apk update && apk add git make && apk add postgresql-client # Добавляем установку make здесь
#WORKDIR /app
#COPY go.mod ./
#COPY go.sum ./
##COPY .env ./
#RUN go mod download
#COPY . .
#RUN make build
#
#FROM alpine:latest
#WORKDIR /app
#COPY --from=BuildStage /app/ /app/
##COPY --from=BuildStage /app/freshFriends freshFriends
###COPY --from=BuildStage /app/migrations migrations
##COPY --from=BuildStage /app/.env .env
##COPY --from=BuildStage /app/Makefile Makefile
##COPY --from=BuildStage /app/migrations migrations
#ENTRYPOINT ["/app/freshFriends"]

# Используйте официальный образ Golang как базовый для сборки приложения
FROM golang:1.21.6-alpine3.18 AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /backend ./cmd/sso/main.go

FROM alpine:3.18.0
WORKDIR /
COPY --from=builder /backend /backend
COPY --from=builder /app/.env /.env
COPY --from=builder /app/migrations /migrations
#EXPOSE 8080
CMD ["/backend"]
