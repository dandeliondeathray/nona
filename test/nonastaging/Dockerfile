ARG INTERFACE_TAG=staging
FROM erikedin/nonainterface:${INTERFACE_TAG}

ENV PYTHONUNBUFFERED=0

WORKDIR /app
COPY requirements.txt /app/
RUN pip install -r requirements.txt
COPY . /app/

RUN useradd -m user

CMD [ "python", "/app/nonastaging.py" ]