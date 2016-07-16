cat << EOF > test/rc/mysql.yml
mysql:
  host: $MYSQL_PORT_3306_TCP_ADDR
  port: $MYSQL_PORT_3306_TCP_PORT
  user: $MYSQL_ENV_MYSQL_USER
  password: $MYSQL_ENV_MYSQL_PASSWORD
  database: $MYSQL_ENV_MYSQL_DATABASE
EOF
