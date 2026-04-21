FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/admin ./cmd/gatus-admin

FROM golang:1.25-alpine AS dev
WORKDIR /src
RUN apk add --no-cache docker-cli
RUN go install github.com/air-verse/air@latest
EXPOSE 8000
CMD ["air", "-c", ".air.toml"]

FROM golang:1.25-alpine AS hashpw-build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/hashpw ./cmd/hashpw

FROM alpine:3.20 AS prod
RUN apk add --no-cache docker-cli ca-certificates
COPY --from=build /out/admin /admin
COPY --from=hashpw-build /out/hashpw /hashpw
EXPOSE 8000
ENTRYPOINT ["/admin"]
