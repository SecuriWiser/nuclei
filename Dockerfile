FROM golang:1.20.0-alpine as build-env
RUN apk add build-base
WORKDIR /app
COPY . .
RUN cd v2/cmd/nuclei/ && go build -o /app/nuclei main.go

FROM alpine:3.17.1
RUN apk add --no-cache bind-tools ca-certificates chromium
COPY --from=build-env /app/nuclei /usr/local/bin/nuclei
ENTRYPOINT ["nuclei"]
