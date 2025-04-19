CREATE TABLE ratings (
                         id UUID PRIMARY KEY,
                         book_id UUID NOT NULL,
                         user_id UUID NOT NULL,
                         rating INTEGER NOT NULL CHECK (rating >= 0 AND rating <= 10),
                         comment TEXT,
                         created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                         UNIQUE (book_id, user_id)
);

-- Индекс для быстрого поиска рейтингов по книгам с сортировкой по дате
CREATE INDEX idx_ratings_book_id_created_at ON ratings(book_id, created_at DESC);

-- Индекс для быстрого поиска рейтингов по книгам с сортировкой по значению рейтинга
CREATE INDEX idx_ratings_book_id_rating ON ratings(book_id, rating DESC);

-- Индекс для быстрого поиска рейтингов по пользователям с сортировкой по дате
CREATE INDEX idx_ratings_user_id_created_at ON ratings(user_id, created_at DESC);

-- Составной индекс для быстрого поиска по книге+пользователю
CREATE UNIQUE INDEX idx_ratings_book_user ON ratings(book_id, user_id);

-- Простое представление для рейтингов
CREATE OR REPLACE VIEW book_average_ratings AS
SELECT
    book_id,
    COALESCE(AVG(rating), 0) AS average_rating,
    COUNT(id) AS rating_count
FROM
    ratings
GROUP BY
    book_id;