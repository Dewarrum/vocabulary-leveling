FROM node:20.14.0 as web-builder
WORKDIR /app

COPY ./web/package.json web/package-lock.json ./
RUN npm ci

COPY ./web/ ./
RUN npm run build

FROM golang:1.22 as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY ./ .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./main ./cmd/api

FROM alpine:3.20.1
ARG PORT
WORKDIR /app

COPY ./db/migrations ./db/migrations
COPY --from=web-builder /app/build ./web/build
COPY --from=builder /app/main ./main

EXPOSE ${PORT}
CMD ["./main"]