# Build
FROM golang:1.20.3-alpine AS build-env
RUN apk add build-base
WORKDIR /app
COPY . /app
WORKDIR /app/v2
RUN go mod download
RUN go build ./cmd/nuclei

# Release
FROM alpine:3.17.3
RUN apk -U upgrade --no-cache \
    && apk add --no-cache bind-tools chromium ca-certificates git
COPY --from=build-env /app/v2/nuclei /usr/local/bin/
RUN git clone https://github.com/projectdiscovery/nuclei-templates ~/nuclei-templates

ENTRYPOINT ["nuclei"]
