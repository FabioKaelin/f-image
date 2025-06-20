FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS build
LABEL org.opencontainers.image.source="https://github.com/FabioKaelin/f-image"
WORKDIR /src
ARG TARGETOS TARGETARCH
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/app .

FROM alpine  AS production
LABEL org.opencontainers.image.source="https://github.com/FabioKaelin/f-image"
LABEL org.opencontainers.image.authors="FabioKaelin"
LABEL org.opencontainers.image.title="f-image"
COPY --from=build /out/app /bin/app

RUN mkdir -p /public

COPY ./public /public

RUN mkdir -p /public/dynamic
RUN mkdir -p /public/dynamic/profiles
RUN mkdir -p /public/dynamic/gallery

CMD ["/bin/app"]
