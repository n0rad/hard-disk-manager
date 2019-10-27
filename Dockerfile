FROM golang:1.13 as builder

WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -o /hdm

#####
FROM archlinux
RUN pacman -Sy --noconfirm smartmontools

#FROM debian
#RUN apt-get update &&apt-get install -y --no-install-recommends smartmontools

COPY --from=builder /hdm /hdm

CMD [ "/hdm", "agent", "-L", "debug" ]
