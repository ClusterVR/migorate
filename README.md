migorate
===

Simple database migration tool for Go.

[![wercker status](https://app.wercker.com/status/d56e8c9b3a5e5aa6d81d4f8c9c74a4ff/m "wercker status")](https://app.wercker.com/project/bykey/d56e8c9b3a5e5aa6d81d4f8c9c74a4ff)

## Features
- Migration management
- Execute SQL command from \*.sql files
- Generate SQL file with Rails-like commands

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
$ migorate generate [migration id] [col:type]...
```

```shell
$ migorate generate create_users id:id name:string login_count:integer last_login_at:datetime created_at:timestamp
2016/07/17 22:55:30 Generated: db/migrations/20160717225530_create_users.sql

$ cat db/migrations/20160717225530_create_users.sql
-- +migrate Up
CREATE TABLE users(id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(255), login_count INT, last_login_at DATETIME, created_at TIMESTAMP);

-- +migrate Down
DROP TABLE users;
```

Currently, only `CREATE TABLE` migration can be generated with migration id `create_(tablename)`.

All other migration id generates empty migration file.

#### Type conversion
| in command | SQL |  |
|---|---|---|
| id | INT PRIMARY KEY AUTO_INCREMENT |  |
| integer | INT |  |
| datetime | DATETIME |  |
| string | VARCHAR(255) |  |
| timestamp | TIMESTAMP |  |
| references | INT | Adds foreign key to `id` in `(column_name)s` table |

### Plan migration
```shell
$ migorate plan
2016/07/17 21:16:53 Planned migrations:
2016/07/17 21:16:53   1: 20160717225530_create_users
```

### Execute migration
```shell
$ migorate exec
2016/07/17 21:17:49 Executing 20160717225530_create_users
2016/07/17 21:17:49   CREATE TABLE users(id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(255), login_count INT, last_login_at DATETIME, created_at TIMESTAMP);
```

### Rollback
```shell
$ migorate rollback [migration id]
```

```shell
$ migorate rollback 20160717225530_create_users
2016/07/17 21:21:09 Executing 20160717225530_create_users
2016/07/17 21:21:09   DROP TABLE users;
```
