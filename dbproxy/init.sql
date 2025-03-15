create database if not exists db;

create table if not exists db.orders_v2 (name text not null) engine = innodb;

insert into
  db.orders_v2 (name)
values
  ('order1');
