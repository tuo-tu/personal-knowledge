# Redis高级特性和应用(发布 订阅、Stream)

## 发布和订阅

Redis提供了基于“发布/订阅”模式的消息机制，此种模式下，消息发布者和订阅者**不进行直接通信**,发布者客户端向指定的频道( channel)发布消息，订阅该频道的每个客户端都可以收到该消息。

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/9f6b447fb8024a3595352326d792ba95.png" alt="image.png" style="zoom: 67%;" />

### 操作命令

Redis主要提供了发布消息、**订阅频道**、取消订阅以及按照模式订阅和取消订阅等命令。

#### 发布消息

```bash
publish channel message
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/a04198c93ac0416eb677631af2adf868.png)

返回值是接收到信息的**订阅者数量**，如果是0说明没有订阅者，这条消息就丢了（再启动订阅者也不会收到）。

#### 订阅消息

```bash
subscribe channel [channel ...]
```

订阅者可以订阅一个或多个频道，如果此时另一个客户端发布一条消息，当前订阅者客户端会收到消息。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/1177c507c9d1491bb9a654a600604fab.png)

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/ede17eb77f294947b942524de78c651d.png)

如果有多个客户端同时订阅了同一个频道，都会收到消息。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/2f812f0fb14b45008410122fb7000a28.png)

客户端在执行订阅命令之后进入了订阅状态（类似于监听），只能接收subscribe、psubscribe、unsubscribe、 punsubscribe的四个命令。

#### 查询订阅情况

##### 查看活跃的频道

```
pubsub channels [pattern]
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/6a08c97b389e4807a69152504a2f5c80.png)

Pubsub 命令用于查看订阅与发布系统状态，包括活跃的频道（是指**当前频道至少有一个订阅者**），其中[pattern]是可以指定具体的模式，类似于通配符。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/87a515a74bb444e3b2a7a71d23caa0dc.png)

##### 查看频道订阅数

```bash
pubsub numsub channel
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/c42cb62a659e456991363af93a8dac0a.png)

最后也可以通过 help看具体的参数运用

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/1a4ab497c7f0489b93bd26fff1b3420d.png)

#### 使用场景和缺点

需要消息解耦又并不关注消息可靠性的地方都可以使用发布订阅模式。

PubSub 的生产者传递过来一个消息，Redis会直接找到相应的消费者传递过去。如果一个消费者都没有，那么消息直接丢弃。如果开始有三个消费者，一个消费者突然挂掉了，生产者会继续发送消息，另外两个消费者可以持续收到消息。但是挂掉的消费者重新连上的时候，这断连期间生产者发送的消息，对于这个消费者来说就是彻底丢失了。

所以和很多专业的消息队列系统（例如Kafka、RocketMQ)相比，**Redis 的发布订阅很粗糙**，例如无法实现消息堆积和回溯。但胜在足够简单，如果当前场景可以容忍的这些缺点，也不失为一个不错的选择。

正是因为 PubSub 有这些缺点，它的应用场景其实是非常狭窄的。从Redis5.0 新增了 Stream 数据结构，这个功能给 Redis 带来了**持久化消息队列**，我们马上将要学习到。

## Redis Stream

> Redis Stream 主要用于消息队列（MQ，Message Queue），Redis 本身是有一个 Redis 发布订阅 (pub/sub) 来实现消息队列的功能，但它有个缺点就是消息无法持久化，如果出现网络断开、Redis 宕机等，消息就会被丢弃。
>
> 简单来说发布订阅 (pub/sub) 可以分发消息，但无法记录历史消息。
>
> 而 Redis Stream 提供了消息的持久化和主备复制功能，可以让任何客户端访问任何时刻的数据，并且能记住每一个客户端的访问位置，还能保证消息不丢失。

Redis5.0 最大的新特性就是多出了一个数据结构 Stream，它是一个新的强大的**支持多播的可持久化的消息队列**，Redis的作者声明Redis Stream地借鉴了 Kafka 的设计。

### Stream总述

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/78e6d284fc6c4f8ab774287aad02f501.png)

Redis Stream 的结构如上图所示，每一个Stream都有一个**消息链表**，将所有加入的消息都串起来，每个消息都有一个唯一的 ID 和对应的内容。**消息是持久化的**，Redis 重启后，内容还在。

**具体的玩法如下：**

1、每个 Stream 都有唯一的名称，它就是 Redis 的 key，在我们首次使用xadd指令追加消息时自动创建。

> **XADD**
>
> 使用 XADD 向队列添加消息，如果指定的队列不存在，则创建一个队列，XADD 语法格式：
>
> ```
> XADD key ID field value [field value ...]
> ```
>
> - **key** ：队列名称，如果不存在就创建
> - **ID** ：消息 id，我们使用 * 表示**由 redis 生成**，可以自定义，但是要自己保证递增性。
> - **field value** ： 记录。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/f6287b28a1604bf29df9eb2b41388cfc.png)

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/baa99ab58cd4432bbf4848788a6dee74.png)

消息 ID 的形式是timestampInMillis-sequence，例如1527846880572-5，它表示当前的消息在毫米时间戳1527846880572时产生，并且是该毫秒内产生的第 5 条消息。消息 ID 可以由服务器自动生成（***代表默认自动**），也可以由客户端自己指定，但是形式必须是**整数-整数**，而且必须是后面加入的消息的 ID 要大于前面的消息 ID。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/d1b29c3c392c4c329226699acb889c37.png)

消息内容就是键值对，形如 hash 结构的键值对，这没什么特别之处。

2、每个 Stream 都可以挂多个**消费组**，每个消费组会有个游标**last_delivered_id**在 Stream 数组之上往前移动，表示当前消费组已经消费到哪条消息了。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/48efbf4656844581bad5dd7470eb4955.png)

每个消费组都有一个Stream 内唯一的名称，消费组不会自动创建，它需要单独的指令xgroup create进行创建，需要指定从 Stream 的某个消息 ID 开始消费，这个 ID 用来初始化last_delivered_id变量。

3、**每个消费组 (Consumer Group) 的状态都是独立的**，相互不受影响。也就是说同一份 Stream 内部的消息会被每个消费组都消费到。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/cb27e53f2451443f96925f58689f7cbb.png)

4、同一个消费组 (Consumer Group) 可以挂接多个消费者 (Consumer)，这些消费者之间是**竞争关系**，任意一个消费者读取了消息都会使游标last_delivered_id往前移动。每个消费者有一个组内唯一名称。

5、消费者 (Consumer) 内部会有个状态变量pending_ids，它记录了当前**已经被客户端读取，但是还没有 ack的消息**。如果客户端没有 ack，这个变量里面的消息 ID 会越来越多，一旦某个消息被 ack，它就开始减少。这个 pending_ids 变量在 Redis 官方被称之为**PEL**，也就是**Pending Entries List**，这是一个很核心的数据结构，它用来**确保客户端至少消费了消息一次**，而不会在网络传输的中途丢失了没处理。

> 在 Redis Streams 中，**Pending Entries List（PEL）** 是与消费者组相关联的一种内部数据结构。它用于跟踪**被消费者读取但尚未确认（`XACK`）的消息。**
>
> PEL 是在 Redis Streams 的消费者组中维护的。

### 常用操作命令

#### 生产端

**xadd 追加消息**

`XADD` 用于向流（Stream）**添加新消息**。每条消息由 **唯一 ID** 和 **键值对字段** 组成，流中的消息按 ID 排序。

```
XADD key [MAXLEN [~|=] count] [NOMKSTREAM] *|ID field value [field value ...]
```

**参数说明**

1. **`key`**: 流的名称。如果不存在，会自动创建流（除非使用 `NOMKSTREAM` 参数）。
2. **`MAXLEN [~|=] count`**: 控制流的最大长度。支持精确匹配（`=`）和近似匹配（`~`）。
   - `~`: 近似修剪，提升性能。
   - `=`: 精确修剪。
3. **`NOMKSTREAM`**: 如果流不存在，则命令失败，而不是创建一个新的流。
4. **`\*|ID`**:
   - 使用 `*` 时，Redis 自动生成消息 ID（格式为 `<ms>-<seq>`，如 `1688629596925-0`）。
   - 提供具体 ID（如 `1688629596925-1`），需确保不与现有消息 ID 冲突。
5. **`field value [field value ...]`**: 消息的键值对字段。

xadd第一次对于一个stream使用可以生成一个stream的结构

```bash
xadd streamtest * name lijin age 18
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/9ed2e1ca3e144bb784eae61ca4b53552.png)

*号表示服务器自动生成 ID，后面顺序跟着一堆 key-value

1626705954593-0 则是生成的消息 ID，由两部分组成：**时间戳-序号**。时间戳时毫秒级单位，是生成消息的Redis服务器时间，它是个64位整型。序号是在这个毫秒时间点内的**消息序号**。它也是个64位整型。

为了保证消息是有序的，因此Redis生成的ID是单调递增有序的。由于ID中包含时间戳部分，为了避免服务器时间错误而带来的问题（例如服务器时间延后了），Redis的每个Stream类型数据都维护一个**latest_generated_id**属性，用于记录**最后一个消息的ID**。若发现当前时间戳退后（小于latest_generated_id所记录的），则采用时间戳不变而序号递增的方案来作为新消息ID（这也是序号为什么使用int64的原因，保证有足够多的的序号），从而保证ID的单调递增性质。

强烈建议使用Redis的方案生成消息ID，因为这种时间戳+序号的单调递增的ID方案，几乎可以满足你全部的需求。但ID是支持自定义的。

**xrange 获取消息列表，会自动过滤已经删除的消息**

```bash
xrange streamtest - +
```

其中-表示最小值 , + 表示最大值

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/af90e1953f8b45f6bb2a41bf5776916c.png" alt="image.png" style="zoom:67%;" />

或者我们可以指定消息 ID 的列表：

```bash
xrange streamtest - 1665646270814-0
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/867166efa9ef440ea7e4dd8c8dd79960.png)

**xlen 消息长度**

```bash
xlen streamtest
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/9f115fc376dd4fd18ecf98d177e22d70.png)

**del 删除 Stream**

- del streamtest  删除整个 Stream


![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/20af39a3b407484e8c7d3608453f56f9.png)

- xdel可以删除指定的消息(指定ID)

  ```bash
  XDEL key ID [ID ...]
  ```

  该命令返回删除的消息 ID。

  **理解删除条目的底层细节**

  Redis 流以一种使其内存高效的方式表示：使用基数树来索引包含线性数十个 Stream 条目的宏节点。**通常，当你从 Stream 中删除一个条目的时候，条目并没有真正被驱逐，只是被标记为删除。**

  最终，如果宏节点中的所有条目都被标记为删除，则会销毁整个节点，并回收内存。这意味着如果你从 Stream 里删除大量的条目，比如超过 50% 的条目，则每一个条目的内存占用可能会增加，因为 Stream 将会开始变得碎片化。然而，流的表现将保持不变。

  在 Redis 未来的版本中，当一个宏节点内删除条目达到一定数量的时候，我们有可能会触发节点垃圾回收机制。目前，根据我们对这种数据结构的预期用途，还不太适合增加这样的复杂度。


![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/30f40f46b0b047568f9b1c8d8098f592.png)

#### 消费端

##### 单消费者

虽然Stream中有消费者组的概念，但是可以在不定义消费组的情况下进行 Stream 消息的独立消费，当 Stream 没有新消息时，甚至可以阻塞等待。Redis 设计了一个单独的消费指令xread，可以将 Stream 当成普通的消息队列 (list) 来使用。**使用 xread 时，我们可以完全忽略消费组 (Consumer Group) 的存在**，就好比 Stream 就是一个普通的列表 (list)。

使用 XREAD 以阻塞或非阻塞方式**获取消息列表** ，语法格式：

```bash
XREAD [COUNT count] [BLOCK milliseconds] STREAMS key [key ...] id [id ...]
```

- **count** ：指定要读取的消息数量，默认为1。
- **milliseconds** ：可选，阻塞毫秒数，没有设置就是非阻塞模式
- **key** ：指定要读取的Stream，可以指定多个Stream。
- **id** ：指定读取的起始位置，即要获取的消息ID。

**实例：**

```bash
# 从 Stream 头部读取两条消息
> XREAD COUNT 2 STREAMS mystream writers 0-0 0-0
1) 1) "mystream"
   2) 1) 1) 1526984818136-0
         2) 1) "duration"
            2) "1532"
            3) "event-id"
            4) "5"
            5) "user-id"
            6) "7782813"
      2) 1) 1526999352406-0
         2) 1) "duration"
            2) "812"
            3) "event-id"
            4) "9"
            5) "user-id"
            6) "388234"
2) 1) "writers"
   2) 1) 1) 1526985676425-0
         2) 1) "name"
            2) "Virginia"
            3) "surname"
            4) "Woolf"
      2) 1) 1526985685298-0
         2) 1) "name"
            2) "Jane"
            3) "surname"
            4) "Austen"
```

再举例：

```bash
xread count 1 streams stream2 0-0
```

表示从 Stream 头部读取1条消息，0-0指从每个流的起点开始读取消息。

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/aff55f71cbbb490aad07a1ab9c631f3c.png" alt="image.png" style="zoom:80%;" />



也可以指定从streams的消息id开始（**不包括命令中的消息id**）

```
xread count 2 streams stream1 1665644057564-0
```

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/428bdd838ad44d279756f4946adf9596.png" alt="image.png" style="zoom:80%;" />

```bash
xread count 1 streams stream1 $
```

- **`$`**：表示从流的当前尾部读取消息，就是从尾部读取最新的一条消息。

> 注意：如果流在 `$` 时刻没有新消息，此命令不会返回任何消息。



<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/04dd6f0da5a647c6b5b4f85767e975ff.png" alt="image.png" style="zoom:80%;" />

所以应该以阻塞的方式读取尾部最新的一条消息，直到新的消息的到来。

```bash
xread block 0 count 1 streams stream1 $
```

block后面的数字代表阻塞时间，单位毫秒，**0代表一直阻塞，直到有新消息可供读取。**

此时我们新开一个客户端，往stream1中写入一条消息。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/c59a0638fe51419f84bfd4a4f8f2a8d0.png)

可以看到看到阻塞解除了，返回了新的消息内容，而且还显示了一个等待时间，这里我们等待了10.82s

一般来说客户端如果想要使用 xread 进行顺序消费，一定要记住当前消费到哪里了，也就是返回的消息 ID。下次继续调用 xread 时，将上次返回的最后一个消息 ID 作为参数传递进去，就可以继续消费后续的消息。不然很容易重复消息，基于这点单**消费者基本上没啥运用场景**，本课也不深入去讲。

##### 消费组

###### 创建消费组

Stream 通过**xgroup create**指令创建消费组 (Consumer Group)，需要传递起始消息 ID 参数用来初始化last_delivered_id变量。

0-表示从头开始消费

```bash
xgroup create stream1 c1 0-0
```

$ 表示从尾部开始消费，**只接受新消息，当前 Stream 消息会全部忽略。**

```bash
xgroup create stream1 c2 $
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/ccd3dd2307a043c1a127621f2d923848.png)

现在我们可以用xinfo命令来看看stream1的情况：

```bash
xinfo stream stream1
```

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/2a4ff34210e84e86ba857a91eb2ce7e9.png" alt="image.png" style="zoom:80%;" />

查看stream1的消费组的情况：

```bash
xinfo groups stream1
```

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/78bda19da7044cfba93a37ce65292146.png" alt="image.png" style="zoom:80%;" />

###### 消息消费

有了消费组，自然还需要消费者，Stream提供了 xreadgroup 指令可以进行消费组的组内消费，需要提供**消费组名称、消费者名称和起始消息 ID。**

它同 xread 一样，也可以阻塞等待新消息。读到新消息后，对应的消息 ID 就会进入消费者的PEL(正在处理的消息) 结构里，客户端处理完毕后使用 xack 指令通知服务器，本条消息已经处理完毕，该消息 ID 就会从 PEL 中移除。

命令：

```bash
XREADGROUP GROUP groupname consumer [COUNT count] [BLOCK milliseconds] STREAMS key [key ...] ID [ID ...]
```

**参数说明**

1. **`GROUP groupname consumer`**:
   - **`groupname`**: 消费组的名称，需先使用 `XGROUP CREATE` 或 `XGroupCreateMkStream` 创建。
   - **`consumer`**: 消费者名称，用于标识消息的具体消费者。
2. **`COUNT count`** (可选): 每次读取的最大消息条数。
3. **`BLOCK milliseconds`** (可选): 如果没有消息，阻塞指定毫秒时间，等待消息到达。
4. **`STREAMS key [key ...]`**: 指定流的名称，可以是一个或多个。
5. **`ID [ID ...]`**: 每个流的起始 ID：
   - **`>`**: 只读取未被其他消费者读取的消息。
   - **消息 ID**: 从指定的消息 ID 开始读取，包括该 ID。

```bash
xreadgroup GROUP c1 consumer1 count 1 streams stream1 >
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/75a5be4cc4c84ccc959407fe6744e0ad.png)

consumer1代表消费者的名字。

">"表示从当前消费组的 **last_delivered_id** 后面开始读，每当消费者读取一条消息，**last_delivered_id 变量就会前进**。前面我们定义cg1的时候是从头开始消费的，自然就获得stream1中第一条消息，再执行一次上面的命令，自然就读取到了下条消息。我们将Stream1中的消息读取完，很自然就没有消息可读了。

然后设置阻塞等待![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/80a52389edda4d25acf8d099756515b6.png)

我们新开一个客户端，发送消息到stream1回到原来的客户端，发现阻塞解除，收到新消息

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/5899b8f323634b7482c0ce9a2975ed7f.png)

我们来观察一下消费组状态

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/3cf7520617144b309fa02c6ef6b70ca1.png" alt="image.png" style="zoom:80%;" />

如果同一个**消费组有多个消费者**，我们还可以通过 xinfo consumers 指令观察组内所有消费者的状态

`XINFO CONSUMERS` 命令用于获取指定消费组中所有消费者的详细信息，包括消费者名称、未确认消息的数量等。

返回的是一个列表，包含消费组中每个消费者的详细信息

语法：

```bash
XINFO CONSUMERS key group
```

实例：这个例子中消费组c1目前只有一个消费者

```bash
xinfo consumers stream1 c1
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/75b6f8e647764a8598370fdbef4bf631.png)

可以看到目前c1消费组中consumer1这个消费者有 7 条待ACK的消息，空闲了2086176ms 没有读取消息。

如果我们确认一条消息

命令：

```bash
XACK key group ID [ID ...]
```

**参数说明**

1. **`key`**: 流的名称（Stream）。
2. **`group`**: 消费组的名称。
3. **`ID [ID ...]`**: 一个或多个消息的 ID，标记为已处理。

```bash
xack stream1 c1 1665647371850-0
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/21bcc8499e4a487c848afde464afc79a.png)

就可以看到待确认消息变成了6条

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/16985be5de824675aa72f21af08baa9d.png)

**xack允许带多个消息id。**

同时Stream还提供了命令 XPENDING 用来获取消费组或消费组内消费者的未处理完毕的消息。

 `XPENDING` 命令用于查看特定消费组中待确认的消息状态。这些消息是被消费者读取但尚未通过 `XACK` 确认的消息。

**命令格式：**

```bash
XPENDING <stream> <group>
```

- **`stream`**: 消息流的名称。
- **`group`**: 消费组的名称。

举例：

```bash
xpending stream1 c1
```

模拟结果：

`XPENDING` 命令返回一个摘要信息，结果如下：

```bash
1) (integer) 10
2) "1689328420951-0"
3) "1689328504225-0"
4) 1) 1) "consumer-1"
      2) (integer) 5
   2) 1) "consumer-2"
      2) (integer) 5
```

**字段解释**

1. **`10`**: 待确认消息的总数。
2. **`1689328420951-0`**: 待确认消息的最小消息 ID。
3. **`1689328504225-0`**: 待确认消息的最大消息 ID。
4. **消费者统计信息**:
   - `consumer-1` 待确认的消息数量为 `5`。
   - `consumer-2` 待确认的消息数量为 `5`。

具体操作细节可以参考：[xpending 命令 -- Redis中国用户组（CRUG）](http://www.redis.cn/commands/xpending.html)

**命令XCLAIM**用以进行消息转移的操作，将某个消息转移到自己的Pending列表中。通常在 Redis Streams 中，当某些消息在一个消费者中未被及时处理（例如消费者意外崩溃）时，可以使用此命令重新分配这些消息给其他消费者。

需要设置组、转移的目标消费者和消息ID，同时需要提供IDLE（已被读取时长），只有超过这个时长，才能被转移。

命令格式：

```bash
XCLAIM <stream> <group> <consumer> <min-idle-time> <id> [id ...] [OPTIONS]
```

**参数说明**

- **`stream`**: 消息流的名称。
- **`group`**: 消费组的名称。
- **`consumer`**: 新的消费者名称。
- **`min-idle-time`**: 指定消息的最小空闲时间（以毫秒为单位）。仅空闲时间超过此值的消息才会被转移。
- **`id`**: 要转移的消息 ID，可以指定一个或多个。

**可选参数**

- **`IDLE <ms>`**: 设置消息的空闲时间（以毫秒为单位）。如果未指定，空闲时间会根据实际计算。
- **`TIME <ms-timestamp>`**: 设置消息的接收时间戳。
- **`RETRYCOUNT <count>`**: 设置消息的重试次数。
- **`FORCE`**: 即使消息当前未处于待确认状态（未被任何消费者读取或已确认），也强制转移。
- **`JUSTID`**: 只返回消息 ID，而不是消息的完整内容。

**实例：**基本用法

将 `mystream` 中 ID 为 `1689328420951-0` 的消息，从 `consumer-1` 转移给 `consumer-2`。

```
XCLAIM mystream mygroup consumer-2 5000 1689328420951-0
```

**解释**:

- `mystream`: 消息流。
- `mygroup`: 消费组。
- `consumer-2`: 新的消费者。
- `5000`: 最小空闲时间，消息必须在 5000 毫秒内未被确认才会被转移。
- `1689328420951-0`: 要转移的消息 ID。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/a584a8bb1f274d61a6dc1b3976316fc0.png)[]

具体操作细节可参考：[xclaim 命令 -- Redis中国用户组（CRUG）](http://www.redis.cn/commands/xclaim.html)

### 在Redis中实现消息队列

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/03c4e0d96c40423ca834515d73c9dabb.png)

#### 基于pub/sub

##### go语言实现

1. **发布者代码**：将消息发布到指定频道

   ```go
   package main
   
   import (
   	"context"
   	"fmt"
   	"github.com/redis/go-redis/v9"
   )
   
   func main() {
   	// 创建 Redis 客户端
   	rdb := redis.NewClient(&redis.Options{
   		Addr: "localhost:6379", // Redis 地址
   		DB:   0,                // 使用默认数据库
   	})
   
   	// 发布消息
   	channel := "my_channel"
   	message := "Hello, Subscribers!"
   	ctx := context.Background()
   
   	err := rdb.Publish(ctx, channel, message).Err()
   	if err != nil {
   		fmt.Printf("Error publishing message: %v\n", err)
   		return
   	}
   
   	fmt.Printf("Message '%s' published to channel '%s'\n", message, channel)
   }
   
   ```

2. **订阅者代码**：监听指定频道并接收消息

   ```go
   package main
   
   import (
   	"context"
   	"fmt"
   	"github.com/redis/go-redis/v9"
   )
   
   func main() {
   	// 创建 Redis 客户端
   	rdb := redis.NewClient(&redis.Options{
   		Addr: "localhost:6379", // Redis 地址
   		DB:   0,                // 使用默认数据库
   	})
   
   	// 创建 PubSub 对象
   	ctx := context.Background()
   	pubsub := rdb.Subscribe(ctx, "my_channel")
   
   	// 等待订阅确认
   	_, err := pubsub.Receive(ctx)
   	if err != nil {
   		fmt.Printf("Error subscribing to channel: %v\n", err)
   		return
   	}
   
   	fmt.Println("Subscribed to channel 'my_channel'")
   
   	// 接收消息
   	ch := pubsub.Channel()
   	for msg := range ch {
   		fmt.Printf("Received message from channel '%s': %s\n", msg.Channel, msg.Payload)
   	}
   }
   ```

##### java实现（不用看）

注意必须继承JedisPubSub这个抽象类

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/81da1646a67b4486b26021ff672a53eb.png)

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/b3b97617ed984b3dae34b83b582e00bb.png)

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/d5f7698eed88443682061d7b137451b7.png)

#### 基于Stream

##### go语言实现

1. **发布者代码**：发布者负责将消息添加到 Redis Stream 中。（使用xadd命令添加消息）

   ```go
   package main
   
   import (
   	"context"
   	"fmt"
   	"github.com/redis/go-redis/v9"
   	"time"
   )
   
   func main() {
   	// 创建 Redis 客户端
   	rdb := redis.NewClient(&redis.Options{
   		Addr: "localhost:6379", // Redis 地址
   		DB:   0,                // 使用默认数据库
   	})
   
   	ctx := context.Background()
   	streamName := "my_stream"
   
   	// 发布消息
   	for i := 1; i <= 5; i++ {
   		message := fmt.Sprintf("Message %d", i)
   		_, err := rdb.XAdd(ctx, &redis.XAddArgs{
   			Stream: streamName,
   			Values: map[string]interface{}{"content": message},
   		}).Result()
   		if err != nil {
   			fmt.Printf("Error adding message to stream: %v\n", err)
   			return
   		}
   		fmt.Printf("Added message: %s to stream: %s\n", message, streamName)
   		time.Sleep(1 * time.Second) // 模拟延时
   	}
   }
   ```

2. **消费者代码：**消费者从 Redis Stream 中读取消息，并确认处理完成。

   `XGroupCreateMkStream` 是 Redis 提供的命令，用于创建一个新的消费组，并在**流不存在时创建一个新的流。**

   ```bash
   XGROUP CREATE key groupname id [MKSTREAM]
   ```

   **参数说明：**

   - **`key`**: 流的名称。
   - **`groupname`**: 要创建的消费组的名称。
   - **`id`**: 消费组开始读取的消息 ID。
     - 常用值：
       - **`$`**: 从最新的消息开始。
       - **`0`**: 从流的第一个消息开始。
   - **`MKSTREAM`** (可选): 如果流不存在，则创建流。

   ```go
   package main
   
   import (
   	"context"
   	"fmt"
   	"github.com/redis/go-redis/v9"
   )
   
   func main() {
   	// 创建 Redis 客户端
   	rdb := redis.NewClient(&redis.Options{
   		Addr: "localhost:6379", // Redis 地址
   		DB:   0,                // 使用默认数据库
   	})
   
   	ctx := context.Background()
   	streamName := "my_stream"
   	groupName := "my_group"
   	consumerName := "consumer1"
   
   	// 1.确保消费组存在
   	err := rdb.XGroupCreateMkStream(ctx, streamName, groupName, "$").Err()
   	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
   		fmt.Printf("Error creating group: %v\n", err)
   		return
   	}
   
   	fmt.Printf("Consumer '%s' started in group '%s'\n", consumerName, groupName)
   
   	// 2.监听 Stream 消息
   	for {
   		messages, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
   			Group:    groupName,
   			Consumer: consumerName,
   			Streams:  []string{streamName, ">"},
   			Block:    0, // 阻塞直到消息到达
   			Count:    1, // 每次读取一条消息
   		}).Result()
   
   		if err != nil {
   			fmt.Printf("Error reading from stream: %v\n", err)
   			continue
   		}
          
           // 3.处理消息，并确认（ack）
   		for _, message := range messages {
   			for _, value := range message.Messages {
   				fmt.Printf("Consumer '%s' received message: %v\n", consumerName, value.Values)
   				// 确认消息
   				err = rdb.XAck(ctx, streamName, groupName, value.ID).Err()
   				if err != nil {
   					fmt.Printf("Error acknowledging message: %v\n", err)
   				}
   			}
   		}
   	}
   }
   ```

##### java实现（不用看）

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/77b699b89e7246449fd2e86e2a8d5acb.png)![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/98c6d17fd9e74c8c92190126da1b1fda.png)

java封装了两个类用于处理消息及消息的元数据。

StreamEntry和StreamEntryID

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/523f8ac3bb1f4c198f40c1496aa22fe5.png)

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/1d017710e9724326867d898a8935d5a0.png)

#### Redis中几种消息队列实现的总结

##### 基于List的 LPUSH+BRPOP 的实现

足够简单，消费消息延迟几乎为零，但是需要处理空闲连接的问题。

如果线程一直阻塞在那里，Redis客户端的连接就成了闲置连接，闲置过久，服务器一般会主动断开连接，减少闲置资源占用，这个时候blpop和brpop或抛出异常，所以在编写客户端消费者的时候要小心，如果捕获到异常，还要重试。

其他缺点包括：

做消费者确认ACK麻烦，不能保证消费者消费消息后是否成功处理的问题（宕机或处理异常等），通常需要维护一个Pending列表，保证消息处理确认；不能做广播模式，如pub/sub，消息发布/订阅模型；不能重复消费，一旦消费就会被删除；不支持分组消费。

##### 基于Sorted-Set的实现

多用来实现延迟队列，当然也可以实现有序的普通的消息队列，但是**消费者无法阻塞的获取消息**，只能轮询，不允许重复消息。

##### PUB/SUB，订阅/发布模式

优点：

典型的广播模式，一个消息可以发布到多个消费者；多信道订阅，消费者可以同时订阅多个信道，从而接收多类消息；消息即时发送，消息不用等待消费者读取，消费者会自动接收到信道发布的消息。

缺点：

消息一旦发布，没有及时接收会丢失。换句话就是**发布时若客户端不在线，则消息丢失，不能寻回**；不能保证每个消费者接收的时间是一致的；若消费者客户端出现消息积压，到一定程度，会被强制断开，导致消息意外丢失。通常发生在消息的生产远大于消费速度时；可见，Pub/Sub 模式不适合做消息存储，消息积压类的业务，而是擅长处理广播，即时通讯，即时反馈的业务。

##### 基于Stream类型的实现

基本上已经有了一个消息中间件的雏形，**可以考虑在生产过程中使用。**

### 消息队列问题

从我们上面对Stream的使用表明，Stream已经具备了一个消息队列的基本要素，生产者API、消费者API，消息Broker，消息的确认机制等等，所以在使用消息中间件中产生的问题，这里一样也会遇到。

#### Stream 消息太多怎么办?

要是消息积累太多，Stream 的链表岂不是很长，内容会不会爆掉？**xdel 指令又不会删除消息，它只是给消息做了个标志位。**

Redis 自然考虑到了这一点，所以它提供了一个定长 Stream 功能。在 xadd 的指令提供一个**定长长度 maxlen**，就可以将老的消息干掉，确保最多不超过指定长度。

#### 消息如果忘记 ACK 会怎样?

Stream 在每个消费者结构中保存了**正在处理中**的消息 ID 列表 PEL，如果消费者收到了消息处理完了但是没有回复 ack，就会导致 PEL 列表不断增长，如果有很多消费组的话，那么这个 PEL 占用的内存就会放大。所以**消息要尽可能的快速消费并确认。**

#### PEL如何避免消息丢失?

在客户端消费者读取 Stream 消息时，Redis 服务器将消息发给客户端的过程中，客户端突然断开了连接，消息就丢失了。但是 PEL 里已经保存了发送者发出去的消息 ID。待客户端重新连上之后，可以再次收到 PEL 中的消息 ID 列表。不过此时 xreadgroup 的起始消息 ID 不能为参数>，而必须是任意有效的消息 ID，一般将参数设为 0-0，表示**读取所有的 PEL 消息**以及自last_delivered_id之后的新消息。

#### 死信问题

**什么是死信（Dead Letter Message）？**

在消息队列或流处理系统（例如 Redis Streams）中，**死信**是指由于某些原因**未被成功处理或无法被消费者消费的消息**。这些消息通常被移动到一个专门的队列或流中，供开发人员进行进一步的排查或处理。

每个Pending的消息有4个属性：

- 消息ID
- 所属消费者
- IDLE，已读取时长
- delivery counter，消息被读取次数

如果某个消息，不能被消费者处理，也就是不能被XACK，这是要长时间处于Pending列表中，即使被反复的转移给各个消费者也是如此。此时该消息的delivery counter（通过XPENDING可以查询到）就会累加，当累加到某个我们预设的临界值时，我们就认为是坏消息（也叫死信，Dead Letter，无法投递的消息），由于有了判定条件，我们将坏消息处理掉即可（删除即可）。删除一个消息，使用XDEL语法。注意，这个命令并没有真正删除Pending中的消息，因此查看Pending，消息还会在。可以在执行XDEL之后，再XACK这个消息，标识其处理完毕。（这里的处理方式似乎不妥，看下面的处理方式）

> **为什么会出现死信？**
>
> 1. **消费者未确认消息**：
>    - 消费者组中的消费者从 Stream 读取消息后，如果没有发送 `XACK` 命令确认消息成功处理，那么消息会被保留在 Pending Entries List (PEL) 中。
>    - 如果该消息滞留在 PEL 中的时间超过预期，就可能需要处理为死信。
> 2. **消费者崩溃或超时**：
>    - 消费者在处理消息时崩溃或因网络问题断开连接，导致消息无法被处理和确认。
> 3. **重复处理失败**：
>    - 多次重试消费某条消息失败时，可以将其标记为死信。
> 4. **处理超时**：
>    - 如果消费者未能在指定时间内完成消息处理，消息可能需要被转移到死信流。
>
> **Redis Streams 的死信处理**
>
> Redis Streams 没有内置的死信队列功能，但可以通过以下方式实现。
>
> **1.创建死信流**
>
> 死信流是一个独立的 Stream，用于存储被标记为死信的消息。
>
> **2. 监控 PEL 中的消息**
>
> 使用 `XPENDING` 命令查看滞留在 PEL 中的消息：可以根据消息滞留时间，判断是否需要将其转移到死信流。
>
> **3. 转移死信**
>
> 将超时未确认的消息转移到死信流：①读取原stream中的死信；②使用XADD将该消息加入到deadLetterStream；③在原stream中使用xack确认掉该消息；

**Go 实现 Redis Streams 死信处理**

以下是一个示例程序，使用 Go 将死信转移到死信流中。

```go
go复制编辑package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v8"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx := context.Background()

	stream := "mystream"
	group := "mygroup"
	consumer := "consumer-1"
	deadLetterStream := "dead_letter_stream"

	// 获取 PEL 中的消息
	pending, err := rdb.XPending(ctx, stream, group).Result()
	if err != nil {
		fmt.Println("Error fetching pending messages:", err)
		return
	}

	// 遍历每条消息，判断是否超时
	for _, message := range pending.Consumers {
		for _, entry := range message.Entries {
			// 假设超时时间为 1 分钟
			if entry.Idle > time.Minute {
				// 转移到死信流
				err = rdb.XAdd(ctx, &redis.XAddArgs{
					Stream: deadLetterStream,
					Values: map[string]interface{}{
						"original_id": entry.ID,
						"reason":      "processing timeout",
					},
				}).Err()
				if err != nil {
					fmt.Println("Error adding to dead letter stream:", err)
					continue
				}

				// 确认原始流消息，避免重复转移
				err = rdb.XAck(ctx, stream, group, entry.ID).Err()
				if err != nil {
					fmt.Println("Error acknowledging message:", err)
					continue
				}

				fmt.Printf("Moved message %s to dead letter stream\n", entry.ID)
			}
		}
	}
}
```

> **5. 消费死信**
>
> 死信流可以被消费者组消费，用于分析、告警或重试：
>
> ```bash
> XREADGROUP GROUP dlq_group consumer-1 STREAMS dead_letter_stream >
> ```

#### Stream 的高可用

Stream 的高可用建立在主从复制的基础上，它和其它数据结构的复制机制没有区别，也就是说在 Sentinel 和 Cluster 集群环境下 Stream 是可以支持高可用的。不过鉴于 Redis 的指令复制是异步的，在 failover 发生时，Redis 可能会丢失极小部分数据，这点 Redis 的其它数据结构也是一样的。

#### 分区 Partition

Redis 的服务器没有原生支持分区能力，如果想要使用分区，那就**需要分配多个 Stream**，然后在客户端使用一定的策略来生产消息到不同的 Stream。

### Stream小结

Stream 的消费模型借鉴了Kafka 的消费分组的概念，它弥补了 Redis Pub/Sub 不能持久化消息的缺陷。但是它又不同于 kafka，Kafka 的消息可以分 partition，而 Stream 不行。如果非要分 parition 的话，得在客户端做，提供不同的 Stream 名称，对消息进行 hash 取模来选择往哪个 Stream 里塞。

关于 Redis 是否适合做消息队列，业界一直是有争论的。很多人认为，要使用消息队列，就应该采用 Kafka、RabbitMQ 这些专门面向消息队列场景的软件，而 **Redis 更加适合做缓存。**

根据这些年做 Redis 研发工作的经验，我的看法是：Redis 是一个非常轻量级的键值数据库，部署一个 Redis 实例就是启动一个进程，部署 Redis 集群，也就是部署多个 Redis 实例。而 Kafka、RabbitMQ 部署时，涉及额外的组件，例如 Kafka 的运行就需要再部署ZooKeeper。**相比 Redis 来说，Kafka 和 RabbitMQ 一般被认为是重量级的消息队列。**

所以，关于是否用 Redis 做消息队列的问题，不能一概而论，我们需要考虑业务层面的数据体量，以及对性能、可靠性、可扩展性的需求。如果分布式系统中的组件消息通信量不大，那么Redis只需要使用有限的内存空间就能满足消息存储的需求，而且，Redis 的高性能特性能支持快速的消息读写，不失为消息队列的一个好的解决方案。

## Redis的Key和Value的数据结构组织

### 全局哈希表

为了实现从键到值的快速访问，**Redis 使用了一个哈希表来保存所有键值对。**哈希表的存储结构本质上是一个数组，数组的每个元素称为一个**哈希桶**，每个桶可以存储多个键值对。所以，我们常说，一个哈希表是由多个哈希桶组成的，每个哈希桶中保存了键值对数据。

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/1eb5ca015a3b4389ad4ede842d98df1f.png" alt="image.png" style="zoom: 67%;" />

**Entry（链地址法）**
 每个数组位置指向一个链表节点（`entry`），链表节点用于存储具体的键值对。

- **冲突解决**：
  - 当多个键通过哈希函数映射到同一个数组索引时，冲突的键值对会被存储在该索引的链表中。
  - 图中展示了一个链表结构，通过指针（`*next`）链接多个节点。
- **示例**：
  - 如果两个键 `key1` 和 `key2` 计算出的哈希值相同，则它们都会存储在数组索引 3 的链表中，链表头是 `entry`。

 **整体工作流程**

- **插入操作**：
  1. 键通过哈希函数计算出哈希值。
  2. 哈希值对应数组的某个索引。
  3. 如果该索引已有数据，新的键值对会以链表的形式插入到该索引的链表中。
- **查找操作**：
  1. 根据键计算出哈希值，找到对应的数组索引。
  2. 在该索引的链表中遍历，查找匹配的键。
- **删除操作**：
  1. 查找到对应的链表节点。
  2. 删除该节点，同时调整链表的指针。

哈希桶中的 entry 元素中保存了key和value指针，分别指向了实际的键和值，这样一来，即使值是一个集合，也可以通过*next指针被查找到。因为这个哈希表保存了所有的键值对，所以，我也把它称为全局哈希表。

哈希表的最大好处很明显，就是让我们可以用 **O(1)** 的时间复杂度来快速查找到键值对：我们只需要计算键的哈希值，就可以知道它所对应的哈希桶位置，然后就可以访问相应的 entry 元素。

但当你往 Redis 中写入大量数据后，就可能发现操作有时候会突然变慢了。这其实是因为你忽略了一个潜在的风险点，那就是哈希表的冲突问题和 **rehash** 可能带来的操作阻塞。

当你往哈希表中写入更多数据时，哈希冲突是不可避免的问题。这里的哈希冲突，两个 key 的哈希值和哈希桶计算对应关系时，正好落在了同一个哈希桶中。

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/f2a1f73d63f1428cad1324a7b67283b4.png" alt="image.png" style="zoom:67%;" />

Redis 解决哈希冲突的方式，就是链式哈希。链式哈希也很容易理解，就是指同一个哈希桶中的多个元素用一个链表来保存，它们之间依次用指针连接。

当然如果这个数组一直不变，那么hash冲突会变很多，这个时候检索效率会大打折扣，所以Redis就需要把**数组进行扩容**（一般是扩大到原来的两倍），但是问题来了，扩容后每个hash桶的数据会分散到不同的位置，这里涉及到元素的移动，必定会阻塞IO，所以这个rehash过程会导致很多请求阻塞。

### 渐进式rehash

为了避免这个问题，Redis 采用了渐进式 rehash。

首先，Redis 默认使用了两个全局哈希表：哈希表 1 和哈希表 2。一开始，当你刚插入数据时，默认使用哈希表 1，此时的哈希表 2 并没有被分配空间。随着数据逐步增多，Redis 开始执行 rehash。

1. 给哈希表 2 分配更大的空间，例如是当前哈希表 1 大小的两倍；
2. 把哈希表 1 中的数据重新映射并拷贝到哈希表 2 中；

3. 释放哈希表 1 的空间。

在上面的第二步涉及大量的数据拷贝，如果一次性把哈希表 1 中的数据都迁移完，会造成 Redis 线程阻塞，无法服务其他请求。此时，Redis 就无法快速访问数据了。

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665642245080/2406a8206e944d449b03f1f390bedf0d.png" alt="image.png"  />

在Redis开始执行 rehash时，Redis仍然正常处理客户端请求，但是要加入一个额外的处理：

1. 处理第1个请求时，把哈希表 1中的第1个索引位置上的所有 entries 拷贝到哈希表 2 中；

2. 处理第2个请求时，把哈希表 1中的第2个索引位置上的所有 entries 拷贝到哈希表 2 中；

3. 如此循环，直到把所有的索引位置的数据都拷贝到哈希表 2 中。


**这样就巧妙的把一次性大量拷贝的开销，分摊到了多次处理请求的过程中**，避免了耗时操作，保证了数据的快速访问。所以这里基本上也可以确保根据key找value的操作在O(1)左右。

不过这里要注意，如果Redis中有海量的key值的话，这个rehash过程会很长很长，虽然采用渐进式Rehash，但在Rehash的过程中还是会导致请求有不小的卡顿。并且像一些统计命令也会非常卡顿：比如keys

按照Redis的配置每个实例（每个**实例**表示一个运行的 Redis 服务进程）能存储的最大的key的数量为2的32次方，即2.5亿，但是尽量把key的数量控制在千万以下，这样就可以避免rehash导致的卡顿问题，如果数量确实比较多，建议采用分区hash存储。
