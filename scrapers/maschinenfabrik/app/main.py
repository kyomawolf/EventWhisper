import threading
import os

from scraper import Scraper


def main():
    try:
        scraper = Scraper()
        scraper.Run()

    except Exception as e:
        print(e)

    print("Done!")
    print("Starting again in 10 minutes...")
    threading.Timer(60.0 * 10, main).start()


if __name__ == "__main__":
    print("Starting scraper...")
    main()
