# github-webhook-listener

Go webhook listener for continous deployment on your server.

## How to use it

1. Create a .env that follows the .env.example schema.
2. Create a config, example:

```yaml
projects:
  - repository: websocket-chat
    branch: main
    command: "cd ~/projects/websocket-chat && just redeploy"
```

3. Now set your github webhook to `{your-server-url}/webhook`.
