# DatabaseSystemProject

用户信息存在于userinfo的info表中。第一次启动前需要数据库提前创建并分配一个用户，

数据库定义如下

```sql
create database userinfo;
use userinfo;
create table info(
username varchar(255) primary key,
password char(60) not null,
level int,
conformed int);
```

session存在于statics/session_data文件夹内，需要预先定义好static/session_data这两个文件夹，不然可能会报错。