# rgstr
**rgstr** (_re-gis-ter_) automatically registers and de-registers `rkt` pods on [Consul](), [etcd]().

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
