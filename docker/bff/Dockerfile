FROM golang:1.24-alpine3.21 AS base
WORKDIR /base
COPY ../../src/bff/go.mod .
RUN go mod download
COPY ../../src/bff .
RUN go build -o ./main ./main.go

FROM alpine:3.21
COPY --from=base /base/main /main
CMD ["/main"]