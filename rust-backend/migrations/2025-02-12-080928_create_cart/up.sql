-- Your SQL goes here
CREATE TABLE Cart (
    id SERIAL PRIMARY KEY,
    user_id INTEGER CONSTRAINT fk_users REFERENCES Users(id),
    product_id INTEGER CONSTRAINT fk_products REFERENCES Products(id)
)