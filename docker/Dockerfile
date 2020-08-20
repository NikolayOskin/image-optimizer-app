FROM golang:latest

RUN apt-get update && apt-get install -y --no-install-recommends g++ make autoconf automake libtool nasm wget pngquant=2.12.2-1

RUN wget https://github.com/mozilla/mozjpeg/releases/download/v3.2-pre/mozjpeg-3.2-release-source.tar.gz && \
    tar -xvzf mozjpeg-3.2-release-source.tar.gz && \
    rm mozjpeg-3.2-release-source.tar.gz && \
    cd mozjpeg && \
    ./configure && \
    make install && \
    cd / && rm -rf mozjpeg && \
    ln -s /opt/mozjpeg/bin/jpegtran /usr/local/bin/jpegtran && \
    ln -s /opt/mozjpeg/bin/cjpeg /usr/local/bin/cjpeg

WORKDIR /app

