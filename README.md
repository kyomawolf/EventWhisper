# EventWhisper

EventWhis


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