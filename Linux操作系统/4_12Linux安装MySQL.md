# Python 全栈文档

## 第十二章  Linux安装MySQL

1、下载并安装MySQL官方的 Yum Repository

​        由于CentOS 的yum源中没有mysql，需要到mysql的官网下载

![1604715200059](C:\Users\root\AppData\Roaming\Typora\typora-user-images\1604715200059.png)

  2、然后进行rpm的安装

​      rpm -ivh mysql-community-server-xxxx

​      rpm -e mariadb-libs --nodeps 卸载一个软件包

3、使用yum命令安装依赖：net-tools

   

```bash
[root@localhost yum.repos.d]#   yum install net-tools

```

4、启动MySQL

```bash
mysqld --initialize 先初始化
chown mysql:mysql mysql -R 修改mysql数据文件的拥有者
[root@localhost yum.repos.d]#  systemctl start mysqld.service
# 注意这里的mysqld是服务名，一般默认是mysqld
```

5、获取安装时的临时密码（在第一次登录时就是用这个密码）

```bash
[root@localhost yum.repos.d]# grep 'temporary password' /var/log/mysqld.log
2020-09-09T09:55:29.051013Z 1 [Note] A temporary password is generated for root@localhost: E?/=tU<k!6Ws
```

![](imgs/04_133.png)

6、登录mysql

```bash
[root@localhost yum.repos.d]# mysql -u root -p
Enter password: 
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 5
Server version: 5.7.31

Copyright (c) 2000, 2020, Oracle and/or its affiliates. All rights reserved.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> 


# 注意密码这里没有回显
```

7、登录后修改密码

```
注意：这里会进行密码强度校验（密码设置时必须包含大小写字母、特殊符号、数字，并且长度大于8位）
而且建议使用强口令，特别是从事商业开发项目，养成这个习惯比较好
```

```bash
mysql> ALTER USER 'root'@'localhost' IDENTIFIED BY 'Msb123.618';
Query OK, 0 rows affected (0.00 sec)
```

8、验证密码

```bash
mysql> exit
Bye

[root@localhost yum.repos.d]# mysql -u root -p
Enter password: 
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 7
Server version: 5.7.31 MySQL Community Server (GPL)

Copyright (c) 2000, 2020, Oracle and/or its affiliates. All rights reserved.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql>
```

