# (Moved) Documentation generator for GraphQL

Markdown generator for documenting GraphQL schema

Download binary in [Releases](https://github.com/mvochoa/graphqldoc/releases) or download docker image [mvochoa/graphqldoc](https://hub.docker.com/r/mvochoa/graphqldoc)

## Use

Generate dir `doc/` with markdown files

```bash
$ graphqldoc http://localhost:8080/graphql
```

## Compile

Is need `go-bindata` [https://github.com/go-bindata/go-bindata](https://github.com/go-bindata/go-bindata)

```
go get github.com/go-bindata/go-bindata/...
```

```bash
$ go-bindata -o assets.go template/
$ sed -i '' 's/package\ main/package\ graphqldoc/g' assets.go # Change package of assets.go file
$ cd cmd/
$ go install -v
```
