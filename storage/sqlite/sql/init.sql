-- 数据库初始化脚本

-- 代理地址表
create table if not exists proxy_addresses (
    -- 主键
    id integer primary key autoincrement,
    -- ip地址，可以是ipv4和ipv6
    ip text not null,
    -- 端口
    port integer not null,
    -- 协议
    protocol text not null,
    -- 上一次检测可用性的时间，为纳秒级时间戳
    detect_time integer default 0,
    -- 检测失败次数
    failure_number integer default 0,
    -- 上次使用的时间，为纳秒级时间戳
    used_time integer default 0,
    -- 禁止重复ip 端口 协议
    unique(ip, port, protocol)
);

-- 为检测时间字段添加索引
create index if not exists idx_detect_time on proxy_addresses(detect_time);

-- 为使用时间创建索引
create index if not exists idx_used_time on proxy_addresses(used_time);
