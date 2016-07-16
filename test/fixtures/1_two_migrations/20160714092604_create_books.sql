-- +migrate Up
CREATE TABLE books(id INT PRIMARY KEY AUTO_INCREMENT, title VARCHAR(255), author VARCHAR(255), created_at TIMESTAMP);


-- +migrate Down
DROP TABLE books;
