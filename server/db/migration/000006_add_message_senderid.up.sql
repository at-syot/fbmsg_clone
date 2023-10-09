ALTER TABLE messages
ADD COLUMN senderId UUID NOT NULL;

ALTER TABLE messages
ADD CONSTRAINT fk_chan_msg_user
    FOREIGN KEY (senderId) REFERENCES users(id);