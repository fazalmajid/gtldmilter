FROM alpine:latest

RUN mkdir /etc/postfix
RUN apk add postfix
RUN apk add go
RUN apk add git

RUN mkdir /src

ADD gtldmilter.go /src
RUN (cd /src;go mod init gtldmilter; go mod tidy; go build gtldmilter.go)

RUN ls -l /etc/postfix
RUN echo master.cf
RUN cat /etc/postfix/master.cf
RUN echo main.cf
RUN cat /etc/postfix/main.cf
ADD main.cf /etc/postfix/main.cf
ADD master.cf /etc/postfix/master.cf
ADD gtlds.bad /etc/postfix
ADD dests.bad /etc/postfix
ADD aliases /etc/postfix
RUN postmap /etc/postfix/aliases
RUN echo 'gtld:x:666:101:GTLD milter user:/var/spool/postfix:/bin/sh' >> /etc/passwd
RUN mkdir -p /var/spool/postfix/milter
RUN chown postfix:postfix /var/spool/postfix/milter
RUN chmod 775 /var/spool/postfix/milter
ADD run.sh /

EXPOSE 25
CMD ["/run.sh"]
