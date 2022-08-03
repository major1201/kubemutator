ARG GO_VERSION=1.18.5
FROM golang:${GO_VERSION}-alpine as builder

WORKDIR /src
RUN apk --no-cache add ca-certificates git make
COPY . .
RUN make linux/amd64

FROM alpine:3.16.1
ENV TLS_CERT_FILE=/etc/kubemutator/kubemutator.crt \
    TLS_PRIVATE_KEY_FILE=/etc/kubemutator/kubemutator.key
RUN apk --no-cache add ca-certificates
COPY --from=builder /src/release/kubemutator* /bin/kubemutator
COPY ./examples/conf /etc/kubemutator
ENTRYPOINT [ "kubemutator" ]
VOLUME /etc/kubemutator
EXPOSE 443
