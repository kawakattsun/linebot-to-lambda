# linebot-to-lambda

LINE にメッセージを送ると Google Calendar に登録してくれる bot 用スクリプト

## 環境変数

GOOGLE_CALENDAR_ID: Google Calenedar のカレンダーID
CHANNEL_SECRET: LINE Messaging API のシークレット
CHANNEL_TOKEN: LINE Messaging API のトークン
## フォーマット

```
予定登録
[タイトル]
[場所]
[詳細]
[開始時間(2018-01-02 12:30)]
[終了時間(2018-01-03 20:30)]
```
