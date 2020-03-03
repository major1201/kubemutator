ARG GO_VERSION=1.12.6
FROM golang:${GO_VERSION}-alpine

RUN apk --no-cache add ca-certificates git

WORKDIR /src
COPY . .
RUN go build

FROM alpine:3.9.2
ENV TLS_CERT_FILE=/etc/kubemutator/kubemutator.crt \
    TLS_PRIVATE_KEY_FILE=/etc/kubemutator/kubemutator.key
RUN apk --no-cache add ca-certificates
COPY --from=0 /src/kubemutator /bin/kubemutator
COPY ./examples/conf /etc/kubemutator
ENTRYPOINT [ "kubemutator" ]
VOLUME /etc/kubemutator
EXPOSE 443
