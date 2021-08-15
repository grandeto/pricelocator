# About

Pricelocator fetches the prices from given urls and sends a summary by email.

Due to the some websites restrictions like CloudFlare DDoS and/or non-domestic IPs blocking the suggested use is on personal PCs where the given urls are opened at least once in a browser instead of use on remote servers.

# Get Started

- edit urls map in `func main()` according to your needs
- navigate to `pricelocator` dir and run `go build`
- set the following env variables (e.g. in `$HOME/.profile`)
```
export PRICELOCATOR_MAIL_FROM="yourmail@gmail.com"
export PRICELOCATOR_MAIL_PASS="yourmailpass"
export PRICELOCATOR_MAIL_TO="yournotificationmail@gmail.com"
export PRICELOCATOR_MAIL_HOST="smtp.gmail.com"
```
- create a cronjob

Linux example:

`0 * * * * . /home/youruser/.profile; cd /home/youruser/path/to/pricelocator/ && /home/youruser/path/to/go run pricelocator >> /home/youruser/pricelocator.log 2>&1`