#!/usr/bin/env python
# pylint: disable=unused-argument
import json

import requests
import logging
import os

import openai
from openai import OpenAI

from flask import Flask, request

from telegram import ReplyKeyboardMarkup, ReplyKeyboardRemove, Update
from telegram.ext import (
    Application,
    CommandHandler,
    ContextTypes,
    ConversationHandler,
    MessageHandler,
    filters,
)


class channel:
    def __init__(self, channel_id):
        self.id = None
        self.channelname = "telegram"
        self.type = "directmessage"
        self.specifics = { "chatId": str(channel_id) }


class user:
    def __init__(self, sub):
        self.sub = sub
        self.name = None
        self.location = None
        self.interest = []
        self.channels = []


user_list = []

def find_by_id(userid):
    for user in user_list:
        if user.sub == userid:
            return user
    return None

# Enable logging
logging.basicConfig(
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s", level=logging.INFO
)
# set higher logging level for httpx to avoid all GET and POST requests being logged
logging.getLogger("httpx").setLevel(logging.WARNING)

logger = logging.getLogger(__name__)

INTERESTS, LOCATION, NAME = range(3)

openaiClient = OpenAI(api_key=os.environ['OPENAI_API_KEY'])

interests = ['Musik', 'Kunst', 'Sport', 'Theater', 'Film', 'Bildung', 'Mode', 'Literatur', 'Technologie', 'Business', 'Religion', 'Wohltätigkeit', 'Gastronomie', 'Outdoor', 'Umwelt', 'Markt', 'Spiele', 'Comedy', 'Wissenschaft', 'Politik', 'Festivals', 'Kinder', 'Handwerk']
interests_as_string = str(', '.join(interests))


def ask_gpt_location(message: str) -> str:
    logger.info(message)
    # return message
    response = openaiClient.chat.completions.create(messages=[{"role": "user",
                                                               "content": "Sag mir zu welcher Stadt oder Kreis die PLZ gehört, gib mir die Antwort im format: \"Stadt, Bundesland\", hier ist die deutsche Postleitzahl: " + message}],
                                                    model='gpt-4')
    return response.choices[0].message.content


def ask_gpt_name(message: str) -> str:
    # return message
    logger.info(message)
    response = openaiClient.chat.completions.create(messages=[{"role": "user",
                                                               "content": "Wie heisst die Person: " + message}],
                                                    model='gpt-4')
    return response.choices[0].message.content



def ask_gpt_interests(message: str):
    # return message.split(',')
    logger.info(message)
    response = openaiClient.chat.completions.create(messages=[{"role": "user", "content": ": welche interessen hat diese Person? Gib mir deine Antwort passend zu der List hier:" + interests_as_string + ". Hier Antwort der Person: " + message}], model='gpt-4')

    return response.choices[0].message.content


async def start(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    user_list.append(user(update.message.chat_id))
    await update.message.reply_text("Hey ich bin EventWhisper, ich kann dir Events vorschlangen. Basierend auf deinen Interessen und deinem Standort. Wie heisst du denn?")
    return NAME


async def name(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    name = update.message.from_user
    local_user = find_by_id(update.message.chat_id)
    if local_user is None:
       return ConversationHandler.END
    local_user.name = name
    await update.message.reply_text("Willst du mir auch deine Postleitzahl geben? Dann kann ich dir Events in deiner Nähe vorschlagen.")
    print(local_user.name)
    return LOCATION


async def location(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    filtered_message = ask_gpt_location(update.message.text)
    logger.info("%s's location is: %s", update.message.from_user, filtered_message)
    local_user = find_by_id(update.message.chat_id)
    if local_user is None:
       return ConversationHandler.END
    local_user.location = filtered_message
    await update.message.reply_text("Willst du mir auch deine Interessen geben? Dann kann ich dir Events vorschlagen die dich interessieren.")
    return INTERESTS


def send_data(data):

    data.channels.append(channel(data.sub))
    data.channels[0].id = data.sub

    print(type(data.sub))
    print(type(data.name.first_name))
    print(type(data.interest))
    print(type(data.location))
    print(type(data.channels[0].id))
    print(type(data.channels[0].channelname))
    print(type(data.channels[0].type))
    print(type(data.channels[0].specifics))
    local_dict = {
        'sub': str(data.sub),
        'name': str(data.name.first_name),
        'interests': data.interest,
        'location': data.location,
        'channels': [{
            'id': str(data.channels[0].id),
            'channelname': data.channels[0].channelname,
            'type': data.channels[0].type,
            'specifics': data.channels[0].specifics
        }],
        "announcedEvents": []
    }
    dump = json.dumps(local_dict)
    print(dump)
    res = requests.post('https://api.eventwhisper.de/identity', headers={'Authorization': os.environ['WHISPER_API_TOKEN'], 'Content-Type': 'application/json'}, data=dump)
    print ('response from server:', res.text)


async def interests(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    mapped_interests = ask_gpt_interests(update.message.text)

    logger.info("%s's interests are: %s", update.message.from_user, mapped_interests)
    local_user = find_by_id(update.message.chat_id)
    if local_user is None:
        return ConversationHandler.END
    local_user.interest = mapped_interests
    print("chat id: %i, location: %s, interests: %s", local_user.sub, local_user.location, local_user.interest)
    await update.message.reply_text("Danke für deine Angaben. Wenn wir tolle neue Evetns finden, dann Benachrichtigen wir dich sofort!")
    
    print("sending data to api")
    send_data(local_user)
    print("data sent")

    user_list.remove(local_user)

    return ConversationHandler.END

async def cancel(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    """Cancels and ends the conversation."""
    user = update.message.from_user
    logger.info("User %s canceled the conversation.", user.first_name)
    await update.message.reply_text(
        "Bye! I hope we can talk again some day.", reply_markup=ReplyKeyboardRemove()
    )
    local = find_by_id(update.message.chat_id)
    if local is not None:
        user_list.remove(local)
    return ConversationHandler.END




def main() -> None:
    """Run the bot."""
    # Create the Application and pass it your bot's token.
    application = Application.builder().token(os.environ['TELEGRAM_API_TOKEN']).build()

    # Add conversation handler with the states GENDER, PHOTO, LOCATION and BIO
    conv_handler = ConversationHandler(
        entry_points=[CommandHandler("start", start)],
        states={
            INTERESTS: [MessageHandler(filters.TEXT & ~filters.COMMAND, interests)],
            NAME: [MessageHandler(filters.TEXT & ~filters.COMMAND, name)],
            LOCATION: [
                MessageHandler(filters.TEXT & ~filters.COMMAND, location)
                # CommandHandler("skip", skip_location),
            ],
        },
        fallbacks=[CommandHandler("cancel", cancel)],
    )

    application.add_handler(conv_handler)

    # Run the bot until the user presses Ctrl-C
    application.run_polling(allowed_updates=Update.ALL_TYPES)


app = Flask(__name__)


@app.route("/new-event", methods=["POST"])
def hello():
    if request.method == "POST":
        formatted = request.json

        for idx in formatted["identity"]["channels"]:
            if idx == "telegram":
                print(formatted["message"])
                url = 'https://api.telegram.org/bot/sendMessage?chat_id=' + formatted['identity']['channels']['specifics']['chatId'] + '&text=' + formatted['message']
                ret = requests.post(url)
                print("Response from telegram: " + ret.text)
                return "OK", 200
    return 403 # Forbidden



if __name__ == "__main__":
    main()
