-- CREATE DATABASE gotodos;
-- CREATE USER gotodos WITH PASSWORD 'gotodos';
-- GRANT ALL PRIVILEGES ON DATABASE "gotodos" to gotodos;
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  first_name VARCHAR(255),
  last_name VARCHAR(255),
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL
);

CREATE TABLE todos (
  id SERIAL PRIMARY KEY,
  title TEXT,
  note TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  modified_at TIMESTAMP NOT NULL DEFAULT NOW(),
  due_at TIMESTAMP,
  completed_at TIMESTAMP,
  user_id INTEGER REFERENCES users(id),
  is_done BOOLEAN NOT NULL DEFAULT FALSE
);

-- @TODO: Implement a todo tagging system
-- 
-- CREATE TABLE tags(
--   id SERIAL PRIMARY KEY,
--   name VARCHAR(50) NOT NULL
-- );
-- 
-- CREATE TABLE todos_tags(
--   tag_id INTEGER REFERENCES tags(id),
--   todo_id INTEGER REFERENCES todos(id)
-- );
