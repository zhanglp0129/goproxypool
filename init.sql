-- 数据库初始化脚本

-- 代理地址表
create table proxy_address (
    -- 主键
    id integer primary key autoincrement,
    -- ip地址，可以是ipv4和ipv6
    ip text not null,
    -- 端口
    port integer not null,
    -- 协议，一般为http或socks
    protocol text not null,
    -- 上一次检测可用性的时间，为纳秒级时间戳
    detect_time integer null,
    -- 检测失败次数
    failure_number integer default 0,
    -- 禁止重复ip 端口 协议
    unique(ip, port, protocol)
);
