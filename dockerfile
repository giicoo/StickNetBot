FROM golang:1.20

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN apt-get update
RUN apt-get install -y libvips libvips-dev
RUN go get -u gopkg.in/h2non/bimg.v1
