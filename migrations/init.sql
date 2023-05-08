CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

drop table if exists data, processed_data;

create table if not exists data (
    id serial primary key,
    n integer not null,
    mqtt text,
    invid varchar(8) not null,
    unit_guid uuid not null,
    msg_id text not null,
    msg_text text not null,
    context text,
    class varchar(7) not null,
    level integer not null,
    area varchar(5) not null,
    addr text not null,
    block text,
    type text,
    bit text,
    invert_bit text
);

create table if not exists processed_data (
    id serial primary key,
    file_name text not null,
    is_processed bool not null default false,
    date_processed timestamp
)