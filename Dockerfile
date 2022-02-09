FROM golang:alpine as builder
RUN apk update && apk add --no-cache git
COPY . $GOPATH/src/webwol/
WORKDIR $GOPATH/src/webwol/
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /go/bin/main .
FROM scratch
COPY --from=builder /go/bin/main /app/
WORKDIR /app
CMD ["./main"]
