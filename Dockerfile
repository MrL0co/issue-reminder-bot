FROM golang:alpine AS build-env

LABEL maintainer="Richard Jacobse <beheer.loco@gmail.com>"

RUN apk add --update --no-cache ca-certificates git

RUN mkdir /app 

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/main ./server

RUN chmod a+x /go/bin/main

FROM scratch

COPY --from=build-env /go/bin/main /go/bin/main

EXPOSE 8080

ENTRYPOINT ["/go/bin/main"]
