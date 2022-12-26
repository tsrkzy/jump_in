-- チェック用コマンド
--
-- $ docker container rm -f db_pg && docker image rm jump_in_pg -f && sh docker_up.sh

drop database if exists "jump_in";
create database "jump_in" encoding ='UTF8' lc_collate ='C' lc_ctype ='C' template "template0";

\c "jump_in"


drop table if exists "account" cascade;
create table "account" (
    "id"         bigserial                not null
        primary key,
    "name"       "varchar"(20)            not null,
    "created_at" timestamp with time zone not null default now(),
    "updated_at" timestamp with time zone not null default now()
);


-- event_series にしたかったが、sql_boiler が不可算名詞に対応していない (-ies) のでgroupにした
drop table if exists "event_group" cascade;
create table "event_group" (
    "id"         bigserial                not null
        primary key,
    "name"       "varchar"(100)           not null
        unique,
    "created_at" timestamp with time zone not null default now(),
    "updated_at" timestamp with time zone not null default now()
);


drop table if exists "event" cascade;
create table "event" (
    "id"             bigserial                not null
        primary key,
    "name"           "varchar"(100)           not null,
    "account_id"     bigserial                not null references "account"("id") on delete set null,
    "event_group_id" bigserial                references "event_group"("id") on delete set null,
    "created_at"     timestamp with time zone not null default now(),
    "updated_at"     timestamp with time zone          default null
);


drop table if exists "mail_account" cascade;
create table "mail_account" (
    "id"           bigserial                not null
        primary key,
    "account_id"   bigserial                not null
        references "account"("id") on delete cascade,
    "mail_address" "varchar"(100)           not null
        unique,
    "created_at"   timestamp with time zone not null default now(),
    "updated_at"   timestamp with time zone          default null
);

drop table if exists "invitation" cascade;
create table "invitation" (
    "id"                  bigserial                not null
        primary key,
    "uri_hash"            "varchar"(100)           not null
        unique,
    "choco_chip"          "varchar"(100)           not null
        unique,

    "mail_account_id"     bigserial                not null
        references "mail_account"("id") on delete cascade,
    "ip_address"          "varchar"(100)           not null,
    "redirect_uri"        "varchar"(100)           not null,
    "expired_datetime"    timestamp with time zone not null,
    "authorised"          "bool"                   not null default false,
    "authorised_datetime" timestamp with time zone          default null,
    "created_at"          timestamp with time zone not null default now(),
    "updated_at"          timestamp with time zone          default null
);

drop table if exists "attend" cascade;
create table "attend" (
    id         BIGSERIAL                NOT NULL
        PRIMARY KEY,
    account_id BIGSERIAL                NOT NULL
        REFERENCES account(id) ON DELETE CASCADE,
    event_id   BIGSERIAL                NOT NULL
        REFERENCES event(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE (account_id, event_id)
);

DROP TABLE IF EXISTS candidate CASCADE;
CREATE TABLE candidate (
    id         BIGSERIAL                NOT NULL
        PRIMARY KEY,
    event_id   BIGSERIAL                NOT NULL
        REFERENCES attend(id) ON DELETE CASCADE,
    open_at    TIMESTAMP WITH TIME ZONE NOT NULL,
    close_at   TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);


DROP TABLE IF EXISTS event_create_mail CASCADE;
CREATE TABLE event_create_mail (
    id         BIGSERIAL                NOT NULL
        PRIMARY KEY,
    event_id   BIGSERIAL                NOT NULL
        REFERENCES event(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

DROP TABLE IF EXISTS event_announce_mail CASCADE;
create table "event_announce_mail" (
    "id"         bigserial                not null
        primary key,
    "event_id"   bigserial                not null
        references "event"("id") on delete cascade,
    "created_at" timestamp with time zone not null default now(),
    "updated_at" timestamp with time zone not null default now()
);

