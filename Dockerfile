FROM golang:1.17.5-alpine AS builder
RUN mkdir /build
ADD go.mod go.sum main.go template.html /build/
WORKDIR /build
RUN go build

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/my-form /app/
COPY --from=builder /build/template.html /app/
WORKDIR /app
CMD ["./my-form"]
