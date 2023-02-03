FROM golang:1.18-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY * ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o dt-runner-amd64

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /dt-runner-amd64 /dt-runner

EXPOSE 9001

USER nonroot:nonroot

ENTRYPOINT ["/dt-runner"]
