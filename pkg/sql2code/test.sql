create table user
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime        null,
    updated_at datetime        null,
    deleted_at datetime        null,
    name       char(50)        not null comment 'username',
    password   char(100)       not null comment 'password',
    email      char(50)        not null comment 'email',
    phone      bigint unsigned not null comment 'phone number',
    age        tinyint         not null comment 'age',
    gender     tinyint         not null comment 'gender, 1:male, 2:female, 3:unknown',
    constraint user_email_uindex
        unique (email)
);
