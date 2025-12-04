# Redis高并发高可用

## 主从复制

在分布式系统中为了解决单点问题，通常会把数据复制多个副本部署到其他机器，满足故障恢复和负载均衡等需求。Redis也是如此，它为我们提供了复制功能，实现了相同数据的多个Redis 副本。复制功能是高可用Redis的基础，后面章节的哨兵和集群都是在复制的基础上实现高可用的。

默认情况下，Redis都是主节点。每个从节点只能有一个主节点，而主节点可以同时具有多个从节点。**复制的数据流是单向的**，只能由主节点复制到从节点。

### 复制的拓扑结构

Redis的复制拓扑结构可以支持单层或多层复制关系，根据拓扑复杂性可以分为以下三种：一主一从、一主多从、树状主从结构，下面分别介绍。

#### 一主一从结构

一主一从结构是最简单的复制拓扑结构，用于主节点出现宕机时从节点提供故障转移支持。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/d634f64c7b6342c9b82d306ad4845e23.png)

当应用写命令并发量较高且需要持久化时，可以**只在从节点上开启AOF**，这样既保证数据安全性同时也避免了持久化对主节点的性能干扰。但需要注意的是，当主节点关闭持久化功能时，如果主节点脱机要避免自动重启操作。

因为主节点之前没有开启持久化功能自动重启后数据集为空，这时从节点如果继续复制主节点会导致从节点数据也被清空的情况，丧失了持久化的意义。安全的做法是在从节点上执行**slaveof no one**断开与主节点的复制关系，再重启主节点从而避免这一问题。

#### 一主多从结构

一主多从结构(又称为星形拓扑结构）使得应用端可以利用多个从节点实现**读写分离。**

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/b7dedbfd19d24db48e82c91719f59114.png)

对于读占比较大的场景，可以把读命令发送到从节点来分担主节点压力。同时在日常开发中如果需要执行一些比较耗时的读命令，如：keys、sort等，可以在其中一台从节点上执行，防止慢查询对主节点造成阻塞从而影响线上服务的稳定性。对于写并发量较高的场景，多个从节点会导致主节点写命令的多次发送从而过度消耗网络带宽，同时也加重了主节点的负载影响服务稳定性。

#### 树状主从结构

树状主从结构(又称为树状拓扑结构）使得从节点不但可以复制主节点数据，同时可以作为其他从节点的主节点继续向下层复制。通过引入复制中间层，可以有效降低主节点负载和需要传送给从节点的数据量。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/8f8edd042ff34531932e699be31bf9c5.png)

数据写入节点A后会同步到B和C节点，B节点再把数据同步到D和E节点，数据实现了一层一层的向下复制。当主节点需要挂载多个从节点时为了避免对主节点的性能干扰,可以采用树状主从结构降低主节点压力。

### 复制的配置

#### 建立复制

参与复制的Redis实例划分为主节点(master)和从节点(slave)。默认情况下，Redis都是主节点。每个从节点只能有一个主节点，而主节点可以同时具有多个从节点。复制的数据流是单向的，只能由主节点复制到从节点。

**配置复制的方式有以下三种**

1. 在配置文件中加入slaveof{masterHost } {masterPort}随 Redis启动生效。
2. 在redis-server启动命令后加入--slaveof{masterHost} {masterPort }生效。

3. 直接使用命令：slaveof {masterHost} { masterPort}生效。

综上所述，slaveof命令在使用时，可以运行期动态配置，也可以提前写到配置文件中。

比如：我在机器上启动2台Redis, 分别是6379 和6380 两个端口。

> flushall 清空数据库并执行持久化操作。
>
> flushdb 清空数据库，但是不执行持久化操作。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/ce991b66b69947159848f8cc047abed2.png)

**slaveof本身是异步命令**，执行slaveof命令时，节点只保存主节点信息后返回，后续复制流程在节点内部异步执行，具体细节见之后。主从节点复制成功建立后，可以使用**info replication**命令查看复制相关状态。

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/58e7f3a1d88243bda30e9206226f70f6.png" alt="image.png" style="zoom:80%;" />

#### 断开复制

slaveof命令不但可以建立复制，还可以在从节点执行**slaveof no one**来断开与主节点复制关系。例如在6380节点上执行slaveof no one来断开复制。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/bf21031a0cf243e68ec0d9c722c9855c.png)

**断开复制主要流程：**

1. 断开与主节点复制关系。
2. 从节点晋升为主节点。

从节点断开复制后并**不会抛弃原有数据**，只是无法再获取主节点上的数据变化。

通过slaveof命令还可以实现**切主操作**，所谓切主是指把当前从节点对主节点的复制切换到另一个主节点。

执行`slaveof{ newMasterIp} { newMasterPort}`命令即可，例如把6880节点从原来的复制6879节点变为复制6881节点。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/e904361290e54bb98ee01b66dd246dcb.png)

切主内部流程如下:

1. 断开与旧主节点复制关系。
2. 与新主节点建立复制关系。

3. **删除从节点当前所有数据。**
4. 对新主节点进行复制操作。

#### 只读

默认情况下，从节点使用**slave-read-only=yes**配置为只读模式。由于复制只能从主节点到从节点，对于从节点的任何修改主节点都无法感知，修改从节点会造成主从数据不一致。因此建议设置线上**不要修改从节点的只读模式。**

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/aec38b7f9d8545cbbdc85bea9183fae0.png)

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/92dbcb842b794158b611f085bd61211b.png)

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/2b0c0be29d5c44dcb1d2c3ffd1c05866.png)

#### 传输延迟

**主从节点一般部署在不同机器上**，复制时的网络延迟就成为需要考虑的问题，Redis为我们提供了`repl-disable-tcp-nodelay`参数用于控制是否关闭`TCP_NODELAY`，**默认关闭**，说明如下：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/fea28b954f964c35ad77dde58ea34020.png)

- **当关闭时**，主节点产生的命令数据无论大小都会及时地发送给从节点，这样主从之间延迟会变小，但增加了网络带宽的消耗。**适用于主从之间的网络环境良好的场景**，如同机架或同机房部署。

- 当开启时，主节点会**合并较小的TCP数据包**从而节省带宽。默认发送时间间隔取决于Linux的内核，一般默认为40毫秒。这种配置节省了带宽但增大主从之间的延迟。适用于主从网络环境复杂或带宽紧张的场景，如跨机房部署。


### Redis主从复制原理

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/9c23120f28e74d54bed1bc455d78e471.png)

在从节点执行slaveof命令后，复制过程便开始运作。

#### 1）保存主节点信息

执行slaveof后从节点**只保存主节点的地址信息**便直接返回，这时建立复制流程还没有开始。

#### 2）建立主从socket连接

从节点(slave)内部通过**每秒运行**的定时任务维护复制相关逻辑，当定时任务发现存在新的主节点后，会尝试与该节点建立网络连接。从节点会建立一个socket套接字，专门用于接受主节点发送的复制命令。从节点连接成功后打印日志。

如果从节点无法建立连接，定时任务会**无限重试**直到连接成功或者执行`slaveof no one`取消复制。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/9a4e0f3db914446399ed2be204217411.png)

#### 3）发送ping命令

连接建立成功后从节点发送ping请求进行首次通信，ping请求主要目的：

- 检测主从之间网络套接字是否可用
- 检测主节点当前是否可接受处理命令

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/d15620425e52488eb3f6143565229dcf.png)

从节点发送的ping命令成功返回，Redis打印日志，并继续后续复制流程。

#### 4）权限验证

如果主节点设置了`requirepass`参数，则需要密码验证，从节点必须配置`masterauth`参数保证与主节点**相同的密码**才能通过验证；如果验证失败复制将终止，从节点重新发起复制流程。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/46f6e8b768b14506b81d424d61f956d2.png)

#### 5) 同步数据集

主从复制连接正常通信后，对于**首次建立**复制的场景，主节点会把持有的数据**全部**发送给从节点，这部分操作是耗时最长的步骤。Redis在2.8版本以后采用新复制命令 psync进行数据同步，原来的sync命令依然支持，保证新旧版本的兼容性。新版同步划分两种情况：**全量同步（第一次）和部分同步。**

> **全量同步：**
>
> -  主服务器会在一个 RDB 文件中保存当前数据集的快照，然后**将这个 RDB 文件发送给从服务器**。从节点接收到 RDB 文件后，会加载这个文件，将自己的数据集替换成主服务器的数据集。 
> -   在 RDB 文件传输的过程中，主服务器会将在传输期间的写操作记录下来，称为**命令传播**（command propagation）。这样一来，主服务器就能够在发送完 RDB 文件后，将这期间的写操作重新发送给从服务器，以保证从服务器的数据集与主服务器保持一致。
>
> 第一次全量同步完成后，后续的同步操作主要是通过 **部分同步** （见后续）和 **实时增量数据同步**。（某种情况也可以理解为同一个东西）
>
>  **增量数据同步：**在完成全量复制后，主从服务器之间会保持一个 TCP 连接，主服务器会将自己的**写操作**发送给从服务器，主节点的每个写操作（如 SET、DEL 等）会记录在 **复制积压缓冲区（Replication Backlog）** 中，从服务器执行这些写操作，从而很好的**保持了主从数据一致性**。增量复制的数据同步是**异步**的，但通过记录写操作，主从服务器之间的数据最终会达到一致状态。

#### 6) 命令持续复制

当主节点把当前的数据同步给从节点后，便完成了复制的建立流程。接下来主节点会持续地把**写命令**发送给从节点，保证主从数据一致性。

### Redis数据同步

Redis早期支持的复制功能只有全量复制（sync命令），它会把主节点全部数据一次性发送给从节点，当数据量较大时，会对主从节点和网络造成很大的开销。

Redis在2.8版本以后采用新复制命令psync进行数据同步，原来的sync命令依然支持，保证新旧版本的兼容性。新版同步划分两种情况:全量复制和部分复制。

#### 全量同步

全量复制：一般用于**初次复制**场景，Redis早期支持的复制功能只有全量复制，它会把主节点全部数据一次性发送给从节点，当数据量较大时，会对主从节点和网络造成很大的开销。

全量复制是Redis最早支持的复制方式，也是主从**第一次建立复制**时必须经历的阶段。触发全量复制的命令是sync和psync。

psync全量复制流程，它与2.8以前的sync全量复制机制基本一致。

**`PSYNC` 命令的基本格式：**

```bash
PSYNC <replicationid> <offset>
```

**参数说明**：

1. **`replicationid`**：
   - 主节点的唯一标识符。
   - 当从节点**首次连接**或不确定主节点的状态时，传入 `?`。
   - 如果从节点之前已与主节点同步过，则使用记录的主节点 ID。
2. **`offset`**：
   - 从节点的复制偏移量，用于标识从节点已经同步的最后一个数据位置。
   - 如果是首次同步或不确定，传入 `-1`。

**示例：**

**初次同步：**

- 从节点发送：`PSYNC ? -1`

- 主节点返回：`+FULLRESYNC c82e4f47b3b9ed0d9e9a8345c6d57b6f 0`，表示主节点通知从节点需要进行**全量同步**，并提供新的 `replicationid` 和偏移量。

  > 如果返回**`-ERR`**，表示主节点拒绝请求，可能由于版本不兼容或其他错误。

然后发送 RDB 文件和增量数据。

**部分同步：**

- 从节点发送：`PSYNC c82e4f47b3b9ed0d9e9a8345c6d57b6f 12345`
- 主节点返回：`+CONTINUE`，表示主节点通知从节点可以进行部分同步，直接发送增量数据。

然后发送偏移量之后的增量数据。

##### **流程说明**

> 下面的客户端缓冲区指的是replication buffer。也称复制缓冲区。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/23ebbf1f77714febb3a5fbcaf9b1ca11.png)

> 注意上图中bgrewriteaof是手动重写，会fork出一个子进程。
>
> 复制积压缓冲区主要用于主从复制的场景，确保从节点可以在重新连接时获取到主节点的最新数据。而客户端缓冲区则是用于管理与每个客户端之间的数据传输。这俩压根不是一个概念。
>
> **1.区分replication_buffer（复制缓冲区） 和 repl_backlog_buffer（复制积压缓冲区）**
>
> - **replication_buffer：**对于客户端或从库与redis通信，redis都会分配一个内存buffer进行数据交互，redis先把数据写入这个buffer中，然后再把buffer中的数据发送出去，所以主从在增量同步时，保证主从数据一致。（既然是内存buffer，那如果持续增大buffer的大小，会消耗大量的资源，所以redis提供了断开这个client的连接，“client-output-buffer-limit”,可以设置限制，当从库处理慢导致主库内存buffer到达限制后，主库会强制断开从库的连接。配置请看后续）
>
> - **repl_backlog_buffer:** 为了解决从库断连后找不到主从差异数据而设立的**环形缓冲区**，从而避免全量同步带来的性能开销。在redis.conf配置文件中可以设置大小，如果从库断开时间过长，repl_backlog_buffer环形缓冲区会被主库的写命令覆盖，那么从库重连后只能全量同步，所以repl_backlog_size配置尽量大一点可以降低从库连接后全量同步的频率。
>
> **2.replication_buffer的被首次使用**
>
> 当且仅当slave与master首次或者出于某种原因，需要全量rdb传输数据后，然后会把replication_buffer中的数据，再次全量传给slave。
>
> 注：此阶段称作主从复制的第一阶段，全量rdb + replication_buffer。
>
> 第二阶段(命令传播)，主要是增量传输，此时replication_backlog_buffer出场。
>
> **总结：**首次主从同步的时候， repl_backlog_buffer尚用不到，首次全量rdb同步后，如果有新的未同步数据产生，这些数据会被写入replication buffer，之后会传给slave。上述步骤完成后，就进入了后期小批量主从增量数据同步阶段，此阶段master产生的新数据会通过repl_backlog_buffer增量传输给slave。
> 简而言之：
> 第一阶段：主从rdb同步
> 第二阶段：master的replication buffer数据传输给slave
> 第三阶段：常规数据增量同步阶段，master的repl_backlog_buffer数据传输给slave (可能会被频繁操作)
> 第四阶段：如果无法通过增量同步时，主从会重复阶段一、二步骤（实在无法进行，则slave会被抛弃，直至恢复主从关系）。

1）发送psync命令进行数据同步，由于是第一次进行复制，从节点没有复制偏移量和主节点的运行ID，所以发送`psync ? -1`。

2）主节点根据psync ? -1解析出当前为全量复制，回复 +FULLRESYNC响应，从节点接收主节点的响应数据保存**运行ID和偏移量offset**，并打印日志。

3）主节点执行bgsave保存RDB文件到本地。

4）主节点发送RDB文件给从节点，从节点把接收的RDB文件保存在本地，并直接作为从节点的数据文件，接收完RDB后，从节点打印相关日志，可以在日志中查看主节点发送的数据量。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/17d37667b1854aeb92371a1ee5a4e600.png)

5）对于从节点开始接收RDB快照到接收完成期间，主节点仍然响应读写命令，因此主节点会把这期间**写命令**数据保存在**客户端缓冲区**内（即replication buffer），当**从节点加载完RDB文件后**，主节点再把缓冲区内的数据发送给从节点，保证主从之间数据一致性。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/a8315cedc2d54bf2a624c0cbb5296bc4.png)

需要注意，对于数据量较大的主节点，比如生成的RDB文件超过6GB以上时要格外小心。传输RDB文件这一步操作非常耗时，速度取决于主从节点之间网络带宽。

> 网络带宽是指在单位时间（一般指的是1秒钟）内能传输的数据量。网络和高速公路类似，带宽越大，就类似高速公路的车道越多，其通行能力越强。带宽的单位有bps，Kbps，Mbps，Gbps，Tbps等。

##### 问题

通过分析全量复制的所有流程，会发现全量复制是一个非常耗时费力的操作。它的时间开销主要包括:

1、主节点bgsave时间。

2、RDB文件网络传输时间。

3、从节点清空旧数据时间。

4、从节点加载RDB的时间。

5、可能的AOF重写时间。

因此当数据量达到一定规模之后，由于全量复制过程中将进行多次持久化相关操作和网络数据传输，这期间会大量消耗主从节点所在服务器的CPU、内存和网络资源。

**另外最大的问题，复制还是有可能会失败！！！**

例如我们主节点的线上数据量在6G左右，从节点发起全量复制的总耗时在2分钟左右。

1、如果复制的总时间超过**repl-timeout**所配置的值（默认60秒)，从节点将放弃接受RDB文件并清理已经下载的临时文件，导致全量复制失败。

> repl-timeout在Redis 配置文件（通常是 `redis.conf`）里面。`repl-timeout` 是一个**从节点端配置项**，如果要修改，需修改从节点的 `redis.conf` 文件。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665215626018/08db81e074244dd6b7a7e27c54625f90.png)

2、如果主节点创建和传输RDB的时间过长，对于高流量写入场景非常容易造成主节点**客户端缓冲区**（replication buffer）溢出。默认配置为下面。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665215626018/99b4f1f238014df0ad4c29e3fac5a97a.png)

意思是如果**60秒内缓冲区消耗持续大于64MB**或者直接超过256MB时，主节点将直接关闭复制客户端连接，造成全量同步失败。

所以除了第一次复制时采用全量复制在所难免之外，对于其他场景应该规避全量复制的发生。正因为全量复制的成本问题。

#### 部分同步

部分复制主要是Redis针对全量复制的过高开销做出的一种优化措施。

使用`psync {runId} {offset}`  命令实现

当从节点(slave)正在复制主节点(master)时，如果出现**网络闪断**或者**命令丢失**等异常情况时，从节点会向主节点要求补发**丢失的命令数据**，如果主节点的**复制积压缓冲区**内存在这部分数据，则直接发送给从节点，这样就可以保持主从节点复制的一致性。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/d284a4e590e5419cb092771ea302b10f.png)

**千万注意！！！从上下文可以看出，在部分复制中，主节点传给从节点的是命令，而不是具体的数据。**

##### 流程说明

1)当主从节点之间**网络出现中断**时，如果超过repl-timeout时间，主节点会认为从节点故障并中断复制连接，打印日志。如果此时从节点没有宕机，也会打印与主节点连接丢失日志。

2）主从连接中断期间，主节点依然响应命令，但因复制连接中断命令无法发送给从节点，不过主节点内部存在的**复制积压缓冲区**，依然可以保存最近一段时间的写命令数据，**默认最大缓存1MB。**

3)当主从节点网络恢复后，从节点会再次连上主节点，打印日志。

4）当主从连接恢复后，由于从节点之前保存了自身已复制的**偏移量**和主节点的**运行ID**。因此会把它们当作psync参数发送给主节点，要求进行**部分复制**操作。

5)主节点接到psync命令后，首先核对参数runId（运行id）是否与自身一致，如果一致，说明之前复制的是当前主节点；之后根据参数offset在自身**复制积压缓冲区**查找，如果偏移量之后的数据存在缓冲区中，则对从节点发送**+CONTINUE**响应，表示可以进行部分复制。如果不在，则**退化为全量复制。**

6）主节点根据偏移量把复制积压缓冲区里的数据发送给从节点，保证主从复制进入正常状态。发送的数据量可以在主节点的日志，传递的数据远远小于全量数据。

> **全量同步和部分同步触发场景**
>
> **全量同步（Full Resynchronization）**
>
> 全量同步是指从节点完全清空自己的数据，并重新从主节点接收完整的数据快照（RDB 文件）和增量数据进行同步。这是一个高成本操作。
>
> **触发场景**：
>
> 1. **首次连接**
>     从节点首次与主节点建立复制关系时，没有任何数据需要同步所有内容。
>    - 从节点发送：`PSYNC ? -1`
> 2. **主节点 ID 变化**
>     主节点重启或数据被清空时，其 `replication ID` 会发生变化，从节点的旧 `replication ID` 不再有效，无法进行部分同步。
> 3. **复制偏移量超出范围**
>     如果主节点的复制积压缓冲区（replication backlog buffer）中不再包含从节点需要的增量数据，则必须执行全量同步。
> 4. **手动触发**
>     当人为清空从节点数据或重启从节点时，也会触发全量同步。
>
> **部分同步（Partial Resynchronization）**
>
> 部分同步是指**主从节点只传输断开连接期间的增量数据**，而无需传输完整的数据快照。这大大减少了网络带宽的使用和同步时间。
>
> **触发场景**：
>
> 1. **短暂的断开连接后重新连接**
>     如果从节点的 `replication ID` 与主节点一致，并且复制偏移量在主节点的复制积压缓冲区范围内，则可以进行部分同步。

#### 心跳

主从节点在建立复制后，它们之间维护着**长连接**并彼此发送心跳命令。

**主从心跳判断机制：**

1. 主从节点彼此都有心跳检测机制，各自模拟成对方的客户端进行通信，通过**client list命令**查看复制相关客户端信息，主节点的连接状态为flags=M，从节点连接状态为flags=S。

2. 主节点默认每隔10秒对从节点发送ping命令，判断从节点的存活性和连接状态。可通过参数**repl-ping-slave-period**控制发送频率。

3. 从节点在主线程中每隔1秒发送**replconf ack {offset}**命令，给主节点上报自身当前的复制偏移量。replconf命令主要作用如下:

   - 实时监测主从节点网络状态；


   - 上报自身复制偏移量，检查复制数据是否丢失，如果从节点数据丢失，再从主节点的复制缓冲区中拉取丢失数据；


   - 实现保证从节点的数量和延迟性功能，通过min-slaves-to-write、min-slaves-max-lag参数配置定义；


主节点根据replconf命令判断**从节点超时时间**，体现在**info replication**统计中的**lag**信息中，lag表示与从节点最后一次通信延迟的秒数，正常延迟应该在0和1之间。如果超过repl-timeout配置的值(（默认60秒)，则判定从节点下线并断开复制客户端连接。主节点判定从节点下线后，如果从节点重新恢复，心跳检测会继续进行。

#### 异步复制机制

主节点不但负责数据读写，还负责把写命令同步给从节点。**写命令的发送过程是异步完成**。也就是说主节点自身处理完写命令后直接发送给客户端，并不需要等待从节点复制完成。

由于主从复制过程是异步的，就会造成从节点的数据相对主节点存在延迟。具体延迟多少字节，我们可以在主节点执行**info replication**命令查看相关指标获得。

在统计信息中可以看到从节点slave信息，分别记录了从节点的ip和 port，从节点的状态，offset表示当前从节点的复制偏移量，master_repl_offset表示当前主节点的复制偏移量，**两者的差值就是当前从节点复制延迟量**。Redis 的复制速度取决于主从之间网络环境，repl-disable-tcp-nodelay，命令处理速度等。正常情况下，延迟在1秒以内。

## 集群模式

Redis Cluster是Redis的分布式解决方案，在3.0版本正式推出，有效地解决了Redis分布式方面的需求。当遇到单机内存、并发、流量等瓶颈时，可以采用Cluster架构方案达到负载均衡的目的。之前，Redis分布式方案一般有两种：

1、客户端分区方案，优点是分区逻辑可控，缺点是需要自己处理数据路由、高可用、故障转移等问题。

2、代理方案，优点是简化客户端分布式逻辑和升级维护便利,缺点是加重架构部署复杂度和性能损耗。

现在官方为我们提供了专有的集群方案：Redis Cluster，它非常优雅地解决了Redis集群方面的问题，因此理解应用好Redis Cluster将极大地解放我们使用分布式Redis 的工作量。

### 集群前置知识

#### 数据分布理论

分布式数据库首先要解决把整个数据集按照分区规则映射到多个节点的问题，即**把数据集划分到多个节点上**，每个节点负责整体数据的一个子集。

需要重点关注的是数据分区规则。常见的分区规则有**哈希分区和顺序分区**两种，哈希分区离散度好、数据分布业务无关、无法顺序访问，顺序分区离散度易倾斜、数据分布业务相关、可顺序访问。

##### 节点取余分区

使用特定的数据，如Redis的键或用户ID，再根据节点数量N使用公式：
hash(key)%N计算出哈希值，用来决定数据映射到哪一个节点上。这种方案存在一个问题：当节点数量变化时，如扩容或收缩节点，数据节点映射关系需要**重新计算**，会导致数据的重新迁移。

这种方式的突出优点是简单性，常用于数据库的分库分表规则，一般采用预分区的方式，提前根据数据量规划好分区数，比如划分为512或1024张表，保证可支撑未来一段时间的数据量，再根据负载情况将表迁移到其他数据库中。扩容时通常采用翻倍扩容，避免数据映射全部被打乱导致全量迁移的情况，如图10-2所示。

##### 一致性哈希分区

一致性哈希分区（Distributed Hash Table）的实现思路是为系统中每个节点分配一个token，范围一般在0~23，这些token构成一个**哈希环。**数据读写执行节点查找操作时，先根据key计算hash值，然后**顺时针**找到第一个大于等于该哈希值的token节点。例如：

集群中有三个节点（Node1、Node2、Node3），五个键（key1、key2、key3、key4、key5），其路由规则为：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666925611079/bab557dc5d9c4932bc2682481a50e8b7.png)

当集群中增加节点时，比如当在Node2和Node3之间增加了一个节点Node4，此时再访问节点key4时，不能在Node4中命中，更一般的，介于Node2和Node4之间的key均失效，这样的失效方式太过于“集中”和“暴力”，更好的方式应该是“平滑”和“分散”地失效。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666925611079/9de42d7a40a94ca5ba61daff1d0cd894.png)

这种方式相比节点取余最大的好处在于加入和删除节点只影响哈希环中相邻的节点，对其他节点无影响。但一致性哈希分区存在几个问题：

1、当使用少量节点时，节点变化将大范围影响哈希环中数据映射，因此这种方式不适合少量数据节点的分布式方案。

2、增加节点只能对下一个相邻节点有比较好的负载分担效果，例如上图中增加了节点Node4只能够对Node3分担部分负载，对集群中其他的节点基本没有起到负载分担的效果；类似地，删除节点会导致下一个相邻节点负载增加，而其他节点却不能有效分担负载压力。

正因为一致性哈希分区的这些缺点，一些分布式系统采用**虚拟槽**对一致性哈希进行改进，比如**虚拟一致性哈希分区。**

##### 虚拟一致性哈希分区

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666925611079/70c7b32ede304da4bf8498c39e836bc1.png)

为了在**增删节点**的时候，各节点能够保持动态的均衡，将每个真实节点虚拟出若干个**虚拟节点**，再将这些虚拟节点随机映射到环上。此时每个真实节点不再映射到环上，真实节点只是用来存储键值对，它负责接应各自的一组环上虚拟节点。当对键值对进行存取路由时，首先路由到虚拟节点上，再由虚拟节点找到真实的节点。

如下图所示，三个节点真实节点：Node1、Node2和Node3，每个真实节点虚拟出三个虚拟节点：X#V1、X#V2、X#V3，这样每个真实节点所负责的hash空间不再是连续的一段，而是分散在环上的各处，这样就可以将局部的压力均衡到不同的节点，虚拟节点越多，分散性越好，理论上负载就越倾向均匀。

##### 虚拟槽分区

> 一个 **Redis 集群** 是由多个 Redis 实例组成的，集群中的每个 Redis 实例称为一个节点，一个节点可以是 **主节点（Master）** 或 **从节点（Replica）**。Redis 集群总共固定有 **16384 个哈希槽**（编号从 `0` 到 `16383`）。数据通过哈希槽（Hash Slot）分布在多个主节点上（每个主节点负责一定范围的哈希槽）。注意，从节点复制对应主节点的数据（但是不参与槽相关操作）。

Redis则是利用了**虚拟槽分区，**可以算上面虚拟一致性哈希分区的变种，它使用**分散度良好的哈希函数**把所有数据映射到一个固定范围的整数集合中，**整数定义为槽(slot)**。这个范围一般**远远大于节点数**，比如Redis Cluster槽范围是0 ～16383。**槽是集群内数据管理和迁移的基本单位。**采用大范围槽的主要目的是为了方便数据拆分和集群扩展。每个节点会负责一定数量的槽。

比如集群有3个节点，则每个节点平均大约负责5460个槽。由于采用高质量的哈希算法，每个槽所映射的数据通常**比较均匀**，将数据平均划分到5个节点进行数据分区。Redis Cluster就是采用虚拟槽分区，下面就介绍Redis数据分区方法。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666925611079/cbd9c96c04f84d9795c24c5f76cf5234.png)

##### 为什么槽的范围是0 ～16383？

为什么槽的范围是0 ～16383，也就是说槽的个数在16384个？redis的作者在github上有个回答：[https://github.com/redis/redis/issues/2576](https://github.com/redis/redis/issues/2576)

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666925611079/1d9939c506b34a9aa5cd004fb83c34c7.png)

**这个意思是：**

Redis集群中，在握手成功后，两个节点之间会定期发送ping/pong消息，交换数据信息，集群中节点数量越多，消息体内容越大，比如说10个节点的状态信息约1kb，同时redis集群内节点，每秒都在发ping消息。例如，一个总节点数为200的Redis集群，默认情况下，这时ping/pong消息占用带宽达到25M。

那么如果槽位为65536，发送心跳信息的消息头达8k，发送的心跳包过于庞大，非常浪费带宽。

其次redis的集群主节点数量基本不可能超过1000个。集群节点越多，心跳包的消息体内携带的数据越多。如果节点过1000个，也会导致网络拥堵。因此redis作者，不建议redis cluster节点数量超过1000个。

那么，对于节点数在1000以内的redis cluster集群，16384个槽位够用了，可以以确保每个 master 有足够的插槽，没有必要拓展到65536个。

再者Redis主节点的配置信息中，它所负责的哈希槽是通过一张**bitmap**的形式来保存的，在传输过程中，会对bitmap进行压缩，但是如果bitmap的填充率slots / N很高的话(N表示节点数)，也就是节点数很少，而哈希槽数量很多的话，bitmap的压缩率就很低，也会浪费资源。

所以Redis作者决定取16384个槽，作为一个比较好的设计权衡。

总而言之，实践出真知。

#### Redis数据分区

Redis Cluser采用虚拟槽分区，**所有的键根据哈希函数映射到0 ~16383整数槽内**，计算公式：`slot = CRC16(key) & 16383`。每一个节点负责维护―部分槽以及槽所映射的**键值数据。**

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666925611079/d8360eb94ffc458294f0f003a7f94380.png)![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666925611079/a3a1240b7e74478abba33cd4a6e5821d.png)

##### Redis 虚拟槽分区的特点

1、解耦数据和节点之间的关系，简化了节点扩容和收缩难度。

2、节点自身维护槽的映射关系，不需要客户端或者代理服务维护槽分区元数据。可支持节点、槽、键之间的映射查询，用于数据路由、在线伸缩等场景。

3、**数据分区是分布式存储的核心**，理解和灵活运用数据分区规则对于掌握Redis Cluster非常有帮助。

##### 集群功能限制

Redis集群相对单机在功能上存在一些限制，需要开发人员提前了解，在使用时做好规避。限制如下：

1、 key批量操作支持有限。如mset、mget，目前只支持**具有相同slot值**的key执行批量操作。对于映射为不同slot值的key由于执行mget、mget等操作可能存在于多个节点上因此不被支持。

2、key事务操作支持有限。同理只支持多**key在同一节点上**的事务操作，当多个key分布在不同的节点上时无法使用事务功能。

3、key作为数据分区的最小粒度，因此**不能将一个大的键值对象如hash、list等映射到不同的节点。**

4、不支持多数据库空间。单机下的Redis可以支持16个数据库空间，集群模式下**只能使用一个数据库空间**，即 **db 0。**

5、复制结构只支持一层，**从节点只能复制主节点**，不支持嵌套树状复制结构。

### 搭建集群

介绍完Redis集群分区规则之后，下面我们开始搭建Redis集群。搭建集群有几种方式：

1）依照Redis协议**手工搭建**，使用cluster meet、cluster addslots、cluster replicate命令。

2）5.0之前使用由ruby语言编写的redis-trib.rb，在使用前需要安装ruby语言环境。

3）5.0及其之后redis摒弃了redis-trib.rb，将搭建集群的功能合并到了**redis-cli。**

我们简单点，采用第三种方式搭建。集群中至少应该有**奇数个节点**，所以至少有三个节点，官方推荐**三主三从**的配置方式，我们就来搭建一个三主三从的集群。

#### 节点配置

我们现在规定，主节点的端口为6900、6901、6902，从节点的端口为6930、6931、6932。

首先需要配置节点的conf文件，这个比较统一，所有的节点的配置文件都是类似的，我们以端口为6900的节点举例：

```bash
# 指定该Redis实例监听的端口为 6900
port 6900

# 这个部分是为了在一台服务上启动多台Redis服务，相关的资源要改
pidfile /var/run/redis_6900.pid
logfile "/home/chenpeng/redis/redis/log/6900.log"
# 数据存储目录，数据文件（如 RDB 文件、AOF 文件）存储的目录。多个实例可以共用一个目录，但文件名需要区分
dir "/home/chenpeng/redis/redis/data/" 
dbfilename dump-6900.rdb

# Cluster Config
# 设置为yes表示以守护进程（后台）方式运行
daemonize yes
# 设置为yes表示启用Redis集群模式
cluster-enabled yes
# 指定集群配置文件，文件中记录了节点的角色、槽位分布等信息。每个实例需要独立的配置文件，如nodes-6900.conf。
cluster-config-file nodes-6900.conf
# 集群节点之间的通信超时时间，单位为毫秒。
cluster-node-timeout 15000
# 开启AOF持久化。
appendonly yes
appendfilename "appendonly-6900.aof"
```

在上述配置中，以下配置是集群相关的：

```bash
cluster-enabled yes # 是否启动集群模式(集群需要修改为yes)
cluster-node-timeout 15000  #指定集群节点超时时间(打开注释即可)
 # 指定集群节点的配置文件(打开注释即可)。这个文件不需要手工编辑，它由Redis节点创建和更新。每个Redis集群节点都需要一个集群配置文件，确保在同一系统中运行的实例没有重叠集群配置文件名。
cluster-config-file nodes-6900.conf
appendonly yes  # 指定redis集群持久化方式(默认是rdb，但建议使用aof方式,此处是否修改不影响集群的搭建)
```

#### 集群创建

##### 创建集群随机主从节点

```bash
./redis-cli --cluster create 127.0.0.1:6900 127.0.0.1:6901 127.0.0.1:6902 127.0.0.1:6930 127.0.0.1:6931
127.0.0.1:6932 --cluster-replicas 1
```

说明：--cluster-replicas 参数为数字，**1表示每个主节点需要1个从节点。**

通过该方式创建的带有从节点的机器不能够自己手动指定主节点，**不符合我们的要求。**所以如果需要指定的话，需要自己手动指定，先创建好主节点后，再添加从节点。

##### 指定主从节点

###### 创建集群主节点

```bash
./redis-cli --cluster create 127.0.0.1:6900 127.0.0.1:6901 127.0.0.1:6902
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666925611079/4bff1e9d3ab84d4ab09ad629946dff27.png)

**注意：**

1. 请记录下每个M后形如“dcd818ab48166ccea9563544839187ffa5d79f62”的字符串，在后面添加从节点时有用；
2. 如果服务器存在着防火墙，那么在进行安全设置的时候，除了redis服务器本身的端口，比如6900要加入允许列表之外，Redis服务在集群中还有一个叫**集群总线端口**，其端口为客户端连接端口加上10000，即 6900 + 10000 = 16900 。所以需开放每个集群节点的客户端端口和集群总线端口才能成功创建集群！

###### 添加集群从节点

命令类似：

```bash
./redis-cli --cluster add-node 127.0.0.1:6930 127.0.0.1:6900 --cluster-slave --cluster-master-id dcd818ab48166ccea9563544839187ffa5d79f62
```

说明：上述命令把6930节点作为从节点加入到6900节点的集群中，并且当做node_id为dcd818ab48166ccea9563544839187ffa5d79f62的从节点。如果不指定 --cluster-master-id 会随机分配到任意一个主节点

效果如下：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666925611079/da6debed19fa42ce956ad2fe918b19ab.png)

第二、三个从节点的配置类似。

#### 集群管理

##### 检查集群

```bash
./redis-cli --cluster check 127.0.0.1:6900 --cluster-search-multiple-owners
```

说明：**任意**连接一个集群节点，进行集群状态检查。

##### 集群信息查看

```bash
./redis-cli --cluster info 127.0.0.1:6900
```

说明：检查key、slots、从节点个数的分配情况

##### 修复集群

```bash
redis-cli --cluster fix 127.0.0.1:6900 --cluster-search-multiple-owners
```

说明：修复集群，并且修复**槽的重复分配**问题。

##### 设置集群的超时时间

```bash
redis-cli --cluster set-timeout 127.0.0.1:6900 10000
```

说明：连接到集群的任意一个节点来设置集群的超时时间参数cluster-node-timeout。表示集群中**节点间通信**（心跳检测）的超时时间，单位为毫秒。

##### 集群配置

> `--cluster call`表示对整个集群的所有节点执行命令。

```bash
# 为Redis集群中的所有节点设置访问密码和主从节点的认证密码，并将更改写入持久化配置文件
redis-cli --cluster call 127.0.0.1:6900 config set requirepass cc
redis-cli --cluster call 127.0.0.1:6900 config set masterauth cc
# 使用config rewrite将上述配置写入到每个节点的持久化配置文件
redis-cli --cluster call 127.0.0.1:6900 config rewrite
```

说明：连接到集群的任意一节点来**对整个集群的所有节点**进行设置。

#### redis-cli –cluster参数参考

```bash
redis-cli --cluster help
Cluster Manager Commands:
  create         host1:port1 ... hostN:portN   #创建集群
                 --cluster-replicas <arg>      #从节点个数
  check          host:port                     #检查集群
                 --cluster-search-multiple-owners #检查是否有槽同时被分配给了多个节点
  info           host:port                     #查看集群状态
  fix            host:port                     #修复集群
                 --cluster-search-multiple-owners #修复槽的重复分配问题
  reshard        host:port                     #指定集群的任意一节点进行迁移slot，重新分slots。它允许将部分哈希槽从一个主节点转移到另一个主节点，从而实现负载均衡或调整集群结构。
                 --cluster-from <arg>          #需要从哪些源节点上迁移slot，可从多个源节点完成迁移，以逗号隔开，传递的是节点的node id，还可以直接传递--from all，这样源节点就是集群的所有节点，不传递该参数的话，则会在迁移过程中提示用户输入
                 --cluster-to <arg>            #slot需要迁移的目的节点的node id，目的节点只能填写一个，不传递该参数的话，则会在迁移过程中提示用户输入
                 --cluster-slots <arg>         #需要迁移的slot数量，不传递该参数的话，则会在迁移过程中提示用户输入。
                 --cluster-yes                 #指定迁移时的确认输入
                 --cluster-timeout <arg>       #设置migrate命令的超时时间
                 --cluster-pipeline <arg>      #定义cluster getkeysinslot命令一次取出的key数量，不传的话使用默认值为10。即迁移key时，一次取出的key数量
                 --cluster-replace             #是否直接replace到目标节点，即指定在迁移哈希槽时，是否替换目标节点已经存在的槽（谨慎使用）
  rebalance      host:port                     #指定集群的任意一节点，进行平衡集群节点slot数量，即重新平衡集群中的哈希槽分布，使每个主节点承载的哈希槽尽可能均匀
                 --cluster-weight <node1=w1...nodeN=wN>         #指定集群节点的权重
                 --cluster-use-empty-masters                    #设置可以让没有分配slot的主节点参与，默认不允许。默认情况下，rebalance 会忽略空主节点（没有分配哈希槽的主节点）。使用此选项后，空主节点也会被纳入重新分配范围。
                 --cluster-timeout <arg>                        #设置migrate命令的超时时间
                 --cluster-simulate                             #模拟rebalance操作，不会真正执行迁移操作
                 --cluster-pipeline <arg>                       #定义cluster getkeysinslot命令一次取出的key数量，默认值为10
                 --cluster-threshold <arg>                      #迁移的slot阈值超过threshold，执行rebalance操作，此选项设置一个不平衡阈值（百分比值）。当集群中主节点之间的哈希槽分布的不均衡程度超过该阈值时，rebalance命令会触发哈希槽的重新分配操作
                 --cluster-replace                              #是否直接replace到目标节点（指覆盖已有槽）
  add-node       new_host:new_port existing_host:existing_port  #添加节点，把新节点加入到指定的集群，默认添加主节点
                 --cluster-slave                                #新节点作为从节点，默认随机一个主节点
                 --cluster-master-id <arg>                      #给新节点指定主节点
  del-node       host:port node_id                              #删除给定的一个节点，成功后关闭该节点服务
  call           host:port command arg arg .. arg               #在集群的所有节点执行相关命令
  set-timeout    host:port milliseconds                         #设置cluster-node-timeout
  import         host:port                                      #将外部redis（单实例Redis）数据导入集群
                 --cluster-from <arg>                           #将指定实例的数据导入到集群
                 --cluster-copy                                 #migrate时指定copy。将源单实例Redis的数据拷贝到 Redis集群中，而不删除源实例的数据
                 --cluster-replace                              #migrate时指定replace。使用此选项，源实例的数据会覆盖目标集群中与源实例键冲突的数据。
```

以下是 Redis `redis-cli --cluster` 命令及其选项的表格化说明：

| **命令**        | **参数**                                 | **说明**                                                     |
| :-------------- | :--------------------------------------- | ------------------------------------------------------------ |
| **create**      | `host1:port1 ... hostN:portN`            | 创建集群，指定节点列表。                                     |
|                 | `--cluster-replicas <arg>`               | 指定从节点个数。                                             |
| **check**       | `host:port`                              | 检查集群。                                                   |
|                 | `--cluster-search-multiple-owners`       | 检查是否有哈希槽被多个节点重复分配。                         |
| **info**        | `host:port`                              | 查看集群状态。                                               |
| **fix**         | `host:port`                              | 修复集群。                                                   |
|                 | `--cluster-search-multiple-owners`       | 修复哈希槽的重复分配问题。                                   |
| **reshard**     | `host:port`                              | 重新分配哈希槽，可用于负载均衡或调整集群结构。               |
|                 | `--cluster-from <arg>`                   | 指定迁移源节点（node ID），可用逗号分隔多个节点，或使用 `--from all` 表示集群所有节点。 |
|                 | `--cluster-to <arg>`                     | 指定迁移目标节点（node ID）。                                |
|                 | `--cluster-slots <arg>`                  | 指定迁移的哈希槽数量。                                       |
|                 | `--cluster-yes`                          | 自动确认迁移操作。                                           |
|                 | `--cluster-timeout <arg>`                | 设置迁移操作的超时时间。                                     |
|                 | `--cluster-pipeline <arg>`               | 定义一次迁移操作中获取的键数量，默认值为 10。                |
|                 | `--cluster-replace`                      | 指定是否覆盖目标节点已有哈希槽。                             |
| **rebalance**   | `host:port`                              | 平衡集群中主节点的哈希槽分布，使其尽可能均匀。               |
|                 | `--cluster-weight <node1=w1...nodeN=wN>` | 指定节点权重，权重高的节点会分配更多哈希槽。                 |
|                 | `--cluster-use-empty-masters`            | 允许未分配哈希槽的主节点参与重新分配（默认忽略空主节点）。   |
|                 | `--cluster-timeout <arg>`                | 设置迁移操作的超时时间。                                     |
|                 | `--cluster-simulate`                     | 模拟重新平衡操作，不真正执行迁移。                           |
|                 | `--cluster-pipeline <arg>`               | 定义一次迁移操作中获取的键数量，默认值为 10。                |
|                 | `--cluster-threshold <arg>`              | 设置不平衡阈值（百分比）。当主节点之间的哈希槽分布不均衡程度超过阈值时触发平衡操作。 |
|                 | `--cluster-replace`                      | 指定是否覆盖目标节点已有哈希槽。                             |
| **add-node**    | `new_host:new_port existing_host:port`   | 添加节点到集群中，默认作为主节点。                           |
|                 | `--cluster-slave`                        | 将新节点作为从节点。                                         |
|                 | `--cluster-master-id <arg>`              | 指定新节点的主节点。                                         |
| **del-node**    | `host:port node_id`                      | 删除指定节点，成功后关闭该节点服务。                         |
| **call**        | `host:port command arg arg...`           | 在集群所有节点上执行指定命令。                               |
| **set-timeout** | `host:port milliseconds`                 | 设置 `cluster-node-timeout` 参数。                           |
| **import**      | `host:port`                              | 将外部 Redis（单实例）数据导入到集群。                       |
|                 | `--cluster-from <arg>`                   | 指定单实例 Redis 的地址作为数据源。                          |
|                 | `--cluster-copy`                         | 数据拷贝模式，保留源实例数据。                               |
|                 | `--cluster-replace`                      | 数据替换模式，覆盖目标集群中已有冲突键。                     |

 **reshard 命令的all选项示例：**

> **示例场景**
>
> 假设 Redis 集群包含以下三个主节点：
>
> | Node ID                                    | Address          | 哈希槽范围    |
> | ------------------------------------------ | ---------------- | ------------- |
> | `abcd1234efgh5678ijkl9012mnop3456qrst7890` | `127.0.0.1:7000` | `0-5460`      |
> | `qrst7890mnop3456ijkl9012efgh5678abcd1234` | `127.0.0.1:7001` | `5461-10922`  |
> | `ijkl9012mnop3456qrst7890efgh5678abcd1234` | `127.0.0.1:7002` | `10923-16383` |
>
> #### 需求：
>
> 将 **300 个哈希槽** 平均从所有主节点迁移到目标节点 `127.0.0.1:7001`，而不需要手动指定具体的源节点。
>
> ------
>
> **执行命令**
>
> ```
> redis-cli --cluster reshard 127.0.0.1:7000
> ```
>
> ------
>
> **交互步骤**
>
> 1. **输入需要迁移的槽位数量**：
>
>    ```
>    How many slots do you want to move (from 1 to 16384)?
>    300
>    ```
>
>    - 表示迁移 300 个哈希槽。
>
> 2. **输入目标节点 ID**：
>
>    ```
>    What is the receiving node ID?
>    qrst7890mnop3456ijkl9012efgh5678abcd1234
>    ```
>
>    - 指定目标节点 `127.0.0.1:7001` 的节点 ID，表示将 300 个槽位迁移到这个节点。
>
> 3. **指定源节点**：
>
>    ```
>    Source node(s) from which to take slots, or all to use all nodes:
>    all
>    ```
>
>    - 输入 `all` 表示 Redis 自动选择源节点（`127.0.0.1:7000` 和 `127.0.0.1:7002`），根据当前哈希槽的分布**自动平衡负载。**
>
> 4. **确认迁移计划**：
>    Redis 会生成迁移计划，并显示如下信息：
>
>    ```
>    Moving 300 slots from all nodes to qrst7890mnop3456ijkl9012efgh5678abcd1234
>    Source nodes: abcd1234efgh5678ijkl9012mnop3456qrst7890, ijkl9012mnop3456qrst7890efgh5678abcd1234
>    Do you want to proceed with the proposed reshard plan (yes/no)?
>    yes
>    ```
>
>    - 确认迁移计划后，输入 `yes` 开始迁移。
>
> 5. **迁移完成**：
>    Redis 会逐步从 `127.0.0.1:7000` 和 `127.0.0.1:7002` 迁移槽位到目标节点 `127.0.0.1:7001`。

### 集群伸缩

Redis集群提供了灵活的节点扩容和收缩方案。在不影响集群对外服务的情况下，可以为集群添加节点进行扩容也可以下线部分节点进行缩容。Redis集群可以实现对节点的灵活上下线控制。其中原理可抽象为**槽和对应数据在不同节点之间灵活移动。**首先来看我们之前搭建的集群槽和数据与节点的对应关系。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666925611079/020a4fcfb89345f480bd047a77cd372d.png)

三个主节点分别维护自己负责的槽和对应的数据，如果希望加入1个节点实现集群扩容时，需要通过相关命令把一部分槽和数据迁移给新节点。

#### 集群扩容

##### 节点配置和启动节点

我们加入两个节点，主节点的端口为6903，从节点的端口为6933。配置与前面的6900类似，不再赘述。

启动这两个节点。

```bash
./redis-server ../conf/cluster_m_6903.conf
./redis-server ../conf/cluster_s_6933.conf
```

##### 加入集群

先执行命令查看待加入的集群信息：

```bash
./redis-cli --cluster info 127.0.0.1:6900
```

执行命令查看集群的节点信息：

`cluster nodes` 是 Redis 集群管理的一条命令，用于显示当前节点所属集群中**所有节点**的详细信息。

每个节点的信息包含：

- 节点的唯一 ID。
- 节点的地址（IP 和端口）。
- 节点角色（主节点或从节点）。
- 哈希槽范围（如果是主节点）。
- 节点状态等。

```bash
./redis-cli -p 6900 cluster nodes
```

可以得出，6903和6933还属于孤立节点，需要将这两个实例节点加入到集群中。

###### 将主节点6903加入集群

执行命令加入新的节点

```bash
./redis-cli --cluster add-node 127.0.0.1:6903 127.0.0.1:6900
```

执行命令查看集群信息

```
./redis-cli --cluster info 127.0.0.1:6900
```

执行命令查看节点信息

```
./redis-cli -p 6900 cluster nodes
```

###### 将从节点6933加入集群

执行命令

```bash
./redis-cli --cluster add-node 127.0.0.1:6933 127.0.0.1:6900 --cluster-slave --cluster-master-id 67dd0e8160a5bf8cd0ca02c2c6268bb9cc17884c
```

将刚刚加入的节点6903作为从节点6933的主节点。

##### 迁移槽和数据

上面的图中可以看到，6903和6933已正确添加到集群中，接下来就开始分配槽位。我们将6900、6901、6902三个节点中的槽位分别迁出一些槽位给6903，假设分配后的每个节点槽位平均，那么应该分出（16384/4）=4096个槽位。

执行命令

```bash
./redis-cli --cluster reshard 127.0.01:6900
```

Redis会提问要迁移的**槽位数**和**接受槽位的节点id**，我们这里输入4096和67dd0e8160a5bf8cd0ca02c2c6268bb9cc17884c（新增的主节点的nodeid）。此处可参考上面举例all的过程。

接下来，Redis会提问从哪些源节点进行迁移，我们输入“all”。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666925611079/3df8fdb055e0420a8f6048c4f0fd3c36.png)

Redis会显示一个分配（迁移）计划：填入“yes”。

Redis会开始进行迁移

**这个时间会比较长.........................**

稍等一会，等待Redis迁移完成。

迁移完成后，执行命令

```bash
./redis-cli -p 6900 cluster nodes
```

```bash
./redis-cli --cluster info 127.0.0.1:6900
```

可以看到槽位确实被迁移到了节点6903之上。这样就实现了集群的扩容。

#### 集群缩容

##### 迁移槽和数据

命令语法：

```bash
redis-cli --cluster reshard --cluster-from 要迁出节点ID --cluster-to 接收槽节点ID --cluster-slots 迁出槽数量已存在节点ip 端口
```

例如：

迁出1365个槽位到6900节点：

```bash
./redis-cli --cluster reshard --cluster-from 67dd0e8160a5bf8cd0ca02c2c6268bb9cc17884c
--cluster-to 7353cda9e84f6d85c0b6e41bb03d9c4bd2545c07 --cluster-slots 1365
127.0.0.1:6900
```

迁出1365个槽位到6901节点：

```bash
./redis-cli --cluster reshard --cluster-from 67dd0e8160a5bf8cd0ca02c2c6268bb9cc17884c
--cluster-to 41ca2d569068043a5f2544c598edd1e45a0c1f91 --cluster-slots 1365
127.0.0.1:6900
```

迁出1366个槽位到6902节点：

```bash
./redis-cli --cluster reshard --cluster-from 67dd0e8160a5bf8cd0ca02c2c6268bb9cc17884c
--cluster-to d53bb67e4c82b89a8d04d572364c07b3285e271f --cluster-slots 1366
127.0.0.1:6900
```

稍等片刻，等全部槽迁移完成后，执行命令查看

```bash
./redis-cli -p 6900 cluster nodes
```

```bash
./redis-cli --cluster info 127.0.0.1:6900
```

可以看到6903上不再存在着槽了。

##### 下线节点

执行命令格式redis-cli --cluster del-node 已存在节点ip:端口 要删除的节点ID

例如：

```bash
./redis-cli --cluster del-node 127.0.0.1:6900 67dd0e8160a5bf8cd0ca02c2c6268bb9cc17884c
```

```bash
./redis-cli --cluster del-node 127.0.0.1:6900 23c0ca7519a181f6ff61580eca014dde209f7a67
```

可以看到这两个节点确实脱离集群了，这样就完成了集群的缩容

再关闭节点即可。（？）

#### 迁移相关

##### 在线迁移slot

在线把集群的一些slot从集群原来slot节点迁移到新的节点。其实在前面扩容集群的时候我们已经看到了相关的用法

直接连接到集群的任意一节点

```bash
redis-cli --cluster reshard XXXXXXXXXXX:XXXX
```

按提示操作即可。

##### 平衡（rebalance）slot

1）平衡集群中各个节点的slot数量

```bash
redis-cli --cluster rebalance XXXXXXXXXXX:XXXX
```

2）还可以根据集群中各个节点设置的权重来平衡slot数量

```bash
./redis-cli --cluster rebalance --cluster-weight 117457eab5071954faab5e81c3170600d5192270=5
815da8448f5d5a304df0353ca10d8f9b77016b28=4
56005b9413cbf225783906307a2631109e753f8f=3 
--cluster-simulate
127.0.0.1:6900
```

### 请求路由（Dummy客户端）

目前我们已经搭建好Redis集群并且理解了通信和伸缩细节，但还没有**使用客户端去操作集群。**Redis集群对客户端通信协议做了比较大的修改，为了追求性能最大化，并没有采用代理的方式而是采用**客户端直连节点**的方式。因此对于希望从单机切换到集群环境的应用需要修改客户端代码。

#### 请求重定向

在**集群模式下**，Redis接收任何键相关命令时首**先计算键对应的槽**，再根据槽找出所对应的节点，如果节点是自身，则处理键命令；否则回复MOVED重定向错误，**通知客户端请求正确的节点**。这个过程称为**MOVED重定向。**

例如，在之前搭建的集群上执行如下命令：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666925611079/639f236673424cc4a13e454fafe70ca3.png" alt="image.png" style="zoom:80%;" />

执行set命令成功，因为键hello对应的槽正好位于6900节点负责的槽范围内，可以借助`cluster keyslot {key}`命令**返回key所对应的槽**，如下所示：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666925611079/767247687d254534b8b37d5ebee9ffbe.png" alt="image.png" style="zoom:80%;" />

再执行以下命令：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666925611079/290fcc2910b14457b32ee0d2c72dd38f.png" alt="image.png" style="zoom:80%;" />

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666925611079/87b18bc6d57345cba8c721f89637a791.png" alt="image.png" style="zoom:80%;" />

由于键对应槽是5798，不属于6900节点，则回复 `MOVED {slot} {ip:port}`格式重定向信息，重定向信息包含了键所对应的槽以及负责该槽的节点地址，根据这些信息客户端就可以向**正确的节点**发起请求。

需要我们在6901节点上成功执行之前的命令：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666925611079/c897cd9ff947402fbefed12d339c84b7.png" alt="image.png" style="zoom:80%;" />

使用redis-cli命令时，可以加入-c参数支持**自动重定向**，简化手动发起重定向操作，如下所示：

> `redis-cli -p` 命令用于Redis服务器运行在本地的主机上时，非常适合快速测试或交互式操作Redis实例。如果 Redis 服务器运行在非本地的主机上，可以结合 `-h` 参数指定主机名或 IP 地址：
>
> ```bash
> redis-cli -h <hostname> -p <port>
> ```

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666925611079/ead78dd9382c4ff78920855f22b02380.png" alt="image.png" style="zoom:80%;" />

redis-cli -c自动帮我们连接到正确的节点执行命令，这个过程是在redis-cli内部维护，实质上是client端接到MOVED信息之后**再次发起请求**，并不在Redis节点中完成请求转发。

**同节点对于不属于它的键命令只回复重定向响应，并不负责转发。**正因为集群模式下把解析发起重定向的过程放到客户端完成，所以集群客户端协议相对于单机有了很大的变化。

键命令执行步骤主要分两步：①计算槽（哈希映射）；②查找槽所对应的节点。

##### 计算槽

Redis首先需要计算键所对应的槽。根据键的有效部分使用**CRC16（哈希）函数**计算出散列值，再取对16383的余数，使每个键都可以映射到0 ~16383槽范围内。

##### 槽节点查找

Redis计算得到键对应的槽后，需要**查找槽所对应的节点。**集群内通过**消息交换**，每个节点都会知道所有节点的槽信息。

根据MOVED重定向机制，客户端可以随机连接集群内任一Redis实例来获取键所在节点，这种客户端又叫 **Dummy（傀儡）客户端**，它优点是代码实现简单，对客户端协议影响较小，只需要根据重定向信息再次发送请求即可。但是它的弊端很明显，每次执行键命令前都要到Redis上进行重定向才能找到要执行命令的节点，额外增加了IO开销，这不是Redis集群高效的使用方式。

正因为如此通常集群客户端都采用另一种实现：**Smart（智能）客户端**，我们后面再说。

#### call命令

call命令可以用来在集群的**全部节点执行相同的命令。**call命令也是需要**通过集群的一个节点地址，连上整个集群**，然后在集群的每个节点执行该命令。

```bash
./redis-cli --cluster call 47.112.44.148:6900 get name
```

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666925611079/d3f333557644483dad5abe098afbc018.png)

### Smart客户端

#### smart客户端原理

大多数开发语言的Redis客户端都采用Smart客户端支持集群协议。Smart客户端通过在内部维护**slot→node**的映射关系，**本地就可实现键到节点的查找**，从而保证IO效率的最大化，而MOVED重定向是负责协助Smart客户端更新slot→node映射。Java的Jedis就默认实现了这个功能。

#### ASK 重定向

**ASK 重定向** 是客户端与集群交互时的一种重定向机制，用于处理某些特殊情况下的请求路由问题。

> Redis 集群中的数据通过哈希槽（hash slot）进行分片管理。当某些槽的数据从一个节点迁移到另一个节点时，客户端可能会请求一个**尚未完成迁移**的数据。为了正确处理这种请求，Redis 使用了 ASK 重定向机制。

- 客户端ASK 重定向流程

  Redis集群支持在线迁移槽（slot）数据来完成水平伸缩，当slot对应的数据从源节点到目标节点**迁移过程中**，客户端需要做到**智能识别**，保证键命令可正常执行。例如当一个slot数据从源节点迁移到目标节点时，期间可能出现一部分数据在源节点，而另一部分在目标节点。

  当出现上述情况时，客户端键命令执行流程将发生变化：

  1. 客户端根据本地slots缓存（槽位映射表）发送命令到**源节点**，如果存在键对象则直接执行并返回结果给客户端。

  2. 如果键对象不存在，则**可能存在于目标节点**，这时源节点会回复**ASK重定向异常**。格式如下：(error) ASK (slot} {targetIP}:{targetPort}。用于指示客户端到新节点查询数据。

  3. 客户端从ASK重定向异常提取出目标节点信息，发送asking命令（表示这是一个临时请求，目标节点不需要检查槽位的所有权）到目标节点打开客户端连接标识。
  4. 客户端向目标节点重新发送原始请求（即执行键命令）。目标节点如果**存在则执行，不存在则返回不存在信息。**

ASK与MOVED虽然都是对客户端的重定向控制，但是有着本质区别。

- ASK重定向说明集群正在进行**slot数据迁移**，客户端无法知道什么时候迁移完成，因此只能是**临时性的重定向，客户端不会更新slots缓存。**

- 但是MOVED重定向说明键对应的槽已经明确指定到新的节点，**因此需要更新slots缓存。**

#### 集群下的Jedis客户端

参见模块redis-cluster。

同时集群下的Jedis客户端只能支持有限的有限的批量操作，必须要求所有key的slot值相等。这时可以考虑使用hash tags。

##### Hash tags

> 在 Redis 集群中，**Hash Tags** 是一种机制，用于显式指定某些键被分配到相同的哈希槽中。它允许开发者控制数据的分布，特别是在需要保证某些键位于同一节点的场景下（例如，为了减少跨节点操作的开销）。
>
> **Hash Tags 的定义**
> 如果一个键包含 用大括号 `{}` 包裹的子字符串（例如 `user:{123}:name`），那么只有大括号内的内容（`123`）会被用于计算哈希槽。
>
> **不使用 Hash Tags 的默认行为**
> 如果键中没有大括号 `{}`，则整个键名会被用于计算哈希槽。

集群支持hash tags功能，即可以**把一类key定位到同一个slot**，tag的标识目前不支持配置，只能使用{}，redis处理hash tag的逻辑也很简单，redis只计算从第一次出现{，到第一次出现}的substring（字串）的hash值，如果substring为空，则仍然计算整个key的值。

比如这两个键{user1000}.following 和 {user1000}.followers 会被哈希到同一个哈希槽里，因为只有 user1000 这个子串会被用来计算哈希值。对于 foo{}{bar} 这个键，整个键都会被用来计算哈希值，因为第一个出现的 { 和它右边第一个出现的 } 之间没有任何字符。对于 foo{bar}{zap} 这个键，用来计算哈希值的是 bar 这个子串。

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666925611079/5c9a4d9d9cb14d2d8a1e7dd3d15419a2.png" alt="image.png" style="zoom:80%;" />

我们在使用hashtag特性时，一定要注意，不能把key的离散性变得非常差。

比如，没有利用hashtag特性之前，key是这样的：mall:sale:freq:ctrl:860000000000001，很明显这种key由于与用户相关，所以离散性非常好。而使用hashtag以后，key是这样的：mall:sale:freq:ctrl:{860000000000001}，这种key还是与用户相关，所以离散性依然非常好。但是我们千万不要这样来使用hashtag特性，例如将key设置为：mall:{sale:freq:ctrl}:860000000000001。这样的话，无论有多少个用户多少个key，其{}中的内容完全一样都是sale:freq:ctrl，也就是说，所有的key都会落在同一个slot上，导致整个Redis集群出现**严重的倾斜问题。**

### 集群原理

#### 节点通信

##### 通信流程

在分布式存储中需要提供维护节点**元数据信息**的机制，所谓元数据是指：节点负责哪些数据，是否出现故障等**状态信息**。常见的元数据维护方式分为：**集中式和P2P方式**。Redis集群采用P2P的**Gossip（流言）协议**，Gossip协议工作原理就是**节点彼此不断通信交换信息**，一段时间后所有的节点都会知道集群完整的信息，这种方式类似**流言传播。**

通信过程说明：

1. 集群中的每个节点都会**单独开辟一个TCP通道**，用于节点之间彼此通信，通信端口号在基础端口上加10000。
2. 每个节点在固定周期内通过特定规则选择**几个节点**发送ping消息。

3. 接收到ping消息的节点用pong消息作为响应。

集群中每个节点通过一定规则挑选要通信的节点，每个节点可能知道全部节点，也可能仅知道部分节点，只要这些节点彼此可以正常通信，最终它们会达到一致的状态。当节点出故障、新节点加入、主从角色变化、槽信息变更等事件发生时，通过不断的ping/pong消息通信，经过一段时间后所有的节点都会知道整个集群全部节点的最新状态，从而达到集群状态同步的目的。

##### Gossip 消息

Gossip协议的主要职责就是**信息交换**。信息交换的载体就是节点彼此发送的Gossip消息，了解这些消息有助于我们理解集群如何完成信息交换。常用的Gossip消息可分为：**ping消息、pong消息、meet消息、fail消息等。**

- **meet消息**：用于**通知新节点加入。**消息发送者通知接收者加入到当前集群，meet消息通信正常完成后，接收节点会加入到集群中并进行周期性的ping、pong消息交换。

- **ping消息**：集群内交换最频繁的消息，集群内每个节点**每秒向多个其他节点**发送ping消息，用于检测节点是否在线和交换彼此状态信息。ping消息发送封装了**自身节点和部分其他节点的状态数据。**

- **pong消息**：当接收到ping、meet消息时，作为响应消息回复给发送方确认消息正常通信。pong消息内部封装了**自身状态数据。**节点也可以向集群内**广播**自身的pong消息来通知整个集群对自身状态进行更新。

- **fail消息**：当节点判定集群内另一个节点下线时，会向集群内**广播一个fail消息**，其他节点接收到fail消息之后**把对应节点更新为下线状态。**

所有的消息格式划分为：消息头和消息体。**消息头包含发送节点自身状态数据**，接收节点根据消息头就可以获取到发送节点的相关数据。

集群内所有的消息都采用相同的**消息头结构clusterMsg**，它包含了**发送节点的关键信息**，如节点id、槽映射、节点标识（主从角色，是否下线）等。**消息体**在Redis内部采用**clusterMsg Data结构**声明。

消息体clusterMsg Data定义发送消息数据。其中ping、meet、pong都采用**clusterMsgDataGossip数组**作为消息体数据结构，实际消息类型使用消息头的**type属性**区分。每个消息体包含该节点的多个clusterMsgDataGossip结构数据，用于信息交换。

当接收到ping、meet消息时，接收节点会解析消息内容并根据自身的识别情况做出相应处理。

##### 节点选择

虽然Gossip协议的信息交换机制具有天然的分布式特性，但它是有成本的。由于内部需要频繁地进行节点信息交换，而ping/pong消息会携带当前节点和部分其他节点的状态数据，势必会加重带宽和计算的负担。Redis集群内节点通信采用**固定频率（定时任务每秒执行10次）。**

因此节点每次选择需要通信的节点列表变得非常重要。通信节点选择过多虽然可以做到信息及时交换但成本过高。节点选择过少会降低集群内所有节点彼此信息交换的频率，从而影响故障判定、新节点发现等需求的速度。因此Redis集群的Gossip协议需要兼顾信息交换的**实时性和成本开销。**

消息交换的成本主要体现在单位时间选择发送消息的**节点数量**和**每个消息携带的数据量。**

1. **选择发送消息的节点数量**

   集群内每个节点维护定时任务默认间隔1秒，每秒执行10次，定时任务里**每秒随机选取5个节点**，找出**最久没有通信**的节点发送ping消息，用于保证 Gossip信息交换的随机性。同时每100毫秒都会扫描**本地节点列表**，如果发现节点（自己）最近一次接受pong消息的时间大于cluster_node_timeout/2，则立刻发送ping消息，防止该节点信息太长时间未更新。

   根据以上规则得出每个节点每秒需要发送ping消息的数量 = 1+10

   * num(node.pong_received > cluster_node_timeout/2)，因此**cluster_node_timeout**参数对消息发送的节点数量影响非常大。当我们的带宽资源紧张时，可以适当调大这个参数，如从默认15秒改为30秒来降低带宽占用率。过度调大cluster_node_timeout 会影响消息交换的频率从而影响故障转移、槽信息更新、新节点发现的速度。因此需要根据业务容忍度和资源消耗进行平衡。同时整个集群消息总交换量也跟节点数成正比。


2. **消息数据量**

   每个ping消息的数据量体现在消息头和消息体中，其中消息头主要占用空间的字段是`myslots [CLUSTER_SLOTS/8]`，占用2KB，这块空间占用相对固定。消息体会携带一定数量的**其他节点信息**用于信息交换。消息体携带的数据量跟集群的节点数息息相关，更大的集群每次消息通信的成本也就更高，因此对于Redis集群来说并不是大而全的集群更好。

#### 故障转移

Redis集群自身实现了高可用。高可用首先需要解决集群部分失败的场景：当集群内少量节点出现故障时通过自动故障转移保证集群可以正常对外提供服务。

##### 故障发现

当集群内某个节点出现问题时，需要通过一种健壮的方式保证识别出节点是否发生了故障。Redis集群内节点通过ping/pong消息实现节点通信，消息不但可以传播节点槽信息，还可以传播其他状态如：主从状态、节点故障等。因此故障发现也是通过消息传播机制实现的，主要环节包括：**主观下线(pfail)和客观下线(fail)**。

- **主观下线：**指某个节点认为另一个节点不可用，即下线状态，这个状态并不是最终的故障判定，只能代表一个节点的意见，可能存在误判情况。

- **客观下线：**指标记一个节点**真正的下线**，集群内多个节点都认为该节点不可用，从而达成共识。如果是**持有槽的主节点故障**，需要为该节点进行故障转移。

###### 主观下线

集群中每个节点都会定期向其他节点发送ping消息，接收节点回复pong消息作为响应。如果在`cluster-node-timeout`时间内通信一直失败，则发送节点会认为接收节点存在故障，把接收节点标记为主观下线(pfail)状态。

**流程说明：**

1. 节点a发送ping消息给节点b，如果通信正常将接收到pong消息，节点a更新最近一次与节点b的通信时间。
2. 如果节点a与节点b通信出现问题则断开连接，下次会进行重连。如果一直通信失败，则节点a记录的与节点b最后通信时间将无法更新。

3. 当节点a内的定时任务检测到与节点b**最后通信时间**超过`cluster-node-timeout`时，更新本地对节点b的状态为**主观下线(pfail)。**

主观下线简单来讲就是：当cluster-note-timeout时间内某节点无法与另一个节点顺利完成ping消息通信时，则将该节点标记为主观下线状态。每个节点内的**clusterstate结构**都需要保存其他节点信息，用于从自身视角判断其他节点的状态。

Redis集群对于节点最终是否故障判断非常严谨，只有一个节点认为主观下线并不能准确判断是否故障。比如节点6379与6385通信中断，导致6379判断6385为主观下线状态，但是6380与6385节点之间通信正常，这种情况不能判定节点6385发生故障。因此对于一个健壮的故障发现机制，需要集群内**大多数节点**都判断6385故障时，才能认为6385确实发生故障，然后为6385节点进行故障转移。而这种多个节点协作完成故障发现的过程叫做客观下线。

###### 客观下线

当某个节点判断另一个节点主观下线后，相应的**节点状态**会跟随消息在集群内**传播**。

ping/pong消息的**消息体**会携带集群1/10的其他节点的**状态数据**，当接收节点发现消息体中含有主观下线的节点状态时，会在本地找到故障节点的ClusterNode结构，保存到**下线报告链表**中。

通过Gossip消息传播，集群内节点**不断收集到故障节点的下线报告**。当**半数以上持有槽的主节点**都标记某个节点是主观下线时。**触发客观下线**流程。这里有两个问题：

1. 为什么必须是持有槽的主节点参与故障发现决策？

   因为集群模式下只有处理槽的**主节点**才负责读写请求和集群槽等**关键信息维护**，而**从节点只进行主节点的数据和状态信息的复制。**

2. 为什么是半数以上处理槽的主节点？

   必须半数以上是为了应对网络分区等原因造成的**集群分割**情况，被分割的小集群无法（参与）完成从主观下线到客观下线这一关键过程，从而防止小集群完成故障转移之后继续对外提供服务。

**尝试客观下线**

集群中的节点每次接收到其他节点的pfail状态，都会**尝试触发客观下线。**

流程说明：

1. 首先统计有效的下线报告数量，如果小于集群内持有槽的主节点总数的一半则退出。
2. 当下线报告大于槽主节点数量一半时，标记对应故障节点为客观下线状态。

3. 向集群**广播**一条fail消息，通知所有的节点将故障节点标记为客观下线，**fail消息的消息体只包含故障节点的ID。**

广播fail消息是客观下线的最后一步，它承担着非常重要的职责：

- 通知集群内所有的节点标记故障节点为客观下线状态并立刻生效。

- 通知故障节点的从节点触发故障转移流程。


##### 故障恢复

故障节点变为客观下线后，如果下线节点是持有槽的主节点，则需要**在它的从节点中选出一个替换它**，从而保证集群的高可用。下线主节点的所有从节点承担故障恢复的义务，当从节点通过内部定时任务发现自身复制的主节点进入客观下线时，将会触发故障恢复流程。

###### 资格检查

每个从节点都要检查最后与主节点**断线时间**，判断是否有资格替换故障的主节点。如果从节点与主节点断线时间超过`cluster-node-time * cluster-slave-validity-factor`，则当前从节点不具备故障转移资格。参数`cluster-slave-validity-factor`用于从节点的有效因子，**默认为10。**

###### 准备选举时间

当从节点符合故障转移资格后，更新触发故障选举的时间，只有到达该时间后才能执行后续流程。这里之所以采用**延迟触发机制**，主要是通过对多个从节点使用不同的**延迟选举时间**来支持**优先级**问题。**复制偏移量越大**说明从节点延迟越低，那么它应该具有**更高的优先级**来替换故障主节点。**所有的从节点中复制偏移量最大的将提前触发故障选举流程。**

主节点b进入客观下线后，它的三个从节点根据自身复制偏移量设置**延迟选举时间**，如复制偏移量最大的节点slave b-1延迟1秒执行，保证复制延迟低的从节点优先发起选举。

###### 从节点发起选举

> 一个从节点在确认其主节点已失效后，尝试成为新的主节点。它会向其他主节点发送 `FAILOVER_AUTH_REQUEST` 请求，请求投票支持。

当从节点定时任务检测到达故障选举时间（failover_auth_time）到达后，发起选举流程如下：

- 更新配置纪元(递增)


> Redis 使用**配置纪元**（config epoch）来管理集群的状态变更，确保整个集群的数据一致性和冲突解决。
>
> **配置纪元的定义**
>
> - 配置纪元是一个全局递增的数值，用来标识集群状态的变更。
> - 每次有新的主节点被选举时，会分配一个新的配置纪元，作为该节点的当前配置版本。

配置纪元是一个**只增不减的整数**，每个主节点自身维护一个**配置纪元（clusterNode.configEpoch）**标识当前主节点的版本，所有主节点的配置纪元都不相等，从节点会复制主节点的配置纪元。整个集群又维护一个**全局的配置纪元（clusterstate.currentEpoch）**，用于记录集群内**所有主节点配置纪元的最大版本**。执行cluster info命令可以查看配置纪元信息：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666925611079/7d2573420ab14414912068c2c0469879.png" alt="image.png" style="zoom:80%;" />

- 配置纪元的主要作用：标识集群内每个主节点的不同版本和当前集群最大的版本。


每次集群发生重要事件时（这里的重要事件指**出现新的主节点**（新加入的或者由从节点转换而来）），从节点竞争选举，都会使集群全局的配置纪元递增，并赋值给相关主节点，用于记录这一关键事件。

主节点具有更大的配置纪元代表了更新的集群状态，因此当节点间进行ping/pong消息交换时，如出现slots等关键信息不一致时，**以配置纪元更大的一方为准**，防止过时的消息状态污染集群。

配置纪元的**应用场景**有：

1. 新节点加入。
2. 槽节点映射冲突检测。
3. 从节点投票选举冲突检测。

> **配置纪元使用过程细节**
>
> 1. **分配新的配置纪元**
>    - 当某个从节点被选举为新的主节点后，它会向负责生成配置纪元的节点（一般是集群中的其他活跃主节点）请求一个新的配置纪元。
>    - 生成的配置纪元比当前集群中所有节点的配置纪元都大。
> 2. **更新配置纪元**
>    - 新的主节点将其新的配置纪元广播给集群中的所有节点。
>    - 集群中的节点接收此更新后，更新其本地元数据。
> 3. **配置纪元的作用**
>    - 确保新主节点的主从关系明确，防止其他从节点争夺主节点身份。
>    - 配置纪元是 Redis 集群防止脑裂的重要机制，因为较大的配置纪元表示当前的最新状态。
>
> **配置纪元的特点**
>
> - **唯一性**：**每个配置纪元只能对应一次选举结果。**
> - **递增性**：配置纪元始终递增，标志着集群的变更顺序。
> - **一致性**：所有节点在完成配置纪元同步后会有一致的状态视图。

###### 选举投票

> 每个主节点只能在一个选举周期内投票给一个候选从节点。
>
> 选票的数量需要超过集群中活跃主节点的半数（`majority`）才能当选成功。
>
> 如果多个从节点同时发起选举，只有第一个收到多数票的从节点会成为新的主节点。

只有持有槽的主节点才会处理**故障选举消息(FAILOVER_AUTH_REQUEST)**，因为每个持有槽的节点在**一个配置纪元内**（即一个选举周期）都有**唯一的一张选票**，当接到第一个请求投票的从节点消息时**回复`FAILOVER_AUTH_ACK`消息作为投票**，之后相同配置纪元内其他从节点的选举消息将忽略。

投票过程其实是一个领导者选举的过程，如集群内有N个持有槽的主节点代表有N张选票。由于在每个配置纪元（即每个选举周期）内持有槽的主节点只能投票给一个从节点，因此只能有一个从节点获得 N/2+1的选票，保证能够找出唯一的从节点。

> Redis集群没有直接使用从节点进行领导者选举，主要因为从节点数必须大于等于3个才能保证凑够N/2+1个节点，将导致从节点资源浪费。使用集群内所有持有槽的主节点进行领导者选举，即使只有一个从节点也可以完成选举过程。

当从节点收集到N/2+1个持有槽的主节点的投票时，从节点可以执行替换主节点操作，例如集群内有5个持有槽的主节点，主节点b故障后还有4个，当其中一个从节点收集到3张投票时，代表获得了足够的选票，可以进行**替换主节点**操作。

**投票作废：每个配置纪元代表了一次选举周期。**如果在开始投票之后的`cluster-node-timeout*2`时间内从节点没有获取足够数量的投票，则本次选举作废。从节点**对配置纪元自增**并发起下一轮投票，**直到选举成功为止。**

###### 替换主节点

当从节点收集到足够的选票之后，触发**替换主节点操作**：

1. 当前从节点取消复制变为主节点。
2. 执行 `clusterDelslot` 操作**撤销故障主节点负责的槽**，并执行 `clusterAddSlot` 把这些槽委派给自己，即接管其原主节点负责的槽和请求。
3. 向集群广播自己的pong消息（Gossip 协议），通知集群内其他所有的节点当前从节点变为了主节点，并接管了故障主节点的槽信息。其他节点更新元数据以识别新的主节点。

##### 故障转移时间

在介绍完故障发现和恢复的流程后，这时我们可以估算出故障转移时间：

1. 主观下线（pfail）识别时间 = cluster-node-timeout。
2. 主观下线状态消息传播时间 <= cluster-node-timeout/2。消息通信机制会对超过cluster-node-timeout/2未通信节点会发起ping消息，消息体在选择包含哪些节点时会优先选取下线状态节点，所以通常这段时间内能够收集到半数以上主节点的 pfail 报告从而完成故障发现（客观下线）。

3. 从节点转移时间 <= 1000毫秒。由于存在延迟发起选举机制，偏移量最大的从节点会**最多延迟1秒**发起选举。**通常第一次选举就会成功**，所以从节点执行转移时间在1秒以内。

根据以上分析可以预估出故障转移时间，如下:

`failover-time(毫秒) ≤ cluster-node-timeout + cluster-node-timeout/2 + 1000`

因此，故障转移时间跟`cluster-node-timeout`参数息息相关，默认15秒。配置时可以根据业务容忍度做出适当调整，但不是越小越好。

#### 集群不可用的判定

> 注意：Redis 集群的槽数量是 **16384**，编号范围为 **0 到 16383**。

为了保证集群完整性，默认情况下当集群16384个槽**任何一个**没有被指派到节点时整个集群都不可用。执行任何键命令都会返回`(error)CLUSTERDOWN Hash slot not served`错误。**这是对集群完整性的一种保护措施，保证所有的槽都指派给在线的节点**。

当持有槽的主节点下线时，从故障发现到自动完成转移期间整个集群是不可用状态，对于大多数业务无法容忍这种情况，因此可以将参数`cluster-require-full-coverage`配置为no，当主节点故障时只影响它负责槽的相关命令执行，不会影响其他主节点的可用性。

但是从集群的故障转移的原理来说，集群会出现不可用，当：

1. 当访问一个 Master 和 Slave 节点都挂了的时候，cluster-require-full-coverage=yes，会报告槽无法获取。
2. 集群主库半数宕机（根据 failover 原理，fail（主观下线）掉一个主节点需要一半以上其他主节点都投票通过才可以）。

另外，当集群 Master 节点个数小于 3 个的时候，或者集群可用节点个数为偶数的时候，基于 fail 的这种选举机制的自动主从切换过程可能会不能正常工作，一个是标记 fail 的过程，一个是选举新的 master 的过程，都有可能异常。

> **cluster-require-full-coverage工作机制**
>
> 1. **`yes` 模式下：默认**
>    - 集群要求所有的 16384 个哈希槽都必须有主节点负责。
>    - 故障状态下，所有客户端请求都会被拒绝，直到管理员修复问题（例如，通过手动重新分配槽或恢复节点）。
> 2. **`no` 模式下：**
>    - 即使部分槽未分配或无主节点负责，集群仍然处于工作状态。
>    - 集群会尽可能地处理来自客户端的请求，客户端仍然可以访问其他健康的节点和槽。
>    - 不过，对于未覆盖的槽的请求，客户端会收到 `MOVED` 或 `ASK` 错误。

#### 集群读写分离

1. **只读连接**

集群模式下**从节点不接受任何读写请求**，发送过来的键命令会重定向到负责槽的主节点上(其中包括它的主节点)。当需要使用从节点分担主节点**读压力**时，可以使用**readonly命令**打开客户端连接只读状态。**之前的复制配置slave-read-only在集群模式下无效**。

当开启只读状态时，从节点接收**读命令处理流程**变为：如果对应的槽属于自己正在复制的主节点则直接执行读命令，否则返回重定向信息。

readonly命令是**连接级别生效**，因此每次新建连接时都需要执行readonly开启只读状态。执行**readwrite命令**可以关闭连接只读状态。

2. **读写分离**

集群模式下的读写分离，同样会遇到：复制延迟、读取过期数据、从节点故障等问题。针对从节点故障问题，客户端需要维护**可用节点列表**，集群提供了cluster slaves {nodeld}命令（已弃用），返回nodeId对应的主节点下所有从节点信息，命令如下:

```bash
cluster slave 41ca2d569068043a5f2544c598edd1e45a0c1f91
```

> `cluster slave` 命令已经被弃用，取而代之的是更现代化的 `CLUSTER REPLICAS` 命令。
>
> ```
> CLUSTER REPLICAS <node_id>
> ```
>
> 返回与指定主节点相关的从节点列表。每个从节点的信息包括：
>
> - 节点 ID
> - 节点地址（IP:端口）
> - 节点角色（从节点）
> - 节点状态

解析以上从节点列表信息，排除fail状态节点，这样客户端对从节点的故障判定可以**委托给集群处理**，简化维护可用从节点列表难度。

同时集群模式下读写分离涉及对客户端修改如下：

1. 维护每个主节点可用从节点列表。
2. 针对读命令维护请求节点路由。

3. 从节点新建连接开启readonly状态。

**集群模式下读写分离成本比较高**，可以直接扩展主节点数量提高集群性能，**一般不建议集群模式下做读写分离。**



## 哨兵模式（Redis Sentinel）

Redis的主从复制模式下，一旦主节点由于故障不能提供服务，需要人工将从节点晋升为主节点，同时还要通知应用方更新主节点地址，对于很多应用场景这种故障处理的方式是无法接受的。

Redis 从 2.8开始正式提供了Redis Sentinel(哨兵）架构来解决这个问题。

### 主从复制的问题

Redis 的主从复制模式可以将主节点的数据改变同步给从节点，这样从节点就可以起到两个作用

- 第一，作为主节点的一个备份，一旦主节点出了故障不可达的情况，从节点可以作为后备“顶”上来，并且保证数据尽量不丢失(主从复制是最终一致性)。

- 第二，从节点可以扩展主节点的读能力，一旦主节点不能支撑住大并发量的读操作，从节点可以在一定程度上**帮助主节点分担读压力。**


但是主从复制也带来了以下问题：

1. 一旦主节点出现故障，需要**手动**将一个从节点晋升为主节点，同时需要修改应用方的主节点地址，还需要命令其他从节点去复制新的主节点，整个过程都需要人工干预。
2. 主节点的写能力受到单机的限制。
3. 主节点的存储能力受到单机的限制。

### Redis Sentinel

Redis Sentinel是一个分布式架构，其中包含若干个Sentinel节点和Redis数据节点，**每个Sentinel节点会对数据节点和其余Sentinel节点进行监控**，当它发现节点不可达时，会对节点做**下线标识**。如果被标识的是主节点，它还会和其他Sentinel节点进行**“协商”**，当大多数Sentinel节点都认为主节点不可达时，它们会**选举出一个Sentinel节点**来完成**自动故障转移**的工作，同时会将这个变化实时通知给Redis应用方。整个过程完全是自动的，不需要人工来介入，所以这套方案很有效地解决了Redis的高可用问题。

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/0bdb4b558f974743ba1172af64f81dd9.png" alt="image.png" style="zoom:80%;" />

### Redis Sentinel的搭建

我们以以3个Sentinel节点、1个主节点、2个从节点组成一个Redis Sentinel进行说明。

启动主从的方式和普通的主从没有不同。

#### 启动Sentinel节点

Sentinel节点的启动方法有两种：

方法一，使用redis-sentinel命令：

```bash
./redis-sentinel   ../conf/reids.conf
```

方法二，使用redis-server命令加--sentinel参数:

```bash
./redis-server ../conf/reids.conf  --sentinel
```

两种方法本质上是—样的。

##### 确认

Sentinel节点本质上是一个**特殊的Redis节点**，所以也可以通过**info命令**来查询它的相关信息。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/a0ee43231155484a9ccf54dc5edf3ceb.png)

### 实现原理

Redis Sentinel的基本实现中包含以下：

Redis Sentinel 的**定时任务、主观下线和客观下线、Sentinel领导者选举、故障转移**等等知识点，学习这些可以让我们对Redis Sentinel的高可用特性有更加深入的理解和认识。

#### 三个定时监控任务

一套合理的监控机制是Sentinel节点判定节点不可达的重要保证，Redis Sentinel通过三个定时监控任务完成对各个节点发现和监控：

##### 1、每隔10秒的定时监控（主从信息同步）

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/8a1dfef71fe74479a2b4a4998255652c.png)

每隔10秒，每个Sentinel节点会向主节点（和从节点？）发送**info命令**获取最新的拓扑结构，Sentinel节点通过对上述结果进行解析就可以找到相应的从节点。

这个定时任务的作用具体可以表现在三个方面：

1. 通过向主节点执行info命令，获取从节点的信息，这也是为什么Sentinel节点不需要显式配置监控从节点。
2. 当有新的从节点加入时都可以立刻感知出来。
3. 节点不可达或者故障转移后，可以通过info命令实时更新节点拓扑信息。

##### 2、每隔2秒的定时监控（哨兵节点间的协作）

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/fddd1423ec6e4c22978546cd2e428a8c.png)

每隔2秒，每个Sentinel节点会向Redis数据节点的**sentinel_:hello频道**上发送该Sentinel节点对于主节点的判断，以及当前Sentinel节点的信息，同时每个Sentinel节点也会**订阅**该频道，来了解其他Sentinel节点以及它们对主节点的判断，所以这个定时任务可以完成以下两个工作:

- 发现新的Sentinel节点：通过订阅主节点的**sentinel:hello**了解其他的Sentinel节点信息，如果是新加入的Sentinel节点，将该Sentinel节点信息保存起来,并与该 Sentinel节点创建连接。

- Sentinel节点之间交换主节点的状态，作为后面客观下线以及领导者选举的依据。


##### 3、每隔1秒的定时监控（监控和心跳检测）

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/905d356c71714e7084583296eeac1350.png)

每隔1秒，每个Sentinel节点会向**主节点、从节点、其余Sentinel节点**发送一条ping命令做一次心跳检测，来确认这些节点当前是否可达。通过这个定时任务，Sentinel节点对主节点、从节点、其余Sentinel节点都建立起连接，实现了对每个节点的监控，这个定时任务是节点失败判定的重要依据。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/9bb97de253d24fa69d4035942175c872.png)

#### 主观下线和客观下线

##### 主观下线

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/a17f946612e24d7f8e95c42f9a16e855.png" alt="image.png" style="zoom:80%;" />

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/ac050de6f3874fbb9b1e79503d69b713.png)

上一小节介绍的第三个定时任务，每个Sentinel节点会每隔1秒对主节点、从节点、其他Sentinel节点发送ping命令做心跳检测。当这些节点超过**down-after-milliseconds**时间没有进行有效回复，Sentinel节点就会对该节点做失败判定，这个行为叫做**主观下线**。从字面意思也可以很容易看出主观下线是当前Sentinel节点的一家之言，存在误判的可能。

##### 客观下线

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/238a751dc2f84431b84b78584a2fcdcd.png)

当Sentinel主观下线的节点是**主节点**时，该Sentinel节点会通过**sentinel is-master-down-by-addr**命令向其他Sentinel节点询问对主节点的判断，当超过**quorum**（法定人数）个数的Sentinel节点认为主节点确实有问题，这时**该Sentinel节点会做出客观下线**（O_DOWN）的决定，这样客观下线的含义是比较明显了，也就是大部分Sentinel节点都对主节点的下线做了同意的判定，那么这个判定就是客观的。

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665215626018/78caadf29ad3402fae0d7d52dc49cb7e.png)

#### 领导者Sentinel节点选举

领导者Sentinel节点选举，是为了后续的故障转移。

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/4708aea0fae3469dae8bec534554749a.png" alt="image.png" style="zoom: 50%;" />

假如Sentinel节点对于主节点已经做了客观下线，那么是不是就可以立即进行**故障转移**了？当然不是，实际上故障转移的工作只需要一个Sentinel节点来完成即可，所以 Sentinel节点之间会做一个领导者选举的工作，**选出一个Sentinel节点作为领导者进行故障转移的工作。**Redis使用了Raft算法实现领导者选举，Redis Sentinel进行领导者选举的大致思路如下：

1. 每个在线的Sentinel节点都有资格成为领导者，当它确认主节点主观下线时候，会向其他Sentinel节点发送`sentinel is-master-down-by-addr`命令，**要求将自己设置为领导者。**
2. 收到命令的Sentinel节点，如果没有同意过其他Sentinel节点的sentinel is-master-down-by-addr命令，将同意该请求，否则拒绝。

3. 如果该Sentinel节点发现自己的票数已经大于等于`max(quorum，num(sentinels)/2+1)`，那么它将成为领导者。

4. 如果此过程没有选举出领导者，将进入下一次选举。

选举的过程非常快，**基本上谁先完成客观下线，谁就是领导者。**

Raft协议的详细版本：（后续再看）

[raft-zh_cn/raft-zh_cn.md at master · maemual/raft-zh_cn · GitHub](https://github.com/maemual/raft-zh_cn/blob/master/raft-zh_cn.md)

如果你想手写一个Raft协议，可以看下蚂蚁金服的开发生产的raft算法组件

[GitHub - sofastack/sofa-jraft: A production-grade java implementation of RAFT consensus algorithm.](https://github.com/sofastack/sofa-jraft)

选举很快的！！

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1665215626018/5191985ef85f49c9a8efd144f5272b1a.png)

#### 故障转移

经过选举得出的领导者Sentinel节点负责故障转移，具体步骤如下：

<img src="https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/a8b8e35e659b4521b7571b5395b0808d.png" alt="image.png" style="zoom: 67%;" />

1. 选举出的Sentinel在**从节点列表中选出一个节点作为新的主节点**，选择方法如下：

   1. 过滤“不健康”(主观下线、断线等)、5秒内没有回复过Sentinel节点ping响应、与主节点失联超过down-after-milliseconds*10秒的节点。

   2. 选择slave-priority(从节点优先级)最高的从节点列表，如果存在则返回，不存在则继续。

   3. 选择**复制偏移量最大**的从节点(**复制的最完整**)，如果存在则返回，不存在则继续。

   3. 选择runid最小的从节点。

2. Sentinel领导者节点会对第一步选出来的从节点执行**slaveof no one**命令让其成为主节点。

   > `SLAVEOF NO ONE` 用于让当前节点从从节点（Replica）切换回主节点（Master）。执行该命令后，当前节点将停止与主节点的复制，并成为独立的主节点。以下是详细解释：
   >
   > **命令作用**
   >
   > 1. 停止复制：执行后，该 Redis 实例会断开与原主节点的复制连接，不再同步数据。
   >
   > 2. 切换角色：从节点（Replica）变为独立的主节点（Master），可接受客户端写入操作。
   >
   > 3. 数据保留：保留切换前已复制的数据，后续写入的数据独立存储，与原主节点不再关联。

3. Sentinel领导者节点会向**剩余的从节点**发送命令，让它们成为新主节点的从节点，复制规则和**parallel-syncs**参数有关。

4. Sentinel节点集合会将原来的主节点更新为从节点，并**保持着对其关注**，当其恢复后命令它去复制新的主节点。

### Redis Sentinel的客户端

如果主节点挂掉了，虽然Redis Sentinel可以完成故障转移，但是客户端无法获取这个变化，那么使用Redis Sentinel的意义就不大了，所以各个语言的客户端需要对Redis Sentinel进行显式的支持。

Sentinel节点集合具备了监控、通知、自动故障转移、配置提供者等若干功能，也就是说实际上最了解主节点信息的就是Sentinel节点集合，而各个主节点可以通过&#x3c;host-name>进行标识的，所以，无论是哪种编程语言的客户端，如果需要正确地连接Redis Sentinel，必须有Sentinel节点集合和masterName两个参数。

#### Java实现（不用看）

我们依然使用Jedis 作为Redis 的 Java客户端，Jedis能够很好地支持Redis
Sentinel，并且使用Jedis连接Redis Sentinel也很简单，按照Redis Sentinel的原理，需要有masterName和Sentinel节点集合两个参数。Jedis针对Redis Sentinel给出了一个 JedisSentinelPool。

具体代码可以参见redis-sentinel：

![image.png](https://fynotefile.oss-cn-zhangjiakou.aliyuncs.com/fynote/fyfile/5983/1666336676066/37b4bfd4f0cc4ba8b2484f5be610ff5e.png)

实现一个Redis Sentinel客户端一般来说需要：

1. 遍历Sentinel节点集合获取一个可用的Sentinel节点，Sentinel节点之间可以共享数据，所以从任意一个Sentinel节点获取主节点信息都是可以的。
2. 通过sentinel的`get-master-addr-by-name host-name`这个API来获取对应主节点的相关信息。

3. 验证当前获取的“主节点”是真正的主节点，**这样做的目的是为了防止故障转移期间主节点的变化。**

4. 保持和 Sentinel节点集合的“联系”，时刻获取关于主节点的相关“信息”。

但是注意，JedisSentinel的实现是不支持读写分离的，所有的连接都是连接到Master上面，Slave就完全当成Master的备份，存在着性能浪费。因此如果想支持读写分离，需要自行实现，这里给一个参考

[基于Spring 的 Redis Sentinel 读写分离 Slave 连接池 (jack-yin.com)](https://www.jack-yin.com/coding/spring-boot/2683.html)

#### Golang实现（暂未实现）

### 高可用读写分离

#### 从节点的作用

- 第一，当主节点出现故障时，作为主节点的后备“顶”上来实现故障转移，Redis Sentinel已经实现了该功能的自动化，实现了真正的高可用。

- 第二，扩展主节点的**读能力**，尤其是在**读多写少**的场景非常适用。


但上述模型中，从节点不是高可用的。

如果slave-1节点出现故障，首先客户端client-1将与其失联，其次Sentinel节点只会对该节点做主观下线，因为**Redis Sentinel的故障转移是针对主节点的。**所以很多时候，Redis Sentinel中的从节点仅仅是作为主节点一个热备，不让它参与客户端的读操作，就是为了保证整体高可用性，但实际上这种使用方法还是有一些浪费，尤其是在有很多从节点或者确实需要读写分离的场景，所以如何实现从节点的高可用是非常有必要的。

#### Redis Sentinel读写分离设计思路参考

Redis Sentinel在对各个节点的监控中，如果有对应事件的发生，都会发出相应的事件消息，其中和从节点变动的事件有以下几个：

- **+switch-master**：切换主节点(原来的从节点晋升为主节点)，说明减少了某个从节点。
- **+convert-to-slave**：切换从节点(原来的主节点降级为从节点)，说明添加了某个从节点。
- **+sdown**：主观下线，说明可能某个从节点**可能不可用**（因为对从节点不会做客观下线），所以在实现客户端时可以采用自身策略来实现类似主观下线的功能。
- **+reboot**：重新启动了某个节点，如果它的角色是slave，那么说明添加了某个从节点。

所以在设计Redis Sentinel的从节点高可用时，只要能够实时掌握所有从节点的状态，**把所有从节点看做一个资源池**，无论是上线还是下线从节点，客户端都能及时感知到(将其从资源池中添加或者删除)，这样从节点的高可用目标就达到了。
