# Initial stage: download modules
FROM golang:1.15 as modules
LABEL stage=builder

ADD go.mod go.sum /m/
RUN cd /m && go mod download

# Intermediate stage: Build the binary
FROM golang:1.15 as builder
LABEL stage=builder

COPY --from=modules /go/pkg /go/pkg

# add a non-privileged user
RUN useradd -u 10001 tinyimage

RUN mkdir -p /tinyimage
ADD . /tinyimage
WORKDIR /tinyimage

# Build the binary with go build
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -ldflags="-w -s" -o ./bin/tinyimage ./

# Final stage: Run the binary
FROM alpine

COPY --from=builder /etc/passwd /etc/passwd

RUN apk add --update \
    bash \
    lcms2-dev \
    libpng-dev \
    gcc \
    g++ \
    make \
    autoconf \
    automake \
    pngquant \
    wget \
    nasm \
  && rm -rf /var/cache/apk/*

RUN wget https://github.com/mozilla/mozjpeg/releases/download/v3.2-pre/mozjpeg-3.2-release-source.tar.gz && \
    tar -xvzf mozjpeg-3.2-release-source.tar.gz && \
    rm mozjpeg-3.2-release-source.tar.gz && \
    cd mozjpeg && \
    ./configure && \
    make install && \
    cd / && rm -rf mozjpeg && \
    ln -s /opt/mozjpeg/bin/jpegtran /usr/local/bin/jpegtran && \
    ln -s /opt/mozjpeg/bin/cjpeg /usr/local/bin/cjpeg

# and finally the binary
COPY --from=builder /tinyimage /tinyimage

WORKDIR /tinyimage

RUN chown tinyimage templates && chown tinyimage images

USER tinyimage

CMD ["./bin/tinyimage"]