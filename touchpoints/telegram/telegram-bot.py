#!/usr/bin/env python
# pylint: disable=unused-argument

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


class user:
    def __init__(self, new_chatid):
        self.chatid = new_chatid
        self.name = None
        self.location = None
        self.interests = None


user_list = []


# Enable logging
logging.basicConfig(
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s", level=logging.INFO
)
# set higher logging level for httpx to avoid all GET and POST requests being logged
logging.getLogger("httpx").setLevel(logging.WARNING)

logger = logging.getLogger(__name__)

INTERESTS, LOCATION = range(2)
openaiClient = OpenAI(api_key=os.environ['OPENAI_API_KEY'])

interests = ['Musik', 'Kunst', 'Sport', 'Theater', 'Film', 'Bildung', 'Mode', 'Literatur', 'Technologie', 'Business', 'Religion', 'Wohltätigkeit', 'Gastronomie', 'Outdoor', 'Umwelt', 'Markt', 'Spiele', 'Comedy', 'Wissenschaft', 'Politik', 'Festivals', 'Kinder', 'Handwerk']
interests_as_string = str(', '.join(interests))


def ask_gpt_location(message: str) -> str:
    logger.info(message)
    response = openaiClient.chat.completions.create(messages=[{"role": "user",
                                                               "content": "Sag mir zu welcher Stadt oder Kreis die PLZ gehört, gib mir die Antwort im format: \"Stadt, Bundesland\", hier ist die deutsche Postleitzahl: " + message}],
                                                    model='gpt-4')
    return response.choices[0].message.content


def ask_gpt_interests(message: str) -> str:
    logger.info(message)
    response = openaiClient.chat.completions.create(messages=[{"role": "user", "content": ": welche interessen hat diese Person? Gib mir deine Antwort passend zu der List hier:" + interests_as_string + ". Hier Antwort der Person: " + message}], model='gpt-4')

    return response.choices[0].message.content


async def start(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    user_list.append(user(update.message.chat_id))
    print(user_list[-1].chatid)
    await update.message.reply_text("Hey gib mir deine Postleitzahl, damit ich deine persönlichen noch profitabler verkaufen kann!")

    return LOCATION


async def location(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    filtered_message = ask_gpt_location(update.message.text)
    await update.message.reply_text("Hey gib mir deine Interessen!")
    logger.info("%s's location is: %s", update.message.from_user, filtered_message)
    for user in user_list:
        if user.chatid == update.message.chat_id:
            user.location = filtered_message
    return INTERESTS


async def interests(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    mapped_interests = ask_gpt_interests(update.message.text)

    logger.info("%s's interests are: %s", update.message.from_user, mapped_interests)
    for user in user_list:
        if user.chatid == update.message.chat_id:
            user.interests = mapped_interests
    print("chat id: %i, location: %s, interests: %s", user_list[-1]. chatid, user_list[-1].location, user_list[-1].interests)
    return ConversationHandler.END

async def cancel(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    """Cancels and ends the conversation."""
    user = update.message.from_user
    logger.info("User %s canceled the conversation.", user.first_name)
    await update.message.reply_text(
        "Bye! I hope we can talk again some day.", reply_markup=ReplyKeyboardRemove()
    )

    return ConversationHandler.END


def main() -> None:
    """Run the bot."""
    # Create the Application and pass it your bot's token.
    application = Application.builder().token(os.environ['API_TOKEN']).build()

    # Add conversation handler with the states GENDER, PHOTO, LOCATION and BIO
    conv_handler = ConversationHandler(
        entry_points=[CommandHandler("start", start)],
        states={
            INTERESTS: [MessageHandler(filters.TEXT & ~filters.COMMAND, interests)],
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
                return "OK", 200
    return 403 # Forbidden


if __name__ == "__main__":
    main()
