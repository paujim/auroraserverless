CREATE DATABASE TestDB;

CREATE TABLE IF NOT EXISTS TestDB.Profiles (
    ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    FullName VARCHAR (250) NOT NULL,
    Email VARCHAR (255) NOT NULL UNIQUE,
    Phone VARCHAR (255) NOT NULL
);
