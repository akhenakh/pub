# pub, a tiny ActivityPub to Mastodon bridge
    
[![Go Reference](https://pkg.go.dev/badge/github.com/davecheney/pub.svg)](https://pkg.go.dev/github.com/davecheney/pub) [![Go Report Card](https://goreportcard.com/badge/github.com/davecheney/pub)](https://goreportcard.com/report/github.com/davecheney/pub)
    
## What is pub?

`pub` is an ActivityPub host indented for a single actor.
To interact with ActivityPub, `pub` implements the Mastodon[^tm] api for use with various apps. 

`pub` is not intended to host an ActivityPub community, rather it is aimed at enabling someone who owns their own domain, and thus controls their identity, to participate in the fediverse. 

[^tm]: Mastodon is a trademark of [Mastodon gGmbh](https://joinmastodon.org/trademark). `pub` is not affiliated with Mastodon gGmbH. `pub` is not a Mastodon server.

## What _doesn't_ it do?

`pub` doesn't have much of a web interface beyond what ActivityPub requires, you're expected to interact with it via a Mastodon compatible app.

## Getting started

_Warning: `pub` is still in development, if it breaks, you can keep the pieces._

### Pre-requisites

- [Go][go]
- [MariaDB][mariadb]

### Installation for MySQL

Create a database and user for `pub`:

```sql
CREATE DATABASE pub;
CREATE USER 'pub'@'localhost' IDENTIFIED BY 'pub';
GRANT ALL PRIVILEGES ON pub.* TO 'pub'@'localhost';
```
Install `pub`:

```bash
go install github.com/davecheney/pub@latest
```
Create/migrate the database:

```bash
pub --dsn 'pub:pub@/pub' auto-migrate
```

### Setup for MySQL

Create an instance for `pub`:

```bash
pub --dsn 'pub:pub@/pub' create-instance --domain domain.com --title "Something cool" --description "Something witty" --admin-email admin@domain.com
```

This will create an instance, and an admin account for that instance.

Create your first user

```bash
pub --dsn 'pub:pub@/pub' create-account --email you@domain.com --name you --domain domain.com --password sssh
```

This will create an account for you to act as `acct:you@domain.com`

### Running for MySQL

Start `pub`:

```bash
pub --log-http --dsn 'pub:pub@/pub' serve 
```    

## Setup for SQLite

`pub  --driver sqlite --dsn db.sqlite auto-migrate`

`pub --driver sqlite --dsn db.sqlite create-instance --domain domain.com --title "yeah" --description "Geeking" --admin-email amdin@mydomain`

### Getting online

`pub` doesn't have a web interface, so you'll need to use a Mastodon app to interact with it.
You'll need to put `pub` behind a reverse proxy, and configure the reverse proxy to forward requests to `pub`.
TLS is also required, so you'll need to configure TLS for your reverse proxy, probably using [Let's Encrypt](https://letsencrypt.org/).

See the [examples](examples) directory for sample configurations for [nginx](examples/nginx).

## Acknowledgements 

`pub` would not be possible without these amazing projects:

- [Chi][chi]
- [GORM][gorm]
- [Kong][kong]
- [Requests][requests]


## Contributions

`pub` is open source, but not open for contributions _just yet_.
That may change in the future, but at the moment please do not send pull requests.

In the meantime, if you have a feature request, or a bug report, please open an issue.
If you're _really_ adventurous, you can contact me via [`@dfc@cheney.net`](acct:dfc@cheney.net).

Thank you in advance for your understanding.

[chi]: https://github.com/go-chi/chi
[kong]: https://github.com/alecthomas/kong
[gorm]: https://gorm.io/
[requests]: https://github.com/carlmjohnson/requests/
[go]: https://golang.org/doc/install
[mariadb]: https://mariadb.org/download/

## notes

https://github.com/foolish15/example-gorm-enum/blob/master/main.go