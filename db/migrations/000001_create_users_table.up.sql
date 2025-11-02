CREATE TABLE users (
   id UUID PRIMARY KEY,
   name TEXT NOT NULL,
   surname TEXT NOT NULL,
   nickname TEXT NOT NULL,
   age INT NOT NULL CHECK (age >= 0),
   email TEXT NOT NULL UNIQUE,
   password TEXT NOT NULL,
   created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
   updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
   deleted_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_users_email ON users(email);
