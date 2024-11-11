# go-hexagonal-example
Example Repository

## Test
```shell
go test ./... -v -coverprofile=coverage.out
```

## Coverage
```shell
go tool cover -html=coverage.out -o tmp/coverage.html
```