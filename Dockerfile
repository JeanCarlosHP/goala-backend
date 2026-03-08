FROM golang:1.26.0-alpine AS build

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

COPY . .

RUN go build -ldflags="-s -w" -o app ./cmd/api/main.go
RUN apk add --no-cache ca-certificates

FROM scratch AS final

ENV GO_ENV=production

WORKDIR /app

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=build ./app/app ./app

EXPOSE 8080

USER 1001

ENTRYPOINT ["./app"]