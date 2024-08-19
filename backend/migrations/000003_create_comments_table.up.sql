CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY,
    content TEXT NOT NULL,
    author_id UUID NOT NULL,
    post_id UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (author_id)
    REFERENCES users(id)
    ON DELETE CASCADE,
    FOREIGN KEY (post_id)
    REFERENCES posts(id)
    ON DELETE CASCADE
);
