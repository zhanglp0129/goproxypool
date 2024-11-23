package sqlite

import _ "embed"

// 初始化SQL语句
//
//go:embed sql/init.sql
var initSql string
