# PiTraefikHole

Do you use [PiHole] for your local DNS, and [Traefik] as your reverse proxy? Do you love that [Traefik] can be configured using labels in Docker? Do you hate that you can't do the same with [PiHole]? This is the solution for you!

PiTraefikHole polls your [Traefik] and if it finds any hosts that are not defined as a CNAME in [PiHole], it will add it. It will not remove any unused CNAMEs, as I would prefer destructive actions to be manual.

## How does it work?

It sits in a container next to [Traefik] and [PiHole] and every 30 seconds (or a configurable time) polls the [Traefik] API and looks at the rules for each router. If the rule is defined as ``Host(`my.domain.com`)`` it will strip the ```Host(``)``` part and keep the domain. This does mean that this only works with Host rules, and cannot have anything that is not just plain text in there. If anything else gets returned by the API, it will add it wrong, as this does not process what is returned using anything built into [Traefik]. You must provide in the configuration the CNAME record you would like the host to CNAME to. It will then use the [PiHole] API to add the CNAME record.

## Installation

This should probably be deployed side-by-side with [Traefik] in Docker. If you are using Kubernetes, you are probably also using a different solution.

```yml
---
services:
    traefik:
        ...
    
    pitraefikhole:
        image: ghcr.io/m50/pitraefikhole:main
        networks:
        - traefik
        volumes:
        - ./data/pitraefikhole.yml:/config.yml
        restart: always
```

Then you need to configure it, either by creating the config file (like above), or configuring it with environment variables.

Example config file:

```yaml
cname-address: domain.com
pihole-password: "super-secure-password"
log-level: INFO
poll-frequency-seconds: 30
pihole-address: http://pihole/
traefik-address: http://traefik:8080/
```

Example environment variables:

```env
PITRAEFIKHOLE_CNAME_ADDRESS=domain.com
PITRAEFIKHOLE_PIHOLE_PASSWORD=super-secure-password
PITRAEFIKHOLE_LOG_LEVEL=INFO
PITRAEFIKHOLE_POLL_FREQUENCY_SECONDS=30
PITRAEFIKHOLE_PIHOLE_ADDRESS=http://pihole/
PITRAEFIKHOLE_TRAEFIK_ADDRESS=http://traefik:8080/
```

Alternatively, if you would like to pass your [PiHole] password in as a secret, you can use the `pihole-password-file` (`PITRAEFIKHOLE_PIHOLE_PASSWORD_FILE`) config option and pass the file path in instead.

[PiHole]: https://pi-hole.net/
[Traefik]: https://traefik.io/traefik/
