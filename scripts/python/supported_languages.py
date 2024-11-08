import os
import psycopg2
from google.cloud import translate_v2 as translate

conn_string = os.environ['TLTV_DB_DSN']


def insert_languages(lang_list):
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
                # Insert each record from JSON data
                cursor.execute(insert_query, (lang['name'], lang['language']))

                # Commit the transaction
                connection.commit()

                print("Data inserted successfully")

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

    for language in results:
        print("{name} ({language})".format(**language))

    return results


if __name__ == "__main__":
    r = list_languages()
    insert_languages(r)
