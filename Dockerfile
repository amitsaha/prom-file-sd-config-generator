FROM golang:1.12 as build

WORKDIR /
COPY . .
RUN go build

FROM gcr.io/distroless/base
COPY --from=build /prom-file-sd-config-generator /
ENTRYPOINT ["/prom-file-sd-config-generator"]