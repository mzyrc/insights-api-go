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

CREATE TABLE tweet (
                       id              bigint      NOT NULL
                           CONSTRAINT tweet_pk
                               PRIMARY KEY,
                       text            text        NOT NULL,
                       user_id         bigint      NOT NULL,
                       created_at      timestamptz NOT NULL,
                       favourite_count int             DEFAULT NULL,
                       retweet_count   int             DEFAULT NULL,
                       positive_score  numeric(20, 19) DEFAULT NULL,
                       negative_score  numeric(20, 19) DEFAULT NULL,
                       neutral_score   numeric(20, 19) DEFAULT NULL,
                       mixed_score     numeric(20, 19) DEFAULT NULL
);

CREATE INDEX tweet_user_id_index
    ON tweet (user_id);

ALTER TABLE tweet
    OWNER TO insights_user;

CREATE TABLE tweet_synchronisation (
                                       last_tweet_id   bigint                    NOT NULL
                                           CONSTRAINT tweet_synchronisation_tweet_id_fk
                                               REFERENCES tweet,
                                       user_id         bigint                    NOT NULL,
                                       synchronised_at timestamptz DEFAULT NOW() NOT NULL
);

CREATE UNIQUE INDEX tweet_synchronisation_last_tweet_id_uindex
    ON tweet_synchronisation (last_tweet_id);

ALTER TABLE tweet_synchronisation
    OWNER TO insights_user;