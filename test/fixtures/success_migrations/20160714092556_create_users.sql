-- +migrate Up
CREATE TABLE users (
  id         INT PRIMARY KEY AUTO_INCREMENT,
  name       VARCHAR(255),
  email      VARCHAR(255),
  created_at TIMESTAMP);

-- user table index of email
ALTER TABLE users
  ADD INDEX index_users_email(email);

-- +migrate Down
DROP TABLE users;

-- TODO: drop test comment
