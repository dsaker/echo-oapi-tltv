import os
import requests
import insert_voices

api_key = os.environ['API_KEY']


def main():
    # Set up the API endpoint
    url = f'https://texttospeech.googleapis.com/v1/voices?key={api_key}'

    # Make the GET request
    response = requests.get(url)

    # Check if the request was successful
    if response.status_code == 200:
        voices = response.json()
        # print(json.dumps(voices, indent=2))
        insert_voices.insert_voices(voices)
    else:
        print(f"Failed to retrieve voices. Status code: {response.status_code}")
        print(response.text)


if __name__ == "__main__":
    main()
