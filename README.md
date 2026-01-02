# go-moneytree

go-moneytree is a Go HTTP API Client library for the [Moneytree LINK API](https://docs.link.getmoneytree.com/docs/product-and-tech-overview).

## Install

```bash
go get github.com/sho-hata/go-moneytree
```

## Usage

### Authentication

Moneytree LINK API uses OAuth 2.0 for authentication. You need to obtain an access token through the OAuth flow.

For detailed authentication flow, please refer to the [Moneytree LINK API documentation](https://docs.link.getmoneytree.com/docs/product-and-tech-overview).

```go
// Initialize client with account name (e.g., "jp-api-staging" or "jp-api")
client, err := moneytree.NewClient("jp-api-staging")
if err != nil {
    log.Fatal(err)
}

// Use access token obtained through OAuth flow
accessToken := "your-access-token"
```

### API call

#### Get profile

```go
// Initialize client
client, err := moneytree.NewClient("jp-api-staging")
if err != nil {
    log.Fatal(err)
}

// Get user profile
profile, err := client.GetProfile(ctx, accessToken)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Moneytree ID: %s, Email: %s\n", profile.MoneytreeID, profile.Email)
```

#### Get transactions

```go
// Initialize client
client, err := moneytree.NewClient("jp-api-staging")
if err != nil {
    log.Fatal(err)
}

// Get transactions for a specific account
response, err := client.GetPersonalAccountTransactions(ctx, accessToken, "account_key_123")
if err != nil {
    log.Fatal(err)
}

for _, transaction := range response.Transactions {
    fmt.Printf("Date: %s, Amount: %v, Description: %s\n", 
        transaction.Date, transaction.Amount, *transaction.DescriptionPretty)
}
```

## API availability

| Category           | API                   | Availability    |
| ------------------ | --------------------- | --------------- |
| OAuth              | Retrieve Token        | Available       |
| OAuth              | Revoke Token          | Available       |
| Profile            | Get Profile           | Available       |
| Profile            | Revoke Profile        | Available       |
| Profile            | Refresh Profile       | Available       |
| Profile            | Get Account Groups    | Available       |
| Profile            | Refresh Account Group | Available       |
| Personal Accounts  | Get Accounts          | Available       |
| Personal Accounts  | Get Balances          | Available       |
| Personal Accounts  | Get Term Deposits     | Available       |
| Personal Accounts  | Get Transactions      | Available       |
| Personal Accounts  | Update Transaction    | Available       |
| Corporate Accounts | Get Accounts          | Available       |
| Corporate Accounts | Get Balances          | Available       |
| Corporate Accounts | Get Transactions      | Available       |
| Corporate Accounts | Update Transaction    | Available       |
| Point Accounts     | Get Accounts          | Available       |
| Point Accounts     | Get Transactions      | Available       |
| Point Accounts     | Get Expirations       | Available       |
| Institutions       | Get Institutions      | Available       |
| Category           | All                   | Not Implemented |
| 2FA                | All                   | Not Implemented |

## Authentication availability

| Type      | Availability |
| --------- | ------------ |
| OAuth 2.0 | Available    |

## Contributing

Contributions are welcome! Please ensure that `make test` and `make lint` succeed before submitting a pull request.

For more details, please refer to the development guidelines in the codebase.
