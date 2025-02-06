create table posts (
  id serial primary key,
  title text not null,
  content text not null,
  user_id integer not null references users(id),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
)

