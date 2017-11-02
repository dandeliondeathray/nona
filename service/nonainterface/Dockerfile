FROM erikedin/librdkafka:0.11.1 as librdkafkabin

FROM python:3.6.1
LABEL maintainer="Erik Edin"

ARG LIBRDKAFKA_VERSION=0.11.1

RUN mkdir /dist
COPY --from=librdkafkabin /librdkafka-${LIBRDKAFKA_VERSION}-bin.tar.gz /dist
WORKDIR /
RUN tar xf /dist/librdkafka-${LIBRDKAFKA_VERSION}-bin.tar.gz

RUN mkdir /nonainterface
COPY requirements.txt /nonainterface/
RUN pip install -r /nonainterface/requirements.txt
COPY . /nonainterface

WORKDIR /nonainterface
RUN mkdir -p /egg/
RUN python setup.py bdist_egg --dist-dir /dist
RUN python setup.py install
