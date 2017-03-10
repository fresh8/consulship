FROM alpine:3.5

ADD ./consulship .

ENTRYPOINT [ "./consulship" ]
