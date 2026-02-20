create table if not exists users (
    id serial primary key,
    name varchar(255) not null,
    email varchar(255) not null,
    age int not null,
    city varchar(255) not null
);

insert into users (name, email, age, city)
values ('John Doe', 'john@example.com', 30, 'New York');