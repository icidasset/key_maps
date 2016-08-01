# Authentication

Passwordless authentication documentation.



## Flow

1. {client} User gives email and submits sign-in form
2. {client} Sends request to Auth0 `/passwordless/start`
              (happens automatically using the Auth0 lib)
3. {user} Receives email and clicks on the authentication link
4. {client} Detects Auth0 token that needs to be exchanged
5. {client} Sends request to server `/auth/exchange`
6. {server} Sends request to Auth0 `/tokeninfo`
7. {server} Determines if given token is valid,
              if so, make new user record and new Guardian token
              (that stores the Auth0 token)
8. {client} Receives Guardian token that can be used
              to authenticate other requests



## Validating existing Guardian token

1. {client} Sends request to server `/auth/validate`
2. {server} Sends request to Auth0 `/tokeninfo`
3. {server} Returns true or false
