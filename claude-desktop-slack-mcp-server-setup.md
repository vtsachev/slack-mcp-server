# Running slack-mcp-server with Claude Desktop on a MacBook

This guide outlines the steps to configure and run `slack-mcp-server` with Claude Desktop on your MacBook.

## Prerequisites

Before you begin, ensure you have the following installed on your MacBook:

*   **Node.js and npm**: `slack-mcp-server` is a Node.js application. You can download and install Node.js (which includes npm) from [https://nodejs.org/](https://nodejs.org/). You can verify the installation by opening Terminal and running `node -v` and `npm -v`.

## Authentication Token Setup

To allow `slack-mcp-server` to connect to your Slack workspace, you need to obtain two authentication tokens: `SLACK_MCP_XOXC_TOKEN` and `SLACK_MCP_XOXD_TOKEN`.

### Obtaining `SLACK_MCP_XOXC_TOKEN`

1.  Open your Slack workspace in a web browser (e.g., Chrome, Firefox).
2.  Open the browser's developer console.
    *   **Chrome**: View > Developer > JavaScript Console (or Option + Command + J)
    *   **Firefox**: Tools > Browser Tools > Web Developer Tools (or Option + Command + I), then select the "Console" tab.
3.  Paste the following JavaScript snippet into the console and press Enter:

    ```javascript
    prompt("Copy your xoxc token:", JSON.parse(localStorage.getItem("localConfig_v2")).teams[JSON.parse(localStorage.getItem("localConfig_v2")).lastActiveTeamId].token)
    ```
4.  A dialog box will appear with your `SLACK_MCP_XOXC_TOKEN`. Copy this token and save it securely.

### Obtaining `SLACK_MCP_XOXD_TOKEN`

1.  Open your Slack workspace in the same web browser.
2.  Open the browser's developer tools.
    *   **Chrome**: View > Developer > Developer Tools (or Option + Command + I)
    *   **Firefox**: Tools > Browser Tools > Web Developer Tools (or Option + Command + I)
3.  Go to the "Application" tab (Chrome) or "Storage" tab (Firefox).
4.  Find the "Cookies" section in the left-hand sidebar and select the URL for your Slack workspace (e.g., `https://app.slack.com`).
5.  Search for the cookie named `d`. The value of this cookie is your `SLACK_MCP_XOXD_TOKEN`.
6.  Copy the value of the `d` cookie (it will start with `xoxd-`) and save it securely.

## Claude Desktop Configuration

Once you have obtained both tokens, you need to configure Claude Desktop to use `slack-mcp-server`.

1.  Locate or create your Claude Desktop configuration file. This is typically named `claude_desktop_config.json`. The location might vary, but it's often in a configuration directory related to Claude Desktop (e.g., `~/.config/claude_desktop/claude_desktop_config.json` or `~/Library/Application Support/ClaudeDesktop/claude_desktop_config.json`). Please refer to the Claude Desktop documentation for the exact location on your system.

2.  Add the following JSON configuration block to your `claude_desktop_config.json` file. If the file already contains configurations, ensure you merge this block correctly into the existing JSON structure (e.g., within a `providers` array or a similar top-level object, depending on Claude Desktop's expected format).

    ```json
    {
      "name": "slack",
      "human_name": "Slack",
      "warning": "Slack responses are not controlled by Anthropic and may contain harmful content.",
      "on_enable": "Ensure you have slack-mcp-server installed and configured. See documentation for details.",
      "transports": [
        {
          "name": "slack-mcp-server",
          "type": "stdio",
          "command": ["npx", "slack-mcp-server"],
          "env": {
            "SLACK_MCP_XOXC_TOKEN": "PASTE_YOUR_SLACK_MCP_XOXC_TOKEN_HERE",
            "SLACK_MCP_XOXD_TOKEN": "PASTE_YOUR_SLACK_MCP_XOXD_TOKEN_HERE",
            "SLACK_MCP_LOG_LEVEL": "info"
          },
          "timeout_seconds": 300,
          "max_concurrent_requests": 2,
          "max_retries": 5
        }
      ],
      "default_model": ["claude-2"],
      "fetch_models_supported": false,
      "status_url": "",
      "allow_custom_models": false,
      "enabled_by_default": true
    }
    ```

3.  **Important**:
    *   Replace `"PASTE_YOUR_SLACK_MCP_XOXC_TOKEN_HERE"` with the actual `SLACK_MCP_XOXC_TOKEN` you obtained.
    *   Replace `"PASTE_YOUR_SLACK_MCP_XOXD_TOKEN_HERE"` with the actual `SLACK_MCP_XOXD_TOKEN` you obtained.
    *   This configuration uses `npx slack-mcp-server` to run the server. `npx` will automatically download and run the latest version of `slack-mcp-server` if it's not already installed globally or in your current project.
    *   The `type` is set to `stdio`, meaning Claude Desktop will communicate with `slack-mcp-server` over standard input/output.

4.  Save the `claude_desktop_config.json` file.

After completing these steps, restart Claude Desktop. It should now be able to connect to Slack via `slack-mcp-server` using the configured tokens. You might need to select "Slack" as a provider or model source within Claude Desktop if it doesn't do so automatically.
