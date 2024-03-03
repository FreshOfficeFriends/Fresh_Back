CREATE TABLE IF NOT EXISTS users
(
  id serial PRIMARY KEY,
  first_name VARCHAR(50) NOT NULL,
  second_name VARCHAR(50) NOT NULL,
  email VARCHAR(50) unique NOT NULL,
  password varchar default NULL,
  birthday date NOT NULL,
  registered_at timestamp default NOW()
);

create index if not exists idx_email on users (email);

CREATE TABLE IF NOT EXISTS refresh_tokens
(
  id serial PRIMARY KEY,
  user_id int NOT NULL,
  token varchar NOT NULL,
  expires_at timestamp NOT NULL
);
