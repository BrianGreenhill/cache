FROM golang:1.23-alpine AS build

WORKDIR /app
COPY go.* /app

RUN go mod download

COPY *.go /app

RUN go build -o /app/cash

FROM scratch

COPY --from=build /app/cash /app/cash

CMD ["/app/cash"]
