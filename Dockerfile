FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS build
LABEL org.opencontainers.image.source https://github.com/fabiokaelin/f-image
WORKDIR /src
ARG TARGETOS TARGETARCH
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/app .

FROM alpine
COPY --from=build /out/app /bin/app

CMD ["/bin/app"]
