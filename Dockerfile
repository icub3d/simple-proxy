FROM golang AS builder

WORKDIR /
COPY . .
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o simple-proxy

FROM scratch

COPY --from=builder /simple-proxy /simple-proxy
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT [ "/simple-proxy" ]