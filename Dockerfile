FROM alpine:3.5

RUN mkdir /app
WORKDIR /app

ADD ./consulship .

ENTRYPOINT [ "./consulship" ]
