FROM golang:1.16.3 as build

ARG TARGETPLATFORM
ARG BUILDPLATFORM

WORKDIR /api/bin
COPY ./bin ./

RUN if [ "$BUILDPLATFORM" = "linux/amd64" ]; then mv api.linux.amd64 api ; fi
RUN if [ "$BUILDPLATFORM" = "linux/arm64" ]; then mv api.linux.arm64 api ; fi


FROM --platform=$BUILDPLATFORM gcr.io/distroless/base

ARG TARGETPLATFORM
ARG BUILDPLATFORM

COPY --from=build /api/bin/api /app/api
ENTRYPOINT ["app/api"]
