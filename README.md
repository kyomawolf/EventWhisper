# EventWhisper

## What is EventWhisper?

We collect information about Events by Scrapers and send over to the WhisperCore over REST.
The WhisperCore matches them with the interests of Identitys and sends them through touchpoints to the users.
Touchpoints are different messaging-apps, webfrontends or other services.

## ToDo

* [ ] implement OAuth2 for endpoints and identity
* [ ] Add a mongo db for the data
* [ ] Add WhatsApp as a touchpoint
* [ ] Add a webfrontend as a touchpoint
* [ ] Improve Input Validation




## data model for a identity

    // Sample object for a identity
    {
        "sub": "string",
        "name": "string",
        "interest": [
            "string", 
            "string"
        ],
        "channels": [
            {
                "id": "string",
                "channelname": "telegram",
                "type": "directmessage", // "group"
                ""
                "specifics": {
                    "chatId": "string"
                },
            }
        ],
        "announcedEvents": [
            {
                "id": "string", // the id of the entry
                "eventid": "string", // the id of the event
                "announced_at": "string", // the timestamp when the entry was announced
                "delete_at": "string" // the timestamp when the entry could be deleted
            }
        ]
    }