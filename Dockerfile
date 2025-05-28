FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o fieltorcedorbot ./cmd/bot

FROM alpine:3.18

RUN apk --no-cache add ca-certificates tzdata

ENV TZ=America/Sao_Paulo

WORKDIR /app

COPY --from=builder /app/fieltorcedorbot .

CMD ["./fieltorcedorbot"]
