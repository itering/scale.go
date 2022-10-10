# Scale Codec in Golang

Scale Codec written by golang, reference from [polkascan/py-scale-codec](https://github.com/polkascan/py-scale-codec).

## Features

- [x] Decode
- [x] Encode
- [x] Metadata decode(V11,V12,V13,V14)
- [x] Extrinsic decode
- [x] EventRecord decode
- [x] Custom type reg
- [x] Reg special version type

## Installation

```bash
go get -u github.com/itering/scale.go
```

## Test

```bash
go test ./... --cover -v
```

## Resources

- [Subscan](https://www.subscan.io)
- [Subscan Open-source](https://github.com/itering/subscan-essentials)
- [Polkadot.js](http://polkadot.js.org/)

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/itering/scale.go

## License

The package is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).
