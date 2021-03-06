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

    FOREIGN KEY (userID) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS vendors (
    id SERIAL NOT NULL PRIMARY KEY,
    vendorName TEXT NOT NULL UNIQUE,
    descript TEXT DEFAULT '' NOT NULL,
    picture TEXT DEFAULT '' NOT NULL,
    coordinates GEOGRAPHY NOT NULL,
    service_radius INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL NOT NULL PRIMARY KEY,
    vendorID INTEGER,
    productName TEXT DEFAULT '' NOT NULL,
    descript TEXT DEFAULT '',
    price INTEGER NOT NULL,
    picture TEXT DEFAULT '' NOT NULL,

    FOREIGN KEY (vendorID) REFERENCES vendors (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS categories (
    category TEXT NOT NULL UNIQUE
);

INSERT INTO categories (category) VALUES
    ('Завтраки'),
    ('Обеды'),
    ('Супы'),
    ('Десерты');

CREATE TABLE IF NOT EXISTS vendor_categories (
    vendorID INTEGER NOT NULL,
    category TEXT NOT NULL,

    FOREIGN KEY (vendorID) REFERENCES vendors (id) ON DELETE CASCADE,
    FOREIGN KEY (category) REFERENCES categories (category) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS vendor_partner (
    partnerID INTEGER NOT NULL,
    vendorID INTEGER NOT NULL,

    FOREIGN KEY (partnerID) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (vendorID) REFERENCES vendors (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS carts (
    userID INTEGER NOT NULL,
    productID INTEGER NOT NULL,
    vendorID INTEGER NOT NULL,

    FOREIGN KEY (userID) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (productID) REFERENCES products (id) ON DELETE CASCADE,
    FOREIGN KEY (vendorID) REFERENCES vendors (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL NOT NULL PRIMARY KEY,
    userID INTEGER NOT NULL,
    vendorID INTEGER NOT NULL,
    vendorName TEXT NOT NULL,
    createdAt TIMESTAMPTZ NOT NULL,
    clientAddress TEXT NOT NULL,
    orderStatus TEXT DEFAULT '' NOT NULL,
    price INTEGER NOT NULL,
    reviewed BOOLEAN DEFAULT false NOT NULL,

    FOREIGN KEY (userID) REFERENCES users (id),
    FOREIGN KEY (vendorID) REFERENCES vendors (id)
);

CREATE TABLE IF NOT EXISTS products_in_order (
    orderID INTEGER NOT NULL,
    productName TEXT NOT NULL,
    price INTEGER NOT NULL,
    picture TEXT DEFAULT '' NOT NULL,

    FOREIGN KEY (orderID) REFERENCES orders (id)
);

CREATE TABLE IF NOT EXISTS reviews (
    userID INTEGER NOT NULL,
    orderID INTEGER NOT NULL,
    vendorID INTEGER NOT NULL,
    rating INTEGER NOT NULL CHECK (rating > 0 AND rating < 6),
    review_text TEXT DEFAULT '' NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,

    FOREIGN KEY (userID) REFERENCES users (id),
    FOREIGN KEY (orderID) REFERENCES orders (id),
    FOREIGN KEY (vendorID) REFERENCES vendors (id)
);

CREATE TABLE IF NOT EXISTS messages (
    orderID INTEGER NOT NULL,
    userID INTEGER NOT NULL,
    message_text TEXT NOT NULL,
    sent_at TIMESTAMPTZ NOT NULL,

    FOREIGN KEY (orderID) REFERENCES orders (id),
    FOREIGN KEY (userID) REFERENCES users (id)
);
