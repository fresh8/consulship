FROM alpine:3.5

RUN mkdir /app
WORKDIR /app

ADD ./consulship /usr/bin/

ENTRYPOINT [ "/usr/bin/consulship" ]
