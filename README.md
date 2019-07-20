# Documentation generator for GraphQL

Markdown generator for documenting GraphQL schema

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