-- +migrate Up
CREATE TABLE users (
  id         INT PRIMARY KEY AUTO_INCREMENT,
  name       VARCHAR(255),
  email      VARCHAR(255),
  created_at TIMESTAMP);

ALTER TABLE users
  ADD INDEX index_users_email(email);

-- +migrate Down
DROP TABLE users;
