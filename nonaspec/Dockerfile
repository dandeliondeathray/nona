FROM python:3.6.1

WORKDIR /nonaspec
COPY requirements.txt /nonaspec
RUN pip install -r requirements.txt
COPY . /nonaspec

RUN mkdir /dist
RUN python setup.py bdist_egg --dist-dir /dist