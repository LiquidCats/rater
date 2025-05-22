FROM golang:1.24.1-alpine AS build

RUN apk update --no-cache ca-certificates

WORKDIR /app

ADD ./ /app

RUN go mod download
RUN go build -o main ./cmd/rater/main.go

FROM scratch

WORKDIR /

EXPOSE 8080

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/main /main

CMD ["/main"]



