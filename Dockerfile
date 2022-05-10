# build
FROM golang:alpine as build
RUN apk add --no-cache ca-certificates git
WORKDIR /crawler
COPY go.mod go.sum ./
RUN go mod download
COPY src/ ./src
WORKDIR /crawler/src
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o app

# image
FROM scratch
WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /crawler/src/app .
EXPOSE 8080
CMD ["./app"]