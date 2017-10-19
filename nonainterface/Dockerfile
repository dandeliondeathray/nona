FROM erikedin/librdkafka:0.11.1 as librdkafkabin

FROM python:3.6.1
LABEL maintainer="Erik Edin"

ARG LIBRDKAFKA_VERSION=0.11.1
COPY --from=librdkafkabin /librdkafka-${LIBRDKAFKA_VERSION}-bin.tar.gz /tmp
WORKDIR /
RUN tar xf /tmp/librdkafka-${LIBRDKAFKA_VERSION}-bin.tar.gz

RUN mkdir /nonainterface
COPY requirements.txt /nonainterface/
RUN pip install -r /nonainterface/requirements.txt
COPY . /nonainterface

WORKDIR /nonainterface
RUN python setup.py install
