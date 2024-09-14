FROM golang:1.22.4-alpine3.19 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp ./cmd

FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/myapp .
COPY .env /app
EXPOSE 3030
ENTRYPOINT [ "/app/myapp" ]