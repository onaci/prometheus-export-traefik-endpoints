# prometheus-export-traefik-endpoints

Queries traefik's frontend rules to find what Hosts its listening for, and then requests all the https SSL certificates, so as to list their details for prometheus to scrape

For example:

```
# HELP traefik_endpoints Traefik https hostname endpoints and certificate info
# TYPE traefik_endpoints gauge
traefik_endpoints{dns="*.yoga260.alho.st",host="keycloak.yoga260.alho.st",ip="10.10.10.94",isurer="Let's Encrypt Authority X3",notafter="2020-01-08 08:43:30 +0000 UTC",notbefore="2019-10-10 08:43:30 +0000 UTC",subject="*.yoga260.alho.st"} 1.57847301e+09
traefik_endpoints{dns="*.yoga260.alho.st",host="monitor.yoga260.alho.st",ip="10.10.10.94",isurer="Let's Encrypt Authority X3",notafter="2020-01-08 08:43:30 +0000 UTC",notbefore="2019-10-10 08:43:30 +0000 UTC",subject="*.yoga260.alho.st"} 1.57847301e+09
traefik_endpoints{dns="*.yoga260.alho.st",host="traefik.yoga260.alho.st",ip="10.10.10.94",isurer="Let's Encrypt Authority X3",notafter="2020-01-08 08:43:30 +0000 UTC",notbefore="2019-10-10 08:43:30 +0000 UTC",subject="*.yoga260.alho.st"} 1.57847301e+09
traefik_endpoints{dns="yoga260.alho.st",host="yoga260.alho.st",ip="10.10.10.94",isurer="Let's Encrypt Authority X3",notafter="2020-01-08 08:43:51 +0000 UTC",notbefore="2019-10-10 08:43:51 +0000 UTC",subject="yoga260.alho.st"} 1.578473031e+09
```

Used in https://github.com/onaci/swarm-infra - prometheus.yml:

```
  traefik_certs:
    image: onaci/prometheus-export-traefik-endpoints:latest
    networks:
      - net
      - infra_traefik
    command:  -listen-address "0.0.0.0:8080" -traefik-api http://server:8686/api
    deploy:
      mode: replicated
      replicas: 1
      resources:
        limits:
          memory: 128M
        reservations:
          memory: 64M
    logging:
      options:
        tag: infra.monitoring
```
