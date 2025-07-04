create table users (
  id SERIAL primary key,
  name VARCHAR(25) not null,
  email VARCHAR(320) not null unique,
  username VARCHAR(25) not null unique,
  hashed_password VARCHAR(100) not null
);

create index email_idx_users on users (email);

create index username_idx_users on users (username);
