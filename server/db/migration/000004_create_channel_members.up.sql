CREATE TABLE IF NOT EXISTS channel_members (
    channelId UUID NOT NULL,
    userId UUID NOT NULL,
    creator BOOLEAN DEFAULT FALSE,
    createdAt timestamp DEFAULT NOW(),

    CONSTRAINT fk_user
        FOREIGN KEY (userId) REFERENCES users(id),
    CONSTRAINT fk_channel
        FOREIGN KEY (channelId) REFERENCES channels(id)
)
