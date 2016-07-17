migorate
===

Simple database migration tool for Go.

[![wercker status](https://app.wercker.com/status/d56e8c9b3a5e5aa6d81d4f8c9c74a4ff/m "wercker status")](https://app.wercker.com/project/bykey/d56e8c9b3a5e5aa6d81d4f8c9c74a4ff)

## Features
- Migration management
- Execute SQL command from \*.sql files

## Usage
### Create .migoraterc
```
mysql:
  host: localhost
  port: 3306
  user: migorate
  password: migorate
  database: migorate
```

### Generate migration file
```shell
$ migorate generate create_users
2016/07/17 21:13:45 Generated: db/migrations/20160717211345_create_users.sql
```

```sql:db/migrations/20160717211345_create_users.sql
-- +migrate Up
-- Edit your migration SQL
CREATE TABLE users(id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(255));

-- +migrate Down
-- Edit your migration SQL for rollback
DROP TABLE users;
```

### Plan migration
```shell
$ migorate plan
2016/07/17 21:16:53 Planned migrations:
2016/07/17 21:16:53   1: 20160717211345_create_users
```

### Execute migration
```shell
$ migorate exec
2016/07/17 21:17:49 Executing 20160717211345_create_users
2016/07/17 21:17:49   CREATE TABLE users(id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(255));
```

### Rollback
```shell
$ migorate rollback 20160717211345_create_users
2016/07/17 21:21:09 Executing 20160717211345_create_users
2016/07/17 21:21:09   DROP TABLE users;
```
