import os
import psycopg2
import requests
from google.cloud import translate_v2 as translate

conn_string = os.environ['TLTV_DB_DSN']
api_key = os.environ['API_KEY']


def voices_sql(data):
    # Connect to PostgreSQL database
    try:
        connection = psycopg2.connect(conn_string)
        cursor = connection.cursor()

        voice_id_query = """
        SELECT id from voices where name = %s;
        """

        lang_id_query = """
        SELECT id from languages where tag = %s;
        """
        # SQL insert statement
        insert_query = """
        INSERT INTO voices (language_id, language_codes, ssml_gender, name, natural_sample_rate_hertz)
        VALUES (%s, %s, %s, %s, %s)
        """

        for voice in data['voices']:
            # check if voice already exists in the db
            cursor.execute(voice_id_query, (voice['name'],))
            voice_id = cursor.fetchone()
            if not voice_id:
                lang_tag = voice["languageCodes"][0].split("-")
                cursor.execute(lang_id_query, (lang_tag[0],))
                lang_id = cursor.fetchone()

                if lang_id:
                    codes_string = "array['"
                    language_codes = voice["languageCodes"]
                    for i in range(len(language_codes) - 1):
                        codes_string += language_codes[i] + "', "
                    codes_string += language_codes[len(language_codes) - 1] + "']"
                    print("INSERT INTO voices (language_id, language_codes, ssml_gender, name, natural_sample_rate_hertz) VALUES (" + str(lang_id[0]) + ", " + codes_string + ", '" + voice['ssmlGender'] + "', '" + voice['name'] + "', " + str(voice['naturalSampleRateHertz']) + ");")

                    # cursor.execute(insert_query, (lang_id, voice['languageCodes'], voice['ssmlGender'], voice['name'], voice['naturalSampleRateHertz']))

                    # Commit the transaction
                    # connection.commit()

    except Exception as e:
        print("Error inserting data:", e)
    finally:
        # Close the database connection
        if connection:
            cursor.close()
            connection.close()


def languages_sql(lang_list):
    # Connect to PostgreSQL database
    try:
        connection = psycopg2.connect(conn_string)
        cursor = connection.cursor()

        lang_id_query = """
        SELECT id from languages where tag = %s;
        """
        # SQL insert statement
        insert_query = """
        INSERT INTO languages (language, tag)
        VALUES (%s, %s)
        """

        for lang in lang_list:
            cursor.execute(lang_id_query, (lang["language"],))
            lang_id = cursor.fetchone()

            if not lang_id:
                print("INSERT INTO languages (language, tag) VALUES ('" + lang['name'] + "', '" + lang['language'] + "');")

    except Exception as e:
        print("Error inserting data:", e)
    finally:
        # Close the database connection
        if connection:
            cursor.close()
            connection.close()


def list_languages() -> list:
    """Lists all available languages."""

    translate_client = translate.Client()

    results = translate_client.get_languages()

    # for language in results:
    #     print("{name} ({language})".format(**language))

    return results


if __name__ == "__main__":
    r = list_languages()
    languages_sql(r)
    # Set up the API endpoint
    url = f'https://texttospeech.googleapis.com/v1/voices?key={api_key}'

    # Make the GET request
    response = requests.get(url)

    # Check if the request was successful
    if response.status_code == 200:
        voices = response.json()
        # print(json.dumps(voices, indent=2))
        voices_sql(voices)
    else:
        print(f"Failed to retrieve voices. Status code: {response.status_code}")
        print(response.text)
