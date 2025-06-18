FROM alpine:latest

RUN apk update && apk upgrade
RUN apk add libgcc jansson libstdc++ zstd-libs binutils-gold
RUN wget https://go.dev/dl/go1.24.1.linux-arm64.tar.gz
RUN tar -C /usr/local -xzf go1.24.1.linux-arm64.tar.gz
RUN export PATH=$PATH:/usr/local/go/bin

RUN mkdir /home/database
WORKDIR /home/database

COPY . .

RUN /usr/local/go/bin/go mod tidy

ENTRYPOINT ["/usr/local/go/bin/go", "run", "."]
