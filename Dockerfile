FROM node:20.14.0 as web-builder
WORKDIR /app

COPY ./web/package.json web/package-lock.json ./
RUN npm ci

COPY ./web/ ./
RUN npm run build

FROM golang:1.22
ARG PORT
WORKDIR /app

COPY --from=web-builder /app/build ./web/build

COPY go.mod go.sum ./
RUN go mod download
COPY ./ .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./main ./cmd/api

EXPOSE ${PORT}
CMD ["./main"]