Ayd? Slack Alert
================

Slack alert sender for [Ayd?](https://github.com/macrat/ayd) status monitoring service.


## Usage

1. Get your [webhook URL](https://api.slack.com/messaging/webhooks).

2. Set environment variables.

``` shell
$ export AYD_URL="https://ayd-external-url.example.com"
$ export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/......"
```

3. Start Ayd? server.

``` shell
$ ayd -a exec:ayd-slack-alert ping:your-target.example.com
```


## Options

Set all options through environment variable.

|Variable           |Default                |Description                |
|-------------------|-----------------------|---------------------------|
|`SLACK_WEBHOOK_URL`|                       |Slack Incoming Webhook URL.|
|`AYD_URL`          |`http://localhost:9000`|Ayd? server address.       |

Below options is set by Ayd? server.

|Variable        |Default|Description                                  |
|----------------|-------|---------------------------------------------|
|`ayd_target`    |       |The alerting target address.                 |
|`ayd_status`    |       |The status of target. "FAILURE" or "UNKNOWN".|
|`ayd_checked_at`|       |The timestamp of alert firing.               |
