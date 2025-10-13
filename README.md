# gomoji-updater

Public helper module used by the gomoji CI to fetch emoji data from unicode.org and generate the `data.go` file for the main gomoji library.

## Usage

```
go run ./cmd/updater/main.go --output ./data.go
```

This will fetch the latest Unicode emoji definitions and render them using the same structure as the gomoji repository expects.

## License

MIT
