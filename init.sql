CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    "to" TEXT NOT NULL,
    content TEXT NOT NULL CHECK (char_length(content) <= 160),
    sent BOOLEAN NOT NULL DEFAULT FALSE,
    sent_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT now()
);

INSERT INTO messages ("to", content) VALUES
('+905551111111', 'Sample message 1'),
('+905551111112', 'Sample message 2');
