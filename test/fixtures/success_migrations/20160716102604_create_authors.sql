-- +migrate Up
CREATE TABLE authors (
  id         INT PRIMARY KEY AUTO_INCREMENT,
  name       VARCHAR(255),
  created_at TIMESTAMP
);

ALTER TABLE books
  DROP COLUMN author;

ALTER TABLE books
  ADD (author_id INT NOT NULL);

ALTER TABLE books
  ADD CONSTRAINT fk_books_author_id FOREIGN KEY (author_id) REFERENCES authors (id);

-- +migrate Down
DROP TABLE authors;
ALTER TABLE books
  ADD (author VARCHAR(255));

