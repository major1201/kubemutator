ARG GO_VERSION=1.12.6
FROM golang:${GO_VERSION}-alpine

RUN apk --no-cache add ca-certificates git

WORKDIR /src
COPY . .
RUN go build

FROM alpine:3.9.2
ENV TLS_CERT_FILE=/etc/k8s-mutator/k8s-mutator.crt \
    TLS_PRIVATE_KEY_FILE=/etc/k8s-mutator/k8s-mutator.key
RUN apk --no-cache add ca-certificates
COPY --from=0 /src/k8s-mutator /bin/k8s-mutator
COPY ./examples/conf /etc/k8s-mutator
ENTRYPOINT [ "k8s-mutator" ]
VOLUME /etc/k8s-mutator
EXPOSE 443
