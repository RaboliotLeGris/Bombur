FROM golang:1.16.2-buster@sha256:5a6302e91acb152050d661c9a081a535978c629225225ed91a8b979ad24aafcd as builder

WORKDIR /go/src/bombur
COPY . .

RUN go build -o build/bombur

FROM debian:10.8-slim@sha256:13f0764262a064b2dd9f8a828bbaab29bdb1a1a0ac6adc8610a0a5f37e514955

WORKDIR /bombur
COPY --from=builder /go/src/bombur/build/bombur /bombur
COPY --from=builder /go/src/bombur/static /bombur/static

CMD ["./bombur"]