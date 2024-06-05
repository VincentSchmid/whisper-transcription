FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go main.go
COPY pkg/ pkg/
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

FROM alpine:edge
WORKDIR /app
COPY --from=build /app/myapp .

RUN apk --no-cache add ca-certificates tzdata ffmpeg

ENTRYPOINT ["/app/myapp"]