FROM golang:1.14.4 AS build
WORKDIR /go/src/github.com/amikhalitsin/SIPServer
COPY app .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /go/src/github.com/amikhalitsin/SIPServer/app .
COPY --from=build /go/src/github.com/amikhalitsin/SIPServer/regs .
CMD ["./app"]
