# Let it out!

<img src="docs/logo.png" />

# Overview

**Let it out** makes it easier to quickly expose your local services to the internet without
messing with `inlets` server tokens, and DNS records.

While we just automate the DNS management on your CloudFlare account and the `inlets` command,
all the tunneling magic comes from [inlets](https://github.com/inlets/inlets) itself. You're still
required to have an exit server up and running in order to expose it to the internet.

# Requirements

- Inlets [`inlets`](https://github.com/inlets/inlets#install-the-cli) executable installed on your machine
- One or more [exit server](https://docs.inlets.dev/#/?id=exit-servers) [installed](https://github.com/inlets/inlets/blob/master/docs/vps.md)
- A domain managed by [CloudFlare](https://www.cloudflare.com/)

## Exit servers

You can use any system you want to run exit servers if they are publicly accessible.

The easiest way to do it is using [`inletsctl`](https://github.com/inlets/inletsctl) to create an exit server, but if you want a quick and easy setup on a provider not supported by `inletsctl`, here it is how you can do it on a fresh Ubuntu machine (adjust as needed):

- Install inlets
  - `sudo curl -sLS https://get.inlets.dev | sudo sh`
  - `sudo curl -sLO https://raw.githubusercontent.com/inlets/inlets/master/hack/inlets-operator.service`
  - `sudo mv inlets-operator.service /etc/systemd/system/inlets.service`
  - `echo AUTHTOKEN=$(head -c 16 /dev/urandom | shasum | cut -d" " -f1) > ~/inlets`
  - `echo CONTROLPORT=8000 >> ~/inlets`
  - `sudo mv ~/inlets /etc/default/inlets`
  - `sudo systemctl start inlets`
  - `sudo systemctl enable inlets`
- Grab your token
  - `source /etc/default/inlets`
  - `echo $AUTHTOKEN`

# Installation

> TODO:

# Configuration

To run **let it out** you need a configuration file that has your domains and exit servers.

The configuration is done through a `.letitout.yml` file that can be placed in your home directory or the current working directory when you run `letitout`. If the file is found in the current working directory, it takes priority over the file placed on your home folder.

## Example contents

```yaml
cloudflare:
  mydomain.com:
    token: <your_cloudflare_token_with_edit_zone_permission>

  anotherdomain.com:
    token: <another_cloudflare_token_with_edit_zone_permission>

servers:
  my_digitalocean_vps:
    address: 123.123.123.123:8000
    token: <the_digitalocean_exit_server_token>

  some_aws_vps:
    address: 213.213.213.213:8000
    token: <the_aws_exit_server_token>

  some_random_place:
    address: 222.222.222.222:8000
    token: <the_token_on_random_place>

projects:
  apache:
    server: my_digitalocean_vps
    hostname: sub.mydomain.com
    upstream: http://127.0.0.1:80

  my_vue_frontend:
    server: some_aws_vps
    hostname: project2.anotherdomain.com
    upstream: http://127.0.0.1:3000
```

You can configure any number of domains, servers and projects in the configuration file. You just need to make sure projects references a valid exit server. If the hostname cannot be found under CloudFlare section, the DNS management will be skipped.

# Running

You can quickly expose your projects by running `letitout <project_name>` command.

For example, using the example configuration file, running `letitout apache` will create an `A` record `sub.mydomain.com` that points to the DigitalOcean machine IP. Note that these records have CloudFlare proxy enabled by default.

If the `exit server` IP  is an IPv6 address, the record is created as an `AAAA` type.

## Quickly (temporarily) changing settings

If you want to quickly change the hostname, the target server, or the upstream service without editing the configuration file, you can specify it using `--hostname <host>`, `--server <server>` or `--upstream <upstream>`.

# Examples

These examples assumes a configuration like this:

Server:

- `server1` at address `111.111.111.111`
- `server2` at address `222.222.222.222`

Domains:

- `domain1.com`
- `ddomain2.com`

Projects:

- `project1` on upstream `http://localhost:8000` with domain `project.domain1.com`
- `project2` on upstream `http://localhost:3000` with domain `project.domain2.com`

## Expose `project1`

### Command

`letitout project1`

### Result

- Creates an `A` DNS record `project.domain1.com` pointing to `111.111.111.111`
- Tunnels `http://localhost:8000` to `server1`, accessible through `https://project.domain1.com`

## Expose `project2`

### Command

`letitout project2`

### Result

- Creates an `A` DNS record `project.domain2.com` pointing to `222.222.222.222`
- Tunnels `http://localhost:3000` to `server2`, accessible through `https://project.domain1.com`

## Expose `project1` on `hey.domain2.com`

### Command

`letitout project1 --hostname hey.domain2.com`

### Result

- Creates an `A` DNS record `hey.domain2.com` pointing to `111.111.111.111`
- Tunnels `http://localhost:3000` to `server1`, accessible through `https://hey.domain2.com`

# Backlog

### Complete

- [x] Backup DNS record into a TXT record
- [x] Create DNS record pointing to Exit server address
- [x] Support multiple servers, projects and domains

### Pending

- [ ] Rollback DNS changed based on TXT backups on exit or command
- [ ] Support more DNS providers

# License

MIT
