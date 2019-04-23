CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    distance INT NOT NULL,
    order_status ENUM('UNASSIGNED', 'TAKEN') NOT NULL
);