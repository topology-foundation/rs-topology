# We use builder pattern to optimize the image size.
FROM golang:1.22-alpine AS build

WORKDIR /gram

ENV BUILD_PACKAGES make
RUN apk add --no-cache $BUILD_PACKAGES

# Optimize to use cached image if go.mod and go.sum didn't change.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make gram

# A runtime image has only binary and config file from a build image.
FROM alpine:latest

WORKDIR /root
EXPOSE 1319

COPY --from=build /gram/bin/gram /usr/bin/gram
COPY --from=build /gram/config.json /root/config.json

ENTRYPOINT ["gram"]
