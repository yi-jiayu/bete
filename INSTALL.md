## Installation

### Ubuntu 18.04

1. Install PostgreSQL 11

```
# add the repository
tee /etc/apt/sources.list.d/pgdg.list <<END
deb http://apt.postgresql.org/pub/repos/apt/ bionic-pgdg main
END

# get the signing key and import it
wget https://www.postgresql.org/media/keys/ACCC4CF8.asc
apt-key add ACCC4CF8.asc

# fetch the metadata from the new repo
apt-get update

apt-get install -y postgresql-11
```

2. Install Go

```
wget https://golang.org/dl/go1.14.7.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.14.7.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
```

3. Create the bete direction

```
cd /var
git clone https://github.com/yi-jiayu/bete.git
```

4. Compile bete binaries:

```
cd /var/bete
(cd cmd/bete && go build -o ../../dist/bin/bete -ldflags "-X main.commit=$(git rev-parse --short --verify HEAD)")
(cd cmd/seed && go build -o ../../bin/seed)
GOBIN=$PWD/bin go get -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate
```

5. Prepare database

```
sudo -u postgres createuser bete
sudo -u postgres createdb -O bete bete
sudo -u postgres psql -c 'alter user bete with superuser'
sudo -u bete bin/migrate -path migrations -database 'postgres:///bete?host=/var/run/postgresql' up
sudo -u postgres psql -c 'alter user bete with nosuperuser'
```

6. Prepare bete environment file

```
mkdir /etc/bete
```

Create a file called `environment` (fill in everything except `DATABASE_URL`):

```
DATABASE_URL=postgres:///bete?host=/var/run/postgresql
DATAMALL_ACCOUNT_KEY=
SENTRY_DSN=
STREET_VIEW_STATIC_API_KEY=
TELEGRAM_BOT_TOKEN=
```

7. Seed database

```
cp bete-seed.service /etc/systemd/system/
systemctl start bete-seed
```

8. Create systemd timer to update bus stops every day

```
cp bete-seed.timer /etc/systemd/system
systemctl enable bete-seed.timer
systemctl start bete-seed.timer
```

9. Prepare application service

```
cp bete-server.service /etc/systemd/system
systemd enable bete-server
systemd start bete-server
```

10. Prepare caddy config

```
mkdir /etc/caddy
cp /var/bete/Caddyfile /etc/caddy
```

Make sure to set the correct host and email.

11. Prepare Caddy

```
cd $(mktemp -d)
wget https://github.com/caddyserver/caddy/releases/download/v2.1.1/caddy_2.1.1_linux_amd64.tar.gz
mv caddy /usr/bin/
groupadd --system caddy
useradd --system \
    --gid caddy \
    --create-home \
    --home-dir /var/lib/caddy \
    --shell /usr/sbin/nologin \
    --comment "Caddy web server" \
    caddy
wget https://raw.githubusercontent.com/caddyserver/dist/master/init/caddy.service
mv caddy.service /etc/systemd/system/
```

12. Start caddy

```
systemctl enable caddy
systemctl start caddy
```
