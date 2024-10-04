CREATE TABLE IF NOT EXISTS languages (
    id smallserial PRIMARY KEY,
    language text NOT NULL,
    tag text NOT NULL
);

INSERT INTO languages (id, language, tag) VALUES (-1, 'Not a Language', 'NaL');
INSERT INTO languages (language, tag) VALUES
                                          ('Afrikaans', 'af'),
                                          ('Albanian', 'sq'),
                                          ('Amharic', 'am'),
                                          ('Arabic', 'ar'),
                                          ('Armenian', 'hy'),
                                          ('Assamese', 'as'),
                                          ('Aymara', 'ay'),
                                          ('Azerbaijani', 'az'),
                                          ('Bambara', 'bm'),
                                          ('Basque', 'eu'),
                                          ('Belarusian', 'be'),
                                          ('Bengali', 'bn'),
                                          ('Bhojpuri', 'bho'),
                                          ('Bosnian', 'bs'),
                                          ('Bulgarian', 'bg'),
                                          ('Catalan', 'ca'),
                                          ('Cebuano', 'ceb'),
                                          ('Chinese (Simplified)', 'zh'),
                                          ('Chinese (Traditional)', 'zh-TW'),
                                          ('Corsican', 'co'),
                                          ('Croatian', 'hr'),
                                          ('Czech', 'cs'),
                                          ('Danish', 'da'),
                                          ('Dhivehi', 'dv'),
                                          ('Dogri', 'doi'),
                                          ('Dutch', 'nl'),
                                          ('English', 'en'),
                                          ('Esperanto', 'eo'),
                                          ('Estonian', 'et'),
                                          ('Ewe', 'ee'),
                                          ('Filipino (Tagalog)', 'fil'),
                                          ('Finnish', 'fi'),
                                          ('French', 'fr'),
                                          ('Frisian', 'fy'),
                                          ('Galician', 'gl'),
                                          ('Georgian', 'ka'),
                                          ('German', 'de'),
                                          ('Greek', 'el'),
                                          ('Guarani', 'gn'),
                                          ('Gujarati', 'gu'),
                                          ('Haitian Creole', 'ht'),
                                          ('Hausa', 'ha'),
                                          ('Hawaiian', 'haw'),
                                          ('Hebrew', 'he'),
                                          ('Hindi', 'hi'),
                                          ('Hmong', 'hmn'),
                                          ('Hungarian', 'hu'),
                                          ('Icelandic', 'is'),
                                          ('Igbo', 'ig'),
                                          ('Ilocano', 'ilo'),
                                          ('Indonesian', 'id'),
                                          ('Irish', 'ga'),
                                          ('Italian', 'it'),
                                          ('Japanese', 'ja'),
                                          ('Javanese', 'jv'),
                                          ('Kannada', 'kn'),
                                          ('Kazakh', 'kk'),
                                          ('Khmer', 'km'),
                                          ('Kinyarwanda', 'rw'),
                                          ('Konkani', 'gom'),
                                          ('Korean', 'ko'),
                                          ('Krio', 'kri'),
                                          ('Kurdish', 'ku'),
                                          ('Kyrgyz', 'ky'),
                                          ('Lao', 'lo'),
                                          ('Latin', 'la'),
                                          ('Latvian', 'lv'),
                                          ('Lingala', 'ln'),
                                          ('Lithuanian', 'lt'),
                                          ('Luganda', 'lg'),
                                          ('Luxembourgish', 'lb'),
                                          ('Macedonian', 'mk'),
                                          ('Maithili', 'mai'),
                                          ('Malagasy', 'mg'),
                                          ('Malay', 'ms'),
                                          ('Malayalam', 'ml'),
                                          ('Maltese', 'mt'),
                                          ('Maori', 'mi'),
                                          ('Marathi', 'mr'),
                                          ('Meiteilon', 'mni-Mtei'),
                                          ('Mizo', 'lus'),
                                          ('Mongolian', 'mn'),
                                          ('Myanmar', 'my'),
                                          ('Nepali', 'ne'),
                                          ('Norwegian', 'no'),
                                          ('Nyanja', 'ny'),
                                          ('Odia', 'or'),
                                          ('Oromo', 'om'),
                                          ('Pashto', 'ps'),
                                          ('Persian', 'fa'),
                                          ('Polish', 'pl'),
                                          ('Portuguese', 'pt'),
                                          ('Punjabi', 'pa'),
                                          ('Quechua', 'qu'),
                                          ('Romanian', 'ro'),
                                          ('Russian', 'ru'),
                                          ('Samoan', 'sm'),
                                          ('Sanskrit', 'sa'),
                                          ('Scots Gaelic', 'gd'),
                                          ('Sepedi', 'nso'),
                                          ('Serbian', 'sr'),
                                          ('Sesotho', 'st'),
                                          ('Shona', 'sn'),
                                          ('Sindhi', 'sd'),
                                          ('Sinhala', 'si'),
                                          ('Slovak', 'sk'),
                                          ('Slovenian', 'sl'),
                                          ('Somali', 'so'),
                                          ('Spanish', 'es'),
                                          ('Sundanese', 'su'),
                                          ('Swahili', 'sw'),
                                          ('Swedish', 'sv'),
                                          ('Tagalog', 'tl'),
                                          ('Tajik', 'tg'),
                                          ('Tamil', 'ta'),
                                          ('Tatar', 'tt'),
                                          ('Telugu', 'te'),
                                          ('Thai', 'th'),
                                          ('Tigrinya', 'ti'),
                                          ('Tsonga', 'ts'),
                                          ('Turkish', 'tr'),
                                          ('Turkmen', 'tk'),
                                          ('Twi', 'ak'),
                                          ('Ukrainian', 'uk'),
                                          ('Urdu', 'ur'),
                                          ('Uyghur', 'ug'),
                                          ('Uzbek', 'uz'),
                                          ('Vietnamese', 'vi'),
                                          ('Welsh', 'cy'),
                                          ('Xhosa', 'xh'),
                                          ('Yiddish', 'yi'),
                                          ('Yoruba', 'yo'),
                                          ('Zulu', 'zu');
