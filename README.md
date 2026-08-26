# FizzBuzz API

A REST API implementing a generalized version of FizzBuzz, plus a statistics
endpoint reporting the most frequently requested parameters.

## Requirements

- Go 1.22 or later (uses the stdlib `net/http` method+path routing introduced
  in 1.22). No external dependencies are required.

## Running

```sh
task run
# or: go run ./cmd/server
```

The server listens on port `8080` by default. Override with the `PORT`
environment variable:

```sh
PORT=9090 task run
```

## Testing

```sh
task test
# or: go test ./... -race -cover
```

## Docker

```sh
task docker-build
task docker-run
```

## API

### `GET /fizzbuzz`

Generates numbers from `1` to `limit`. Multiples of `int1` are replaced by
`str1`, multiples of `int2` by `str2`, and multiples of both by `str1+str2`.

| Parameter | Type   | Constraint                 |
|-----------|--------|-----------------------------|
| `int1`    | int    | > 0                          |
| `int2`    | int    | > 0                          |
| `limit`   | int    | 0 < limit <= 1,000,000       |
| `str1`    | string | non-empty                   |
| `str2`    | string | non-empty                   |

```sh
curl "localhost:8080/fizzbuzz?int1=3&int2=5&limit=16&str1=fizz&str2=buzz"
```

```json
{"result":["1","2","fizz","4","buzz","fizz","7","8","fizz","buzz","11","fizz","13","14","fizzbuzz","16"]}
```

Invalid or missing parameters return `400 Bad Request`:

```json
{"error":"limit must be a positive integer"}
```

### `GET /statistics`

Returns the most frequently requested set of parameters and its hit count.
Takes no parameters. Every valid `/fizzbuzz` call is counted, keyed by its
exact `(int1, int2, limit, str1, str2)` combination.

```sh
curl "localhost:8080/statistics"
```

```json
{"request":{"int1":3,"int2":5,"limit":16,"str1":"fizz","str2":"buzz"},"hits":3}
```

Before any request has been made, `request` is `null` and `hits` is `0`.

### `GET /health`

Liveness check, returns `{"status":"ok"}`.
