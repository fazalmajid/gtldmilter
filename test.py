#!/usr/bin/env python3
import smtplib, email.mime.text, subprocess

def one(src, dst, expect_failure=False):
  s = smtplib.SMTP(host="localhost", port=2525)
  msg = email.mime.text.MIMEText('Sopo la pougne')
  msg['Subject'] = 'Sopo'
  msg['From'] = src
  msg['To'] = dst
  try:
    s.sendmail(src, [dst], msg.as_string())
    assert not expect_failure
  except smtplib.SMTPRecipientsRefused:
    assert expect_failure
  s.close()

one('sopo@sopo.cam', 'bad1@example.com', True)
one('sopo@sopo.cam', 'bad2@example.com', True)
one('sopo@sopo.cam', 'good@example.com')

out = subprocess.check_output('grype postfix-milter-test', shell=True)
assert out == b'No vulnerabilities found\n'

print('\033[1;32m%s\033[0m' % 'All tests passed')

