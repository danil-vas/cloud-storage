# первый stage - сборка проекта
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /docker-go-cloud-storage .

# второй stage - запуск приложения в чистой среде
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/docker-go-cloud-storage /docker-go-cloud-storage
RUN echo 'nobody:x:65534:65534:nobody:/:' > /etc/passwd
USER nobody:nogroup
EXPOSE 8080
ENTRYPOINT ["/docker-go-cloud-storage"]