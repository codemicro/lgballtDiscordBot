FROM golang:1.16 as builder
RUN go get -u github.com/hhatto/gocloc/cmd/gocloc
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN gocloc --output-type=json . > internal/buildInfo/jdat
RUN date > internal/buildInfo/currentDate
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static" -s -w' -o main github.com/codemicro/lgballtDiscordBot/cmd/lgballtDiscordBot
FROM alpine
COPY --from=builder /build/main /
WORKDIR /run
CMD ["../main"]
