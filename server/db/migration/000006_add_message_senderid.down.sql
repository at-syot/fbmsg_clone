ALTER TABLE messages
DROP CONSTRAINT fk_chan_msg_user,
DROP COLUMN senderId;
