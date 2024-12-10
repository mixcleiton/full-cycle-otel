FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o server .

FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/server .
COPY .env .

EXPOSE 8081

CMD ["./server"]