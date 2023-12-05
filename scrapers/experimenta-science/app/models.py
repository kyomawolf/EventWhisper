class Location:
    def __init__(self, city, country, email, name, street, telefone, zip):
        self.city = city
        self.country = country
        self.email = email
        self.name = name
        self.street = street
        self.telefone = telefone
        self.zip = zip


class Event:
    def __init__(
        self,
        id: str,
        description: str,
        interests: list[str],
        location: Location,
        organizer: str,
        pricing: str,
        start_date_time: str,
        end_date_time: str,
        title: str,
        url: str,
    ):
        self.id = id
        self.description = description
        # Array of strings
        self.interests = interests
        # Location object
        self.location = location
        self.organizer = organizer
        self.pricing = pricing
        self.start_date_time = start_date_time
        self.end_date_time = end_date_time
        self.title = title
        self.url = url
