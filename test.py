#!/usr/bin/env python3
import smtplib

def one(src, dst):
  s = smtplib.SMTP(host="localhost", port=2525)
  s.ehlo('test.py')
  print(s.ehlo_msg)
  s.mail(src)
  s.rcpt(dst)
  s.close()

one('sopo@sopo.cam', 'bad1@example.com')
one('sopo@sopo.cam', 'bad2@example.com')
one('sopo@sopo.cam', 'good@example.com')
