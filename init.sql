-- 数据库初始化脚本

-- 代理地址表
create table proxy_address (
    -- 主键
    id bigint primary key autoincrement,
    -- ip地址，可以是ipv4和ipv6
    ip varchar(45) not null,
    -- 端口
    port smallint not null,
    -- 协议，一般为http或socks
    protocol varchar(20) not null,
    -- 上一次检测可用性的时间
    detect_time datetime,
    -- 检测失败次数
    failure_number int default 0
);
