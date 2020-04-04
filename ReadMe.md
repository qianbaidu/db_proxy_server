# db proxy server
通过实现Mysql、Mongodb协议完成登录验证，登录后代理到后端真实mysql、mongdob资源地址代理查询。
并实现用户查询权限管控：基于数据库、表、字段级别

### 功能
- MySQL
    - 实现mysql_native_password协议用户创建、登录
    - 实现数据库、表、字段加密级别权限查询管理
    - 实现mysql查询协议及结果封装返回
    - 执行sql语句（暂禁用insert / update / delete / set 类sql语句）
    - 慢sql查询超时终止并kill查询进程
    - 实现mysql链接检查，断开自动重连

- Mongodb
    - 实现SCRAM-SHA-1协议用户创建、登录
    - 实现数据库、表级别权限查询管理
    - 实现mongodb查询协议及结果封装返回
    - 执行语句(暂禁用insert / update / delete 类sql语句）
    - 实现简易连接池资源复用



### 运行
- 1.创建用户: (username替换要创建的用户名，password替换自己的密码)
```
go run create_user.go -u username -p password
```
```
# 示例
➜  go run create_user.go -u test -p test
INFO[0000] int app config                                source="conf.go:63"
INFO[0000] init load app config success.                 source="conf.go:58"
INFO[0000] init db dns : root:@tcp(127.0.0.1:3306)/db_proxy?charset=utf8  source="models.go:24"
inser user success, id : 2
```
- 2.启动sql解析程序: repo: https://github.com/qianbaidu/sql_parser_server
- 3.启动主程序
```
go run main.go
```
- 4.mongodb登录测试
```
# 示例
➜ mongo -u test -p test 127.0.0.1:4000/test
MongoDB shell version v3.4.0
connecting to: mongodb://127.0.0.1:4000/test
MongoDB server version: 3.4.23
>
```
- 查询测试
```
> show dbs;
test  0.000GB
> use test;
switched to db test
...
```
- 5.mysql登录测试
```
# 示例
➜  mysql -utest -ptest --host=127.0.0.1 --port=3307 -Ddb_proxy
Warning: Using a password on the command line interface can be insecure.
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 10002
Server version: 5.7.0 Homebrew

Copyright (c) 2000, 2016, Oracle and/or its affiliates. All rights reserved.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> show databases;
+----------+
| Database |
+----------+
| db_proxy |
+----------+
1 row in set (0.00 sec)

mysql> use db_proxy
Database changed
mysql> show tables;
+--------------------+
| Tables_in_db_proxy |
+--------------------+
| user               |
+--------------------+
1 row in set (0.00 sec)

mysql> select * from user limit 1 \G;
*************************** 1. row ***************************
                           id: 1
                         user: alex1
               mysql_password: Secret Field
        mysql_read_permission: 1
  mysql_read_write_permission: 0
             mongodb_password: Secret Field
      mongodb_read_permission: 1
mongodb_read_write_permission: 0
1 row in set (0.01 sec)

```




### todo
- 执行权限功能实现完善、并开启
- mysql管理客户端连接兼容支持（目前部分支持）
- 用户连接数管控
- mysql连接建立连接池资源复用
- 编写测试用例
- 管理界面
    - 登录
    - 用户添加
        - 添加用户设置初始密码；并生成代理登录时效验用的密码
    - 库权、表、加密字段权限配置

