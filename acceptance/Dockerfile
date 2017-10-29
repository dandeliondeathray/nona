FROM erikedin/nonainterface:latest as nonainterface_builder
FROM erikedin/nonaspec:latest as nonaspec_builder

FROM python:3.6.1

ARG LIBRDKAFKA_VERSION=0.11.1

WORKDIR /dist
COPY --from=nonainterface_builder /dist/nonainterface-0.1.0-py3.6.egg /dist
COPY --from=nonainterface_builder /dist/librdkafka-${LIBRDKAFKA_VERSION}-bin.tar.gz /dist
COPY --from=nonaspec_builder /dist/nonaspec-0.1-py3.6.egg /dist

WORKDIR /
RUN tar xf /dist/librdkafka-${LIBRDKAFKA_VERSION}-bin.tar.gz
RUN easy_install /dist/*.egg

WORKDIR /acceptance
COPY requirements.txt /acceptance
RUN pip install -r requirements.txt
COPY . /acceptance

CMD [ "behave" ]