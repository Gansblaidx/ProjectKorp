# Stage 1: Build
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod tidy && \
    go build -o http-server-projeto .

# Stage 2: Runtime mínimo
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/http-server-projeto .

# Expõe porta 8080 (documentação, não expõe automaticamente)
# Acesso externo é feito via NGINX em porta 80
EXPOSE 8080
ENTRYPOINT ["./http-server-projeto"]
