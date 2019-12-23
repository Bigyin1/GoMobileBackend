FROM golang:latest AS builder
ENV SRC_DIR=/app/backend
WORKDIR $SRC_DIR
COPY go.mod .
COPY go.sum .
RUN go mod download
ADD . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/main cmd/main/main.go

FROM alpine:latest
ENV SRC_DIR=/app/backend
COPY --from=builder $SRC_DIR/bin/main ./
COPY --from=builder $SRC_DIR/local-config.json ./
COPY --from=builder $SRC_DIR/pkg/controllers/mail/templates pkg/controllers/mail/templates
COPY --from=builder $SRC_DIR/pkg/controllers/mail/token pkg/controllers/mail/token
RUN chmod +x ./main
ENTRYPOINT ["./main"]
