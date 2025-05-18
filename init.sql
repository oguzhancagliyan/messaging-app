CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    "to" TEXT NOT NULL,
    content TEXT NOT NULL CHECK (char_length(content) <= 160),
    sent BOOLEAN NOT NULL DEFAULT FALSE,
    sent_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT now()
);

INSERT INTO messages ("to", content,sent) VALUES
('+905551111111', 'Sample message 1',false),
('+905551111112', 'Sample message 2',false),
('+905551111113', 'Sample message 3',false),
('+905551111114', 'Sample message 4',false),
('+905551111115', 'Sample message 5',false),
('+905551111116', 'Sample message 6',false),
('+905551111117', 'Sample message 7',false),
('+905551111118', 'Sample message 8',false),
('+905551111119', 'Sample message 9',false),
('+905551111120', 'Sample message 10',false);
