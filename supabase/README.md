# Supabase OAuth2 Provider

This package provides Supabase authentication support for the [go-zoox/oauth2](https://github.com/go-zoox/oauth2) library.

## Features

- Full OAuth2 integration with Supabase Auth
- Support for custom Supabase project URLs
- Automatic user information extraction
- Token refresh support
- Customizable scopes

## Prerequisites

1. A Supabase project - Create one at [https://supabase.com](https://supabase.com)
2. OAuth2 application configured in your Supabase project

## Configuration

### 1. Set up OAuth2 in Supabase Dashboard

1. Go to your Supabase project dashboard
2. Navigate to **Authentication** > **Settings**
3. In the **Site URL** section, add your application's URL
4. In the **Redirect URLs** section, add your callback URL (e.g., `http://localhost:8080/auth/callback`)
5. Note down your project URL, it will be something like `https://your-project-id.supabase.co`

### 2. Environment Variables

Set the following environment variables:

```bash
export SUPABASE_BASE_URL="https://your-project-id.supabase.co"
export SUPABASE_CLIENT_ID="your-client-id"
export SUPABASE_CLIENT_SECRET="your-client-secret"
export SUPABASE_REDIRECT_URI="http://localhost:8080/auth/callback"
```

## Usage

### Basic Usage

```go
package main

import (
    "log"
    "github.com/go-zoox/oauth2"
    "github.com/go-zoox/oauth2/supabase"
)

func main() {
    // Create Supabase OAuth2 client
    client, err := supabase.New(&supabase.SupabaseConfig{
        BaseURL:      "https://your-project-id.supabase.co",
        ClientID:     "your-client-id",
        ClientSecret: "your-client-secret",
        RedirectURI:  "http://localhost:8080/auth/callback",
        Scope:        "openid email profile",
    })
    if err != nil {
        log.Fatal("Failed to create Supabase client:", err)
    }

    // Start authentication flow
    client.Authorize("state-value", func(loginURL string) {
        // Redirect user to loginURL
        log.Println("Visit:", loginURL)
    })
}
```

### Handle OAuth Callback

```go
// Handle the OAuth callback
client.Callback(code, state, func(user *oauth2.User, token *oauth2.Token, err error) {
    if err != nil {
        log.Printf("Authentication failed: %v", err)
        return
    }

    // Authentication successful
    log.Printf("User: %+v", user)
    log.Printf("Token: %+v", token)
    
    // User information available:
    // user.ID          - User's unique ID
    // user.Email       - User's email address
    // user.Username    - User's username
    // user.Nickname    - User's display name
    // user.Avatar      - User's avatar URL
})
```

### Custom Scopes

```go
client, err := supabase.New(&supabase.SupabaseConfig{
    BaseURL:      "https://your-project-id.supabase.co",
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
    RedirectURI:  "http://localhost:8080/auth/callback",
    Scope:        "openid email profile user_metadata", // Custom scopes
})
```

## Configuration Options

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `BaseURL` | string | Yes | Your Supabase project URL (e.g., `https://your-project.supabase.co`) |
| `ClientID` | string | Yes | OAuth2 client ID from your Supabase project |
| `ClientSecret` | string | Yes | OAuth2 client secret from your Supabase project |
| `RedirectURI` | string | Yes | Callback URL after authentication |
| `Scope` | string | No | OAuth2 scopes (default: `"openid email profile"`) |

## User Information

The following user information is automatically extracted from Supabase:

- **ID**: User's unique identifier
- **Email**: User's email address
- **Username**: User's username (falls back to email if not set)
- **Nickname**: User's display name from `user_metadata.full_name`
- **Avatar**: User's avatar URL from `user_metadata.avatar_url`

## Example Application

See the [example](../example/supabase/) directory for a complete working example with a web server.

To run the example:

1. Set the required environment variables
2. Run the example:
   ```bash
   cd example/supabase
   go run main.go
   ```
3. Visit `http://localhost:8080` in your browser

## Error Handling

The provider handles common OAuth2 errors:

- Invalid credentials
- Missing configuration
- Network errors
- Invalid callback parameters

Always check for errors in the callback function:

```go
client.Callback(code, state, func(user *oauth2.User, token *oauth2.Token, err error) {
    if err != nil {
        // Handle authentication error
        log.Printf("Authentication failed: %v", err)
        return
    }
    // Success case
})
```

## Security Considerations

1. **Environment Variables**: Store sensitive information like client secrets in environment variables
2. **HTTPS**: Always use HTTPS in production
3. **State Parameter**: Use a random state parameter to prevent CSRF attacks
4. **Scope Limitation**: Only request the minimum scopes required for your application
5. **Token Storage**: Store tokens securely and consider encryption for sensitive data

## Troubleshooting

### Common Issues

1. **"Invalid redirect URI"**: Ensure your redirect URI is exactly the same in your code and Supabase dashboard
2. **"Invalid client"**: Check that your client ID and secret are correct
3. **"Base URL required"**: Make sure you provide the full Supabase project URL
4. **CORS errors**: Configure CORS settings in your Supabase dashboard if needed

### Debug Mode

Enable debug logging to see OAuth2 flow details:

```go
import "github.com/go-zoox/logger"

logger.SetLevel(logger.DEBUG)
```

## License

This package is part of the [go-zoox/oauth2](https://github.com/go-zoox/oauth2) library and follows the same license terms.