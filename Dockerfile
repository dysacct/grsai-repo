FROM golang:1.25.10-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY config ./config
COPY grsai ./grsai
COPY handler ./handler
COPY model ./model
COPY router ./router
COPY storage ./storage

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/grsai-newapi .

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /out/grsai-newapi /app/grsai-newapi

ENV SERVER_PORT=3456
EXPOSE 3456

CMD ["/app/grsai-newapi"]
