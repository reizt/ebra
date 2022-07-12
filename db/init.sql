ALTER USER 'root'@'%' PASSWORD EXPIRE;
CREATE DATABASE IF NOT EXISTS `mysql_test_db`;
CREATE USER 'mysql_test_user'@'%' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON mysql_test_db.* TO 'mysql_test_user'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;
