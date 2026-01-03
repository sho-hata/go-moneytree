# go-moneytree

go-moneytree is a Go HTTP API Client library for the [Moneytree LINK API](https://docs.link.getmoneytree.com/docs/product-and-tech-overview).

## Install

```bash
go get github.com/sho-hata/go-moneytree
```

## Usage

Moneytree LINK API uses OAuth 2.0 for authentication. You need to obtain an access token through the OAuth flow.

For detailed authentication flow, please refer to the [Moneytree LINK API documentation](https://docs.link.getmoneytree.com/docs/product-and-tech-overview).

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

// Get user profile using the access token
profile, err := client.GetProfile(ctx, *token.AccessToken)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Moneytree ID: %s, Email: %s\n", profile.MoneytreeID, profile.Email)
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
