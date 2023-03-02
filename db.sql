CREATE USER 'myuser'@'localhost' IDENTIFIED BY 'root66';
CREATE DATABASE mydb;
GRANT ALL PRIVILEGES ON mydb.* TO 'myuser'@'%' IDENTIFIED BY 'root66';
GRANT ALL PRIVILEGES ON mydb.* TO 'myuser'@'localhost' IDENTIFIED BY 'root66';
USE mydb
CREATE TABLE Songs (
  `SongID` int NOT NULL AUTO_INCREMENT,
  `Name` varchar(255),
  `Duration` int,
  PRIMARY KEY (`SongID`)
);
