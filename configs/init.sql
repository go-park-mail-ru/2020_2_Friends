CREATE TABLE IF NOT EXISTS users (
    id SERIAL NOT NULL PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
<<<<<<< HEAD
);

CREATE TABLE IF NOT EXISTS profiles (
    userID INTEGER NOT NULL,
    username TEXT,
    phone TEXT,
    addresses TEXT[],
    points INTEGER

    FOREIGN KEY (userID) REFERENCES users (id)
=======
>>>>>>> a828960e68c444badd0293b31c9c0fda0cb9553a
);
