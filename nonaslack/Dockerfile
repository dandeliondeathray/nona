FROM python:3.6.1

ENV PYTHONUNBUFFERED=0

WORKDIR /app
COPY requirements.txt /app/
RUN pip install -r requirements.txt
COPY . /app/

RUN useradd -m user

CMD [ "python", "/app/main.py" ]