## 1.2.MySQL中的索引

InnoDB存储引擎支持以下几种常见的索引：B+树索引、全文索引、哈希索引，其中**比较关键的是B+树索引**

### 1.2.1.B+树索引

InnoDB中的索引自然也是按照B+树来组织的，前面我们说过B+树的叶子节点用来放数据的，但是放什么数据呢？**索引自然是要放的**，因为B+树的作用本来就是就是为了快速检索数据而提出的一种数据结构，不放索引放什么呢？但是数据库中的表，数据才是我们真正需要的数据，索引只是辅助数据，甚至于一个表可以没有自定义索引。InnoDB中的数据到底是如何组织的？

#### 1.2.1.1.聚集索引（聚簇索引）

InnoDB中使用了聚集索引，它与B+树结合，聚集索引**将表的主键用来构造一棵B+树**，并且将**整张表的行记录数据存放在该B+树的叶子节点中**。也就是说，**叶子节点不仅存储 主键值（Key），还存储整行数据（数据页）**，这就是所谓的**索引即数据，数据即索引**。由于聚集索引是利用表的主键构建的，所以**每张表只能拥有一个聚集索引**。

```less
    [ 10 |  30  |  50 ]
     /      |       \
    ↓       ↓        ↓  //为了好看实际上可以旋转一下
[ id=1, name=A, age=25 ]
[ id=10, name=B, age=30 ]
[ id=30, name=C, age=40 ]
[ id=50, name=D, age=35 ]
```

**聚集索引的叶子节点就是数据页**。换句话说，数据页上存放的是完整的每行记录。因此聚集索引的

一个优点就是：通过过聚集索引能直接获取完整的整行数据。

另一个优点是：对于主键的排序查找和范围查找速度非常快。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1651212459071/e280d5663e534e08a5e5afc2d5e6dbc3.png)

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1651212459071/23a70c89f27240ba82f8b288d30e364b.png)

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1651212459071/7f328f91d2b14a4eb50d19187b9518ad.png" alt="image.png" style="zoom:67%;" />

注意：如果我们没有显示定义主键，InnoDB 如何选择聚集索引呢？MySQL会做以下操作：

1. 如果表中的**存在唯一索引（`UNIQUE`）且该索引列 `NOT NULL`**， InnoDB **会选择这个唯一索引作为聚集索引**，
2. 如果没有唯一索引，InnoDB 会创建一个隐藏列 `RowID`来做主键，然后用这个主键来建立聚集索引。

> 聚集索引的特点：
>
> 1. 每张表只能有一个聚集索引，因为数据行只能按照一种顺序存储。
> 2. 在 InnoDB 存储引擎中，主键索引是默认的聚集索引。
> 3. 如果表没有主键，InnoDB 会选择第一个唯一的非空索引作为聚集索引；如果没有唯一索引，则 MySQL 会隐式生成一个伪列（RowID）作为聚集索引。
>
> 创建聚集索引：只要有上面对应的键，就会自动创建聚集索引。

创建索引的方式有3种：①新建表中添加索引；②在已建表中添加索引；③ 以修改表的方式添加索引。

链接[mysql 中添加索引的三种方法 - MaxBruce - 博客园](https://www.cnblogs.com/bruce1992/p/13958166.html)

#### 1.2.1.2.辅助索引（二级索引，非聚集索引）

**聚集索引只能在搜索条件是主键值时才能发挥作用**，因为B+树中的数据都是按照主键进行排序的。

如果我们**想以别的列作为搜索条件**怎么办？我们一般会建立多个索引，**这些索引（注意可以不止一个）被称为辅助索引**（或二级索引）。

（每建立一个索引，就会生成一棵独立的 B+ 树）

对于辅助索引，其叶子节点存储的**并不是整行数据**。叶子节点除了包含**键值**以外，每个叶子节点中的索引行中还包含了一个书签(bookmark)。该书签用来告诉InnoDB存储引擎哪里可以找到与索引相对应的行数据。因此InnoDB存储引擎的辅助索引的书签就是相应行数据的聚集索引的键。

> chatgpt：在 InnoDB 存储引擎中，**二级索引（Secondary Index）** 的叶子节点存储的**并不是整行数据**，而是**索引列的值 + 聚集索引（Primary Key）的值**。这个 **Primary Key 值就相当于书签（Bookmark）**，用于找到完整的行数据。
>
> **二级索引**是非聚簇索引，数据按主键以外的列排序存储。二级索引是一个笼统的概念，包括唯一索引、普通索引等。

辅助索引可以加速某些列的查询，但不会改变表的数据存储结构。

举例：

```sql
CREATE TABLE users (
    id INT PRIMARY KEY,
    name VARCHAR(50),
    age INT,
    email VARCHAR(100),
    INDEX idx_name (name)
) ENGINE=InnoDB;
```

**聚集索引**（`PRIMARY KEY(id)`）的叶子节点存储完整的行数据。此处也叫主键索引

**辅助索引**（`INDEX idx_name(name)`）的叶子节点存储：**索引列（name）的值 + 聚集索引（Primary Key，即id）的值**

```pgsql
| name  | id  |
|-------|-----|
| Alice | 1   |
| Bob   | 3   |
| Carl  | 2   |
```

当查询 `SELECT * FROM users WHERE name = 'Alice'` 时：

1. 先通过 `idx_name` 找到 `name = 'Alice'` 对应的 `id = 1`。
2. 再通过 `id = 1` 到**聚集索引**中找到完整的行数据。

这就是**二级索引回表（回主键索引）查询**的过程。

**额外优化：覆盖索引（Covering Index）**

如果查询的字段**全部**都在辅助索引的列中，比如：

```sql
SELECT name FROM users WHERE name = 'Alice';
```

此时 InnoDB **不需要回表**，因为 `name` 已经在 `idx_name` 索引中，直接返回结果。

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1651212459071/70e7cdc503e14bc995ac06887c583128.png" alt="image.png"  />

比如辅助索引index(node)，那么叶子节点中包含的数据就包括了(**note和主键**)。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1651212459071/eb80ca41a76c44bdba04fed1f6ec42fc.png)

#### 1.2.1.3.回表

> 定义：通过辅助索引找到聚集索引，再通过聚集索引找到完整的行记录
>

辅助索引的存在并不影响数据在聚集索引中的组织，因此**每张表上可以有多个辅助索引**。当通过辅助索引来寻找数据时，InnoDB存储引擎会遍历辅助索引并通过叶级别的指针获得指向**主键索引的主键**，然后再通过主键索引（聚集索引）来找到一个完整的行记录。这个过程也被称为**回表** 。也就是根据辅助索引的值查询一条完整的用户记录需要使用到2棵B+树：一次辅助索引，一次聚集索引。

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1651212459071/fc1a9aa3e1444b6e82dfec49cec398f4.png" alt="image.png" style="zoom:80%;" />

**为什么我们还需要一次回表操作呢**？直接把完整的用户记录放到辅助索引d的叶子节点不就好了么？如果把完整的用户记录放到叶子节点是可以不用回表，但是太占地方了，相当于每建立一棵B+树都需要把所有的用户记录再都拷贝一遍，这就有点太浪费存储空间了。而且每次对数据的变化要在所有包含数据的索引中全部都修改一次，性能也非常低下。

很明显，**回表的记录越少，性能提升就越高**，需要回表的记录越多，使用二级索引的性能就越低，甚至让某些查询宁愿使用全表扫描也不使用二级索引。

那什么时候采用全表扫描的方式，什么时候使用采用二级索引 + 回表的方式去执行查询呢？这个就是**查询优化器**做的工作，查询优化器会事先对表中的记录计算一些统计数据，然后再利用这些统计数据根据查询的条件来计算一下需要回表的记录数，需要回表的记录数越多，就越倾向于使用全表扫描，反之倾向于使用二级索引 + 回表的方式。

#### 1.2.1.4.联合索引（又称复合索引）

前面我们对索引的描述，隐含了一个条件，那就是构建索引的字段只有一个，但**实践工作中构建索引的完全可以是多个字段**。所以，将表上的多个列组合起来进行索引我们称之为联合索引或者复合索引，比如index(a,b)就是将a,b两个列组合起来构成一个索引。

语法：

```sql
CREATE INDEX index_name ON table_name (column1, column2, column3); -- 推荐
-- 或者
ALTER TABLE table_name ADD INDEX index_name (column1, column2, column3);
```

千万要注意一点，**建立联合索引只会建立1棵B+树**，多个列分别建立索引会分别以每个列则建立B+树，有几个列就有几个B+树，比如下列`index(note)`、`index(b)`，就分别对note,b两个列各构建了一个索引。

而如果是index(note,b)在索引构建上，包含了两个意思：

1. **先把各个记录按照note列进行排序。**

2. **在记录的note列相同的情况下，采用b列进行排序**

从原理可知，为什么有最佳左前缀法则，就是这个道理

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1651212459071/6f15eac444b045179d65e601b588a03f.png" alt="image.png" style="zoom:80%;" />

✅ **联合索引的关键点**

1. **遵循最左前缀匹配原则**：**查询必须按照索引定义的最左列开始匹配**，否则索引不会被完全利用。

   **最左前缀**指的是联合索引中的**前 N 列**（从左往右的连续部分）。假设有如下的一个联合索引

   ```sql
   CREATE INDEX idx_user_status_created ON orders (user_id, status, created_at);
   ```

   则该索引的最左前缀包括：

   1. `(user_id)`
   2. `(user_id, status)`
   3. `(user_id, status, created_at)`

2. **联合索引比多个单列索引更高效**，但要考虑查询模式。

3. **索引顺序很重要**，**高选择性列放在前面**，排序列适当放入索引。

#### 1.2.1.5.覆盖索引（是一个过程）

既然多个列可以组合起来构建为联合索引，那么**辅助索引自然也可以由多个列组成。**

InnoDB存储引擎支持覆盖索引(covering index，或称索引覆盖)，是指**查询所需的所有列都能从辅助索引的 B+ 树叶子节点直接获取**，即从辅助索引中就可以获取查询的记录，而不需要回表（查询聚集索引中的记录）。使用覆盖索引的一个好处是辅助索引不包含整行记录的所有信息，故其大小要远小于聚集索引，因此可以减少大量的IO操作。所以记住，**覆盖索引并不是索引类型的一种，而是一个结果。**

**注意：查询的 SELECT 列、WHERE 条件列、ORDER BY 列等必须全部包含在同一个索引中。**

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1651212459071/37985f1a9c0a4c2c8ba506d89bf14803.png" alt="image.png" style="zoom: 60%;" />

### 1.2.2.哈希索引

InnoDB存储引擎除了我们前面所说的各种索引，还有一种**自适应哈希索引**，我们知道B+树的查找次数,取决于B+树的高度,在生产环境中,B+树的高度一般为3、4层,故需要3、4次的IO查询。

所以在InnoDB存储引擎内部自己去监控索引表，如果监控到**某个索引经常用**，那么就认为是**热数据**，然后内部自己创建一个hash索引，称之为**自适应哈希索引**( Adaptive Hash Index,AHI)，并存储指向数据行的指针。创建以后，如果下次又查询到这个索引，那么直接通过hash算法推导出记录的地址，直接一次就能查到数据，比重复去B+tree索引中查询三四次节点的效率高了不少。

InnoDB存储引擎使用的哈希函数采用除法散列方式，其冲突机制采用链表方式。

注意，对于**自适应哈希索引仅是数据库自身创建并使用的，我们并不能对其进行干预。**

显示 innodb 引擎状态：

```matlab
插入缓冲区和自适应哈希索引
----------------------------
lbuf：大小 1，可用列表 len 0，分段大小 2,0 合并
合并操作：
	插入0,删除标记0,删除0
丢弃操作：
	插入0,删除标记0,删除0
哈希表大小 2267，节点堆有 0 个缓冲区
哈希表大小 2267，节点堆有 0 个缓冲区
哈希表大小 2267，节点堆有 0 个缓冲区
哈希表大小 2267，节点堆有 0 个缓冲区
哈希表大小 2267，节点堆有 0 个缓冲区
哈希表大小 2267，节点堆有 0 个缓冲区
哈希表大小 2267，节点堆有 0 个缓冲区
哈希表大小 2267，节点堆有 0 个缓冲区
0.00 次哈希搜索/秒，0.00 次非哈希搜索/秒
```

```sql
show engine innodb status;

-------------------------------------
INSERT BUFFER AND ADAPTIVE HASH INDEX
-------------------------------------
Ibuf: size 1, free list len 0, seg size 2, 0 merges
merged operations:
 insert 0, delete mark 0, delete 0
discarded operations:
 insert 0, delete mark 0, delete 0
Hash table size 34679, node heap has 1 buffer(s)
Hash table size 34679, node heap has 0 buffer(s)
Hash table size 34679, node heap has 0 buffer(s)
Hash table size 34679, node heap has 0 buffer(s)
Hash table size 34679, node heap has 1 buffer(s)
Hash table size 34679, node heap has 1 buffer(s)
Hash table size 34679, node heap has 2 buffer(s)
Hash table size 34679, node heap has 5 buffer(s)
0.00 hash searches/s, 0.00 non-hash searches/s
```

**哈希索引只能用来搜索等值的查询**，如 `SELECT* FROM table WHERE index co=xxx`。而对于其他查找类型，如范围查找，是不能使用哈希索引的。

因此这里出现了non-hash searches/s（**非哈希索引搜索**）的情况。通过 hash searches（**哈希索引查询**）和non- hash searches可以大概了解使用哈希索引后的效率。

自适应哈希索引的启用：**innodb_adaptive_hash_index**来考虑是禁用或启动此特性，**默认AHI为开启状态。**

```sql
SET GLOBAL innodb_adaptive_hash_index = ON; -- 启用AHI（默认）
SET GLOBAL innodb_adaptive_hash_index = OFF; -- 禁用AHI
```

 哈希表的存储结构及插入和查询流程：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1651212459071/6825291f38f845c88cb7903d690be711.png" alt="image.png" style="zoom:80%;" />

chatgpt的回答：未知真假

1. **存储流程（插入操作）**
   1. `put(key, value)`：插入键值对。
   2. 根据 `key` 计算 `Hash` 值。
   3. 根据 `Hash` 计算数组下标。
   4. 根据数组下标，将数据插入到相应位置。

2. **查询流程（获取操作）**
   1. `get(key)`：获取键对应的值。
   2. 根据 `key` 计算 `Hash` 值。
   3. 根据 `Hash` 计算数组下标。
   4. 根据数组下标查找对应元素。

3. **存储结构**

   - **数组（蓝色）**：哈希表的基础存储结构，每个元素可以是一个 `Node`（节点）。

   - **链表（蓝色）**：当**哈希冲突发生时**，采用链地址法（Linked List）存储多个节点。

   - **红黑树（黑色和红色）**：当链表长度超过一定阈值（如 Java HashMap 设定为 8），链表会转换为红黑树以提高查询效率。

### 1.2.3.全文索引

什么是**全文检索**（Full-Text Search, FTS）？它**是将存储于数据库中的整本书或整篇文章中的任意内容信息查找出来的技术**。它可以根据需要获得全文中有关章、节、段、句、词等信息，也可以进行各种**统计和分析**。我们比较熟知的Elasticsearch、Solr等就是全文检索引擎，底层都是基于Apache Lucene的。

> chatgpt：全文索引允许用户对大量文本数据进行高效的查询和检索，是 MySQL 用于全文搜索的一种特殊索引类型，主要适用于查找文本字段中的关键词。它通过**倒排索引**（Inverted Index）的方式实现。
>

举个例子，现在我们要保存唐宋诗词，数据库中我们们会怎么设计？诗词表我们可能的设计如下：

| 朝代 | 作者   | 诗词年代 | 标题   | 诗词全文                                                                         |
| ---- | ------ | -------- | ------ | -------------------------------------------------------------------------------- |
| 唐   | 李白   |  | 静夜思 | 床前明月光，疑是地上霜。 举头望明月，低头思故乡。                                |
| 宋   | 李清照 |  | 如梦令 | 常记溪亭日暮，沉醉不知归路，兴尽晚回舟，误入藕花深处。争渡，争渡，惊起一滩鸥鹭。 |
| ….  | ….    | …       | ….    | …….                                                                            |

要根据朝代或者作者寻找诗，都很简单，比如`select 诗词全文 from 诗词表 where作者 = ‘李白’`，但是如果数据很多，查询速度很慢，怎么办？我们可以在**对应的查询字段上建立索引加速查询。**

但是如果我们现在有个需求：要求找到包含“望”字的诗词怎么办？用`select 诗词全文 from 诗词表 where 诗词全文 like‘%望%’`，这个意味着要扫描库中的诗词全文字段，逐条比对，找出所有包含关键词“望”字的记录。基本上，数据库中一般的SQL优化手段都是用不上。数量少，大概性能还能接受，如果数据量稍微大点，就完全无法接受了，更何况在互联网这种海量数据的情况下呢？怎么解决这个问题呢，用**倒排索引**。

倒排索引就是，**将文档中包含的关键字全部提取处理，然后再将关键字和文档之间的对应关系保存起来，最后再对关键字本身做索引排序。**用户在检索某一个关键字时，先对关键字的索引进行查找，再通过关键字与文档的对应关系找到所在文档。

于是我们可以这么保存，“望”字在蜀道难、静夜思、春台望、鹤冲天几首诗都有

| 序号 | 关键字 | 蜀道难 | 静夜思 | 春台望 | 鹤冲天 |
| ---- | ------ | ------ | ------ | ------ | ------ |
| 1    | **望** | 有     | 有     | 有     | 有     |

如果查哪个诗词中包含上，怎么办，上述的表格可以继续填入新的记录

| 序号 | 关键字 | 蜀道难 | 静夜思 | 春台望 | 鹤冲天 |
| ---- | ------ | ------ | ------ | ------ | ------ |
| 1    | 望     | 有     | 有     | 有     | 有     |
| 2    | 上     | 有     |        |        | 有     |

从InnoDB 1.2.x版本开始，InnoDB存储引擎开始支持全文检索，对应的MySQL版本是5.6.x系列。不过MySQL从设计之初就是关系型数据库，存储引擎虽然支持全文检索，**整体架构上对全文检索支持并不好而且限制很多**，比如每张表只能有一个全文检索的索引，不支持没有单词界定符( delimiter）的语言，如中文、日语、韩语等。

所以**MySQL中的全文索引功能比较弱鸡，了解即可。**



**chatgpt**：

倒排索引是一种高效的文本搜索数据结构，它能够快速定位包含特定关键词的文档，而不需要遍历整个数据库或文件系统。它是全文检索（Full-Text Search, FTS）的核心技术。

**倒排索引的核心思想是：**

- **传统的正向索引** 是按**文档 ID** 存储对应的内容。
- **倒排索引** 则是按**关键词** 存储包含该词的**文档 ID 列表**。

**对比正向索引与倒排索引：**

| **文档 ID** | **内容**              |
| ----------- | --------------------- |
| 1           | Golang 并发编程很高效 |
| 2           | Python 也支持并发     |
| 3           | Golang 适合高并发场景 |

- **正向索引（Forward Index）**

  ```css
  1 → ["Golang", "并发", "编程", "高效"]
  2 → ["Python", "支持", "并发"]
  3 → ["Golang", "高并发", "适合"]
  ```

- **倒排索引（Inverted Index）**

  ```css
  "Golang" → [1, 3]
  "并发"   → [1, 2, 3]
  "高效"   → [1]
  "Python" → [2]
  ```

### 1.2.4.索引在查询中的使用

索引在查询中的作用到底是什么？在我们的查询中发挥着什么样的作用呢？请记住：

1. **一个索引就是一个B+树**，索引让我们的查询可以快速定位和扫描到我们需要的数据记录上，加快查询的速度 。

2. **一个select查询语句在执行过程中最多使用一个二级索引**，即使在where条件中用了多个二级索引。

   解释：MySQL优化器会评估多个二级索引，**选择最优的一个索引**。

### 1.2.4.高性能的索引创建策略

正确地创建和使用索引是实现高性能查询的基础。前面我们已经了解了索引相关的数据结构，各种类型的索引及其对应的优缺点。现在我们一起来看看如何真正地发挥这些索引的优势。

#### 1.2.4.1.索引列的类型尽量小

我们在定义表结构的时候要显式的指定列的类型，以整数类型为例，有TINYINT（1字节）、MEDIUMINT（3字节）、INT（4字节）、BIGTNT（8字节）这么几种，它们占用的存储空间依次递增，**我们这里所说的类型大小指的就是该类型表示的数据范围的大小。**能表示的整数范围当然也是依次递增，如果我们想要对某个整数列建立索引的话，在表示的整数范围允许的情况下，尽量让索引列使用较小的类型，比如我们**能使用INT就不要使用BIGINT**，能使用MEDIUMINT就不要使用INT，这是因为**数据类型越小，在查询时进行的比较操作越快**（CPU层次)数据类型越小，索引占用的存储空间就越少，在一个数据页内就可以放下更多的记录，从而减少磁盘IO带来的性能损耗，也就意味着可以把更多的数据页缓存在内存中，从而加快读写效率。

**这个建议对于表的主键来说更加适用**，因为不仅是聚簇索引中会存储主键值，其他所有的二级索引的节点处都会存储一份记录的主键值，如果主键适用更小的数据类型，也就意味着节省更多的存储空间和更高效的I/0。

#### 1.2.4.2.索引的选择性

创建索引应该选择**选择性（离散性）高的列。**索引的**选择性**（离散性）是指，**不重复的索引值（也称为基数，cardinality）和数据表的记录总数N的比值**，选择性的范围是 `(0,1]`
$$
索引选择性 = \frac{索引中不同值的个数（即基数 Cardinality）}{表中的总记录数}
$$
假设某`users`表有1000000条数据，并在`age` 字段上建立了索引，`age` 取值范围为 `18-60`，共有 `43` 种不同的年龄值，因此：
$$
\text{age索引的选择性} = \frac{\text{不同值个数}}{\text{总记录数}} = \frac{43}{1000000} = 0.000043
$$


**索引的选择性越高则查询效率越高**，因为选择性高的索引可以让MySQL在查找时过滤掉更多的行。**唯一索引的选择性是1，这是最好的索引选择性，性能也是最好的。**

**很差的索引选择性就是列中的数据重复度很高，比如性别字段**，不考虑政治正确的情况下，只有两者可能，男或女。那么我们在查询时，即使使用这个索引，从概率的角度来说，依然可能查出一半的数据出来。

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1651212459071/f3908cc7cedc4344aa928936bd86c70f.png" alt="image.png" style="zoom: 67%;" />

哪列作为索引字段最好？当然是姓名字段，因为里面的数据没有任何重复，性别字段是最不适合做索引的，因为数据的重复度非常高。

怎么算索引的选择性/离散性？比如`person`这个表：

```sql
SELECT count(DISTINCT name)/count(*) FROM person;
SELECT count(DISTINCT sex)/count(*) FROM person;
SELECT count(DISTINCT age)/count(*) FROM person;
SELECT count(DISTINCT area)/count(*) FROM person;
```

#### 1.2.4.3.前缀索引

**针对`blob`、`text`、`varchar`等字符串类型的列**，MySQL不支持索引他们的全部长度，只针对字段的前 N 个字符，建立前缀索引。

语法：

```sql
ALTER TABLE tableName ADD KEY (columnName(length));
-- 或者
ALTER TABLE tableName ADD INDEX indexName (columnName(length));
```

示例：为 `email` 添加前缀索引（只索引前 10 个字符）

```sql
ALTER TABLE users ADD INDEX idx_email (email(10));
```

案例：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1651212459071/bf3a2bf88d56492a958086486e7d57c7.png" alt="image.png" style="zoom:67%;" />

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1651212459071/1a662f14ed1948f3a2b5dbdb46690091.png" alt="image.png" style="zoom:67%;" />

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1651212459071/5d5c5368694b409f99f0d1a8a238f29a.png" alt="image.png" style="zoom:67%;" />

首先找到最常见的值的列表：通过计算选择性来判定选择前缀。

```sql
SELECT COUNT(DISTINCT LEFT(order_note,3))/COUNT(*) AS sel3,
COUNT(DISTINCT LEFT(order_note,4))/COUNT(*)AS sel4,
COUNT(DISTINCT LEFT(order_note,5))/COUNT(*) AS sel5,
COUNT(DISTINCT LEFT(order_note, 6))/COUNT(*) As sel6,
COUNT(DISTINCT LEFT(order_note, 7))/COUNT(*) As sel7,
COUNT(DISTINCT LEFT(order_note, 8))/COUNT(*) As sel8,
COUNT(DISTINCT LEFT(order_note, 9))/COUNT(*) As sel9,
COUNT(DISTINCT LEFT(order_note, 10))/COUNT(*) As sel10,
COUNT(DISTINCT LEFT(order_note, 11))/COUNT(*) As sel11,
COUNT(DISTINCT LEFT(order_note, 12))/COUNT(*) As sel12,
COUNT(DISTINCT LEFT(order_note, 13))/COUNT(*) As sel13,
COUNT(DISTINCT LEFT(order_note, 14))/COUNT(*) As sel14,
COUNT(DISTINCT LEFT(order_note, 15))/COUNT(*) As sel15,
COUNT(DISTINCT order_note)/COUNT(*) As total
FROM order_exp;
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1651212459071/ffbfccc4e0fa4ff39903505adb5f1a4e.png)

可以看见，从第10个开始**选择性**的增加值很高，随着前缀字符的越来越多，选择度也在不断上升，但是增长到第15时，已经和第14没太大差别了，选择性提升的幅度已经很小了，都非常接近整个列的选择性了。

**那么针对这个字段做前缀索引的话，从第13到第15都是不错的选择**

在上面的示例中，已经找到了合适的前缀长度，如何创建前缀索引:

```sql
ALTER TABLE order_exp ADD KEY (order_note(14));
```

建立前缀索引后查询语句并不需要更改：

```sql
select * from order_exp where order_note = 'xxxx' ;
```

**缺点：**前缀索引是一种能使索引更小、更快的有效办法，但另一方面也有其缺点，MySQL无法使用前缀索引做`ORDER BY`、`GROUP BY`，也无法使用前缀索引做覆盖扫描（覆盖索引）、后缀匹配。

有时候**后缀索引 (suffix index)**也有用途（例如，找到某个域名的所有电子邮件地址）。**MySQL原生不支持直接创建后缀索引**，但可以通过**反转字符串**的方式，**存储反转后的值**，并基于此创建前缀索引，从而间接实现后缀索引。可以通过触发器或者应用程序自行处理来维护索引。

#### 1.2.4.4.只为用于搜索、排序或分组的列创建索引

也就是说，**应该只为那些用于搜索 (`WHERE`)、排序 (`ORDER BY`)、分组 (`GROUP BY`) 的列创建索引**，而出现在查询列表中的列一般就没必要建立索引了，**除非是需要使用覆盖索引**；这句话什么意思呢？比如：

**搜索**

```sql
select order_note from .... where ....
```

只为 where 条件中的列建立索引即可

**排序（ORDER BY）**

```sql
SELECT * FROM order_exp ORDER BY insert_time, order_status,expire_time;
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1651212459071/30ba8ead865741a1bad33637c189fefa.png)

查询的结果集依次按照`insert_time`、`order_status`、`expire_time`来排序。回顾一下**联合索引**的存储结构，`u_idx_day_status`（联合索引）本身就是按照上述规则排好序的，所以**直接从该联合索引生成的B+树的叶子节点中提取包含的列**，不包含的列再进行回表操作获取就好了。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1651212459071/a97a0c6253ef4030bf18088eee4cc974.png)

#### 1.2.4.5.多列索引（联合索引？）

很多人对多列索引的理解都不够。一个常见的错误就是，为每个列创建独立的索引，或者按照错误的顺序创建多列索引。

我们遇到的最容易引起困惑的问题就是索引列的顺序。**正确的顺序依赖于使用该索引的查询**，并且同时需要考虑如何更好地满足排序和分组的需要。反复强调过，在一个多列B-Tree索引中，**索引列的顺序意味着索引首先按照最左列进行排序**，其次是第二列，等等。所以，索引可以按照升序或者降序进行扫描，以满足精确符合列顺序的ORDER BY、GROUP BY和DISTINCT等子句的查询需求。

所以**多列索引的列顺序至关重要**。对于如何选择索引的列顺序有一个经验法则：**将选择性最高的列放到索引最前列**。当不需要考虑排序和分组时，将选择性最高的列放在前面通常是很好的。**这时候索引的作用只是用于优化WHERE条件的查找。**在这种情况下，这样设计的索引确实能够最快地过滤出需要的行，对于在WHERE子句中只使用了索引部分前缀列的查询来说选择性也更高。

然而，性能不只是依赖于索引列的选择性，也和查询条件的有关。可能需要根据那些运行频率最高的查询来调整索引列的顺序，比如排序和分组，让这种情况下索引的选择性最高。

**（涉及排序的查询）**：

```sql
SELECT * FROM orders WHERE user_id = 123 ORDER BY created_at;
```

**索引建议**：`INDEX(user_id, created_at)`，可以支持 `WHERE user_id = ?` 和 `ORDER BY created_at`。

同时，在优化性能的时候，可能需要使用**相同的列但顺序不同**的索引来满足不同类型的查询需求。

#### 1.2.4.6.三星索引

**三星索引概念**

三星索引，顾名思义，是满足了三个星级的索引。对于一个查询而言，一个三星索引，可能是其最好的索引。

那么，这个三个星级是如何给定的呢？满足的条件如下：

1. **相关性星**：索引**将相关的记录放到一起**则获得一星（比重27%）
2. **排序星**：如果**索引中的数据顺序和查找中的排列顺序一致**则获得二星（比重27%）
3. **宽索引星（覆盖索引）**：如果**索引中的列包含了查询中需要的全部列**则获得三星（比重50%）

**第一星：索引行相邻（减少I/O次数）**

一星的意思就是：如果一个查询相关的索引行是相邻的，或者至少相距足够靠近的话，必须扫描的索引片宽度就会缩至最短，也就是说，让索引片尽量变窄，也就是我们所说的**索引的扫描范围越小越好**。

✅ 实现：把`WHERE`条件后的**等值匹配且高选择性的列**放在索引的最左侧（作为开头）。实际使用中，范围条件也可以。

**第二星：排序顺序匹配（避免排序操作）**

**在满足一星的情况下（必须）**，当查询需要排序，group by、 order by，如果查询所需的顺序与索引是一致的（**索引本身是有序的**），是不是就可以不用再另外排序了，一般来说**排序可是影响性能的关键因素**。

✅ 实现：将 ORDER BY 列加入到索引中，保持列的顺序。

**第三星：覆盖索引（避免回表查询）**

在满足了二星的情况下，如果**索引中所包含了这个查询所需的所有列**（包括 where 子句和 select 子句中所需的列，也就是覆盖索引），这样一来，查询就不再需要回表了，减少了查询的步骤和IO请求次数，性能几乎可以提升一倍。

✅ 实现：将查询语句中剩余的列都加入到索引中。

这三颗星，哪颗最重要？**第三颗星**。因为将一个列排除在索引之外可能会导致很多磁盘随机读（回表操作）。第一和第二颗星重要性差不多，可以理解为第三颗星比重是50%，第一颗星为27%，第二颗星为23%，所以在大部分的情况下，会**先考虑第一颗星**，但会根据业务情况调整这两颗星的优先度。

**chatgpt举例：**

**1️⃣ 业务场景**

假设有一个 **订单表（orders）**，存储电商交易信息：

```sql
CREATE TABLE orders (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    order_status ENUM('pending', 'shipped', 'delivered', 'canceled') NOT NULL,
    order_date DATETIME NOT NULL,
    amount DECIMAL(10,2) NOT NULL
);
```

**2️⃣ 查询需求**

```sql
SELECT id, order_date, amount
FROM orders
WHERE user_id = 1001 AND order_status = 'shipped'
ORDER BY order_date DESC
LIMIT 10;
```

3️⃣ 创建三星索引（联合索引）

```sql
CREATE INDEX idx_user_status_date ON orders(user_id, order_status, order_date, id, amount);
```

**4️⃣ 解析为什么它是三星索引**

✅ **⭐ 相关性（27%）**：where中

- `user_id` 是过滤条件的第一列，可以快速查找符合 `user_id = 1001` 的记录。
- `order_status` 进一步缩小范围（比单列索引效率更高）。

✅ **⭐ 排序优化（23%）**：

- `ORDER BY order_date DESC` **可以直接利用索引的有序性**，避免 `filesort` 排序，提高查询效率。

✅ **⭐ 覆盖索引（50%）**：

- `id, order_date, amount` **全部包含在索引中**，查询时 **无需回表** 直接从索引中获取数据。

#### 1.2.4.6.设计三星索引实战

**现在有表，SQL如下**

```sql
CREATE TABLE customer (
	cno INT,
	lname VARCHAR (10),
	fname VARCHAR (10),
	sex INT,
	weight INT,
	city VARCHAR (10)
);
CREATE INDEX idx_cust ON customer (lname, city, fname, cno);
```

对于下面的SQL而言，这是个三星索引

```sql
select cno,fname 
from customer 
where lname=’xx’ and city =’yy’ 
order by fname;
```

来评估下：

**相关性星**：所有等值谓词的列（where条件），是组合索引的开头的列，可以把**索引片缩得很窄**，符合。根据之前讲过的联合索引，我们是知道条件已经把搜索范围搜到很窄了

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1651212459071/f0f05f49cd024c538bf8afa800a71e9b.png" alt="image.png" style="zoom: 70%;" />

**排序星**：order by的fname字段在组合索引中且是索引自动排序好的，符合。

**索引覆盖星**：select中的cno字段、fname字段在组合索引中存在，符合。

**实战：现在有表test，SQL如下：**

```sql
CREATE TABLE `test` (
	`id` INT (11) NOT NULL AUTO_INCREMENT,
	`user_name` VARCHAR (100) DEFAULT NULL,
	`sex` INT (11) DEFAULT NULL,
	`age` INT (11) DEFAULT NULL,
	`c_date` datetime DEFAULT NULL,
	PRIMARY KEY (`id`),

) ENGINE = INNODB AUTO_INCREMENT = 12 DEFAULT CHARSET = utf8;
```

SQL语句如下：

```sql
select user_name, sex, age 
from test 
where user_name like 'test%' and sex =1 
ORDER BY age
```

如果我们建立联合索引`(user_name,sex,age)`：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1651212459071/35954180a81049c68d1e45275b71b19e.png" alt="image.png" style="zoom:67%;" />

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1651212459071/02edf4066ee5456fbd58b658af6f5904.png" alt="image.png" style="zoom:67%;" />

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1651212459071/a6fe8c913c6d4f1ba3b41f7922ecaebb.png" alt="image.png" style="zoom:67%;" />



第一颗星，满足

第二颗星，不满足，`user_name` 采用了范围匹配（获取的查询结果中，范围内的数据是乱序的），进而**导致age列无法保证有序的**。sex 是过滤列。

第三颗星，满足

上述我们看到，此时索引(user_name,sex,age)并不能满足三星索引中的第二颗星（排序）。

```sql
select user_name, sex, age 
from test 
where user_name like 'test%' and sex =1 
ORDER BY age
```

于是我们改改，建立索引`(sex,age,user_name)`：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1651212459071/85c6ee6f25934d5eadef7cca32ff278d.png" alt="image.png" style="zoom: 67%;" />

第一颗星，不满足，只可以匹配到sex，sex选择性很差，意味着是一个宽索引片(同时因为age也会导致排序选择的碎片问题)

第二颗星，满足，**等值sex 的情况下，age是有序的**，

第三颗星，满足，select查询的列都在索引列中，

对于索引(sex,age，user_name)我们可以看到，此时无法满足第一颗星（where后面的第一个要是等值匹配才符合）窄索引片的需求。

以上2个索引，都是无法同时满足三星索引设计中的三个需求的，**我们只能尽力满足2个**。而在多数情况下，能够满足2颗星，已经能缩小很大的查询范围了，具体最终要保留那一颗星（排序星 or 窄索引片星），这个就需要看查询者自己的着重点了，无法给出标准答案。1.3.MySQL性能调优
