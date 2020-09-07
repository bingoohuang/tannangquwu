create database card;
use card;

drop table if exists card;
create table card
(
    id          bigint auto_increment primary key comment '自增ID',
    created     datetime default current_timestamp comment '创建时间',
    updated     datetime on update current_timestamp comment '更新时间',
    num     varchar(32) not null comment '卡号',
    code   varchar(32) not null comment '卡密'
) engine = innodb
  default charset = utf8mb4 comment '重置卡表';

create unique index idx_card_no on card(card_no);
create unique index idx_card_code on card(card_code);

drop table if exists seq;
create table seq 
(
  name        varchar(32) primary key comment '序列名称',
  num         bigint not null default 0 comment '序列值',
  created     datetime default current_timestamp comment '创建时间',
  updated     datetime on update current_timestamp comment '更新时间's
) engine = innodb
  default charset = utf8mb4 comment '序列表';

