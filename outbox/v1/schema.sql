CREATE TABLE orders (
    id uuid PRIMARY KEY,
    product varchar(255) NOT NULL,
    quantity int NOT NULL
);
