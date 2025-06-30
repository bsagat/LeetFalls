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

INSERT INTO Users (Name, Token_ID, Avatar_URL) VALUES
('Dipper Pines', 'token_dipper_123', 'http://127.0.0.1:9090/objects/characters/Dipper_Pines.png'),
('Mabel Pines', 'token_mabel_123', 'http://127.0.0.1:9090/objects/characters/mabel_pines.png'),
('Bill Cipher', 'token_bill_123', 'http://127.0.0.1:9090/objects/characters/Bill_Cipher.png');

INSERT INTO Posts (ImageURL, Title, Content, Author_id) VALUES
('None', 'Journal 3 Found!', 'Guys, I finally found Journal 3! There are things in here that shouldn’t exist...', 1),
('None', 'Glitter Bomb Attack 💥', 'I may or may not have filled Waddles’ food bowl with glitter. Worth it.', 2),
('None', 'The End Is Near', 'Time is an illusion. So is your free will.', 3);

INSERT INTO Comments (Post_id, Author_id, Content, ImageURL) VALUES
(1, 2, 'OMG! Is there anything about unicorns in there?', NULL),
(1, 3, 'Careful with that, kid. Knowledge is... unstable.', '/static/images/comments/eye.png'),
(2, 1, 'Mabel... not again. Poor Waddles!', NULL);

