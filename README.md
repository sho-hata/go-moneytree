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

#### Get personal accounts

```go
// Initialize client
client, err := moneytree.NewClient("jp-api-staging")
if err != nil {
    log.Fatal(err)
}

// Get personal accounts
response, err := client.GetPersonalAccounts(ctx, accessToken)
if err != nil {
    log.Fatal(err)
}

for _, account := range response.Accounts {
    fmt.Printf("Account: %s, Type: %s, Balance: %v\n", 
        account.AccountKey, account.AccountType, account.Balance)
}
```

#### Get personal account transactions

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

#### Update personal account transaction

```go
// Initialize client
client, err := moneytree.NewClient("jp-api-staging")
if err != nil {
    log.Fatal(err)
}

// Update transaction description and category
descriptionGuest := "新しいメモ"
categoryID := int64(123)
request := &moneytree.UpdatePersonalAccountTransactionRequest{
    DescriptionGuest: &descriptionGuest,
    CategoryID:       &categoryID,
}

transaction, err := client.UpdatePersonalAccountTransaction(ctx, accessToken, "account_key_123", 1337, request)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Updated transaction: ID=%d, Description=%s\n", transaction.ID, *transaction.DescriptionGuest)
```

#### Get corporate accounts

```go
// Initialize client
client, err := moneytree.NewClient("jp-api-staging")
if err != nil {
    log.Fatal(err)
}

// Get corporate accounts
response, err := client.GetCorporateAccounts(ctx, accessToken)
if err != nil {
    log.Fatal(err)
}

for _, account := range response.Accounts {
    fmt.Printf("Account: %s, Subtype: %s, Balance: %v\n", 
        account.AccountKey, account.AccountSubtype, account.CurrentBalance)
}
```

#### Get corporate account transactions

```go
// Initialize client
client, err := moneytree.NewClient("jp-api-staging")
if err != nil {
    log.Fatal(err)
}

// Get transactions for a specific corporate account
response, err := client.GetCorporateAccountTransactions(ctx, accessToken, "account_key_123")
if err != nil {
    log.Fatal(err)
}

for _, transaction := range response.Transactions {
    fmt.Printf("Date: %s, Amount: %v, Description: %s\n", 
        transaction.Date, transaction.Amount, *transaction.DescriptionPretty)
}
```

#### Update corporate account transaction

```go
// Initialize client
client, err := moneytree.NewClient("jp-api-staging")
if err != nil {
    log.Fatal(err)
}

// Update transaction description and category
descriptionGuest := "新しいメモ"
categoryID := int64(123)
request := &moneytree.UpdateCorporateAccountTransactionRequest{
    DescriptionGuest: &descriptionGuest,
    CategoryID:       &categoryID,
}

transaction, err := client.UpdateCorporateAccountTransaction(ctx, accessToken, "account_key_123", 1337, request)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Updated transaction: ID=%d, Description=%s\n", transaction.ID, *transaction.DescriptionGuest)
```

#### Get institutions

```go
// Initialize client
client, err := moneytree.NewClient("jp-api-staging")
if err != nil {
    log.Fatal(err)
}

// Get financial institutions list
response, err := client.GetInstitutions(ctx, systemAccessToken)
if err != nil {
    log.Fatal(err)
}

for _, inst := range response.Institutions {
    if inst.Status == "available" {
        fmt.Printf("Available: %s (%s)\n", inst.Name, inst.EntityKey)
    }
}
```

#### Get point accounts

```go
// Initialize client
client, err := moneytree.NewClient("jp-api-staging")
if err != nil {
    log.Fatal(err)
}

// Get point accounts
response, err := client.GetPointAccounts(ctx, accessToken)
if err != nil {
    log.Fatal(err)
}

for _, account := range response.PointAccounts {
    fmt.Printf("Account: ID=%d, Type=%s, Balance=%v\n", 
        account.ID, account.AccountType, account.CurrentBalance)
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
