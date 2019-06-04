# Monzo Go Client

A Go client for interacting with the Monzo API.

## Installation

```
go get -u github.com/tmus/monzo
```

## Usage

> The client API may change at some point, don't rely on this
> for anything serious, and make sure you use version numbers
> in your `go.mod` file.

1. Create a `monzo.Client` and pass it your Monzo access token:

```go
token, _ := os.LookupEnv("MONZO_TOKEN")
c = monzo.NewClient(token)
```

Creating a new client doesn't verify the connection to Monzo,
so you should call the `Ping` method on the new client to ensure
that the client can access Monzo:

```go
if err := c.Ping(); err != nil {
    panic(err)
}
```

2. Call the `Accounts` function on the client to return a slice
of accounts associated with the Monzo token.

```go
accs, _ := c.Accounts()

for _, acc := range accs {
    fmt.Println(acc.ID)
}
```

3. You can call the `Account` function on the client to return
a single account.

```go
acc, _ := c.Account("acc_00000XXXXXXXXXXXXXXXXX")
fmt.Println(acc.ID)
```

**More details coming soon.**