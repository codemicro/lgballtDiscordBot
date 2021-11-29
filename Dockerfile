FROM golang:1 as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static" -s -w' -tags sqlite_omit_load_extension -o main github.com/codemicro/lgballtDiscordBot/cmd/lgballtDiscordBot
FROM alpine
COPY --from=builder /build/main /
WORKDIR /run
LABEL com.centurylinklabs.watchtower.stop-signal="SIGINT"
CMD ["../main"]
