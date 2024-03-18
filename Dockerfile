FROM golang:1.21

WORKDIR $GOPATH/src/github.com/GeorgeGorbanev/vibeshare

COPY go.mod go.sum ./
COPY cmd ./cmd
COPY internal ./internal
RUN go mod download
RUN go build -o /bin/ ./cmd/vibeshare
CMD ["/bin/vibeshare"]