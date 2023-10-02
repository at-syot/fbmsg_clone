CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT uuid_generate_v4(),
    username VARCHAR NOT NULL,
    PRIMARY KEY (id)
)
