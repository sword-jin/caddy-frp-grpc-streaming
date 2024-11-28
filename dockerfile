FROM golang:1.23 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o grpc-server cmd/server/main.go

FROM gcr.io/distroless/base-debian11

WORKDIR /
COPY --from=builder /app/grpc-server .

ENTRYPOINT ["/grpc-server"]