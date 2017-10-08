FROM arm32v7/python:3.6.1
LABEL maintainer="Erik Edin"

RUN mkdir /nonainterface
COPY requirements.txt /nonainterface/
RUN pip install -r /nonainterface/requirements.txt
COPY . /nonainterface

WORKDIR /nonainterface
RUN python setup.py install
