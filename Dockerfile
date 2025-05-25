FROM golang:1.24.2-alpine AS build

RUN apk update --no-cache ca-certificates

WORKDIR /app

ADD ./ /app

ENV CGO_ENABLED=0

RUN go mod download
RUN go build -o main ./cmd/upgrader/main.go

FROM scratch

WORKDIR /

EXPOSE 8080

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/main /main

CMD ["/main"]