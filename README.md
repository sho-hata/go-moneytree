# go-moneytree

go-moneytree is a Go HTTP API Client library for the [Moneytree LINK API](https://docs.link.getmoneytree.com/docs/product-and-tech-overview).

## Install

```bash
go get github.com/sho-hata/go-moneytree
```

## Usage

Moneytree LINK API uses OAuth 2.0 for authentication. You need to obtain an access token through the OAuth flow.

For detailed authentication flow, please refer to the [Moneytree LINK API documentation](https://docs.link.getmoneytree.com/docs/product-and-tech-overview).

### Basic Usage

```go
// Initialize client with account name (e.g., "jp-api-staging" or "jp-api")
client, err := moneytree.NewClient("jp-api-staging")
if err != nil {
    log.Fatal(err)
}

// Set client credentials (obtain from Moneytree)
client.config.ClientID = "your-client-id"
client.config.ClientSecret = "your-client-secret"

// Retrieve access token using authorization code from OAuth flow
grantType := "authorization_code"
code := "authorization-code-from-oauth-flow"
redirectURI := "https://your-app.com/callback"
request := &moneytree.RetrieveTokenRequest{
    GrantType:   &grantType,
    Code:        &code,
    RedirectURI: &redirectURI,
}

token, err := client.RetrieveToken(ctx, request)
if err != nil {
    log.Fatal(err)
}

// Set the token in the client
// The client will automatically refresh the token when it expires
client.SetToken(token)

// Get user profile (no need to pass access token - it's managed automatically)
profile, err := client.GetProfile(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Moneytree ID: %s, Email: %s\n", profile.MoneytreeID, profile.Email)
```

### Token Management

The client automatically manages OAuth tokens:

- **Initial Token**: Use `RetrieveToken()` to get the initial token, then call `SetToken()` to set it in the client.
- **Automatic Refresh**: When a token expires, the client automatically refreshes it using the `refresh_token` grant type.
- **Thread-Safe**: Token refresh is thread-safe and prevents multiple concurrent refresh attempts.

```go
// After setting the initial token, all API calls will automatically use and refresh the token
client.SetToken(token)

// Subsequent API calls don't require passing tokens
accounts, err := client.GetPersonalAccounts(ctx)
balances, err := client.GetPersonalAccountBalances(ctx, "account-key-123")
transactions, err := client.GetPersonalAccountTransactions(ctx, "account-key-123")
```

## API availability

| Category           | API                   | Availability    |
| ------------------ | --------------------- | --------------- |
| OAuth              | Start Authorization   | Not Implemented |
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
| Manual Accounts    | All                   | Not Implemented |
| Institutions       | Get Institutions      | Available       |
| Category           | Get Categories        | Available       |
| Category           | Create Category       | Available       |
| Category           | Get Category          | Available       |
| Category           | Update Category       | Available       |
| Category           | Delete Category       | Available       |
| Category           | Get System Categories | Available       |
| 2FA                | Submit 2FA            | Available       |
| 2FA                | Get Captcha           | Available       |

## Authentication availability

| Type      | Availability |
| --------- | ------------ |
| OAuth 2.0 | Available    |

## Contributing

Contributions are welcome! Please ensure that `make test` and `make lint` succeed before submitting a pull request.

For more details, please refer to the development guidelines in the codebase.
