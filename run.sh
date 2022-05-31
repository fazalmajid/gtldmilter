#!/bin/sh
echo starting milter
su gtld -c /src/gtldmilter &

tail -F /var/spool/postfix/mail.log&

echo starting postfix
/usr/libexec/postfix/master

sleep 5
ps
wait
