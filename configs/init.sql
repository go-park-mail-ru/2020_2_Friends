CREATE TABLE IF NOT EXISTS users (
    id SERIAL NOT NULL PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role INT NOT NULL CHECK (role > 0 AND role < 3)
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
    vendorName TEXT NOT NULL UNIQUE,
    descript TEXT DEFAULT '' NOT NULL,
    picture TEXT DEFAULT '' NOT NULL
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL NOT NULL PRIMARY KEY,
    vendorID INTEGER,
    productName TEXT,
    price TEXT,
    picture TEXT,

    FOREIGN KEY (vendorID) REFERENCES vendors (id)
);

CREATE TABLE IF NOT EXISTS vendor_partner (
    partnerID INTEGER NOT NULL,
    vendorID INTEGER NOT NULL,

    FOREIGN KEY (partnerID) REFERENCES users (id),
    FOREIGN KEY (vendorID) REFERENCES vendors (id)
);

CREATE TABLE IF NOT EXISTS carts (
    userID INTEGER NOT NULL,
    productID INTEGER NOT NULL,
    vendorID INTEGER NOT NULL,

    FOREIGN KEY (userID) REFERENCES users (id),
    FOREIGN KEY (productID) REFERENCES products (id),
    FOREIGN KEY (vendorID) REFERENCES vendors (id)
);
