FROM golang:1.14.0 AS builder
WORKDIR /go/src/archives
COPY . /go/src/archives
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o bin .

FROM node:13 AS assetsbuilder
WORKDIR /go/src/archives
COPY . /go/src/archives
RUN npm install && cd node_modules/@gentoo/tyrian && npm install && npm run dist && cd /go/src/archives
RUN npx webpack

FROM scratch
WORKDIR /go/src/archives
COPY --from=assetsbuilder /go/src/archives/assets /go/src/archives/assets
COPY --from=builder /go/src/archives/bin /go/src/archives/bin
COPY --from=builder /go/src/archives/pkg /go/src/archives/pkg
COPY --from=builder /go/src/archives/web /go/src/archives/web
ENTRYPOINT ["/go/src/archives/bin/archives", "serve"]
