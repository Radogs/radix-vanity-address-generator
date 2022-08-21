FROM golang:1.19-bullseye as golang

RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid 65532 \
  radley

WORKDIR $GOPATH/src/radogs/radix-vanity-address-generator
COPY . .

RUN go mod download
RUN go mod verify

RUN mkdir /service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /service/radix-vanity .

FROM scratch

COPY --from=golang /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=golang /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=golang /etc/passwd /etc/passwd
COPY --from=golang /etc/group /etc/group
COPY --from=golang /service .

USER radley:radley

ENTRYPOINT ["./radix-vanity"]