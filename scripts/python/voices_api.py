import os
import requests
import insert_voices

api_key = os.environ['API_KEY']

"""voices_api.py

voices_api.py makes a request to google text-to-speech and downloads all of the 
supported voices in the cloud text-to-speech api. 

for more info -> https://cloud.google.com/text-to-speech/docs/voices
and https://cloud.google.com/text-to-speech/docs/reference/rest/v1/voices/list

you must get an api key to perform this request ->
https://cloud.google.com/docs/authentication/api-keys

must export api key
export API_KEY=<api key>

db schema is at db/migrations/voices
"""

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
