FROM golang:1.24.3 AS go-builder

WORKDIR /app

RUN apt-get update \
    && apt-get install -y git gcc \
    && rm -rf /var/lib/apt/lists/*

COPY ./ ./
RUN go mod download 

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/...

FROM gcr.io/distroless/cc-debian12

WORKDIR /app

COPY --from=go-builder /app/main ./

EXPOSE 1323
USER nonroot:nonroot
CMD ["./main"]

