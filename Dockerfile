FROM golang:1.14-alpine as builder

RUN apk add git

WORKDIR /app
COPY . ./
RUN ./gomake build -L debug && cp dist/hdm-linux-amd64/hdm /hdm

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
