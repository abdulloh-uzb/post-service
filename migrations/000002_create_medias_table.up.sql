CREATE TABLE IF NOT EXISTS medias(
    id bigserial PRIMARY KEY,
    post_id bigint NOT NULL REFERENCES posts(id),
    name varchar NOT NULL,
    link varchar NOT NULL, 
    type varchar NOT NULL
)
