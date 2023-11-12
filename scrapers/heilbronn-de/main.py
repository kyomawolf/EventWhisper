import requests
from bs4 import BeautifulSoup
import re
import requests
import json
import datetime, threading
import os
from openai import OpenAI


def get_single_event(event_url):
    event_response = requests.get(event_url)
    event_soup = BeautifulSoup(event_response.content, "html.parser")

    date = (
        event_soup.find("header", {"class": "pageTitle"})
        .find("h1")
        .find("span")
        .text.strip()
    )
    event_soup.find("header", {"class": "pageTitle"}).find("h1").span.decompose()
    title = event_soup.find("header", {"class": "pageTitle"}).find("h1").text.strip()
    location = event_soup.find("div", {"class": "address"}).find("p").text.strip()


    pricing = ""
    pricing_div = event_soup.find("div", {"id": "prices"})
    if pricing_div is not None:
        pricing = pricing_div.find("div", {"class": "tab-element"}).text.strip()

    organizer = ""
    organizer_div = event_soup.find("div", {"id": "contributor"})
    if organizer_div is not None:
        organizer = organizer_div.find("div", {"class": "tab-element"}).text.strip()

    # event_soup.find("header", {"class": "pageTitle"}).find("h1").span.decompose()
    event_soup.find("header", {"class": "pageTitle"}).find_next_sibling("div", {"class": "row"}).find("div").find("div", { "class": "images"}).decompose()
    event_soup.find("header", {"class": "pageTitle"}).find_next_sibling("div", {"class": "row"}).find("div").find("div", { "class": "accordion"}).decompose()

    description = event_soup.find("header", {"class": "pageTitle"}).find_next_sibling("div", {"class": "row"}).find("div").text.strip()

    pattern = re.compile(r"\s\s+")
    title = re.sub(pattern, " ", title)
    date = re.sub(pattern, " ", date)
    location = re.sub(pattern, " ", location)
    description = re.sub(pattern, " ", description)
    pricing = re.sub(pattern, " ", pricing)
    organizer = re.sub(pattern, " ", organizer)

    eventObject = {
        "id": "",
        "url": event_url,
        "title": title,
        "start_time": date,
        "end_time": "",
        "organizer": organizer,
        "location": location,
        "description": description,
        "pricing": pricing,
        "interest": [
            "Heilbronn",
        ],
    }

    return eventObject


def load_event_urls(url):
    d = {
        "tx_hbtypo3extimxevents_events[filters][from]": "11.11.2023",
        "tx_hbtypo3extimxevents_events[filters][to]": "30.11.2023",
    }

    response = requests.post(url, data=d)
    soup = BeautifulSoup(response.content, "html.parser")

    event_elements = soup.find_all("a", class_="event-item")

    urls = []

    for event_element in event_elements:
        urls.append("https://www.heilbronn.de" + event_element["href"])

    return urls


def send_event(eventObject):
    # Define the endpoint URL
    url = "https://api.eventwhisper.de/events"

    api_key = os.environ.get("API_KEY")

    # Define the authorization header
    headers = {"Authorization": "Bearer " + api_key, "Content-Type": "application/json"}

    # Convert the event object to JSON
    eventJson = json.dumps(eventObject)

    # Make the POST request
    response = requests.post(url, headers=headers, data=eventJson)

    # Check the response status code
    if response.status_code == 200:
        print("Event created successfully!")
    else:
        print("Error creating event:", response.text)


def fix_data(eventObject):
    promt = '''
You are a helpful assistant designed to output JSON. 
You are given a JSON string representing a part of a calendar. 
You need to parse it into a proper format.

This is the desired output JSON Schema:
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "type": "object",
  "properties": {
    "id": {
      "type": "string"
    },
    "title": {
      "type": "string"
    },
    "description": {
      "type": "string"
    },
    "location": {
      "type": "object",
      "properties": {
        "name": {
          "type": "string"
        },
        "street": {
          "type": "string"
        },
        "zip": {
          "type": "string"
        },
        "city": {
          "type": "string"
        },
        "country": {
          "type": "string"
        },
        "telefone": {
          "type": "string"
        },
        "email": {
          "type": "string"
        }
      },
      "required": [
        "name",
        "street",
        "zip",
        "city",
        "country",
        "telefone",
        "email"
      ]
    },
    "start_date_time": {
      "type": "string"
    },
    "end_date_time": {
      "type": "string"
    },
    "organizer": {
      "type": "string"
    },
    "pricing": {
      "type": "string"
    },
    "url": {
      "type": "string"
    },
    "interests": {
      "type": "array",
      "items": [
        {
          "type": "string"
        },
        {
          "type": "string"
        },
        {
          "type": "string"
        },
        {
          "type": "string"
        },
        {
          "type": "string"
        }
      ]
    }
  },
  "required": [
    "id",
    "title",
    "description",
    "location",
    "start_time",
    "end_time",
    "organizer",
    "pricing",
    "url",
    "interests"
  ]
}

Pricing should be "Eintritt frei" or the actual price in euros.
Please also fix html entities like &auml; or &ouml; or any other encoded or escaped characters.
Use the description to find five interests that fit the event and add them to the interests array. Select only intrests from here 
    interests: ["Konzerte", "Theater", "Kino", "Sport", "Ausstellungen",
    "Familie", "Auto", "Kulinarik", "Popkultur", "Tanz",
    "Job", "Messen", "Weiterbildung", "Sprache", "Kirche",
    "Technik", "KI", "Programmieren", "MAKING", "Handwerk",
    "DIY"]

Please answer only the JSON string without any additional text.
Please try to parse the following JSON string into the previous given output format. 

'''
    
    json_string = json.dumps(eventObject)
    promt += json_string

    client = OpenAI()

    response = client.chat.completions.create(
        model="gpt-4-1106-preview",
        response_format={"type": "json_object"},
        messages=[
            {
                "role": "system",
                "content": "You are a helpful assistant designed to output JSON.",
            },
            {"role": "user", "content": promt},
        ],
    )

    return json.loads(response.choices[0].message.content)


def get_known_urls():
    url = "https://api.eventwhisper.de/events"

    api_key = os.environ.get("API_KEY")

    # Define the authorization header
    headers = {"Authorization": "Bearer " + api_key, "Content-Type": "application/json"}

    reponse = requests.get(url, headers=headers)

    if reponse.status_code == 200:
        print("Events loaded successfully!")

        events = json.loads(reponse.text)

        return list(map(lambda event: event["url"], events))
    
    return []


def main():
    try:
        url = "https://www.heilbronn.de/tourismus/veranstaltungen/veranstaltungskalender.html"
        event_urls = load_event_urls(url)

        print("found " + str(len(event_urls)) + " urls in calendar")

        known_urls = get_known_urls()
        event_urls = list(filter(lambda url: url not in known_urls, event_urls))

        print("new urls: " + str(len(event_urls)))

        # event_urls = event_urls[:2]

        for event_url in event_urls:
            try:
                print("Running for: " + event_url)
                eventObject = get_single_event(event_url)
                eventObject = fix_data(eventObject)

                print(json.dumps(eventObject, sort_keys=True, indent=4))

                send_event(eventObject)
            except Exception as e:
                print(e)

    except Exception as e:
        print(e)

    print("Done!")
    print("Starting again in 10 minutes...")
    threading.Timer(60.0 * 10, main).start()


if __name__ == "__main__":
    OpenAI.api_key = os.environ.get("OPENAI_API_KEY")

    print("Starting scraper...")
    main()
