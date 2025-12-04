# MySQL8新特性底层原理

## 降序索引

### 什么是降序索引

MySQL 8.0开始真正支持降序索引 (descendingindex) 。只有InnoDB存储引擎支持降序索引，只支持BTREE降序索引。

```sql
CREATE INDEX idx_col_desc ON table_name (col_name DESC);
```

也就是说，如果后续的查询需要降序排序的结果（`ORDER BY col_name DESC`），MySQL 直接利用降序索引，无需额外的排序操作。

另外MySQL8.0不再对GROUP BY操作进行隐式排序，如果需要排序结果，必须显式使用 `ORDER BY`。

在MySQL中创建一个t2表

```sql
create table t2(
    c1 int,
    c2 int,
    index idx1(c1 asc,c2 desc) -- 定义联合索引，c1 为升序，c2 为降序。
);
show create table t2\G
```

> 定义联合索引举例：
>
> ```sql
> CREATE TABLE users (
>     id INT AUTO_INCREMENT PRIMARY KEY,
>     name VARCHAR(50),
>     age INT,
>     city VARCHAR(50),
>     INDEX idx_name_age_city (name, age, city) -- 定义联合索引
> );
> ```

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/75d8a1f112174a0eae3bdcc329844c80.png" alt="image.png" style="zoom:80%;" />

如果是5.7中，则没有显示升序还是降序信息

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/44eda4306a624a77ae6f4c9bf3811143.png)

我们插入一些数据，给大家演示下降序索引的使用

```sql
insert into t2(c1,c2) values(1,100),(2,200),(3,150),(4,50);
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/6c688fd5a39342a08158c5cdecc8592b.png)

看下索引使用情况

```sql
explain select * from t2 order by c1,c2 desc;
```

type = index表明需要扫描全部的索引记录，但是不用回表，也就是覆盖索引的意思。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/74a77cd932174740a18b3b6b28e19617.png)

我们在5.7对比一下（using filesort表示需要外部排序）。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/9dfcf7f4f605455b851f3ae64c5b2fa5.png)

这里说明，这里需要一个额外的排序操作，才能把刚才的索引利用上。

我们把查询语句换一下

```sql
explain select * from t2 order by c1 desc,c2;
```

MySQL8中使用了

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/00d9db6d7c71430a955e7107a0678461.png)

**`Backward Index Scan`** 表示 MySQL 正在以**逆序方式扫描索引**。

> 降序索引结合覆盖索引举例：
>
> ```sql
> CREATE TABLE t2 (
>     id INT AUTO_INCREMENT PRIMARY KEY,
>     score INT NOT NULL,
>     age INT NOT NULL,
>     INDEX idx_score_age (score DESC, age ASC) -- 联合索引：score 降序，age 升序
> );
> 
> -- 查询：
> SELECT score, age FROM t2 WHERE score > 50 ORDER BY score DESC; -- age 默认升序
> 
> ```
>
> 索引 `idx_score_age (score DESC, age ASC)` 包含查询所需的所有列，避免了回表操作。

另外还有一点，就是group by语句在 8之后不再默认排序

```sql
select count(*),c2 from t2 group by c2;
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/e51a002e1f034c15bad26b7781532518.png)

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/4bbbf1e99fa847479e232377e5d66a77.png)

在8要排序的话，就需要手动把排序语句加上

```sql
select count(*),c2 from t2 group by c2 order by c2;
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1653544165079/23c4a24d51db40e4aaae81792d4b8a40.png)

到此为止，大家应该对升序索引和降序索引有了一个大概的了解，但并没有真正理解，因为大家并不知道升序索引与降序索引底层到底是如何实现的。

### 降序索引的底层实现

升序索引对应的B+树

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1656320896043/b595577948f94b2b84de3d79557eb9fc.png)

降序索引对应的B+树

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1656320896043/d9b945594d154f63b167b3c6c150a99b.png" alt="image.png" style="zoom:80%;" />

如果没有降序索引，查询的时候要实现降序的数据展示，那么就需要把原来默认是升序排序的数据处理一遍（比如利用压栈和出栈操作），而降序索引的话就不需要，所以在优化一些SQL的时候更加高效。

还有一点，现在 **只有Innodb存储引擎支持降序索引** 。

## Doublewrite Buffer的改进

#### **MySQL5.7**

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1656320896043/d462b4c147bc41148f82bee5564c02b3.png)

#### **MySQL8.0**

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1656320896043/2b08d8cf75e64ceb909b03ce818cb287.png)

**在MySQL 8.0.20 版本之前，doublewrite 存储区位于系统表空间，从 8.0.20 版本开始，doublewrite 有自己独立的表空间文件**，这种变更，能够降低doublewrite的写入延迟，增加吞吐量，为设置doublewrite文件的存放位置提供了更高的灵活性。

因为系统表空间在存储中就是一个文件，那么doublewrite必然会受制于这个文件的读写效率（其他向这个文件的读写操作，比如统计、监控等数据操作）

**系统表空间(system tablespace)**

这个所谓的系统表空间可以对应文件系统上一个或多个实际的文件，默认情况下，InnoDB会在数据目录下创建一个名为ibdata1(在你的数据目录下找找看有木有)、大小为12M的文件，这个文件就是对应的系统表空间在文件系统上的表示。

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1654000409075/6d3767c7708848ddb0b110e20df12388.png" alt="image.png" style="zoom:80%;" />

而单独的文件必然效率比放在系统表空间效率要高！！！

**新增的参数：**

**innodb_doublewrite_dir**

指定doublewrite文件存放的目录，如果没有指定该参数，那么与innodb数据目录一致(innodb_data_home_dir)，如果这个参数也没有指定，那么默认放在数据目录下面(datadir)。

**innodb_doublewrite_files**

指定doublewrite文件数量，**默认**情况下，每个buffer pool实例，对应**2个doublewrite文件**。

**innodb_doublewrite_pages**

一次批量写入的doublewrite页数量的最大值，默认值、最小值与innodb_write_io_threads参数值相同，最大值512。

**innodb_doublewrite_batch_size**

一次批量写入的页数量。默认值为0，取值范围0到256。

## redo log 无锁优化

[MySQL :: MySQL 8.0: New Lock free, scalable WAL design](https://dev.mysql.com/blog-archive/mysql-8-0-new-lock-free-scalable-wal-design/)

## MySQL8中快速添加列的底层实现原理

MySQL 8 中快速添加列的底层实现原理是通过 InnoDB 存储引擎的 **"Fast Index Creation"** 特性实现的。**该特性允许在大型表中高效地添加列，而无需重建整个表。**

Online DDL 操作添加了 instant 算法，使得添加列时不再需要重建整个表，只需要在表的 metadata 中记录新增列的基本信息即可。

新的算法依赖于 MySQL 8.0 对表 metadata 结构做出的一些变更。8.0除了在表的 metadata 信息中新增了 **instant 列的默认值**以及非 instant 列的数量以外，还在数据的物理记录中加入了 info_bit，包括一个 flag 来标记这条记录是否为添加 instant 列之后才更新、插入的，以及 column_num，用来记录行数据总共有多少列。

当使用 instant 算法来添加列的时候，无需 rebuild 表，直接把列的信息记录到 metadata 中即可，对这些行进行操作时，可以读取 metadata 的信息来组合出完整的行数据。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1656320896043/c54c878ca10144eeb2e5bcdc829bdbad.png)


### 使用限制

1. 如果 alter 语句包含了 add column 和其他的操作，其中**有操作不支持 instant 算法的，那么 alter 语句会报错**，所有的操作都不会执行。
2. 添加列时，不能使用 after 关键字控制列的位置，**只能添加在表的末尾（最后一列）。**
3. 开启压缩的 innodb 表无法使用 instant 算法。
4. 不支持包含全文索引的表。
5. 仅支持使用 MySQL 8.0 新表空间格式的表。
6. 不支持临时表。
7. 包含 instant 列的表无法在旧版本的 MySQL 上使用（即物理备份无法恢复）。
8. 在旧版本上，如果表或者表的索引已经 corrupt（损坏），除非已经执行 fix（修复） 或者 rebuild，否则升级到新版本后无法添加 instant 列。
