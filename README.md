# gtldmilter

A milter to add (source, destination) filtering of gTLDs to Postfix

## Problem statement

A few years back ICANN in its [infinite wisdom][1] (and the need to raise revenues since it is paid a fee for each domain name registered) created a bunch of new generic top-level  domains (gTLDs). Some are vanity domains like `.google` or `.bnpparibas`. Some are land grabs like `.app`, `.museum` or `.aero`.

I started recceiving huge amounts of spam from dodgy domains like `.cam` (apparently a domain created specifically to host online video strippers). My initial implementation was to write a script that:

1. fetched the list of all TLDs
2. filtered out known ones like `.com`, `.net`, `.org`, the country-specific (ccTLDs) and a whitelist of non-dodgy ones
3. created a Postfix map rejecting the others, not even accepting the connections

Unfortunately, a number of merchants and other legitimate senders started using these crackpot domains, and maintaining the whitelists started becoming a whack-a-mole exercise, and since it is reactive, the emails would be lost

I generate a different address for each vendor, e.g. Dell would get `dell@example.com` (not the actual domain, but you get the point). That way, when I started receiving pornographic spam addressed to that address (true story), I knew Dell's security was worth jack all. I could live with a scheme where if an email comes from a crackpot domain and is addressed to a non-vendor email address (those most likely to receive spam), it would be rejected.

Unfortunately, Postfix's native facilities do not allow this, but using the [milter][2] interface, you can build it. I used [pf-milters][4] as a guide.

## Installation

Just do a git checkout and run:

    go build gtldmilter.go

## Usage

You need to add the following entry to your Postfix `main.cf`:

    smtpd_milters = unix:milter/gtld

and run, using SMF, daemontools or (shudder) systemd the following program as a user in the `postfix` group:

    mkdir -p /var/spool/postfix/milter
    chown postfix:postfix /var/spool/postfix/milter

finally run, using SMF, daemontools or (shudder) systemd the `gtldmilter` program as a user in the `postfix` group.

[1]: https://www.theregister.com/2012/06/29/domain_land_grab_under_the_microscope/
[2]: http://www.postfix.org/MILTER_README.html
[3]: https://github.com/phalaaxx/milter
[4]: https://github.com/phalaaxx/pf-milters
