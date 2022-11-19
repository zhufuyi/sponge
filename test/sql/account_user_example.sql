
create table user_example
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime        null,
    updated_at datetime        null,
    deleted_at datetime        null,
    name       char(50)        not null comment 'username',
    password   char(100)       not null comment 'password',
    email      char(50)        not null comment 'email',
    phone      char(30)        not null comment 'phone number',
    avatar     varchar(200)    null comment 'avatar',
    age        tinyint         not null comment 'age',
    gender     tinyint         not null comment 'gender, 1:Male, 2:Female, other values:unknown',
    status     tinyint         not null comment 'account status, 1:inactive, 2:activated, 3:blocked',
    login_at   bigint unsigned not null comment 'login timestamp',
    constraint user_email_uindex
        unique (email)
);

create index user_example_deleted_at_index
    on user_example (deleted_at);
