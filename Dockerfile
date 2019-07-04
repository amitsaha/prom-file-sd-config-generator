FROM golang:1.12 as build

WORKDIR /go/src/prom-file-sd-config-generator
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

FROM gcr.io/distroless/base
COPY --from=build /go/bin/prom-file-sd-config-generator /
ENTRYPOINT ["/prom-file-sd-config-generator"]