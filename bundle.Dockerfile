ARG SERVER_VERSION=latest
FROM ghcr.io/gotify/server:${SERVER_VERSION}

# NOTE: Additional plugins will ALSO need to go there, or they won't be detected!
ENV GOTIFY_PLUGINSDIR=/opt/plugins

# NOTE: Need to re-declare this to make available _inside_ image build
ARG TARGETARCH
COPY out/gotify-slack-webhook-linux-${TARGETARCH}.so $GOTIFY_PLUGINSDIR
