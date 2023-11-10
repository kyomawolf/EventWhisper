# EventWhisper

EventWhis


## data model Idendity

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
                "channel": "telegram",
                "type": "directmessage", // "group"
                ""
                "specifics": {
                    "chatId": "string"
                },
            }
        ]

        <!-- "email": "string", -->
        <!-- "password": "string", -->
        <!-- "created": "string", -->
        <!-- "updated": "string" -->
    }