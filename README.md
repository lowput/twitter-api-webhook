# Google Cloud Functions Sample Code

Twitterの[Account Activity API](https://developer.twitter.com/en/docs/accounts-and-users/subscribe-account-activity/api-reference/aaa-premium)
で登録したユーザーのツイートに画像があればSlackへPOSTする


## 環境

- [Cloud Functions](https://cloud.google.com/functions?hl=ja)

## 設定

`.env.yaml.sample` を `.env.yaml` にリネームしてそれぞれの環境に合わせます。

- TWITTER_CONSUMER_SECRET: [Twitter Developer](https://developer.twitter.com/)で作成した `API secret key`
- SLACK_ACCESS_TOKEN: Slack API トークン
- SLACK_CHANNEL_ID: Slack チャンネルID

## デプロイ

```bash
gcloud functions deploy TwitterApi --env-vars-file .env.yaml --entry-point TwitterApi --runtime go111 --trigger-http
```

## 参考

- [Twitter Developer - Subscribe to account activity](https://developer.twitter.com/en/docs/accounts-and-users/subscribe-account-activity/guides/securing-webhooks)

## License

MIT License
