CREATE TABLE ratings (
                         id UUID PRIMARY KEY,
                         book_id UUID NOT NULL,
                         user_id UUID NOT NULL,
                         rating INTEGER NOT NULL CHECK (rating >= 0 AND rating <= 10),
                         comment TEXT,
                         created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                         UNIQUE (book_id, user_id)
);

-- Keep the unique composite index for the constraint
CREATE UNIQUE INDEX idx_ratings_book_user ON ratings(book_id, user_id);

-- Keep index for fast retrieval of book ratings sorted by rating
CREATE INDEX idx_ratings_book_id_rating ON ratings(book_id, rating DESC);

-- View for average ratings
CREATE OR REPLACE VIEW book_average_ratings AS
SELECT
    book_id,
    COALESCE(AVG(rating), 0) AS average_rating,
    COUNT(id) AS rating_count
FROM
    ratings
GROUP BY
    book_id;