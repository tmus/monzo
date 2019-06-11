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

### Creating a Client

Create a `monzo.Client` and pass it your Monzo access token:

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

### Retrieving Accounts

Call the `Accounts` function on the client to return a slice
of accounts associated with the Monzo token.

```go
accs, _ := c.Accounts()

for _, acc := range accs {
    fmt.Println(acc.ID)
}
```

#### Retrieving a Single Account

If you know the account id, you can call the `Account` function
on the client to return a single account.

```go
acc, _ := c.Account("acc_00000XXXXXXXXXXXXXXXXX")
fmt.Println(acc.ID)
```

### Account Balance

The `Account` struct is fluent: it contains a pointer to the
monzo.Client, meaning further API calls can be done directly
from the Account, such as retrieving a balance:

```go
acc, _ := c.Account("acc_00000XXXXXXXXXXXXXXXXX")

b, _ := acc.Balance()
```

The `Balance` struct contains all the information about an
Account's balance:

- `Balance.Balance` returns the amount available to spend
- `Balance.Total` returns the total balance (including money
  in Pots)
- `Balance.WithSavings` includes the total balance including
  money in Savings pots.

**More details coming soon.**