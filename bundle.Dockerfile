# Different source images for each platform
# Stolen from https://github.com/docker/buildx/discussions/1928
FROM --platform=linux/amd64 ghcr.io/gotify/server AS build-amd64
FROM --platform=linux/arm64 ghcr.io/gotify/server-arm64 AS build-arm64
FROM --platform=linux/arm/v7 ghcr.io/gotify/server-arm7 AS build-arm

# Run the actual build for the specific platform
FROM build-$TARGETARCH
# NOTE: Need to re-declare this to make available _inside_ image build
ARG TARGETARCH
COPY out/gotify-slack-webhook-${TARGETARCH}.so /app/data/plugins/
