FROM golang:alpine as build
WORKDIR /go/src/myapp
COPY go.mod .
COPY go.sum .

RUN go mod download

ENV CGO_ENABLED=0

COPY . .

RUN go build -o /go/bin/myapp ./cmd/web

FROM golang:alpine

RUN apk add --no-cache poppler

WORKDIR /go/bin

COPY --from=build /go/bin/myapp .

ENTRYPOINT ["./myapp", "serve"]
