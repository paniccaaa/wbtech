FROM golang:1.23.2-alpine AS builder
# c lib
RUN apk add --no-cache alpine-sdk

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# from stackoverflow 
RUN GOOS=linux GOARCH=amd64 go build -tags musl -o ./bin/app cmd/wbtech/main.go

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /build/bin/app /app/app
COPY ./config/docker.yaml /app/docker.yaml  

EXPOSE 8089
CMD ["ls", "/app/app"]
