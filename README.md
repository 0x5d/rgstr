[![Build Status](https://travis-ci.org/castillobg/rgstr.svg?branch=master)](https://travis-ci.org/castillobg/rgstr)

# rgstr
**rgstr** (_rĕjˈĭ-stər_) automatically registers and de-registers `rkt`
(_rŏkˈĭt_) pods on [Consul](https://www.consul.io/), [etcd](https://coreos.com/etcd/).

## Please Note:
This project relies heavily on
[`rkt`'s API service](https://github.com/coreos/rkt/blob/master/Documentation/subcommands/api-service.md),
which is currently in an experimental phase. Until the API reaches stability, this project is not
suitable for a production environment.


# Run it
Build **rgstr**:
```sh
go build
```
Run rkt's API service:
```sh
rkt api-service
```

Run Consul

Run **rgstr**:
```sh
./rgstr
```

## Inspiration
**rgstr** was inspired by [progrium](https://github.com/progrium)'s and
[gliderlabs](https://github.com/gliderlabs)'s [registrator](https://github.com/gliderlabs/registrator).

## License
MIT
