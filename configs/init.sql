CREATE TABLE IF NOT EXISTS users (
    id SERIAL NOT NULL PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
);

CREATE TABLE IF NOT EXISTS profiles (
    userID INTEGER NOT NULL,
    name TEXT,
    phone TEXT,
    addresses TEXT[],
    points INTEGER,

    FOREIGN KEY (userID) REFERENCES users (id)
);