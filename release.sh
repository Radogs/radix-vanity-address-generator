docker buildx build \
  --platform linux/arm64,linux/amd64,linux/arm/v7,linux/arm/v6 \
  -t radogs/radix-vanity-address-generator:latest \
  --push \
  .