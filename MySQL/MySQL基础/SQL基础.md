## SQL语言的5个部分

### 数据查询语言（Data Query Language，DQL）

DQL主要用于数据的查询，其基本结构是使用SELECT子句，FROM子句和WHERE子句的组合来查询一条或多条数据。

### 数据操作语言（Data Manipulation Language，DML）

DML主要用于对**数据库中的数据**进行增加、修改和删除的操作，其主要包括：

```sql
1) INSERT：增加数据
2) UPDATE：修改数据
3) DELETE：逐行删除数据，可回滚，而 `TRUNCATE` 不能。
```

### 数据定义语言（Data Definition Language，DDL）

DDL主要用针对是数据库对象（数据库、表、索引、视图、触发器、存储过程、函数）进行**创建、修改和删除**操作。DDL 语句包括 **CREATE、ALTER、DROP、TRUNCATE、RENAME** 等。

  1) CREATE：创建数据库对象
  2) **ALTER（修改数据库对象）：**用于修改表的结构，如添加/删除列、修改数据类型等。
  3) **DROP（删除数据库对象）**：删除数据库或数据库对象（如表、索引等），操作不可恢复。
  4) **TRUNCATE（清空表）**：用于删除表中的所有数据，但不删除表结构，**且不可回滚。**
  5) **RENAME**（重命名对象）：用于修改数据库对象的名称，例如重命名表。

### 数据控制语言（Data Control Language，DCL）

DCL用来授予或回收访问数据库的权限，其主要包括：

#### **GRANT：授予用户某种权限**

```sql
GRANT 权限列表 ON 数据库对象 TO 用户 [IDENTIFIED BY 密码] [WITH GRANT OPTION];
```

**参数说明：**

- **权限列表：** 指定要授予的权限，如 `SELECT`、`INSERT`、`UPDATE`、`DELETE` 等。
- **数据库对象：** 指定权限作用的对象，可以是表、视图、存储过程等。
- **用户：** 指定接收权限的用户。
- **WITH GRANT OPTION：** 允许用户将权限授予其他用户。

**示例：**

1. 授予用户 `user1` 对表 `employees` 的查询和插入权限：

   ```sql
   GRANT SELECT, INSERT ON employees TO 'user1';
   ```

2. 授予用户 `admin` 对数据库 `school` 的所有权限，并允许其再分配权限：

   ```sql
   GRANT ALL PRIVILEGES ON school.* TO 'admin' WITH GRANT OPTION;
   ```

#### REVOKE：回收授予的某种权限

`REVOKE` 用于撤销已授予用户的权限。

```sql
REVOKE 权限列表 ON 数据库对象 FROM 用户;
```

**参数说明：**

- **权限列表：** 指定要撤销的权限。
- **数据库对象：** 指定权限作用的对象。
- **用户：** 指定被撤销权限的用户。

**示例：**

1. 撤销用户 `user1` 对表 `employees` 的插入权限：

   ```sql
   REVOKE INSERT ON employees FROM 'user1';
   ```

2. 撤销用户 `admin` 对数据库 `school` 的所有权限：

   ```sql
   REVOKE ALL PRIVILEGES ON school.* FROM 'admin';
   ```

#### **DCL 权限列表**

以下是常见的权限类型：

- **ALL PRIVILEGES:** 授予所有权限。
- **SELECT:** 查询权限。
- **INSERT:** 插入权限。
- **UPDATE:** 更新权限。
- **DELETE:** 删除权限。
- **EXECUTE:** 执行存储过程的权限。
- **ALTER:** 修改表结构的权限。
- **CREATE:** 创建数据库对象的权限。
- **DROP:** 删除数据库对象的权限.

### 事务控制语言（Transaction Control Language，TCL）

TCL用于数据库的事务管理。其主要包括：

 1. START TRANSACTION：开启事务

 2. COMMIT：提交事务

 3. ROLLBACK：回滚事务

 4. SAVEPOINT：设置保存点

       用于在事务中创建一个保存点，允许部分回滚到特定的保存点，而不是回滚整个事务。

       ```sql
       SAVEPOINT 保存点名;
       ```

 5. RELEASE SAVEPOINT：释放保存点，用于删除指定的保存点。

       ```sql
       RELEASE SAVEPOINT 保存点名;
       ```

 6. SET TRANSACTION：设置事务的属性，**用于设置事务的隔离级别。**

       ```sql
       SET TRANSACTION ISOLATION LEVEL 隔离级别;
       ```

       隔离级别是数据库管理系统（DBMS）在处理并发事务时，用来控制不同事务之间可见性的一种机制。**它决定了一个事务对其他事务的可见程度**，从而影响数据的一致性和并发性能。

       1. **四种事务隔离级别**

       SQL 标准定义了四种隔离级别，**从低到高**依次是：

       | 隔离级别                         | 解决的问题           | 允许的问题             | 并发性 |
       | -------------------------------- | -------------------- | ---------------------- | ------ |
       | **READ UNCOMMITTED**（读未提交） | 无                   | 脏读、不可重复读、幻读 | 高     |
       | **READ COMMITTED**（读已提交）   | 解决脏读             | 不可重复读、幻读       | 中等   |
       | **REPEATABLE READ**（可重复读）  | 解决脏读、不可重复读 | 幻读                   | 较低   |
       | **SERIALIZABLE**（可串行化）     | 解决所有问题         | 无                     | 最低   |

       2. **各隔离级别的详细解释**

          - READ UNCOMMITTED（读未提交）

            事务可以读取 **其他事务未提交的数据**，可能会导致**脏读**。

          - READ COMMITTED（读已提交）

            事务只能读取**其他事务已提交的数据**，可以避免**脏读**，但可能会发生**不可重复读**。

          - REPEATABLE READ（可重复读）

            事务在执行过程中，多次读取同一条记录的结果保持一致（防止不可重复读），但可能发生**幻读**。

          - SERIALIZABLE（可串行化）

            **最高级别的隔离性**，事务**完全串行化**执行，所有事务必须按顺序执行，从而**完全避免脏读、不可重复读和幻读**。

### select语句的执行顺序

在 SQL 查询中，`SELECT` 语句的执行顺序（逻辑执行顺序）与它的书写顺序并不完全相同。SQL 查询的执行通常由数据库优化器进行解析、优化和执行，其逻辑顺序如下：

------

#### **1. SQL 语句的书写顺序**

通常，我们书写一个 SQL 语句的顺序如下：

```sql
SELECT column1, column2
FROM table_name
WHERE condition
GROUP BY column
HAVING condition
ORDER BY column
LIMIT number;
```

但 SQL 实际上是按照**逻辑执行顺序**来执行的，并非从 `SELECT` 开始。

------

#### **2. `SELECT` 语句的执行顺序**

SQL 语句的执行顺序如下（从 1 到 7 是逻辑执行顺序）：

| 执行顺序 | 关键字 (`SQL` 语句) | 作用                                          |
| -------- | ------------------- | --------------------------------------------- |
| 1️⃣        | `FROM`              | 确定数据来源的表或视图，并进行连接操作        |
| 2️⃣        | `WHERE`             | 过滤来源数据，仅保留满足条件的记录            |
| 3️⃣        | `GROUP BY`          | 对来源数据进行分组                            |
| 4️⃣        | `HAVING`            | 过滤分组后的数据                              |
| 5️⃣        | `SELECT`            | 选择要返回的列，并进行计算或去重 (`DISTINCT`) |
| 6️⃣        | `ORDER BY`          | 对结果进行排序                                |
| 7️⃣        | `LIMIT`             | 限制返回的行数                                |

------

#### **3. 详细解析 SQL 执行顺序**

##### **① `FROM`（确定数据来源，执行连接）**

数据库首先会确定查询的数据来源，这一步包括：

- 解析表的结构（字段、索引）。
- 处理表的连接（`JOIN`）。
- 生成一个**临时数据集**，供后续步骤使用。

**示例**

```sql
SELECT * FROM employees;
```

如果有多个表连接：

```sql
SELECT * FROM employees e JOIN departments d ON e.dept_id = d.dept_id;
```

------

##### **② `WHERE`（行级过滤）**

在 `FROM` 生成的临时数据集中，`WHERE` 用于**过滤数据**，只保留符合条件的记录。

**示例**

```sql
SELECT * FROM employees WHERE age > 30;
```

💡 **注意**：`WHERE` 不能用于聚合函数（如 `SUM()`、`AVG()`），因为 `WHERE` 在 `GROUP BY` 之前执行。

------

##### **③ `GROUP BY`（分组）**

如果查询包含聚合操作（如 `SUM()`、`COUNT()`），数据库会在 `WHERE` 过滤后的数据基础上**按指定列分组**。

**示例**

```sql
SELECT department, COUNT(*) 
FROM employees
WHERE age > 30
GROUP BY department;
```

💡 **注意**：`GROUP BY` 之后，`SELECT` 语句中只能包含：

- **分组列**
- **聚合函数（SUM、COUNT、AVG 等）**

错误示例：

```sql
SELECT department, age, COUNT(*) 
FROM employees
GROUP BY department; -- ❌ 错误：age 不是聚合列
```

------

##### **④ `HAVING`（对分组数据进行过滤）**

- `HAVING` 和 `WHERE` 类似，都是用来**过滤数据**，但：
  - `WHERE` 在 `GROUP BY` **之前** 执行，作用于**行**。
  - `HAVING` 在 `GROUP BY` **之后** 执行，作用于**分组后的数据**。

**示例**

```sql
SELECT department, COUNT(*) AS employee_count
FROM employees
GROUP BY department
HAVING COUNT(*) > 10;
```

💡 **注意**：`HAVING` 通常用于过滤聚合结果，如 `SUM()`、`COUNT()`。

------

##### **⑤ `SELECT`（选择数据列）**

- 经过前面 `FROM`、`WHERE`、`GROUP BY`、`HAVING` 的处理后，`SELECT` 负责：
  - 选择具体的字段或计算值
  - 进行去重操作（`DISTINCT`）

**示例**

```sql
SELECT department, COUNT(*) AS employee_count
FROM employees
GROUP BY department
HAVING COUNT(*) > 10;
```

💡 **注意**：

- `SELECT` 不能引用 `WHERE` 之后被过滤掉的数据。
- `DISTINCT` 仅影响 `SELECT` 返回的结果。

错误示例：

```sql
SELECT department, COUNT(*) 
FROM employees
WHERE COUNT(*) > 10 -- ❌ 错误，COUNT() 不能用于 WHERE
GROUP BY department;
```

------

##### **⑥ `ORDER BY`（排序）**

- 按指定列进行排序，默认升序（`ASC`），可以使用 `DESC` 降序。

**示例**

```sql
SELECT department, COUNT(*) AS employee_count
FROM employees
GROUP BY department
HAVING COUNT(*) > 10
ORDER BY employee_count DESC;
```

💡 **注意**：

- `ORDER BY` 只能使用 `SELECT` 选择的列或计算列。

------

##### **⑦ `LIMIT`（限制返回的行数）**

- `LIMIT` 用于限制查询结果的返回行数。

**示例**

```sql
SELECT department, COUNT(*) AS employee_count
FROM employees
GROUP BY department
HAVING COUNT(*) > 10
ORDER BY employee_count DESC
LIMIT 5; -- 只返回前 5 行
```

💡 **注意**：

- `LIMIT` 在 `ORDER BY` 之后执行，所以 `LIMIT` 作用于**排序后的数据**。

------

#### **4. 示例：完整 SQL 语句的执行顺序**

```sql
SELECT department, COUNT(*) AS employee_count
FROM employees
WHERE age > 30
GROUP BY department
HAVING COUNT(*) > 10
ORDER BY employee_count DESC
LIMIT 5;
```

##### **实际执行顺序**

 1️⃣ **FROM** `employees`（选择 `employees` 表）
 2️⃣ **WHERE** `age > 30`（过滤出年龄大于 30 的员工）
 3️⃣ **GROUP BY** `department`（按 `department` 分组）
 4️⃣ **HAVING** `COUNT(*) > 10`（过滤分组后 `employee_count > 10` 的部门）
 5️⃣ **SELECT** `department, COUNT(*) AS employee_count`（选择要返回的列）
 6️⃣ **ORDER BY** `employee_count DESC`（按员工数降序排序）
 7️⃣ **LIMIT** `5`（限制返回前 5 行）

------

#### **5. 总结**

- **SQL 语句的执行顺序 ≠ 书写顺序**，真正执行顺序为： 1️⃣ `FROM` → 2️⃣ `WHERE` → 3️⃣ `GROUP BY` → 4️⃣ `HAVING` → 5️⃣ `SELECT` → 6️⃣ `ORDER BY` → 7️⃣ `LIMIT`
- `WHERE` **不能** 过滤 `GROUP BY` 之后的数据，应使用 `HAVING`。
- `ORDER BY` **最后执行**，决定最终的排序结果。
- `LIMIT` **在最后执行**，用于返回指定数量的结果。

这套逻辑对 SQL **查询优化** 和 **调试** 非常重要，有助于提高查询效率和正确性！🚀