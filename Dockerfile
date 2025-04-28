FROM golang:1.24 AS src

WORKDIR /go/src/app/

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

ENV CGO_ENABLED=0

RUN go build -o ./server ./cmd/server/main.go

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=src /go/src/app/server /usr/bin/local/server

EXPOSE 3000

ENTRYPOINT ["/usr/bin/local/server"]
