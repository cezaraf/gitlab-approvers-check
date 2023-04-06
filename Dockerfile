FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/gac

FROM scratch

COPY --from=builder /go/bin/gac /go/bin/gac

ENTRYPOINT ["/go/bin/gac"]