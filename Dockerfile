FROM docker.io/library/golang:1.22.2-alpine@sha256:cdc86d9f363e8786845bea2040312b4efa321b828acdeb26f393faa864d887b0 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY * ./
RUN CGO_ENABLED=0 go build -o nasaso .

FROM docker.io/library/alpine:3.16.2@sha256:1304f174557314a7ed9eddb4eab12fed12cb0cd9809e4c28f29af86979a3c870

COPY --from=build /app/nasaso /nasaso

ENTRYPOINT ["/nasaso"]
