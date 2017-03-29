FROM alpine:3.5

WORKDIR /app

ADD ./consulship .

ENTRYPOINT [ "./consulship" ]
