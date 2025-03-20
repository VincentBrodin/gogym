FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./

RUN go test -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o /backend-app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /backend-app /app/backend-app
COPY /bin/content /app/content
COPY /backend/.env /app/.env

EXPOSE 8080
CMD ["/app/backend-app"]
