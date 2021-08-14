# About

Pricelocator fetches the prices from given urls and sends a summary by email.

Due to the some websites restrictions like CloudFlare DDoS and/or non-domestic IPs blocking the suggested use is on personal PCs where the given urls are opened at least once in a browser instead of use on remote servers.

# Get Started

- `go build`
- set the following env variables in `$HOME/.profile`
```
PRICELOCATOR_MAIL_FROM
PRICELOCATOR_MAIL_PASS
PRICELOCATOR_MAIL_TO
PRICELOCATOR_MAIL_HOST
```
- install cronjob

`0 * * * * . /home/youruser/.profile; cd /home/youruser/path/to/pricelocator/ && /home/youruser/path/to/go run pricelocator >> /home/youruser/pricelocator.log 2>&1`