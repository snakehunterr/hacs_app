create database if not exists `MYSQL_DATABASE`;
use `MYSQL_DATABASE`;

create table if not exists client (
    client_id bigint not null,
    client_name varchar(100),
    is_admin boolean not null,
    last_edited timestamp not null default current_timestamp,
    primary key (client_id)
);

create table if not exists room (
    room_id int not null,
    client_id bigint not null,
    room_people_count int not null,
    room_area float not null,
    last_edited timestamp not null default current_timestamp,
    primary key (room_id),
    foreign key (client_id) references client(client_id) on delete cascade
);

create table if not exists payment (
    payment_id int not null auto_increment,
    client_id bigint not null,
    room_id int not null,
    payment_date timestamp not null default current_timestamp,
    payment_amount float not null,
    last_edited timestamp not null default current_timestamp,
    primary key(payment_id),
    foreign key (client_id) references client(client_id),
    foreign key (room_id) references room(room_id)
);

create table if not exists expense (
    expense_id int not null auto_increment,
    expense_date timestamp not null default current_timestamp,
    expense_amount float not null,
    last_edited timestamp not null default current_timestamp,
    primary key (expense_id)
);
