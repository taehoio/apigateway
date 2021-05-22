FROM golang:1.16.3 as build

ARG TARGETPLATFORM
ARG BUILDPLATFORM

WORKDIR /apigateway/bin
COPY ./bin ./

RUN if [ "$BUILDPLATFORM" = "linux/amd64" ]; then mv apigateway.linux.amd64 apigateway; fi
RUN if [ "$BUILDPLATFORM" = "linux/arm64" ]; then mv apigateway.linux.arm64 apigateway; fi


FROM --platform=$BUILDPLATFORM gcr.io/distroless/base

ARG TARGETPLATFORM
ARG BUILDPLATFORM

ENV ENV=development
ENV IS_GRPC_INSECURE=true
ENV CERT_FILE=/etc/ssl/certs/ca-certificates.crt
ENV BAEMINCRYPTO_GRPC_SERVICE_ENDPOINT=localhost:50051

COPY --from=build /apigateway/bin/apigateway /app/apigateway

EXPOSE 8080

ENTRYPOINT ["app/apigateway"]
