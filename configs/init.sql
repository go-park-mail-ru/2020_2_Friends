CREATE TABLE IF NOT EXISTS users (
    id SERIAL NOT NULL PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS profiles (
    userID INTEGER NOT NULL,
    username TEXT,
    phone TEXT,
    addresses TEXT[],
    points INTEGER,
    avatar TEXT,

    FOREIGN KEY (userID) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS vendors (
    id SERIAL NOT NULL PRIMARY KEY,
    vendorName TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL NOT NULL PRIMARY KEY,
    vendorID INTEGER,
    productName TEXT,
    price TEXT,
    picture TEXT,

    FOREIGN KEY (vendorID) REFERENCES vendors (id)
);
