FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV CGO=0 OS=linux ARCH=amd64

RUN CGO_ENABLED=$CGO GOOS=$OS GOARCH=$ARCH go build -ldflags '-extldflags "-static"' -o news_app cmd/app/main.go

# Run stage
FROM alpine:3.14

RUN apk update && apk upgrade && apk add --no-cache bash=5.1.16-r0 

ENV USER=docker GROUPNAME=dockergr UID=12345 GID=23456

WORKDIR /home/$USER/app

RUN addgroup --gid "$GID" "$GROUPNAME" \
    && adduser \
    --disabled-password \
    --gecos "" \
    --home "$(pwd)" \
    --ingroup "$GROUPNAME" \
    --no-create-home \
    --uid "$UID" \
    $USER

USER $USER

COPY --from=builder /app .

HEALTHCHECK NONE

CMD [ "./news_app" ]