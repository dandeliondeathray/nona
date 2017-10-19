FROM erikedin/nonainterface:staging

WORKDIR /acceptance
COPY requirements.txt /acceptance
RUN pip install -r requirements.txt
COPY . /acceptance
RUN behave