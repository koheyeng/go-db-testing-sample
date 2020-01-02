CREATE TABLE sample
(
    id SERIAL PRIMARY KEY NOT NULL,
    name varchar(32) NOT NULL,
    email varchar(32) NOT NULL,

    deleted_at timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
)