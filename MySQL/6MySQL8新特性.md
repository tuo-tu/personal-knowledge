## 1.7.**MySQL8新特性**

对于 MySQL 5.7 版本，其将于 2023年 10月31日 停止支持。后续官方将不再进行后续的代码维护。

MySQL 8.0 全内存访问可以轻易跑到 200W QPS，I/O 极端高负载场景跑到 16W QPS，除此之外MySQL 8还新增了很多功能，那么我们来一起看一下。

补充：**QPS**（Queries Per Second）是指**每秒查询率**，是衡量数据库或系统性能的一个重要指标，**表示系统每秒能够处理的查询请求数。**

### 1.7.1. 账户与安全

#### 1.7.1.1. 用户创建和授权

到了**MySQL8中，用户创建与授权语句必须是分开执行**，之前版本是可以一起执行。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/de3ffd4b519443439bdbd8c543c14a32.png)

**MySQL8的版本**

```sql
grant all privileges on *.* to 'lijin'@'%' identified by 'Lijin@2022'; -- 会报错
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/69ff061c543848569f28443f48fd3d2a.png)

```sql
 -- 分开执行，正确
create user 'lijin'@'%' identified by 'Lijin@2022';
grant all privileges on *.* to 'lijin'@'%';
```

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/377e9651c03148b09f1255b47c4f2e19.png" alt="image.png" style="zoom:80%;" />

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/b2a6e6d4d2f741838c8ca70c7da75693.png" alt="image.png" style="zoom:80%;" />

**MySQL5.7的版本**

```sql
grant all privileges on *.* to 'lijin'@'%' identified by 'Lijin@2022';
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/cf816534a0734330acfc5d993b55310b.png)

#### 1.7.1.2. 认证插件更新

MySQL8.0中默认的**身份认证插件**是`caching_sha2_password`，替代了之前的`mysql_native_password`。

```sql
show variables like 'default_authentication%';
```

**5.7 版本**

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/015cd62f5aa14b5887464c8562b5dd98.png)

**8 版本**

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/13929213dfb0439fa2373b969d19115d.png)

```sql
select user, host,plugin from mysql.user;
```

由于MySQL 8.0的身份认证插件更新了，带来的问题就是如果客户端没有更新，就连接不上！！

下列报错翻译：客户端不支持 server 请求的认证协议;考虑升级 MySQL 客户端。

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/662464a5e80a4fcb8ad498c5adcc7b80.png" alt="image.png" style="zoom:80%;" />

当然，也可以通过在MySQL的服务端找到my.cnf的文件，把**相关参数进行修改**（不过要MySQL重启后才能生效）

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/d8fc4c99b4d444ca8a2793f6f55cb6b8.png)

如果没办法重启服务，还有一种动态的方式：0000

```sql
alter user 'lijin'@'%' identified with mysql_native_password by 'Lijin@2022';
select host,user from mysql.user;
```

修改用户 `lijin` 的密码，并指定使用 `mysql_native_password` 认证插件。

- **`'lijin'`**：用户名。
- **`'%'`**：允许从任意主机访问。
- **`mysql_native_password`**：使用传统的用户名-密码认证插件。
- **`Lijin@2022`**：新密码。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/4a96e29ae52f46d7815fc1af6fb2f151.png)

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/1abfde3fba0a4945b526d51ba098bb26.png)

使用老的Navicat for MySQL也能访问

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/43c1bfc3251e4fa68354a57d9a15bd5e.png" alt="image.png" style="zoom:80%;" />

#### 1.7.1.3. 密码管理

MySQL 8.0可以**限制重复使用以前的密码**（修改密码时）。并且还加入了密码的修改管理功能。

```sql
show variables like 'password%';
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/06b3320142bd41459b9ace0126211d88.png)

- **修改策略（全局级）**

```sql
set persist password_history=3; -- 修改密码不能和最近3次一致
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/e85e65e80db64732b894b7a116401335.png)

- **修改策略（用户级）**

```sql
alter user 'lijin'@'%' password history 3; -- 用户级，修改密码不能和最近3次一致
```

```sql
select user,host,Password_reuse_history from mysql.user;
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/e7f298fd623a497b9b4af57d5afb0a2c.png)

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/6e3f3b85114d492faa8eadaea67038f4.png)使用重复密码修改

**用户密码(指定lijin用户)**

```sql
alter user 'lijin'@'%' identified by 'Lijin@2022';
```

- **`ALTER USER`**：用于修改现有用户的属性。
- **`'lijin'@'%'`**：指定要修改的用户。`@'%'` 表示该用户可以从**任何主机**连接到 MySQL 服务器。
- **`IDENTIFIED BY 'Lijin@2022'`**：设置新的密码为 `Lijin@2022`。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/e7d506e7e3e044238354b7a1db369012.png)

如果我们把全局的参数（password_history）改为0，则对于root用户可以反复的修改密码。

```sql
alter user 'root'@'localhost' identified by '789456';
```

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/8eac02053e2f4e0e9686d997e2431c60.png" alt="image.png" style="zoom:80%;" />

**password_reuse_interval**   则是按照天数来限定（不允许重复的）

**password_require_current**    是否需要校验旧密码（off 不校验、 on校验（修改时必须提供旧密码））(针对非root用户)

```sql
set persist password_require_current=on;
```

### 1.7.2. 索引增强

#### 1.7.2.1. 隐藏索引

**隐藏索引 (Invisible Index)** 是 MySQL 从 **8.0** 版本开始引入的一项功能，用于管理索引的可见性。隐藏索引的主要作用是将索引标记为“不可见”，从而**使优化器在查询中忽略该索引，但索引仍然存在并维护数据。**

**注意：隐藏索引不会被优化器使用，但仍然需要进行维护。**

**应用场景：软删除、灰度发布。**

- **软删除：**

  - 概念：软删除是指在数据库中，通过标记记录的状态而非实际删除数据的方式来实现“删除”。软删除的实现通常通过在表中添加一个额外的字段，例如 `is_deleted` 或 `deleted_at`，用来标记数据是否被删除。

  - 我们在线上会经常删除和创建索引，如果是以前的版本，我们如果删除了索引，后面发现删错了，我又需要创建一个索引，这样做的话就非常影响性能。**在MySQL8中我们可以将不确定是否需要删除的索引变成隐藏索引**（索引就不可用了，查询优化器也用不上），**最后确定要删除这个索引我们才会进行删除索引操作。**

- **灰度发布：**

  - 概念：灰度发布（Canary Release）是一种系统发布策略，指在新版本上线时，将其**逐步推广**到部分用户或设备，**而不是一次性面向所有用户发布**，以降低上线风险。灰度发布的核心思想是：**在真实环境中进行增量测试，逐步验证变更的安全性和效果，最终实现全量推广。**

  - 也是类似的，我们想在线上进行一些测试，可以**先创建一个隐藏索引，不会影响当前的生产环境**，然后我们通过一些附加的测试，发现这个索引没问题，那么就直接**把这个隐藏索引改成正式索引**，让线上环境生效。

    - **隐藏状态**：隐藏索引不会影响现有的查询执行计划，优化器不会选择它，这等于索引的存在对现有生产环境透明。

    - **动态启用**：通过简单的切换操作，**隐藏索引可以变为正式索引**，而无需重新创建或重启服务。

    - **风险隔离**：在测试期间，隐藏索引的任何问题不会对当前的业务流程造成影响。

**使用案例（灰度发布）：**

```sql
create table t1(i int,j int);  -- 创建一张t1表
create index i_idx on t1(i);  -- 创建一个正常索引
create index j_idx on t1(j) invisible;  -- 在表 t1 的列 j 上创建一个隐藏索引 j_idx。
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/c73a1f8ec0714350aa7ba3cca4b510ef.png)

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/b0acbc1caa464f3082cc44f44838bee7.png)

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/04da87a0e2ed4fb2bf8e2393b67a30f1.png" alt="image.png"  />

```sql
show index from t1\G  -- 查看索引信息
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/11f5ca11849c4323be34c1f3713d9a2a.png)

使用查询优化器看下：

```sql
explain select * from t1 where i=1;
explain select * from t1 where j=1;
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/168794969d3b42a2ad224643c6d1cfc7.png)

这里可以看到隐藏索引不会用上（也就是j列）。

这里可以通过优化器的开关，打开一个设置，方便我们对隐藏索引进行设置。

```sql
select @@optimizer_switch\G;   -- 查看各种参数
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/0ab72a3fea784b3397319fadde1293b0.png)

红色的部分就**是默认查询优化器对隐藏索引不可见**（默认不会考虑使用）。我们可以通过参数进行修改，确保我们可以用隐藏索引进行测试。

```sql
set session optimizer_switch='use_invisible_indexes=on';   -- 在会话级别设置查询优化器可以看到隐藏索引
```

`use_invisible_indexes` 用于控制**是否允许优化器使用隐藏索引**。

- **use_invisible_indexes=on**：允许优化器在查询计划中考虑隐藏索引，这表示：
  - 在**查询计划生成阶段**，优化器**会将隐藏索引视为普通索引一样进行评估。**
  - 如果优化器判断使用隐藏索引比其他索引更优，则它会将隐藏索引纳入查询执行计划中。

- **use_invisible_indexes=off：** 表示查询**优化器一律不会考虑使用隐藏索引**，即便隐藏索引存在且对某些查询有优化作用。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/d4b7225269b64c728ed749f7b2f5d79b.png)

设置为on后，再使用查询优化器看下：

```sql
explain select * from t1 where j=1;
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/dba6db5fe1e04e418fccd1f47718df9f.png)

把隐藏索引变成可见索引（正常索引）

```sql
alter table t1 alter index j_idx visible;   -- 变成可见
alter table t1 alter index j_idx invisible;   -- 变成不可见(隐藏索引)
```

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/e4735a743ec14fcba847a9fac18b0798.png" alt="image.png" style="zoom:80%;" />

最后一点，**不能把主键设置成隐藏索引**（MySQL做了限制）

#### 1.7.2.2. 降序索引

降序索引（Descending Index）是指按 **降序（从大到小）** 排列的索引。它是 MySQL 8.0 开始支持的功能，用于优化涉及降序排列的查询。

**背景**：在 MySQL 8.0 之前，虽然可以通过 `ORDER BY column DESC` 来排序，但即使使用了索引，MySQL 仍会进行额外的排序操作。

**MySQL 8.0 的改进**：

- MySQL 8.0开始真正支持降序索引 (Descending Index) 。
- 降序索引可以直接支持 `ORDER BY column DESC` 的查询，**无需额外的排序操作。**

**创建降序索引**：

- 通过 `CREATE INDEX` 或在 `PRIMARY KEY` 或 `UNIQUE KEY` 中使用 `DESC` 关键字。

- 示例：

  ```sql
  CREATE TABLE employees (
      id INT NOT NULL,
      name VARCHAR(50),
      salary INT,
      PRIMARY KEY (id),
      KEY idx_salary_desc (salary DESC) -- 为 salary 字段创建了一个普通索引（非唯一索引），并指定索引以降序顺序组织 salary 字段的值。
  );
  ```

**降序索引的作用**：

- 当查询使用 `ORDER BY column DESC` 时，优化器会选择降序索引以**避免额外的排序操作。**

- 对于**范围查询**（例如 `BETWEEN` 或 `<`、`>`）**也能有效地利用降序索引**。

**只有InnoDB存储引擎支持降序索引，并且只支持BTREE降序索引**。

- **BTREE** 的全称是 **Balanced Tree**（平衡树）。**降序索引的实现基于 BTREE 索引结构**，它通过调整索引节点的排序方式，将字段值按降序排列存储。因此：

  - 只有 **BTREE 类型索引** 支持降序排序。

  - **HASH 索引** 或其他类型的索引（如 FULLTEXT 索引）不支持降序索引。

另外MySQL8.0**不再对GROUP BY操作进行隐式排序。**

- **MySQL 5.7及更早版本**：

  - `GROUP BY` 会**对结果进行隐式排序**，即使用户没有明确指定 `ORDER BY`。

  - 这在某些情况下导致不必要的性能开销。

- **MySQL 8.0 的改进**：

  - `GROUP BY` 不再隐式排序。

  - 如果需要排序的结果，可以**明确使用** `ORDER BY` 语句。

  - 示例：

    ```sql
    SELECT department_id, COUNT(*)
    FROM employees
    GROUP BY department_id;
    ```

    在 MySQL 8.0 中，结果可能不是按 `department_id` 排序的。如果需要排序，需要添加`ORDER BY`：

    ```sql
    SELECT department_id, COUNT(*)
    FROM employees
    GROUP BY department_id
    ORDER BY department_id; -- 显示排序
    ```


在MySQL中创建一个t2表

```sql
create table t2(
    c1 int,
    c2 int,
    index idx1(c1 asc,c2 desc)
);
show create table t2\G
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/75d8a1f112174a0eae3bdcc329844c80.png)

如果是5.7中，则没有显示升序还是降序信息

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/44eda4306a624a77ae6f4c9bf3811143.png)

我们插入一些数据，给大家演示下降序索引的使用。

```sql
insert into t2(c1,c2) values(1,100),(2,200),(3,150),(4,50);
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/6c688fd5a39342a08158c5cdecc8592b.png)

看下索引使用情况

```sql
explain select * from t2 order by c1,c2 desc;
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/74a77cd932174740a18b3b6b28e19617.png)

我们在5.7对比一下

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/9dfcf7f4f605455b851f3ae64c5b2fa5.png)

Using filesort这里说明，这里需要一个额外的排序操作，才能把刚才的索引利用上。

我们把查询语句换一下

```sql
explain select * from t2 order by c1 desc,c2 ;
```

MySQL8中使用了

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/00d9db6d7c71430a955e7107a0678461.png)

另外还有一点，就是**group by语句在 8 之后不再默认排序。**

```sql
select count(*),c2 from t2 group by c2;
```

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/e51a002e1f034c15bad26b7781532518.png" alt="image.png" style="zoom:80%;" />

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/4bbbf1e99fa847479e232377e5d66a77.png)

在8要排序的话，就需要手动把排序语句加上

```sql
select count(*),c2 from t2 group by c2 order by c2;
```

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/23c4a24d51db40e4aaae81792d4b8a40.png" alt="image.png" style="zoom:80%;" />

#### 1.7.2.3. 函数索引

**函数索引**（Functional Index）是 MySQL 8.0 开始支持的一种索引类型。它允许你**在索引中使用表达式或函数**，而不仅仅是列值。通过函数索引，可以对计算后的表达式结果进行索引，从而优化复杂查询。

之前我们知道，如果在查询中加入了函数，索引不生效，所以MySQL8引入了函数索引。

MySQL 8.0.13开始支持**在索引中使用函数（表达式）**的值。**支持降序索引，支持JSON 数据的索引**，函数索引基于**虚拟列功能**实现。

- **使用函数索引（表达式）**

```sql
create table t3(
    c1 varchar(10),
    c2 varchar(10)
);
create index idx_c1 on t3(c1);   -- 普通索引
create index func_idx on t3((UPPER(c2)));  -- UPPER(c2) 表示索引存储的是 c2 列中每个值转换为大写后的结果。
```

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/47d2fd6adec142cf9dfde5c0cb677495.png" alt="image.png" style="zoom:80%;" />

在 **函数索引** 中，`UPPER()` 和 `LOWER()` 是用来**将某个字符串列的内容统一转化为大写或小写**形式，**从而对变换后的结果进行索引。**

```sql
show index from t3\G
```

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/d8508066bd504ba0897c5a2d25b85bf5.png" alt="image.png" style="zoom: 80%;" />

```sql
explain select * from t3 where upper(c1)='ABC' ;  
explain select * from t3 where upper(c2)='ABC' ;
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/7bfb171eb5c3483ba2d8ec79399c8261.png)

- **使用函数索引（JSON）**

```sql
# 为表中的 data 列指定了一个 函数索引，基于 JSON 数据列中的某个字段 ($.name) 创建索引。
create table t4(
    data json,
    index((CAST(data->>'$.name' as char(25)) ))
);

explain select * from t4 where CAST(data->>'$.name' as char(25)) = 'lijin ';
```

**`data->>'$.name'`**：

- 这是 MySQL 的 JSON 提取语法，用于从 `data` 列的 JSON 文档中提取 `$` 表示的根对象的 `name` 字段。
- `->>` 是 JSON 提取操作符，用于提取并返回字符串值。
- 例如，如果 `data` 的值是 `{ "name": "John", "age": 30 }`，则 `data->>'$.name'` 的结果为 `'John'`。

**`CAST(... AS CHAR(25))`**：函数表达式为 `CAST(expression AS target_data_type)`

- 将提取的 JSON 值转换为一个定长字符串，长度**最多**为 25 个字符。
- 这是为了确保提取的值可以存储在索引中，因为索引要求数据类型为字符串或数值。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/45a6485744a9438da17953c4df4d83fd.png)

- **函数索引基于 虚拟列功能 实现**

**函数索引在MySQL中相当于新增了一个列**，这个虚拟列会**根据你的函数来计算结果**，在使用函数索引的时候就会用这个**动态计算后**的列作为索引。这不会占用额外的存储空间。

### 1.7.3. 通用表表达式（CTE）

**通用表表达式（Common Table Expression, CTE）** 是一种在 SQL 查询中用于简化复杂查询的结构化方法。它是一个**临时的结果集**，可以在查询中像表一样引用。CTE 通常用于提高 SQL 查询的可读性和结构化程度，**尤其是在递归查询或多层嵌套查询中。**

MySQL8.0开始支持通用表表达式（Common Table Expression, CTE），即**WITH**子句。

```sql
WITH cte_name (column1, column2, ...) AS (
    -- CTE 查询部分
    SELECT ...
)
-- 主查询部分
SELECT ... 
FROM cte_name;
```

1. **非递归 CTE**

   非递归 CTE：**定义的结果集直接通过一个查询生成。**

   示例：假设有一个 `employees` 表，列包含 `id`、`name` 和 `manager_id`，我们需要查找所有直接由某个特定经理（如 `manager_id = 2`）管理的员工：

   ```sql
   WITH direct_reports AS (
       SELECT id, name, manager_id
       FROM employees
       WHERE manager_id = 2
   )
   SELECT * FROM direct_reports;
   ```

2. **递归 CTE**

   **递归 CTE 用于处理层级关系**，例如**组织结构、文件目录**等。它包含两个部分：

   - **锚查询**：生成初始结果集。

   - **递归查询**：基于锚查询的结果继续生成新结果。

   示例：查找员工层级关系。

   ```sql
   WITH RECURSIVE employee_hierarchy AS (
       -- 锚查询：找到顶层员工；并且这部分会生成employee_hierarchy的初始结果集
       SELECT id, name, manager_id, 1 AS level
       FROM employees
       WHERE manager_id IS NULL -- 没有上级意味着是老板
       UNION ALL
       -- 递归查询：找到每个员工的直接下属
       SELECT e.id, e.name, e.manager_id, eh.level + 1
       FROM employees e
       INNER JOIN employee_hierarchy eh -- employee_hierarchy是一个虚拟表，它是通过 通用表表达式 (CTE) 创建的临时结果集。
       ON e.manager_id = eh.id
   )
   SELECT * FROM employee_hierarchy; -- 结果的数据是union all格式
   ```

   **锚查询**：

   - 查找没有上级（`manager_id IS NULL`）的员工，通常是公司的顶层领导（例如 CEO）。

   - 为顶层员工设置层级为 `1`。

   **递归查询**：

   - 从锚查询的结果开始，不断查找每个员工的直接下属。

   - 通过 `INNER JOIN`，匹配条件为下属员工的 `manager_id` 等于当前层级的员工 `id`。

   - 每次递归，将层级（`level`）加 `1`，以此表示层次关系。

   **结果说明**

   假设表 `employees` 的数据如下：

   | id   | name    | manager_id |
   | ---- | ------- | ---------- |
   | 1    | Alice   | NULL       |
   | 2    | Bob     | 1          |
   | 3    | Charlie | 1          |
   | 4    | Dave    | 2          |
   | 5    | Eve     | 2          |

   执行上述查询后，生成的 `employee_hierarchy` 结果：

   | id   | name    | manager_id | level |
   | ---- | ------- | ---------- | ----- |
   | 1    | Alice   | NULL       | 1     |
   | 2    | Bob     | 1          | 2     |
   | 3    | Charlie | 1          | 2     |
   | 4    | Dave    | 2          | 3     |
   | 5    | Eve     | 2          | 3     |

**简单入门：**

以下SQL就是一个简单的CTE表达式，类似于递归调用，这段SQL中，首先执行select 1 然后得到查询结果后把这个值n送入 union all下面的 `select n+1 from cte where n < 10`,然后一直这样递归调用union all下面sql语句。

递归表 `cte` 的结构定义为一列 `n`，表示数字。

```sql
WITH recursive cte(n) as (
    select 1
    union ALL -- 合并查询结果，不去重
    select n+1 from cte where n<10
)
select * from cte;
```

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/7f73594320584b3ca192eec1fdee9389.png" alt="image.png" style="zoom:67%;" />

**案例介绍：**一个staff表，里面有id，有name还有一个 m_id，这个是对应的上级id。数据如下：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/faecda299a6e4117b9b60a9f41e09a9b.png" alt="image.png" style="zoom:80%;" />

如果我们想查询出每一个员工的上下级关系，可以使用以下方式递归CTE：挺经典的应用。

```sql
with recursive staff_view(id,name,m_id) as (
    select id ,name ,cast(id as char(200)) 
    from staff where m_id = 0 -- 先查老板
    union ALL 
    select s2.id ,s2.name,concat(s1.m_id,'-',s2.id)
    from staff_view as s1 
    join staff as s2
    on s1.id = s2.m_id
)
select * from staff_view order by id
```

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/ca7f4cfa8b80461c919ab34a5931decb.png" alt="image.png" style="zoom: 67%;" />

使用通用表表达式的好处就是上下级层级就算有4，5，6甚至更多层，都可以帮助我们遍历出来，而老的方式的写法SQL语句就要调整。

**总结：**

通用表表达式与派生表类似，就像语句级别的临时表或视图。**CTE可以在查询中多次引用，可以引用其他CTE，可以递归**。CTE支持`SELECT、INSERT、UPDATE、DELETE`等语句。

### 1.7.4. 窗口函数

窗口函数（Window Function）是 SQL 中的一种高级功能，用于对查询结果集的子集（窗口）进行计算。它与聚合函数类似，但不同之处在于窗口函数不会将结果进行分组，而是保留所有行，同时在指定的窗口内计算结果。

**语法结构**：窗口函数的基本语法如下：

```sql
<窗口函数>(<参数>)
OVER (
    [PARTITION BY <列名>]
    [ORDER BY <列名> ASC|DESC]
    [ROWS|RANGE BETWEEN <范围>]
)
```

**组成部分**

1. **窗口函数**：

   - 常见的窗口函数包括：
     - 聚合函数：`SUM()`, `AVG()`, `COUNT()`, `MAX()`, `MIN()`
     - 排序函数：`ROW_NUMBER()`, `RANK()`, `DENSE_RANK()`
     - 偏移函数：`LAG()`, `LEAD()`
     - 窗口统计函数：`NTILE()`, `PERCENT_RANK()`, `CUME_DIST()`

2. **OVER 子句**：

   - 定义窗口范围。

   - 关键部分包括：

     - **PARTITION BY**：**按某列分区，类似于按组计算。**

     - **ORDER BY**：指定排序顺序。

     - **ROWS|RANGE**：定义窗口的具体行范围（可选）。

       - **ROWS**：按物理行定义范围。

       - **RANGE**：按逻辑值定义范围。

**窗口函数 VS 聚合函数**

| **特性**       | **窗口函数**         | **聚合函数**     |
| -------------- | -------------------- | ---------------- |
| 是否分组保留行 | 保留所有行           | 聚合后只返回一行 |
| 计算范围       | 按窗口定义的范围计算 | 整个分组计算     |
| 使用场景       | 排名、累计和滑动计算 | 汇总统计         |

MySQL 8.0支持窗口函数(Window Function)，也称**分析函数**。窗口函数与分组聚合函数类似，但是**每一行数据都生成一个结果**。聚合窗口函数包括：SUM、AVG、COUNT、MAX、MIN等等。

案例如下：sales表结构与数据如下

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/55a0f2e91f0b4ff99d807c997899bd85.png" alt="image.png" style="zoom: 67%;" />

- 普通的分组、聚合（以国家统计）


```sql
SELECT country,sum(sum)
FROM sales
GROUP BY country
order BY country;
```

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/2f621f4802de415aa8e41152f0e96e67.png" alt="image.png" style="zoom:80%;" />

- 窗口函数（以国家汇总）：注意**每一行数据都生成一个sum结果**


```sql
select year,country,product,sum,
sum(sum) over (
    PARTITION by country -- 按照country列进行分区（类似于分组）
) as country_sum -- 给新的列取名
from sales
order by country,year,product,sum;
```

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/119d39d7e7b4425bb79873cf1e532d30.png" alt="image.png" style="zoom:67%;" />

- 窗口函数（计算平局值）


```sql
select year,country,product,sum,
sum(sum) over (PARTITION by country) as country_sum,
avg(sum) over (PARTITION by country) as country_avg
from sales
order by country,year,product,sum;
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/b1efa7b8456e43949f4b46aec9b12a70.png)

- 专用窗口函数：

  * 序号函数：ROW_NUMBER()、RANK()、DENSE_RANK()

  * 分布函数：PERCENT_RANK()、CUME_DIST()

  * 前后函数：LAG()、LEAD()

  * 头尾函数：FIRST_VALUE()、LAST_VALUE()

  * 其它函数：NTH_VALUE()、NTILE()


| **函数名称**       | **参数**                            | **功能描述**                                                 |
| ------------------ | ----------------------------------- | ------------------------------------------------------------ |
| **ROW_NUMBER()**   | 无                                  | 返回当前行在分组内的序号，不受重复值影响，序号始终连续：1, 2, 3, 4, 5。 |
| **DENSE_RANK()**   | 无                                  | 返回**不间断**的分组排名，可出现重复排名：1, 1, 2, 2。       |
| **RANK()**         | 无                                  | 返回**间断**的分组排名，重复值后排名会跳跃：1, 1, 3, 3, 5。  |
| **PERCENT_RANK()** | 无                                  | 计算分组内累计百分比：`(当前值的前面行数) / (分组总行数 - 1)`，返回范围为 [0, 1]。 |
| **CUME_DIST()**    | 无                                  | 计算累计分布值：`分组值小于等于当前值的行数（注意不是id哦） / 分组总行数`，返回范围为 [0, 1]。 |
| **LAG()**          | `lag(column_name, [N, [default]])`  | 获取当前行往前第 N 行的值，如果N缺失， 默认为 1。如前面不足 N 行，返回默认值 `default`（default默认值为 NULL）。 |
| **LEAD()**         | `lead(column_name, [N, [default]])` | 获取当前行往后第 N 行的值，与 `LAG()` 逻辑相反，其余相同。   |
| **FIRST_VALUE()**  | `first_value(column_name)`          | 返回窗口（分组）中按指定顺序排列的**第一行的值**。           |
| **LAST_VALUE()**   | `last_value(column_name)`           | 返回窗口（分组）中按指定排序规则排列的**最后一行的值**。     |
| **NTH_VALUE()**    | `nth_value(column_name, N)`         | 返回窗口中按指定顺序排列的第 N 行的值。                      |
| **NTILE()**        | `ntile(N)`                          | 将分组数据平均分成 N 个桶，返回当前行所在的桶号，范围从 1 到 N。适用于等分数据的场景。 |

- **窗口函数（排名）**

用于计算分类排名的排名窗口函数，以及获取指定位置数据的取值窗口函数

```sql
SELECT
YEAR,
country,
product,
sum,
row_number() over (ORDER BY sum) AS 'rank', -- 排名不间断不重复
rank() over (ORDER BY sum) AS 'rank_1' -- 排名间断重复
FROM sales;
```

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/5cca667ce4344608a49ea780d7a69ba3.png" alt="image.png" style="zoom:67%;" />

```sql
SELECT
YEAR,
country,
product,
sum,
sum(sum) over (PARTITION by country order by sum rows unbounded preceding) as sum_1
FROM sales order by country,sum;
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/24c5c7cd8cf44dee84e869ee5dfaadc3.png)

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/a96976f74acc4b3ab445bb22e6cdece7.png" alt="image.png" style="zoom: 67%;" />

当然可以做的操作很多，具体见官网：

[https://dev.mysql.com/doc/refman/8.0/en/window-function-descriptions.html]()

### 1.7.5. 原子DDL操作

**原子 DDL（Atomic Data Definition Language）** 是指对数据定义语句（DDL）操作的执行具有原子性，意味着这些操作要么完全成功，要么完全失败，并且**失败时会自动回滚**，不会影响数据库的完整性和一致性。

MySQL 8.0 开始支持原子 DDL 操作，并且是默认开启的，其中**与表相关的原子 DDL 只支持 InnoDB 存储引擎**。

- 一个原子 DDL 操作内容包括：更新数据字典，存储引擎层的操作，在 binlog 中记录 DDL 操作。

- 支持与表相关的 DDL：数据库、表空间、表、索引的 CREATE、ALTER、DROP 以及 TRUNCATE TABLE。

- 支持的其他 DDL ：存储程序、触发器、视图、UDF 的 CREATE、DROP 以及ALTER 语句。
- 支持账户管理相关的 DDL：用户和角色的 CREATE、ALTER、DROP 以及适用的 RENAME，以及 GRANT 和 REVOKE 语句。

举例：

```sql
drop table t1,t2;   
```

上面这个语句，如果只有t1表，没有t2表在MySQL5.7与 8 的表现是不同的：

5.7会删除t1表，而在8中因为报错了，整个是一个原子操作，所以不会删除t1表。

### 1.7.6. JSON增强

具体看官网信息，英文好的直接看，英文不好的找个翻译工具即可看懂

[MySQL :: MySQL 8.0 Reference Manual :: 11.5 The JSON Data Type](https://dev.mysql.com/doc/refman/8.0/en/json.html)

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/1d64831a19c4488089e934b96119686e.png" alt="image.png" style="zoom:80%;" />

### 1.7.7. InnoDB其他改进功能

#### 自增列持久化

- 在 **MySQL 5.7 及早期版本** 中，自增列计数器（**AUTO_INCREMENT**）的值仅**存储在内存**中，因此存在以下问题：

  1. **值丢失风险**：
      如果 MySQL 实例意外崩溃，重启后自增计数器的值无法恢复准确值，可能导致值重复或跳跃。

  2. **高并发风险**：
      在主从复制或高并发场景中，因自增值仅存在内存中，可能引发自增字段值的不一致。

- **MySQL 8.0的改进**：**写入持久化存储**，解决了长期以来的自增字段值可能重复的 bug。
  - 每次自增计数器变化时，InnoDB 会将其 **最大值** 写入 **redo log**。
  - 每次检查点（**checkpoint**）时，将自增计数器的值写入引擎私有的 **系统表**。

**Redo Log** 是 InnoDB 存储引擎用来实现 **崩溃恢复** 和 **事务持久性（Durability，ACID 的 D）** 的一种**物理日志文件**。它**记录了事务对数据库所做的修改**，以保证即使系统崩溃，也能恢复到一致的状态。

#### **死锁检查控制**

MySQL 8.0 （MySQL 5.7.15）增加了一个新的动态变量`innodb_deadlock_detect`，用于控制 InnoDB 是否执行死锁检测。对于高并发的系统，禁用死锁检查可能带来性能的提高。

- `ON`：启用死锁检测（默认值）。

- `OFF`：禁用死锁检测。

#### **锁定语句选项**

在 **MySQL 8.0** 中，`SELECT ... FOR SHARE` 和 `SELECT ... FOR UPDATE` 支持 **`NOWAIT`** 和 **`SKIP LOCKED`** 这两种处理行锁冲突的选项，提供了对锁定行为的更精细控制。

- **`NOWAIT`**：如果**目标行已被其他事务锁定**，查询会**立即失败并返回错误**，无需等待锁释放。

- **`SKIP LOCKED`**：如果目标行已被其他事务锁定，查询会**跳过这些行，并从结果集中移除。**

```sql
-- FOR UPDATE 示例
SELECT * FROM table_name WHERE condition FOR UPDATE [NOWAIT | SKIP LOCKED];

-- FOR SHARE 示例
SELECT * FROM table_name WHERE condition FOR SHARE [NOWAIT | SKIP LOCKED];
```

#### InnoDB 其他改进功能。

* 支持部分快速 DDL：`ALTER TABLE ALGORITHM=INSTANT;`

  - 部分快速 DDL 操作**不会拷贝数据文件，也不需要对表加全表锁。**

  - 例如，添加一个虚拟列或新列（不包含默认值）可以瞬间完成。

    ```sql
    ALTER TABLE employees ADD COLUMN age INT ALGORITHM=INSTANT;
    ```

* InnoDB 临时表使用共享的临时表空间文件 ibtmp1。

  - MySQL 8.0 开始，所有 InnoDB 临时表统一使用一个共享的临时表空间文件 `ibtmp1`。

  - 该文件位于 `datadir` 目录下，**生命周期与服务器启动时间一致。**

  - **临时表空间在重启后自动清空。**

* 新增静态变量 `innodb_dedicated_server`，**自动配置** InnoDB 内存参数，特别适合专用数据库服务器场景。当启用时，会**自动调整**以下参数：

  ```sql
  innodb_buffer_pool_size
  innodb_log_file_size
  innodb_flush_method
  -- 启动方式
  SET GLOBAL innodb_dedicated_server = 1;
  ```

* 默认创建 2 个 UNDO 表空间，不再使用系统表空间。

  - MySQL 8.0 默认创建两个独立的 UNDO 表空间，不再使用系统表空间存储 UNDO 日志。

* 支持 ALTER TABLESPACE ... RENAME TO 重命名通用表空间。

  - 新增对通用表空间重命名的支持。

  - 可以通过简单的 SQL 语句更改表空间名称，无需导出和重新创建。

**表空间**（Tablespace）是数据库管理系统（DBMS）中的一种逻辑存储结构，用于管理和组织物理存储资源。它为表、索引和其他数据库对象提供了**逻辑的存储位置。**

**表空间的特点：**

- **逻辑与物理分离**：
  - 表空间是逻辑存储的概念，**用户在表空间中创建和管理数据库对象**（如表和索引）。
  - **表空间映射到一个或多个物理存储文件**（如磁盘文件）。
