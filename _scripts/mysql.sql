drop database if exists card;
create database card character set utf8mb4 collate utf8mb4 unicode ci;
use card;

drop table if exists card;
create table card
(
    id          bigint auto_increment primary key comment '自增ID',
    created     datetime default current_timestamp comment '创建时间',
    updated     datetime on update current_timestamp comment '更新时间',
    num     varchar(36) not null comment '卡号',
    code   varchar(36) not null comment '卡密',
    state int not null default 0 comment '状态。 0:未用 1:预占 2:已用'
) engine = innodb
  default charset = utf8mb4 comment '重置卡表';

create unique index idx_card_num on card(num);
create unique index idx_card_code on card(code);

drop table if exists seq;
create table seq
(
  name        varchar(32) primary key comment '序列名称',
  num         bigint not null default 0 comment '序列值',
  created     datetime default current_timestamp comment '创建时间',
  updated     datetime on update current_timestamp comment '更新时间'
) engine = innodb
  default charset = utf8mb4 comment '序列表';

