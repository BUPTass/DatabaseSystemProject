[TOC]



#### 程序结构说明

main.go为唯一入口。

controllers文件夹为各个组件并构成controllers包，其中database.go为数据库相关模块接口，uesr.go定义用户注册、登陆相关模块接口。

docs为文档文件夹。

statics存放各个静态数据，其中session_data文件夹下存放各个用户登陆对应session。

#### 用户登陆控制

程序使用session管理登陆。session存在~/DBMS/statics/session_data文件夹下。

用户在登陆时会创建一个对应session，名称即为用户名。具体说明如下表：

| 名称      | 说明                                   |
| --------- | -------------------------------------- |
| 查询索引  | username                               |
| id属性    | username                               |
| level属性 | 用户等级（0代表管理员，1代表普通用户） |
| isLogin   | true                                   |
| MaxAge    | 最大存在时间，默认为7天                |

#### 路由说明

##### /login

路由形如`"/login?username=uname&password=pwd"`（例如 `/login?username=root&password=123456`)。

​	提供用户名密码进行登陆

##### /logout

路由形如`"/logout?username=name"`（例如 `/logout?username=root`)。

​	提供用户名进行登出

##### /show/users

路由形如`"/show/users?adminname=name"` （例如 `/show/users?adminname=root`)

​	管理员在登陆状态下查看所有用户（及其状态）。返回一个json表，其内容为用户状态，表项结构如下

```go
type UserInfo struct {
  UserName  string `db:"username"`
  Password  string `db:"password"`
  Level   int   `db:"level"`   //0 for admin,1 for normal user
  Conformed int   `db:"conformed"` //false for unavailable
}
```

##### /add/user

路由形如`"/add/user?adminname=name&username=name"`（例如 `"/add/user?adminname=root&username=test1"`)。

​	管理员在登陆状态下添加指定用户（将用户的Conformed改为true）。

##### /delete/user

路由形如`"/delete/user?adminname=name&username=name"`（例如 `"/delete/user?adminname=root&username=test1"`)。

​	管理员在登陆状态下删除指定用户（删除其对应记录）。

##### /regist

路由形如`"/regist?username=name&password=password&level=num"`（例如 `"/regist?username=test1&password=123456&level=1"`)。

​	用户注册。提供用户名（唯一），密码，用户等级，从而生成一条Conformed为false的记录，经管理员添加之后成功。

##### /manage/databaseInfo

路由形如`"/manage/databaseInfo?adminname=name&item=itemname&condition=string"`

​	管理员在登陆后获取数据库信息。condition有三种情况，

1. condition=databaseinfo
   返回一个json表，结构如下：

   ```go
   		struct {
   			Username          string
   			UserDbName        string
   			UserinfoTableName string
   			Ip                string
   			DataDbName        string
   			Port              int
   		}
   ```

2. condition为空字符串
   相当于在mysql上运行 `show ...`(其中...的内容即为item字符串)并返回json表，表项如下

   ```go
   type databaseInfo struct {
   	Name  string
   	Value string
   }
   ```

3. 其他
   相当于在mysql上运行 `show ... like ‘condition’`(其中...的内容即为item字符串，condition即为指定过滤条件)并返回json表，表项如下

   ```go
   type databaseInfo struct {
   	Name  string
   	Value string
   }
   ```

   

##### /manage/databaseConnection

路由形如`"/manage/databaseConnection?adminname=name&condition=string"`

​	管理员在登陆后获取数据库连接信息。condition有三种情况，

1. condition=ip
   查看当前连接中各个IP的连接数
   相当于在mysql上运行`select SUBSTRING_INDEX(host,':',1) as ip , count(*) from information_schema.processlist group by ip`。返回各个ip的连接数，为一个json表，表项如下

   ```go
   		type info struct {
   			ipname string
   			cnt    int
   		}
   ```

2. condition=user
   查看当前连接中各个用户的连接数

   相当于在mysql上运行`select USER , count(*) from information_schema.processlist group by USER`。返回各个用户的连接数，为一个json表，表项如下

   ```go
   		type info struct {
   			user string
   			cnt  int
   		}
   ```

3. condition=list
   查看当前数据库的连接情况
   相当于在mysql上运行`show full processlist`。返回连接情况json表，表项如下

   ```go
   		type allinfo struct {
   			id      int
   			user    string
   			host    string
   			db      string
   			command string
   			time    int
   			state   string
   			info    string
   		}
   ```

##### /manage/databse

路由形如`"/manage/database?adminname=name&item=itemname&value=value"`

​	管理员在登陆状态下设置数据库参数，相当于在mysql上运行`set itemname=value`（其中itemname及value给定）