FROM golang:1.19.5-buster AS builder

ARG VERSION=dev

WORKDIR /go/src/app
COPY . .
COPY ./internal/ ./internal/
COPY ./tests/ ./tests/
RUN go build -o main

FROM debian:buster-slim
COPY --from=builder /go/src/app/main /go/bin/main
ENV PATH="/go/bin:${PATH}"
CMD ["main"]
