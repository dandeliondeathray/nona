FROM arm32v7/openjdk:8-jre
LABEL maintainer="Erik Edin"

RUN useradd -m user

ADD kafka_2.11-0.11.0.1.tgz /tmp/
RUN mv /tmp/kafka_2.11-0.11.0.1 /kafka/
COPY zookeeper-server-start.sh /kafka/bin/

RUN mkdir /logs
RUN chown -R user:user /logs

# This is used by /kafka/bin/kafka-run-class.sh to set the log directory.
ENV LOG_DIR=/logs

USER user
