ARG PHP_VERSION

FROM ghcr.io/allincart-org/allincart-cli-base:${PHP_VERSION}

COPY allincart-cli /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/allincart-cli"]
CMD ["--help"]
