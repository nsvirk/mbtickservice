# Moneybots Tick Data Service

## Endpoints

### POST /publish

#### Request

```bash
curl -X POST https://ticks.moneybots.app/publish \
        -H "Authorization: <user_id>:<enctoken>" \
        -H "Content-Type: application/json" \
        -d '{
            "bot_id": "BOT1",
            "ticker_instruments": [
                "NSE:INFY",
                "MCX:GOLDM24DECFUT",
                "MCX:GOLDM24NOVFUT",
                "MCX:SILVERMIC25APRFUT",
                "MCX:SILVERM25JUNFUT",
                "MCX:GOLDM24OCTFUT"
                ]
            }'
```

#### Response

```bash
{
  "status": "ok",
  "data": {
    "published_channel": "CH:TICKS:ABXXXX:BOT1",
    "subscribed_count": 6
}
```

#### Request Parameters

| Parameter          | Type   | Description                                             |
| ------------------ | ------ | ------------------------------------------------------- |
| bot_id             | string | The ID of the bot to publish the ticker instruments for |
| ticker_instruments | array  | The list of ticker instruments to publish               |

#### Response Data

| Data              | Type   | Description                                               |
| ----------------- | ------ | --------------------------------------------------------- |
| published_channel | string | The channel on which the ticker instruments are published |
| subscribed_count  | int    | The number of ticker instruments subscribed to            |
