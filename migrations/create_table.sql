CREATE TABLE wallets
(
    id SERIAL PRIMARY KEY,
    amount NUMERIC(15,2)
);

INSERT INTO wallets
VALUES
    (1, 100),
    (2, 200),
    (3, 300);
