CREATE TABLE ratings (
                         id UUID PRIMARY KEY,
                         book_id UUID REFERENCES books(id) ON DELETE CASCADE,
                         user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                         rating INTEGER NOT NULL CHECK (rating >= 0 AND rating <= 10),
                         comment TEXT,
                         created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                         UNIQUE (book_id, user_id)
);

-- Индекс для быстрого поиска рейтингов по книгам
CREATE INDEX idx_ratings_book_id ON ratings(book_id);

-- Индекс для быстрого поиска рейтингов по пользователям
CREATE INDEX idx_ratings_user_id ON ratings(user_id);


CREATE OR REPLACE VIEW book_average_ratings AS
SELECT
    b.id AS book_id,
    b.name,
    b.author,
    COALESCE(AVG(r.rating), 0) AS average_rating,
    COUNT(r.id) AS rating_count
FROM
    books b
        LEFT JOIN
    ratings r ON b.id = r.book_id
GROUP BY
    b.id, b.name, b.author;