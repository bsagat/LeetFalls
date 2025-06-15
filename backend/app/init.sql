CREATE TABLE IF NOT EXISTS Users (
    ID SERIAL PRIMARY KEY,
    Name Text NOT NULL,
    Token_ID TEXT NOT NULL,
    Avatar_URL TEXT NOT NULL,
    TokenDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    Expires_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP + INTERVAL '2 weeks'
);


CREATE TABLE IF NOT EXISTS Posts (
    ID SERIAL PRIMARY KEY,
    ImageURL TEXT,
    Title TEXT NOT NULL,
    Content TEXT NOT NULL,
    Author_id INT REFERENCES Users(ID) ON DELETE CASCADE,
    Created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    Expires_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP + INTERVAL '15 minutes'
);

-- Добавлен Reply to ,нужно обновить erd 
CREATE TABLE IF NOT EXISTS Comments (
    ID SERIAL PRIMARY KEY,
    Post_id INT REFERENCES Posts(ID) ON DELETE CASCADE,
    Author_id INT REFERENCES Users(ID) ON DELETE CASCADE, 
    Created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    Content TEXT NOT NULL,
    ImageURL TEXT,
    Reply_to INT 
);

CREATE INDEX idx_token_id ON Users(Token_ID);
CREATE INDEX idx_expires_at_users ON Users(Expires_at);

CREATE INDEX idx_author_id_posts ON Posts(Author_id);
CREATE INDEX idx_expires_at_posts ON Posts(Expires_at);

CREATE INDEX idx_post_id ON Comments(Post_id);
CREATE INDEX idx_author_id_comments ON Comments(Author_id);

