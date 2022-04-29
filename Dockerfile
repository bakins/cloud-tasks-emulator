FROM golang:1.17-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV CGO_ENABLED=0

RUN go build -o emulator ./cmd/emulator


FROM gcr.io/distroless/static

LABEL org.opencontainers.image.source=https://github.com/aertje/cloud-tasks-emulator

ENTRYPOINT ["/emulator"]

WORKDIR /

COPY --from=builder /app/emulator .
