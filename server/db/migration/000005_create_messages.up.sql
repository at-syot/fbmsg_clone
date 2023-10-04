CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    channelId UUID NOT NULL,
    content VARCHAR NOT NULL,
    createdAt timestamp DEFAULT NOW(),

    CONSTRAINT fk_chan_msg
        FOREIGN KEY (channelId) REFERENCES channels(id)
)