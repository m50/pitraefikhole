FROM golang:1.24.3 AS go-builder

WORKDIR /app

RUN apt-get update \
    && apt-get install -y git gcc

COPY ./ ./
RUN go mod download 

RUN go build -o main ./cmd/sync/main.go
RUN chmod +x ./main

FROM gcr.io/distroless/cc-debian12

WORKDIR /app

COPY --from=go-builder /app/main ./

USER 1000:1000
CMD ["./main", "sync"]

