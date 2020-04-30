FROM golang:1.14-alpine as builder

WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -o /hdm

#####
FROM alpine
RUN apk add \
    bash \
    smartmontools \
    util-linux \
    rsync \
    sgdisk \
    cryptsetup \
    device-mapper \
    hdparm \
    udev

COPY --from=builder /hdm /usr/bin/hdm

CMD [ "hdm", "agent" ]
