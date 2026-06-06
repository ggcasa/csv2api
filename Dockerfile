FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o server ./cmd/server

FROM scratch
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/cars.db . 
EXPOSE 9977
CMD ["./server"]