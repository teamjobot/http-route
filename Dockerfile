FROM golang:1.13 as build

WORKDIR /go/src/github.com/teamjobot/http-route/
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o http-route .

FROM alpine:3.12.0
RUN apk --no-cache add ca-certificates
COPY --from=build /go/src/github.com/teamjobot/http-route/http-route /usr/local/bin/http-route
CMD http-route
EXPOSE 80
