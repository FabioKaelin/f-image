FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS build
LABEL org.opencontainers.image.source="https://github.com/FabioKaelin/f-image"
WORKDIR /src
ARG TARGETOS TARGETARCH
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/app .

FROM alpine  AS production
LABEL org.opencontainers.image.source="https://github.com/FabioKaelin/f-image"
COPY --from=build /out/app /bin/app

RUN mkdir -p /public

COPY ./public /public

RUN mkdir -p /public/dynamic
RUN mkdir -p /public/dynamic/profiles

CMD ["/bin/app"]
