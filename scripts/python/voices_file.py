import json
import insert_voices


def main():
    # Load JSON data
    with open('../../internal/util/voices.json', 'r') as file:
        data = json.load(file)
    insert_voices.insert_voices(data)


if __name__ == "__main__":
    main()
