Ayd? Slack Alert
================

Slack alert sender for [Ayd?](https://github.com/macrat/ayd) status monitoring service.


## Install

1. Download binary from [release page](https://github.com/macrat/ayd-slack-alert/releases).

2. Save downloaded binary as `ayd-slack-alert` to somewhere directory that registered to PATH.


## Usage

1. Get your [webhook URL](https://api.slack.com/messaging/webhooks).

2. Set environment variables.

``` shell
$ export AYD_URL="https://ayd-external-url.example.com"
$ export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/......"
```

3. Start Ayd? server.

``` shell
$ ayd -a slack: ping:your-target.example.com
```


## Options

Set all options through environment variable.

|Variable           |Default                |Description                |
|-------------------|-----------------------|---------------------------|
|`SLACK_WEBHOOK_URL`|                       |Slack Incoming Webhook URL.|
|`AYD_URL`          |`http://localhost:9000`|Ayd server address.        |
