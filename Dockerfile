FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_DATE=unknown
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath \
    -ldflags="-s -w -extldflags=-static \
    -X 'github.com/Norgate-AV/genlinx/internal/version.Version=${VERSION}' \
    -X 'github.com/Norgate-AV/genlinx/internal/version.GitCommit=${GIT_COMMIT}' \
    -X 'github.com/Norgate-AV/genlinx/internal/version.BuildDate=${BUILD_DATE}'" \
    -o /genlinx .

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

COPY --from=builder /genlinx /usr/local/bin/genlinx

USER nobody
ENTRYPOINT ["genlinx"]
