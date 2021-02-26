CREATE DATABASE insights;

CREATE USER insights_user WITH PASSWORD 'abcd1234';

-- usr
CREATE TABLE usr (
                     id         serial                    NOT NULL
                         CONSTRAINT usr_pk
                             PRIMARY KEY,
                     username   text                      NOT NULL,
                     created_at timestamptz DEFAULT NOW() NOT NULL
);

ALTER TABLE usr
    OWNER TO insights_user;

-- user_following
CREATE TABLE IF NOT EXISTS usr_following (
                                             user_id         integer NOT NULL
                                             CONSTRAINT usr_following_usr_id_fk
                                             REFERENCES usr,
                                             twitter_user_id bigint  NOT NULL,
                                             created_at      timestamp WITH TIME ZONE DEFAULT NOW()
    );

ALTER TABLE usr_following
    OWNER TO insights_user;

CREATE UNIQUE INDEX IF NOT EXISTS usr_following_user_id_twitter_user_id_uindex
    ON usr_following (user_id, twitter_user_id);

INSERT INTO usr (username)
VALUES ('db998');