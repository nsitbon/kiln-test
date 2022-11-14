FROM golang:1.19.3-alpine3.16 as builder
WORKDIR /project
COPY . .
ENV CGO_ENABLED=0
RUN --mount=type=cache,target=/root/.cache/go-build go build -ldflags "-s -w" -installsuffix cgo -o /app ./cmd/main.go

FROM scratch
ADD https://curl.se/ca/cacert.pem /etc/ssl/certs/
COPY --from=builder /app /
ENTRYPOINT ["/app"]