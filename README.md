# Authentication Service

The authentication service provides two endpoints:

1. **Generate Tokens Endpoint**
   - **Method:** GET
   - **Params:** user_id: uuid
   - **Path:** `/generate-tokens/`
   - **Purpose:** Generate a pair of authentication tokens.

2. **Refresh Tokens Endpoint**
   - **Method:** POST
   - **Request Data:** access_token: string, refresh_token: string
   - **Path:** `/refresh-tokens`
   - **Purpose:** Regenerate authentication tokens.

## Basic Commands:

### Start the Service

To launch the authentication service, use the following command:

```bash
make start
```
