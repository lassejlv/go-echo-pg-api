create table users (
  id serial primary key,
  username text unique not null,
  password_hash text not null
  created_at timestamptz not null default now()
  updated_at timestamptz not null default now()
)