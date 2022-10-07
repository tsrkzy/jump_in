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
    "id"         bigserial                not null
        primary key,
    "account_id" bigserial                not null
        references "account"("id") on delete cascade,
    "event_id"   bigserial                not null
        references "event"("id") on delete cascade,
    "created_at" timestamp with time zone not null default now(),
    "updated_at" timestamp with time zone not null default now()
);

drop table if exists "wish_date" cascade;
create table "wish_date" (
    "id"         bigserial                not null
        primary key,
    "attend_id"  bigserial                not null
        references "attend"("id") on delete cascade,
    "created_at" timestamp with time zone not null default now(),
    "updated_at" timestamp with time zone not null default now()
);


drop table if exists "event_create_mail" cascade;
create table "event_create_mail" (
    "id"         bigserial                not null
        primary key,
    "event_id"   bigserial                not null
        references "event"("id") on delete cascade,
    "created_at" timestamp with time zone not null default now(),
    "updated_at" timestamp with time zone not null default now()
);

drop table if exists "event_announce_mail" cascade;
create table "event_announce_mail" (
    "id"         bigserial                not null
        primary key,
    "event_id"   bigserial                not null
        references "event"("id") on delete cascade,
    "created_at" timestamp with time zone not null default now(),
    "updated_at" timestamp with time zone not null default now()
);

