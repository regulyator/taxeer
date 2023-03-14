create table taxeer_user
(
    id               uuid default gen_random_uuid() not null
        constraint taxeer_user_pk
            primary key,
    telegram_user_id varchar                        not null,
    chat_id          BIGINT                         not null
);

create table taxeer_record
(
    id              uuid default gen_random_uuid() not null
        constraint taxeer_record_pk
            primary key,
    taxeer_user_id  uuid                           not null
        constraint taxeer_record_taxeer_user_id_fk
            references taxeer_user,
    date            timestamp                      not null,
    income_value    float8                         not null,
    income_currency varchar                        not null,
    rate            float8                         not null
);