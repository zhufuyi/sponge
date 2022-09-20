
create table user_example
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime        null,
    updated_at datetime        null,
    deleted_at datetime        null,
    name       char(50)        not null comment '用户名',
    password   char(100)       not null comment '密码',
    email      char(50)        not null comment '邮件',
    phone      char(30)        not null comment '手机号码',
    avatar     varchar(200)    null comment '头像',
    age        tinyint         not null comment '年龄',
    gender     tinyint         not null comment '性别，1:男，2:女，其他值:未知',
    status     tinyint         not null comment '账号状态，1:未激活，2:已激活，3:封禁',
    login_at   bigint unsigned not null comment '登录时间戳',
    constraint user_email_uindex
        unique (email)
);

create index user_example_deleted_at_index
    on user_example (deleted_at);
