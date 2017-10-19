FROM python:3.6.1

ARG VERSION=0.11.1

ADD https://github.com/edenhill/librdkafka/archive/v${VERSION}.tar.gz /tmp/
WORKDIR /tmp/librdkafka
RUN tar xf ../v${VERSION}.tar.gz
WORKDIR /tmp/librdkafka/librdkafka-0.11.1
RUN ./configure --prefix=/usr
RUN make

RUN mkdir -p /tmp/install
RUN make install DESTDIR=/tmp/install

WORKDIR /
RUN tar -C /tmp/install -czf librdkafka-${VERSION}-bin.tar.gz usr