CREATE TABLE IF NOT EXISTS posts (
    id bigserial PRIMARY KEY,
    name varchar NOT NULL,
    description varchar NOT NULL,
    created_at TIMESTAMP NULL DEFAULT now(),
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    customer_id bigint NOT NULL 
)