package database

func Schema() string {
	schema := `
        CREATE TABLE IF NOT EXISTS USERS_PROFILE(
            ID          UUID            PRIMARY KEY,
            USERNAME    VARCHAR(30),
            EMAIL       VARCHAR(255)    UNIQUE,
            PWD         VARCHAR(72),
            CREATED_AT  TIMESTAMPTZ,
            EDITED_AT   TIMESTAMPTZ,
            VERIFIED    BOOLEAN
        );

        CREATE TABLE IF NOT EXISTS LANGUAGES(
            ISO         VARCHAR(2)      PRIMARY KEY,
            LANG        VARCHAR(20),
            IMG         VARCHAR
        );

        CREATE TABLE IF NOT EXISTS WORDS(
            WORD        VARCHAR(50)     PRIMARY KEY,
            LANG_ISO    VARCHAR(2)      REFERENCES LANGUAGES(ISO) ON DELETE CASCADE
        );

        CREATE TABLE IF NOT EXISTS LITERACY(
            USER_ID     UUID            REFERENCES USERS_PROFILE(ID) ON DELETE CASCADE,
            WORD        VARCHAR(50)     REFERENCES WORDS(WORD) ON DELETE CASCADE,
            LANG_ISO    VARCHAR(2)      REFERENCES LANGUAGES(ISO) ON DELETE CASCADE,
            KNOWN_LEVEL INT
        );

        CREATE TABLE IF NOT EXISTS CONTENTS(
            ID          UUID            PRIMARY KEY,
            USER_ID     UUID            REFERENCES USERS_PROFILE(ID) ON DELETE CASCADE,
            TITLE       VARCHAR(255),
            LANG_ISO    VARCHAR(20)     REFERENCES LANGUAGES(ISO),
            BODY        TEXT,
            IMG         VARCHAR,
            CREATED_AT  TIMESTAMP,
            EDITED_AT   TIMESTAMP
        );

        CREATE TABLE IF NOT EXISTS STUDY_OVERVIEW(
            USER_ID     UUID            REFERENCES USERS_PROFILE(ID) ON DELETE CASCADE,
            LANG_ISO    VARCHAR(2)      REFERENCES LANGUAGES(ISO) ON DELETE CASCADE,
            WORD        VARCHAR(50)     REFERENCES WORDS(WORD) ON DELETE CASCADE,
            LITERACY    INT             
        );

        CREATE TABLE IF NOT EXISTS INVENTORY(
            USER_ID     UUID            REFERENCES USERS_PROFILE(ID) ON DELETE CASCADE,
            HEAD_ID     UUID            REFERENCES CONTENTS(ID) ON DELETE CASCADE,
            LANG_ISO    VARCHAR(2)      REFERENCES LANGUAGES(ISO) ON DELETE CASCADE
        );

        CREATE TABLE IF NOT EXISTS NOTE(
            USER_ID     UUID            REFERENCES USERS_PROFILE(ID) ON DELETE CASCADE,
            WORD        VARCHAR(50)     REFERENCES WORDS(WORD) ON DELETE CASCADE,
            NOTE        VARCHAR(255),
            LANG_ISO    VARCHAR(2)      REFERENCES LANGUAGES(ISO) ON DELETE CASCADE
        );
    `
    return schema
}
