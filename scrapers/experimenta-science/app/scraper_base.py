import requests
import requests
import json
import os

from models import Event

class BaseScraper:

    def SendEvent(eventObject:Event):
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


    def GetKnownUrls(self):
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
