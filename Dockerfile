FROM golang:alpine as build-env

RUN mkdir /plug
WORKDIR /plug
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/plug

FROM scratch

ENV USER_UID=1001 \
  USER_NAME=plug

COPY --from=build-env /go/bin/plug /usr/local/bin/plug
ENTRYPOINT ["/usr/local/bin/plug"]

# switch to non-root user
USER ${USER_UID}

