# Slack MCP Server

Model Context Protocol (MCP) server for Slack Workspaces. This integration supports both Stdio and SSE transports, proxy settings and does not require any permissions or bots being created or approved by Workspace admins 😏.

### Feature Demo

![ezgif-316311ee04f444](https://github.com/user-attachments/assets/35dc9895-e695-4e56-acdc-1a46d6520ba0)

## Tools

1. `conversations_history`
  - Get messages from the channel by channelID
  - Required inputs:
    - `channel_id` (string): ID of the channel in format Cxxxxxxxxxx.
    - `cursor` (string): Cursor for pagination. Use the value of the last row and column in the response as next_cursor field returned from the previous request.
    - `limit` (string, default: 28): Limit of messages to fetch.
  - Returns: List of messages with timestamps, user IDs, and text content

2. `channels_list`
  - Get list of channels
  - Required inputs:
    - `channel_types` (string): Comma-separated channel types. Allowed values: 'mpim', 'im', 'public_channel', 'private_channel'. Example: 'public_channel,private_channel,im'.
    - `sort` (string): Type of sorting. Allowed values: 'popularity' - sort by number of members/participants in each channel.
    - `limit` (number, default: 100): Limit of channels to fetch.
    - `cursor` (string): Cursor for pagination. Use the value of the last row and column in the response as next_cursor field returned from the previous request.
  - Returns: List of channels

## Setup Guide

### 1. Authentication Setup

Open up your Slack in your browser and login.

#### Lookup `SLACK_MCP_XOXC_TOKEN`

- Open your browser's Developer Console.
- In Firefox, under `Tools -> Browser Tools -> Web Developer tools` in the menu bar
- In Chrome, click the "three dots" button to the right of the URL Bar, then select
`More Tools -> Developer Tools`
- Switch to the console tab.
- Type "allow pasting" and press ENTER.
- Paste the following snippet and press ENTER to execute:
  `JSON.parse(localStorage.localConfig_v2).teams[document.location.pathname.match(/^\/client\/([A-Z0-9]+)/)[1]].token`

Token value is printed right after the executed command (it starts with
`xoxc-`), save it somewhere for now.

#### Lookup `SLACK_MCP_XOXD_TOKEN`

 - Switch to "Application" tab and select "Cookies" in the left navigation pane.
 - Find the cookie with the name `d`.  That's right, just the letter `d`.
 - Double-click the Value of this cookie.
 - Press Ctrl+C or Cmd+C to copy it's value to clipboard.
 - Save it for later.

### 2. Installation

Choose one of these installation methods:

- [npx](#Using-npx)
- [Docker](#Using-Docker)

### 3. Configuration and Usage

You can configure the MCP server using command line arguments and environment variables.

#### Using npx

If you have npm installed, this is the fastest way to get started with `slack-mcp-server` on Claude Desktop.

Open your `claude_desktop_config.json` and add the mcp server to the list of `mcpServers`:
``` json
{
  "mcpServers": {
    "slack": {
      "command": "npx",
      "args": [
        "-y",
        "slack-mcp-server@latest",
        "--transport",
        "stdio"
      ],
      "env": {
        "SLACK_MCP_XOXC_TOKEN": "xoxc-...",
        "SLACK_MCP_XOXD_TOKEN": "xoxd-..."
      }
    }
  }
}
```

<details>
<summary>Or, stdio transport with docker.</summary>

```json
{
  "mcpServers": {
    "slack": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "-e",
        "SLACK_MCP_XOXC_TOKEN",
        "-e",
        "SLACK_MCP_XOXD_TOKEN",
        "ghcr.io/korotovsky/slack-mcp-server",
        "mcp-server",
        "--transport",
        "stdio"
      ],
      "env": {
        "SLACK_MCP_XOXC_TOKEN": "xoxc-...",
        "SLACK_MCP_XOXD_TOKEN": "xoxd-..."
      }
    }
  }
}
```

Please see [Docker](#Using-Docker) for more information.
</details>

#### Using npx with `sse` transport:

In case you would like to run it in `sse` mode, then you  should use `mcp-remote` wrapper for Claude Desktop and deploy/expose MCP server somewhere e.g. with `ngrok` or `docker-compose`.

```json
{
  "mcpServers": {
    "slack": {
      "command": "npx",
      "args": [
        "-y",
        "mcp-remote",
        "https://x.y.z.q:3001/sse",
        "--header",
        "Authorization: Bearer ${SLACK_MCP_SSE_API_KEY}"
      ],
      "env": {
        "SLACK_MCP_SSE_API_KEY": "my-$$e-$ecret"
      }
    }
  }
}
```

<details>
<summary>Or, sse transport for Windows.</summary>

```json
{
  "mcpServers": {
    "slack": {
      "command": "C:\\Progra~1\\nodejs\\npx.cmd",
      "args": [
        "-y",
        "mcp-remote",
        "https://x.y.z.q:3001/sse",
        "--header",
        "Authorization: Bearer ${SLACK_MCP_SSE_API_KEY}"
      ],
      "env": {
        "SLACK_MCP_SSE_API_KEY": "my-$$e-$ecret"
      }
    }
  }
}
```
</details>

#### TLS and Exposing to the Internet

There are several reasons why you might need to setup HTTPS for your SSE.
- `mcp-remote` is capable to handle only https schemes;
- it is generally a good practice to use TLS for any service exposed to the internet;

You could use `ngrok`:

```bash
ngrok http 3001
```

and then use the endpoint `https://903d-xxx-xxxx-xxxx-10b4.ngrok-free.app` for your `mcp-remote` argument.

#### Using Docker

For detailed information about all environment variables, see [Environment Variables](https://github.com/korotovsky/slack-mcp-server?tab=readme-ov-file#environment-variables).

**Note:** For improved security, the Docker container now runs as a non-root user (`nonroot`).

```bash
export SLACK_MCP_XOXC_TOKEN=xoxc-...
export SLACK_MCP_XOXD_TOKEN=xoxd-...

docker pull ghcr.io/korotovsky/slack-mcp-server:latest
docker run -i --rm \
  -e SLACK_MCP_XOXC_TOKEN \
  -e SLACK_MCP_XOXD_TOKEN \
  slack-mcp-server --transport stdio
```

Or, the docker-compose way:

```bash
wget -O docker-compose.yml https://github.com/korotovsky/slack-mcp-server/releases/latest/download/docker-compose.yml
wget -O .env https://github.com/korotovsky/slack-mcp-server/releases/latest/download/default.env.dist
nano .env # Edit .env file with your tokens from step 1 of the setup guide
docker-compose up -d
```

#### Console Arguments

| Argument              | Required ? | Description                                                              |
|-----------------------|------------|--------------------------------------------------------------------------|
| `--transport` or `-t` | Yes        | Select transport for the MCP Server, possible values are: `stdio`, `sse` |

#### Environment Variables

| Variable                       | Required ? | Default            | Description                                                                                                                               |
|--------------------------------|------------|--------------------|-------------------------------------------------------------------------------------------------------------------------------------------|
| `SLACK_MCP_XOXC_TOKEN`         | Yes        | `nil`              | Authentication data token field `token` from POST data field-set (`xoxc-...`).                                                            |
| `SLACK_MCP_XOXD_TOKEN`         | Yes        | `nil`              | Authentication data token from cookie `d` (`xoxd-...`).                                                                                     |
| `SLACK_MCP_DS_COOKIE`          | No         | `"1744415074"`     | The `d-s` cookie value required for Slack API requests. Defaults to a known value if not set.                                               |
| `SLACK_MCP_SERVER_PORT`        | No         | `3001`             | Port for the MCP server to listen on (used with `sse` transport).                                                                         |
| `SLACK_MCP_SERVER_HOST`        | No         | `127.0.0.1`        | Host for the MCP server to listen on (used with `sse` transport).                                                                         |
| `SLACK_MCP_SSE_API_KEY`        | No         | `nil`              | If set, requires clients of the SSE transport to provide this key as a Bearer token in the `Authorization` header for authentication.       |
| `SLACK_MCP_PROXY`              | No         | `nil`              | Proxy URL for the MCP server to use for outbound Slack API requests.                                                                        |
| `SLACK_MCP_SERVER_CA`          | No         | `nil`              | Path to a custom CA certificate file for trusting self-signed certificates (e.g., for a corporate proxy).                                   |
| `SLACK_MCP_SERVER_CA_INSECURE` | No         | `false`            | If `true`, trusts all insecure server certificates. **NOT RECOMMENDED.** Use `SLACK_MCP_SERVER_CA` instead if possible.                     |
| `SLACK_MCP_ENABLE_USER_CACHE`  | No         | `false`            | If `true`, enables on-disk caching of user data (PII). See Security section for implications.                                               |
| `SLACK_MCP_USERS_CACHE`        | No         | `.users_cache.json`| Path to the user cache file. Only used if `SLACK_MCP_ENABLE_USER_CACHE` is `true`.                                                        |

### Debugging Tools

```bash
# Run the inspector with stdio transport
npx @modelcontextprotocol/inspector go run mcp/mcp-server.go --transport stdio

# View logs
tail -n 20 -f ~/Library/Logs/Claude/mcp*.log
```

## Security

- **API Tokens**: Never share your `SLACK_MCP_XOXC_TOKEN` and `SLACK_MCP_XOXD_TOKEN`. Keep `.env` files and any configuration containing these tokens secure and private.
- **SSE API Key**: If you use the `sse` transport and expose the server, it is highly recommended to set `SLACK_MCP_SSE_API_KEY`. This variable enforces Bearer token authentication on incoming SSE connections, preventing unauthorized access. Clients must include this key in the `Authorization` header (e.g., `Authorization: Bearer your-secret-key`).
- **'d-s' Cookie Configuration**: The `d-s` cookie, necessary for Slack API interactions, can be configured using the `SLACK_MCP_DS_COOKIE` environment variable. If not set, it defaults to `"1744415074"`. While this cookie is not as sensitive as the primary auth tokens, its configurability can be useful if the default value becomes outdated.
- **PII (User Data) Caching**:
    - By default, this server **disables** on-disk caching of user data (which includes Personally Identifiable Information like user IDs, names, and email addresses) to enhance privacy and security.
    - If you need to enable on-disk user caching (e.g., to reduce API calls in a trusted environment), set the `SLACK_MCP_ENABLE_USER_CACHE` environment variable to `true`.
    - When enabled, the cache file path can be specified using `SLACK_MCP_USERS_CACHE` (defaults to `.users_cache.json`).
    - **Security Implication**: Enabling user caching means PII will be stored on the filesystem where the server runs. Ensure that this location is adequately secured and that you understand the risks associated with storing such data.
- **Non-Root Docker User**: The Docker container now runs as a non-root user (`nonroot`) by default, reducing the potential impact of a container compromise.

## License

Licensed under MIT - see [LICENSE](LICENSE) file. This is not an official Slack product.
