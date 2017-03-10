FROM alpine:3.5

ADD ./consulship /usr/local/bin

ENTRYPOINT [ "/usr/local/bin/consulship" ]
