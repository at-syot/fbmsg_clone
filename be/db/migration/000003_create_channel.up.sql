CREATE TABLE IF NOT EXISTS channels (
    id UUID DEFAULT uuid_generate_v4(),
    displayname VARCHAR,

    PRIMARY KEY (id)
)