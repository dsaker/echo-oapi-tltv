import json
import insert_voices

"""voices_file.py

This script loads a file from local and calls insert_voices to insert them into
the database,
"""

def main():
    # Load JSON data
    with open('../../internal/util/voices.json', 'r') as file:
        data = json.load(file)
    insert_voices.insert_voices(data)


if __name__ == "__main__":
    main()
