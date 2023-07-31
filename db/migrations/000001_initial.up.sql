begin;
create table if not exists users (
  id serial primary key,
  password varchar(255) not null,
  login varchar(255) not null unique,
  name varchar(255) not null unique,
  age int not null
);
create table if not exists sessions (
  id serial primary key,
  user_id int not null,
  active boolean default true,
  created_at timestamp with time zone default now() not null,
  token varchar(255) not null,

  constraint fk_user_id foreign key(user_id) references users(id) on delete cascade
);
create table if not exists user_phones (
  id serial primary key,
  user_id int not null,
  phone varchar(255) not null,
  description varchar(255) not null,
  is_fax boolean default false,

  constraint fk_user_id foreign key(user_id) references users(id) on delete cascade,
  constraint uq_user_id_phone unique(user_id, phone)
);
commit;
