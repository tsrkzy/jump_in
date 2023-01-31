-- チェック用コマンド
--
-- $ docker container rm -f db_pg && docker image rm jump_in_pg -f && sh docker_up.sh

DROP DATABASE IF EXISTS jump_in;
CREATE DATABASE jump_in ENCODING ='UTF8' LC_COLLATE ='C' LC_CTYPE ='C' TEMPLATE template0;

\c "jump_in"


DROP TABLE IF EXISTS account CASCADE;
CREATE TABLE account (
    id         BIGSERIAL                NOT NULL
        PRIMARY KEY,
    name       VARCHAR(20)              NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);


-- event_series にしたかったが、sql_boiler が不可算名詞に対応していない (-ies) のでgroupにした
DROP TABLE IF EXISTS event_group CASCADE;
CREATE TABLE event_group (
    id         BIGSERIAL                NOT NULL
        PRIMARY KEY,
    name       VARCHAR(100)             NOT NULL
        UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);


DROP TABLE IF EXISTS event CASCADE;
CREATE TABLE event (
    id             BIGSERIAL                NOT NULL
        PRIMARY KEY,
    name           VARCHAR(100)             NOT NULL,
    description    VARCHAR(1000)            NOT NULL DEFAULT '',
    certified      bool                     NOT NULL DEFAULT FALSE,
    is_open        bool                     NOT NULL DEFAULT FALSE,
    account_id     BIGSERIAL                NOT NULL REFERENCES account(id) ON DELETE SET NULL,
    event_group_id BIGSERIAL                REFERENCES event_group(id) ON DELETE SET NULL,
    created_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMP WITH TIME ZONE          DEFAULT NULL
);


DROP TABLE IF EXISTS mail_account CASCADE;
CREATE TABLE mail_account (
    id           BIGSERIAL                NOT NULL
        PRIMARY KEY,
    account_id   BIGSERIAL                NOT NULL
        REFERENCES account(id) ON DELETE CASCADE,
    mail_address VARCHAR(100)             NOT NULL
        UNIQUE,
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP WITH TIME ZONE          DEFAULT NULL
);

DROP TABLE IF EXISTS invitation CASCADE;
CREATE TABLE invitation (
    id                  BIGSERIAL                NOT NULL
        PRIMARY KEY,
    uri_hash            VARCHAR(100)             NOT NULL
        UNIQUE,
    choco_chip          VARCHAR(100)             NOT NULL
        UNIQUE,

    mail_account_id     BIGSERIAL                NOT NULL
        REFERENCES mail_account(id) ON DELETE CASCADE,
    ip_address          VARCHAR(100)             NOT NULL,
    redirect_uri        VARCHAR(100)             NOT NULL,
    expired_datetime    TIMESTAMP WITH TIME ZONE NOT NULL,
    authorised          bool                     NOT NULL DEFAULT FALSE,
    authorised_datetime TIMESTAMP WITH TIME ZONE          DEFAULT NULL,
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP WITH TIME ZONE          DEFAULT NULL
);

DROP TABLE IF EXISTS attend CASCADE;
CREATE TABLE attend (
    id         BIGSERIAL                NOT NULL
        PRIMARY KEY,
    account_id BIGSERIAL                NOT NULL
        REFERENCES account(id) ON DELETE CASCADE,
    event_id   BIGSERIAL                NOT NULL
        REFERENCES event(id) ON DELETE CASCADE,
    comment    VARCHAR(30)              NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE (account_id, event_id)
);

DROP TABLE IF EXISTS candidate CASCADE;
CREATE TABLE candidate (
    id         BIGSERIAL                NOT NULL
        PRIMARY KEY,
    event_id   BIGSERIAL                NOT NULL
        REFERENCES event(id) ON DELETE CASCADE,
    open_at    CHAR(12)                 NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_candidate_event_id_open_at
        UNIQUE (event_id, open_at)
);

DROP TABLE IF EXISTS vote CASCADE;
CREATE TABLE vote (
    id           BIGSERIAL                NOT NULL
        PRIMARY KEY,
    account_id   BIGSERIAL                NOT NULL
        REFERENCES account(id) ON DELETE CASCADE,
    candidate_id BIGSERIAL                NOT NULL
        REFERENCES candidate(id) ON DELETE CASCADE,
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_vote_account_id_candidate_id
        UNIQUE (account_id, candidate_id)
);


DROP TABLE IF EXISTS administrator CASCADE;
CREATE TABLE administrator (
    id            BIGSERIAL                NOT NULL
        PRIMARY KEY,
    account_id    BIGSERIAL                NOT NULL
        REFERENCES account(id) ON DELETE NO ACTION,
    invitation_id BIGSERIAL                NOT NULL
        REFERENCES invitation(id) ON DELETE SET NULL,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_administrator_account_id_invitation_id
        UNIQUE (account_id, invitation_id)
);

DROP TABLE IF EXISTS Consent CASCADE;
CREATE TABLE Consent (
    id               BIGSERIAL                NOT NULL
        PRIMARY KEY,
    administrator_id BIGSERIAL                NOT NULL
        REFERENCES administrator(id) ON DELETE NO ACTION,
    event_id         BIGSERIAL                NOT NULL
        REFERENCES event(id) ON DELETE NO ACTION,
    message          VARCHAR(2000)            NOT NULL DEFAULT '',
    accepted         bool                     NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);