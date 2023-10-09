ALTER TABLE channels
ADD COLUMN createdAt timestamp DEFAULT now();

ALTER TABLE channels
ADD COLUMN createdBy UUID;

ALTER TABLE channels
ADD CONSTRAINT fk_chan_createdby
    FOREIGN KEY (createdBy) REFERENCES users(id)