import os
import psycopg2

conn_string = os.environ['TLTV_DB_DSN']


def insert_voices(data):
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
                for lang_code in voice["languageCodes"]:

                    # get the language id for the voice from the language tag
                    lang_tag = lang_code.split("-")
                    cursor.execute(lang_id_query, (lang_tag[0],))
                    lang_id = cursor.fetchone()

                    # Insert each record from JSON data
                    cursor.execute(insert_query, (lang_id, voice['languageCodes'], voice['ssmlGender'], voice['name'], voice['naturalSampleRateHertz']))

                    # Commit the transaction
                    connection.commit()

    except Exception as e:
        print("Error inserting data:", e)
    finally:
        # Close the database connection
        if connection:
            cursor.close()
            connection.close()