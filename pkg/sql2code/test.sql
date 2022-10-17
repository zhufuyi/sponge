create table user
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime        null,
    updated_at datetime        null,
    deleted_at datetime        null,
    name       char(50)        not null comment '用户名',
    password   char(100)       not null comment '密码',
    email      char(50)        not null comment '邮件',
    phone      bigint unsigned not null comment '手机号码',
    age        tinyint         not null comment '年龄',
    gender     tinyint         not null comment '性别，1:男，2:女，3:未知',
    constraint user_email_uindex
        unique (email)
);