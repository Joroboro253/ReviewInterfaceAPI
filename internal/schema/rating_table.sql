CREATE TABLE review_ratings (
    id SERIAL PRIMARY KEY,
    review_id INT,
    rating FLOAT,
    user_id INT,
    FOREIGN KEY (review_id) REFERENCES reviews(id)
);
