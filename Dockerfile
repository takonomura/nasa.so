FROM docker.io/library/golang:1.19.1-alpine@sha256:ca4f0513119dfbdc65ae7b76b69688f0723ed00d9ecf9de68abbf6ed01ef11bf AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY * ./
RUN CGO_ENABLED=0 go build -o nasaso .

FROM docker.io/library/alpine:3.16.2@sha256:1304f174557314a7ed9eddb4eab12fed12cb0cd9809e4c28f29af86979a3c870

COPY --from=build /app/nasaso /nasaso

ENTRYPOINT ["/nasaso"]
