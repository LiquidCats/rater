FROM golang:1.24.1-alpine AS build

WORKDIR /app

ADD ./ /app

RUN go build -o main ./cmd/rater/main.go

FROM scratch

EXPOSE 8080

COPY --from=build /app/main /main

CMD ["/main"]



