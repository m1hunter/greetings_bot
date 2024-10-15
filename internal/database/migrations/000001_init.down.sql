create table if not exists users
(
    id        serial
        primary key,
    chat_id   bigint not null,
    firstname text   not null,
    lastname  text   not null,
    birthday  date   not null,
    password  text   not null,
    username  text   not null
);

create table if not exists subscriptions
(
    id                   serial
        primary key,
    user_id              integer not null
        references users
            on delete cascade,
    subscribed_to        integer not null
        references users
            on delete cascade,
    is_send_notification boolean not null
);

