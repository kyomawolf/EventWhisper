const TelegramBot = require('node-telegram-bot-api');
const express = require('express')
const OpenAI = require('openai');
const axios = require('axios');
const uuid = require('uuid');
const bodyParser = require('body-parser');

const config = {
    telegram_api_key: process.env.TELEGRAM_API_TOKEN,
    openai_api_key: process.env.OPENAI_API_KEY,
    whisper_api_key: process.env.WHISPER_API_TOKEN
}

const openai = new OpenAI({
    apiKey: config.openai_api_key
})

let bot = new TelegramBot(config.telegram_api_key, { polling: true });

bot.setMyCommands([
    { command: '/start', description: 'Starte den Bot' },
    { command: '/heute', description: 'Zeige Events für heute an' },
    { command: '/morgen', description: 'Zeige Events für morgen an' }
])


const basic_interests = [
    "Konzerte", "Theater", "Kino", "Sport", "Ausstellungen",
    "Familie", "Auto", "Kulinarik", "Popkultur", "Tanz",
    "Job", "Messen", "Weiterbildung", "Sprache", "Kirche",
    "Technik", "KI", "Programmieren", "MAKING", "Handwerk",
    "DIY",
]


const currentFlows = []

const FlowState = {
    Welcome: 0,
    AskedName: 1,
    AskedLocation: 2,
    AskedInterests: 3,
    Finished: 4
}


const sendWelcomeAndAskForName = (chatId) => {
    bot.sendMessage(chatId, 'Ich bin EventWhisper, dein zuverlässiger Begleiter für die Entdeckung aufregender Veranstaltungen in deiner Nähe. Mit meiner Kombination aus Algorithmen und deinen Interessen finde ich die perfekten Events für dich. Lass uns zusammen die aufregendsten Aktivitäten in deiner Umgebung erkunden!')
        .then(() => {
            setTimeout(() => {
                bot.sendMessage(chatId, 'Um zu beginnen, benötige ich ein paar Informationen von dir. Ich würde sagen, wir beginnen mit den Basics.')
                    .then(() => {
                        setTimeout(() => {

                            bot.sendMessage(chatId, 'Kannst du mir sagen, wie du genannt werden möchtest?').then((msg) => {
                                currentFlows.find((flow) => flow.chatId === chatId).state = FlowState.AskedName;
                            });
                        }, 1000);
                    });

            }, 3000);
        });
}

const askForLocation = (chatId, msgText) => {

    flow = currentFlows.find((flow) => flow.chatId === chatId);
    flow.data.name = msgText;

    bot.sendMessage(chatId, 'Hallo ' + msgText + ', schön dich kennenzulernen!').then((msg) => {
        setTimeout(() => {
            bot.sendMessage(chatId, 'Wenn du mir deine Postleitzahl sagst, kann ich dir in zukunft Events aus der Gegend vorschlagen. Wie ist deine PLZ?').then((msg) => {
                currentFlows.find((flow) => flow.chatId === chatId).state = FlowState.AskedLocation;
            });
        }, 3000);
    });
}

const askForInterests = (chatId, msgText) => {
    const germanZipCodeRegex = /^\d{5}$/;
    if (germanZipCodeRegex.test(msgText) == false) {
        bot.sendMessage(chatId, 'Das ist leider keine gültige Postleitzahl. Versuchen wir das ganze noch einmal. Wie ist deine PLZ?').then((msg) => {
            return;
        });

        return;
    }


    flow = currentFlows.find((flow) => flow.chatId === chatId);
    flow.data.zip = msgText;

    bot.sendMessage(chatId, 'Danke für deine PLZ!').then((msg) => {
        setTimeout(() => {
            bot.sendMessage(chatId, 'Was interessiert dich?').then((msg) => {
                currentFlows.find((flow) => flow.chatId === chatId).state = FlowState.AskedInterests;
            });
        }, 3000);
    });
}

const finishFlow = (chatId, msgText) => {

    const chatCompletion = openai.chat.completions.create({
        messages: [
            { role: "system", content: "You are a helpful assistant designed to output JSON." },
            {
                role: 'user', content: 'You are a helpful assistant designed to output JSON. \n' +
                    'You are given a user input. It should contain some things, the user likes to do. Please match this interests to something of the following:. \n' +
                    basic_interests.join(", ") + '  \n' +
                    '  \n' +
                    'Now add me the matches to an array on put them under the key "interests" in the following JSON Schema. Return me only the generated JSON \n' +
                    '  \n' +
                    '{  \n' +
                    '    "$schema": "http://json-schema.org/draft-04/schema#",  \n' +
                    '    "type": "object",  \n' +
                    '    "properties": {  \n' +
                    '      "interests": {  \n' +
                    '        "type": "array",  \n' +
                    '        "items": [  \n' +
                    '          {  \n' +
                    '            "type": "string"  \n' +
                    '          },  \n' +
                    '          {  \n' +
                    '            "type": "string"  \n' +
                    '          }  \n' +
                    '        ]  \n' +
                    '      }  \n' +
                    '    },  \n' +
                    '    "required": [  \n' +
                    '      "interests"  \n' +
                    '    ]  \n' +
                    '  }'
            },
            { role: "user", content: "here are the users inputs: " + msgText }
        ],
        response_format: { "type": "json_object" },
        model: 'gpt-4-1106-preview',
    }).then((response) => {
        console.log("response data: ", response.choices[0].message.content);

        try {
            jsonStr = JSON.parse(response.choices[0].message.content);
            console.log("json: ", jsonStr.interests);

            flow = currentFlows.find((flow) => flow.chatId === chatId);

            jsonStr.interests.forEach((interest) => {
                flow.data.interests.push(interest);
            });

            if (flow.data.interests.length < 3) {
                bot.sendMessage(chatId, 'Ok, das ist ein guter Anfang. Was interessiert dich noch?').then((msg) => {
                    return;
                });

                return;
            }

            bot.sendMessage(chatId, 'Vielen Dank. Da hast du aber ein paar spannende Dinge. Ich glaube da finde ich gute Events. Ich sende dir immer mal wieder Nachrichten mit atuellen Veranstaltungen in der Region. Außerdem kannst du mit /heute und /morgen auch mal spontan etwas unternehmen.').then((msg) => {
                flow = currentFlows.find((flow) => flow.chatId === chatId)
                flow.state = FlowState.Finished;

                axios.post('https://api.eventwhisper.de/identity', {
                    sub: "",
                    name: flow.data.name,
                    interests: flow.data.interests,
                    location: flow.data.zip,
                    channels: [{
                        id: uuid.v4(),
                        channelname: "telegram",
                        type: "directmessage",
                        specifics: {
                            chatId: chatId.toString()
                        }
                    }],
                    announcedEvents: []
                }, {
                    headers: {
                        Authorization: 'Bearer ' + config.whisper_api_key,
                    }
                }).then((response) => {
                    console.log(response.data);

                    currentFlows.splice(currentFlows.indexOf(flow), 1);
                }).catch((error) => {
                    console.error(error);
                });
            });

        } catch (error) {
            console.log("error: ", error);
            bot.sendMessage(chatId, 'Das habe ich leider nicht verstanden. Kannst du es noch einmal versuchen?').then((msg) => {
                return;
            });
        }
    });
}

const initialFlow = (chatId, msg) => {

    console.log("running initial flow")
    let state = currentFlows.find((flow) => flow.chatId == chatId)

    if (state == null) {
        console.log("found no existing flow")

        state = {
            chatId: chatId,
            state: FlowState.Welcome,
            data: { name: "", zip: "", interests: [] }
        };

        currentFlows.push(state);
    }


    switch (state.state) {
        case FlowState.Welcome:
            sendWelcomeAndAskForName(chatId);
            break;
        case FlowState.AskedName:
            askForLocation(chatId, msg);
            break;
        case FlowState.AskedLocation:
            askForInterests(chatId, msg);
            break;
        case FlowState.AskedInterests:
            finishFlow(chatId, msg);
            break;
    }
}


bot.on('message', (msg) => {
    const chatId = msg.chat.id;
    const messageText = msg.text;

    if (messageText == "/start" || currentFlows.find((flow) => flow.chatId == chatId) != null) {
        initialFlow(chatId, messageText);
    }
    else if (messageText == "/heute") {
        const now = new Date();
        const today = now.getFullYear() + "-" + (now.getMonth() + 1) + "-" + now.getDate();

        axios.post('https://api.eventwhisper.de/notify', {
            chatId: chatId.toString(),
            day: today
        }, {
            headers: {
                Authorization: 'Bearer ' + config.whisper_api_key,
            }
        }).then((response) => {
            console.log(response.data);
        }).catch((error) => {
            console.error(error);
        });
    }
    else if (messageText == "/morgen") {
        const now = new Date();
        now.setDate(now.getDate() + 1);
        const today = now.getFullYear() + "-" + (now.getMonth() + 1) + "-" + now.getDate();

        axios.post('https://api.eventwhisper.de/notify', {
            chatId: chatId.toString(),
            day: today
        }, {
            headers: {
                Authorization: 'Bearer ' + config.whisper_api_key,
            }
        }).then((response) => {
            console.log(response.data);
        }).catch((error) => {
            console.error(error);
        });
    }
    else {
        bot.sendMessage(chatId, 'Es freut mich immer von die zu hören. Wenn du Events vorgeschlagen bekommen willst, kannst du mit /heute oder /morgen Events abfragen');
    }

});

// SERVER STUFF BELOW

let server = express();
server.use(bodyParser.urlencoded({ extended: true }));
server.use(bodyParser.json());
server.use(bodyParser.raw());


server.post('/telegram/sendmsg', (req, res) => {


    const identity = req.body.identity;
    const msg = req.body.msg;

    telegram_specific = identity.channels.find((channel) => channel.channelname == "telegram");

    if (telegram_specific == null) {
        console.log("could not find telegram channel");
        res.send('identity does not have a telegram channel')
        return;
    }


    console.log("sending message to telegram chat: ", telegram_specific)

    chatId = parseInt(telegram_specific.specifics.chatId)


    bot.sendMessage(chatId, msg)
    res.send('message sent')
});




server.listen(3000, () => {
    console.log(`Example app listening on port 3000`)
});
