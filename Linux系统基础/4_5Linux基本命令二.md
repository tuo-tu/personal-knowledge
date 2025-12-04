# ·Python 全栈文档

## 第五章  Linux基本命令二

#### 1、head 命令

head

​           **作用：**用于**查看文件的开头部分的内容**，有一个常用的参数 **-n** 用于显示行数，默认为 10，即显示 10 行的内容

​           **语法：**head [参数] [文件]  

​           **命令参数：**

| 参数     | 参数描述     |
| -------- | ------------ |
| -q       | 隐藏文件名   |
| -v       | 显示文件名   |
| -c<数目> | 显示的字节数 |
| -n<行数> | 显示的行数   |

​            （1）显示 1.log 文件中前 20 行

```bash
[root@localhost ~]# head 1.log -n 20
```

​            （2）显示 1.log 文件前 20 字节

```bash
[root@localhost ~]# head -c 20 log2014.log
```

​            （3）显示 t.log最后 10 行

```bash
[root@localhost ~]# head -n -10 t.log
```

扩展：tail 命令，查看文件的末尾

#### 2、which 命令

which

​        在 linux 要查找某个命令或者文件，但不知道放在哪里了，可以使用下面的一些命令来搜索

```
which     查看可执行文件的位置。
whereis   查看文件的位置。
locate    配合数据库查看文件位置。
find      实际搜寻硬盘查询文件名称。
```

​           **作用：**用于查找文件（which指令会在环境变量$PATH设置的目录里查找符合条件的文件。）

​           **语法：**which [文件...]

​           **命令参数：**

| 参数           | 参数描述                                                     |
| -------------- | ------------------------------------------------------------ |
| -n<文件名长度> | 指定文件名长度，指定的长度必须大于或等于所有文件中最长的文件名 |
| -p<文件名长度> | 与-n参数相同，但此处的<文件名长度>包括了文件的路径           |
| -w             | 指定输出时栏位的宽度                                         |
| -V             | 显示版本信息                                                 |

​            （1）查看 ls 命令是否存在，执行哪个

```bash
[root@localhost ~]# which ls
alias ls='ls --color=auto'
        /usr/bin/ls

```

​            （2）查看 which

```bash
[root@localhost ~]# which which
alias which='alias | /usr/bin/which --tty-only --read-alias --show-dot --show-tilde'
        /usr/bin/alias
        /usr/bin/which
```

​            （3）查看 cd

```bash
[root@localhost ~]# which cd
/usr/bin/cd

（注意：显示不存在，因为 cd 是内建命令，而 which 查找显示是 PATH 中的命令）

```

​            （4）查看当前 PATH 配置

```bash
[root@localhost ~]# echo $PATH
/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/root/bin

```



#### 3、whereis命令

whereis

```
whereis 命令只能用于程序名的搜索，而且只搜索二进制文件（参数-b）、man说明文件（参数-m）和源代码文件（参数-s）。如果省略参数，则返回所有信息。
whereis 及 locate 都是基于系统内建的数据库进行搜索，因此效率很高，
而find则是遍历硬盘查找文件
```

​           **作用：**用于查找文件

​           **语法：**whereis \[-bfmsu]\[-B <目录>...]\-M <目录>...]\[-S <目录>...][文件...]

​           **命令参数：**

| 参数     | 参数描述                                                     |
| -------- | ------------------------------------------------------------ |
| -b       | 定位可执行文件                                               |
| -B<目录> | 只在设置的目录下查找可执行文件                               |
| -f       | 不显示文件名前的路径名称                                     |
| -m       | 定位帮助文件                                                 |
| -M<目录> | 只在设置的目录下查找说帮助文件                               |
| -s       | 定位源代码文件                                               |
| -S<目录> | 只在设置的目录下查找源代码文件                               |
| -u       | 搜索默认路径下除可执行文件、源代码文件、帮助文件以外的其它文件 |

​            （1）查找 locate 程序相关文件

```bash
[root@localhost ~]# whereis bash
bash: /usr/bin/bash /usr/share/man/man1/bash.1.gz

```

​            （2）查找 locate 的源码文件

```bash
[root@localhost ~]# whereis -s locate

```

​            （3）查找 lcoate 的帮助文件

```bash
[root@localhost ~]# whereis -m locate
```

#### 4、locate命令

locate

​        需要注意这个命令在我们的最小mini系统里面是没有安装的

```bash
[root@localhost ~]# yum install mlocate
...省略...
[root@localhost ~]# updatedb
```

​           **作用：**用于查找符合条件的文档，他会去保存文档和目录名称的数据库内，查找合乎范本样式条件的文档或目录

​           **语法：**locate \[-d ]\[--help]\[--version][范本样式...]

​           **命令参数：**

| 参数    | 参数描述                                                     |
| ------- | ------------------------------------------------------------ |
| -b      | 仅匹配路径名的基本名称                                       |
| -c      | 只输出找到的数量                                             |
| -d      | 使用 DBPATH 指定的数据库，而不是默认数据库 /var/lib/mlocate/mlocate.db |
| -e      | 仅打印当前现有文件的条目                                     |
| -1      | 如果 是 1．则启动安全模式。在安全模式下，使用者不会看到权限无法看到 的档案。<br />这会始速度减慢，因为 locate 必须至实际的档案系统中取得档案的 权限资料 |
| -0      | 在输出上带有NUL的单独条目                                    |
| -S      | 不搜索条目，打印有关每个数据库的统计信息                     |
| -q      | 安静模式，不会显示任何错误讯息                               |
| -P      | 检查文件存在时不要遵循尾随的符号链接                         |
| -l      | 将输出（或计数）限制为LIMIT个条目                            |
| -n      | 至多显示 n个输出                                             |
| -m      | 被忽略，为了向后兼容                                         |
| -r      | REGEXP -- 使用基本正则表达式                                 |
| --regex | 使用扩展正则表达式                                           |
| -o      | 指定资料库存的名称                                           |
| -h      | 显示帮助                                                     |
| -i      | 忽略大小写                                                   |
| -V      | 显示版本信息                                                 |

​        常用参数：

| 参数 | 参数描述                                   |
| ---- | ------------------------------------------ |
| -l   | num（要显示的行数）                        |
| -f   | 将特定的档案系统排除在外，如将proc排除在外 |
| -r   | 使用正则运算式做为寻找条件                 |

​            （1）查找和 pwd 相关的所有文件(文件名中包含 pwd）

```bash
[root@localhost ~]# locate pwd
/etc/.pwd.lock
/usr/bin/pwd
...省略...
```

​            （2）搜索 etc 目录下所有以 sh 开头的文件

```bash
[root@localhost ~]# locate /etc/sh
/etc/shadow
/etc/shadow-
/etc/shells
```

​            （3）查找 /root 目录下，以 sh 结尾的文件

```bash
[root@localhost ~]# locate -r '^/root.*sh$'
/root/test.sh

```

#### 5、find命令

find

​           **作用：**用于在文件树中查找文件，并作出相应的处理

​           **语法：**

```
find [-H] [-L] [-P] [-Olevel] [-D help|tree|search|stat|rates|opt|exec] [path...] [expression]
```

​           **命令参数：**

| 参数     | 参数描述                                                     |
| -------- | ------------------------------------------------------------ |
| pathname | find命令所查找的目录路径。例如用.来表示当前目录，用/来表示系统根目录 |
| -print   | find命令将匹配的文件输出到标准输出                           |
| -exec    | find命令对匹配的文件执行该参数所给出的shell命令。相应命令的形式为'command' {  } \;，注意{   }和\；之间的空格 |
| -ok      | 和-exec的作用相同，只不过以一种更为安全的模式来执行该参数所给出的shell命令，在执行每一个命令之前，都会给出提示，让用户来确定是否执行 |

​           **命令选项：**

| 选项         | 选项描述                                                     |
| ------------ | ------------------------------------------------------------ |
| -name        | 按照文件名查找文件                                           |
| -perm        | 按文件权限查找文件                                           |
| -user        | 按文件属主查找文件                                           |
| -group       | 按照文件所属的组来查找文件                                   |
| -type        | 查找某一类型的文件，诸如：<br />b - 块设备文件<br />d - 目录<br />c - 字符设备文件<br />l - 符号链接文件<br />p - 管道文件<br />f - 普通文件 |
| -size n :[c] | 查找文件长度为n块文件，带有c时表文件字节大小                 |
| -amin n      | 查找系统中最后N分钟访问的文件                                |
| -atime n     | 查找系统中最后n*24小时访问的文件                             |
| -cmin n      | 查找系统中最后N分钟被改变文件状态的文件                      |
| -ctime n     | 查找系统中最后n*24小时被改变文件状态的文件                   |
| -mmin n      | 查找系统中最后N分钟被改变文件数据的文件                      |
| -mtime n     | 查找系统中最后n*24小时被改变文件数据的文件                   |
| -maxdepth n  | 最大查找目录深度                                             |
| -prune       | 选项来指出需要忽略的目录。在使用-prune选项时要当心，<br />因为如果你同时使用了-depth选项，那么-prune选项就会被find命令忽略 |
| -newer       | 如果希望查找更改时间比某个文件新但比另一个文件旧的所有文件，可以使用-newer选项 |

​            来看例子

​            （1）查找 48 小时内修改过的文件

```bash
[root@localhost ~]# find -atime -2
.
./.bash_profile
...省略...
```

​            （2）在当前目录查找 以 .log 结尾的文件。 **.** 代表当前目录

```bash
[root@localhost ~]# find ./ -name '*.log'
./mydir/text2.log
./test.log

```

​            （3）查找 /opt 目录下 权限为 777 的文件

```bash
[root@localhost ~]# find /opt -perm 777
```

​            （4）查找大于 1K 的文件

```bash
[root@localhost ~]# find -size +1000c
./anaconda-ks.cfg
./.bash_history
./.viminfo

```

​            （5）查找等于 1000 字符的文件

```bash
[root@localhost ~]# find -size 1000c
```

​       **-exec**

```
注意：-exec 参数后面跟的是 command 命令，它的终止是以 ; 为结束标志的，所以这句命令后面的分号是不可缺少的，考虑到各个系统中分号会有不同的意义，所以前面加反斜杠。{} 花括号代表前面find查找出来的文件名
```

​            （6）在当前目录中查找更改时间在10日以前的文件并删除它们(无提醒）

```bash
[root@localhost ~]# find . -type f -mtime +10 -exec rm -f {} \;

# 可以不用操作，了解一下结构即可
```

​            （7）当前目录中查找所有文件名以.log结尾、更改时间在0日以上的文件，并删除它们，

​                     只不过在删除之前先给出提示。 按y键删除文件，按n键不删除

```bash
[root@localhost ~]# find . -name '*.log' -mtime +0 -ok -exec rm {} \;
< -exec ... ./mydir/text2.log > ? n
```

​            （8）用 exec 选项执行 cp 命令

```bash
[root@localhost ~]# find . -name '*.log' -exec cp {} test3 \;
```

​       **-xargs find**

```
-xargs find 命令把匹配到的文件传递给 xargs 命令，而 xargs 命令每次只获取一部分文件而不是全部，不像 -exec 选项那样。这样它可以先处理最先获取的一部分文件，然后是下一批，并如此继续下去。
```

​            （9）查找当前目录下每个普通文件，然后使用 xargs 来判断文件类型

```bash
[root@localhost ~]# find . -type f -print | xargs file
```

​            （10）查找当前目录下所有以 js 结尾的并且其中包含 'editor' 字符的普通文件

```bash
[root@localhost ~]# find . -type f -name "*.js" -exec grep -lF 'ueditor' {} \;
[root@localhost ~]# find -type f -name '*.js' | xargs grep -lF 'editor'
```

​            （11）利用 xargs 执行 mv 命令

```bash
[root@localhost ~]# find . -name "*.log" | xargs -i mv {} test4
```

​            （12）用 grep 命令在当前目录下的所有普通文件中搜索 hostnames 这个词，并标出所在行：

```bash
[root@localhost ~]# find . -name \*(转义） -type f -print | xargs grep -n 'hostnames'
```

​            （13）查找当前目录中以一个小写字母开头，最后是 4 到 9 加上 .log 结束的文件：

```bash
[root@localhost ~]# find . -name '[a-z]*[4-9].log' -print
```

​            （14）在 test 目录查找不在 test4 子目录查找

```bash
[root@localhost ~]# find test -path 'test/test4' -prune -o -print
```

​            （15）实例1：查找更改时间比文件 log2012.log新但比文件 log2017.log 旧的文件

```bash
[root@localhost ~]# find -newer log2012.log ! -newer log2017.log
```

​       **depth** 

```
depth 选项可以使 find 命令向磁带上备份文件系统时，希望首先备份所有的文件，其次再备份子目录中的文件。

实例：find 命令从文件系统的根目录开始，查找一个名为 CON.FILE 的文件。 它将首先匹配所有的文件然后再进入子目录中查找
```

```bash
[root@localhost ~]# find / -name "CON.FILE" -depth -print
```

#### which、whereis、locate、find的区别

`which`用于查找某个命令的**可执行文件路径**。

`whereis`用于查找某个命令的**二进制文件、源代码和帮助文档路径**。

`locate`通过**文件数据库**快速查找文件。速度很快。

`find`在指定目录（文件树中）**逐个扫描查找文件**。

#### 6、chmod命令

​       Linux/Unix 的文件调用权限分为三级 : 文件拥有者、群组、其他。

```
用于改变linux系统文件或目录的访问权限。用它控制文件或目录的访问权限。该命令有两种用法。一种是包含字母和操作符表达式的文字设定法；另一种是包含数字的数字设定法。

每一文件或目录的访问权限都有三组，每组用三位表示，分别为文件属主的读、写和执行权限；与属主同组的用户的读、写和执行权限；系统中其他用户的读、写和执行权限。可使用 ls -l test.txt 查找
```

```bash
[root@localhost ~]# ll
总用量 20
-rw-------. 1 root root 1437 8月  31 15:54 anaconda-ks.cfg
-rw-r--r--. 1 root root    0 9月   8 18:29 file1
lrwxrwxrwx. 1 root root   15 9月   7 16:31 link_text2 -> mydir/text2.log
drwxr-xr-x. 4 root root   92 9月   7 18:08 mydir
-rw-r--r--. 1 root root   13 9月   8 16:35 myfile
-rw-r--r--. 1 root root   36 9月   9 13:16 test3
-rw-r--r--. 1 root root   36 9月   8 18:36 test.log
-rwxr-xr-x. 1 root root   67 9月   8 18:36 test.sh

```

​       这里使用test.log作为例子

```
-rw-r--r--. 1 root root   36 9月   8 18:36 test.log

第一列共有 10 个位置，
第一个字符指定了文件类型。
在通常意义上，一个目录也是一个文件。如果第一个字符是横线，表示是一个非目录的文件。
如果是 d，表示是一个目录。
从第二个字符开始到第十个 9 个字符，3 个字符一组，分别表示了 3 组用户对文件或者目录的权限。
权限字符用横线代表空许可，r 代表只读，w 代表写，x 代表可执行
```

​           **语法：**

```bash
chmod [-cfvR] [--help] [--version] mode file...
```

​           **常用参数：**

| 参数 | 参数描述                           |
| ---- | ---------------------------------- |
| -c   | 当发生改变时，报告处理信息         |
| -R   | 处理指定目录以及其子目录下所有文件 |

​           **权限范围：**

```
u ：目录或者文件的当前的用户
g ：目录或者文件的当前的群组
o ：除了目录或者文件的当前用户或群组之外的用户或者群组
a ：所有的用户及群组
```

​           **权限代号：**

| 代号 | 代号权限              |
| ---- | --------------------- |
| r    | 读权限，用数字4表示   |
| w    | 写权限，用数字2表示   |
| x    | 执行权限，用数字1表示 |
| -    | 删除权限，用数字0表示 |
| s    | 特殊权限              |

```bash
环境：-rw-r--r--. 1 root root   36 9月   8 18:36 test.log
```

​            （1）增加文件 t.log 所有用户可执行权限

```bash
[root@localhost ~]# ls -n test.log
-rwxr-xr-x. 1 0 0 36 9月   8 18:36 test.log
```

​            （2）撤销原来所有的权限，然后使拥有者具有可读权限,并输出处理信息

```bash
[root@localhost ~]# chmod u=r test.log -c
mode of "test.log" changed from 0755 (rwxr-xr-x) to 0455 (r--r-xr-x)
[root@localhost ~]# ls -n test.log
-r--r-xr-x. 1 0 0 36 9月   8 18:36 test.log

```

​            （3）给 file 的属主分配读、写、执行(7)的权限，给file的所在组分配读、执行(5)的权限，给其他用户分配执行(1)的权限

```bash
[root@localhost ~]# chmod 751 test.log -c
或者
[root@localhost ~]# chmod u=rwx,g=rx,o=x t.log -c
```

​            （4）将mydir 目录及其子目录所有文件添加可读权限

```bash
[root@localhost ~]# chmod u+r,g+r,o+r -R text/ -c
```

#### 7、chown命令

```
更改文件或目录的所有者和所属组

chown 将指定文件的拥有者改为指定的用户或组，
用户可以是用户名或者用户 ID；
组可以是组名或者组 ID；文件是以空格分开的要改变权限的文件列表，支持通配符
```

注意：一般来说，这个指令只有是由系统管理者(root)所使用，一般使用者没有权限可以改变别人的文件拥有者，也没有权限把自己的文件拥有者改设为别人。只有系统管理者(root)才有这样的权限。

​           **语法：**

```bash
chown [-cfhvR] [--help] [--version] user[:group] file...
```

​           **常用参数：**

| 参数      | 参数描述                             |
| --------- | ------------------------------------ |
| user      | 新的文件拥有者的使用者 ID            |
| group     | 新的文件拥有者的使用者组(group)      |
| -c        | 显示更改的部分的信息                 |
| -f        | 忽略错误信息                         |
| -h        | 修复符号链接                         |
| -v        | 显示详细的处理信息                   |
| -R        | 处理指定目录以及其子目录下的所有文件 |
| --help    | 显示辅助说明                         |
| --version | 显示版本                             |

​            （1）改变拥有者和群组 并显示改变信息

```bash
[root@localhost ~]# chown -c mail:mail test.log
changed ownership of "test.log" from root:root to mail:mail

-r--r-xr-x. 1 mail mail       36 9月   8 18:36 test.log

```

​            （2）改变文件群

```bash
[root@localhost ~]# chown -c :mail test.sh 
changed ownership of "test.sh" from root:root to :mail

```

​            （3）改变文件夹及子文件目录属主及属组为 mail

```bash
[root@localhost ~]# chown -cR mail: mydir
changed ownership of "mydir/test1/text1.txt" from root:root to mail:mail
changed ownership of "mydir/test1" from root:root to mail:mail
...省略...


```

#### 8、tar 命令

```
用来压缩和解压文件。tar 本身不具有压缩功能，只具有打包功能，有关压缩及解压是调用其它的功能来完成。

弄清两个概念：打包和压缩。打包是指将一大堆文件或目录变成一个总的文件；压缩则是将一个大的文件通过一些压缩算法变成一个小文件

```

​           **作用：**用于备份文件（tar是用来建立，还原备份文件的工具程序，它可以加入，解开备份文件内的文件）

​           **语法：**

```
tar [-ABcdgGhiklmMoOpPrRsStuUvwWxzZ][-b <区块数目>][-C <目的目录>][-f <备份文件>][-F <Script文件>][-K <文件>][-L <媒体容量>][-N <日期时间>][-T <范本文件>][-V <卷册名称>][-X <范本文件>][-<设备编号><存储密度>][--after-date=<日期时间>][--atime-preserve][--backuup=<备份方式>][--checkpoint][--concatenate][--confirmation][--delete][--exclude=<范本样式>][--force-local][--group=<群组名称>][--help][--ignore-failed-read][--new-volume-script=<Script文件>][--newer-mtime][--no-recursion][--null][--numeric-owner][--owner=<用户名称>][--posix][--erve][--preserve-order][--preserve-permissions][--record-size=<区块数目>][--recursive-unlink][--remove-files][--rsh-command=<执行指令>][--same-owner][--suffix=<备份字尾字符串>][--totals][--use-compress-program=<执行指令>][--version][--volno-file=<编号文件>][文件或目录...]

```

​           **命令参数：**

| 参数 | 参数描述                       |
| ---- | ------------------------------ |
| -c   | 建立新的压缩文件               |
| -f   | 定压缩文件                     |
| -r   | 添加文件到已经压缩文件包中     |
| -u   | 添加改了和现有的文件到压缩包中 |
| -x   | 从压缩包中抽取文件             |
| -t   | 显示压缩文件中的内容           |
| -z   | 支持gzip压缩                   |
| -j   | 支持bzip2压缩                  |
| -Z   | 支持compress解压文件           |
| -v   | 显示操作过程                   |

​           有关 gzip 及 bzip2 压缩:

```bash
gzip 实例：压缩 gzip fileName .tar.gz 和.tgz  解压：gunzip filename.gz 或 gzip -d filename.gz
          对应：tar zcvf filename.tar.gz     tar zxvf filename.tar.gz
```

```bash
bz2实例：压缩 bzip2 -z filename .tar.bz2 解压：bunzip filename.bz2或bzip -d filename.bz2
       对应：tar jcvf filename.tar.gz         解压：tar jxvf filename.tar.bz2
```

​            （1）将test.log  test.sh全部打包成 tar 包

```bash
[root@localhost ~]# [root@localhost ~]# tar -cvf log.tar test.log  test.sh
test.log
test.sh

```

​            （2）将 /etc 下的所有文件及目录打包到指定目录或当前目录，并使用 gz 压缩

```bash
[root@localhost ~]# tar -zcvf ./etc.tar.gz /etc
```

​            （3）查看刚打包的文件内容（一定加z，因为是使用 gzip 压缩的）

```bash
[root@localhost ~]# tar -ztvf ./etc.tar.gz
...省略...
```

​            （4）要压缩打包 /home, /etc ，但不要 /home/mashibing ，只能针对文件，不能针对目录

```bash
[root@localhost ~]# tar --exclude /home/mshibing -zcvf myfile.tar.gz /home/* /etc
```

#### 

#### 9、date命令

​           **作用：**用来显示或设定系统的日期与时间

​           **语法：**

```
date [-u] [-d datestr] [-s datestr] [--utc] [--universal] [--date=datestr] [--set=datestr] [--help] [--version] [+FORMAT] [MMDDhhmm[[CC]YY][.ss]]
```

​           **时间参数**

| 参数 | 描述参数                                           |
| ---- | -------------------------------------------------- |
| %    | 印出 %                                             |
| %n   | 下一行                                             |
| %t   | 跳格                                               |
| %H   | 小时(00..23)                                       |
| %I   | 小时(01..12)                                       |
| %k   | 小时(0..23)                                        |
| %l   | 小时(1..12)                                        |
| %M   | 分钟(00..59)                                       |
| %p   | 显示本地 AM 或 PM                                  |
| %r   | 直接显示时间 (12 小时制，格式为 hh:mm:ss [AP]M)    |
| %s   | 从 1970 年 1 月 1 日 00:00:00 UTC 到目前为止的秒数 |
| %S   | 秒(00..61)                                         |
| %T   | 直接显示时间 (24 小时制)                           |
| %X   | 相当于 %H:%M:%S                                    |
| %Z   | 显示时区                                           |

​           **日期参数**

| 参数 | 描述参数                                                 |
| ---- | -------------------------------------------------------- |
| %a   | 星期几 (Sun..Sat)                                        |
| %A   | 星期几 (Sunday..Saturday)                                |
| %b   | 月份 (Jan..Dec)                                          |
| %B   | 月份 (January..December)                                 |
| %c   | 直接显示日期与时间                                       |
| %d   | 日 (01..31)                                              |
| %D   | 直接显示日期 (mm/dd/yy)                                  |
| %h   | 同 %b                                                    |
| %j   | 一年中的第几天 (001..366)                                |
| %m   | 月份 (01..12)                                            |
| %U   | 一年中的第几周 (00..53) (以 Sunday 为一周的第一天的情形) |
| %w   | 一周中的第几天 (0..6)                                    |
| %W   | 一年中的第几周 (00..53) (以 Monday 为一周的第一天的情形) |
| %x   | 直接显示日期 (mm/dd/yy)                                  |
| %y   | 年份的最后两位数字 (00.99)                               |
| %Y   | 完整年份 (0000..9999)                                    |

```
若是不以加号作为开头，则表示要设定时间，而时间格式为 MMDDhhmm[[CC]YY][.ss]，其中 MM 为月份，DD 为日，hh 为小时，mm 为分钟，CC 为年份前两位数字，YY 为年份后两位数字，ss 为秒数。

使用权限：所有使用者。

当您不希望出现无意义的 0 时(比如说 1999/03/07)，则可以在标记中插入 - 符号，比如说 date '+%-H:%-M:%-S' 会把时分秒中无意义的 0 给去掉，像是原本的 08:09:04 会变为 8:9:4。另外，只有取得权限者(比如说 root)才能设定系统时间。

当您以 root 身分更改了系统时间之后，请记得以 clock -w 来将系统时间写入 CMOS 中，这样下次重新开机时系统时间才会持续抱持最新的正确值。
```

​           **语法：**

```
date [-u] [-d datestr] [-s datestr] [--utc] [--universal] [--date=datestr] [--set=datestr] [--help] [--version] [+FORMAT] [MMDDhhmm[[CC]YY][.ss]]
```

​           常见参数

|           |                                          |
| --------- | ---------------------------------------- |
| -d        | 显示 datestr 中所设定的时间 (非系统时间) |
| --help    | 显示辅助讯息                             |
| -s        | 将系统时间设为 datestr 中所设定的时间    |
| -u        | 显示目前的格林威治时间                   |
| --version | 显示版本编号                             |

​            （1）

#### 10、cal 命令

​           **作用：**用户显示公历（阳历）日历

​           **语法：**cal [选项] [[[日] 月] 年]

| 参数 | 参数描述                 |
| ---- | ------------------------ |
| -1   | 只显示当前月份(默认)     |
| -3   | 显示上个月、当月和下个月 |
| -s   | 周日作为一周第一天       |
| -m   | 周一用为一周第一天       |
| -j   | 输出儒略日               |
| -y   | 输出整年                 |
| -V   | 显示版本信息并退出       |
| -h   | 显示此帮助并退出         |

​            （1）显示指定年月日期

```bash
[root@localhost ~]# cal 9 2020
```

​            （2）显示2020年每个月日历

```bash
[root@localhost ~]# cal -y 2020
```

​            （3）将星期一做为第一列,显示前中后三月

```bash
[root@localhost ~]# cal -3m
```



#### 11、grep命令

**grep （文本搜索工具）**

​       强大的文本搜索命令，grep(Global Regular Expression Print) **全局正则表达式搜索**

​       grep 的工作方式是这样的，它在一个或多个文件中搜索字符串模板。如果模板包括空格，则必须被引用，模板后的所有字符串被看作文件名。搜索的结果被送到标准输出，不影响原文件内容

​           **作用：**用于查找文件里符合条件的字符串

```
注意：如果发现某文件的内容符合所指定的范本样式，预设 grep 指令会把含有范本样式的那一列显示出来。若不指定任何文件名称，或是所给予的文件名为 **-**，则 grep 指令会从标准输入设备读取数据
```

​           **语法：**

```
grep [option] pattern file|dir
```

​           **常用参数：**

| 参数 | 参数描述                               |
| ---- | -------------------------------------- |
| -A n | 显示匹配字符后n行                      |
| -B n | 显示匹配字符前n行                      |
| -C n | 显示匹配字符前后n行                    |
| -c   | 计算符合样式的列数                     |
| -i   | 忽略大小写                             |
| -l   | 只列出文件内容符合指定的样式的文件名称 |
| -f   | 从文件中读取关键词                     |
| -n   | 显示匹配内容的所在文件中行数           |
| -R   | 递归查找文件夹                         |

​           **grep 的规则表达式**

```
^               #锚定行的开始 如：'^grep'匹配所有以grep开头的行。 
$               #锚定行的结束 如：'grep$'匹配所有以grep结尾的行。 
.               #匹配一个非换行符的字符 如：'gr.p'匹配gr后接一个任意字符，然后是p。  
*               #匹配零个或多个先前字符 如：'*grep'匹配所有一个或多个空格后紧跟grep的行。
.*              #一起用代表任意字符。  
[]              #匹配一个指定范围内的字符，如'[Gg]rep'匹配Grep和grep。 
[^]             #匹配一个不在指定范围内的字符，如：'[^A-FH-Z]rep'匹配不包含A-R和T-Z的一个字母开头，紧跟rep的行。  
\(..\)          #标记匹配字符，如'\(love\)'，love被标记为1。   
\<              #锚定单词的开始，如:'\<grep'匹配包含以grep开头的单词的行。
\>              #锚定单词的结束，如'grep\>'匹配包含以grep结尾的单词的行。
x\{m\}          #重复字符x，m次，如：'0\{5\}'匹配包含5个o的行。 
x\{m,\}         #重复字符x,至少m次，如：'o\{5,\}'匹配至少有5个o的行。  
x\{m,n\}        #重复字符x，至少m次，不多于n次，如：'o\{5,10\}'匹配5--10个o的行。  
\w              #匹配文字和数字字符，也就是[A-Za-z0-9]，如：'G\w*p'匹配以G后跟零个或多个文字或数字字符，然后是p。  
\W              #\w的反置形式，匹配一个或多个非单词字符，如点号句号等。  
\b              #单词锁定符，如: '\bgrep\b'只匹配grep。
```

​            （1）查找指定进程

```bash
[root@localhost ~]# ps -ef | grep svn
root       6771   9267  0 15:17 pts/0    00:00:00 grep --color=auto svn
```

​            （2）查找指定进程个数

```bash
[root@localhost ~]# ps -ef | grep svn -c
1

```

​            （3）从文件中读取关键词

```bash
[root@localhost ~]# cat test.log | grep -f test.log
马士兵教育：www.mashibing.com

```

​            （4）从文件夹中递归查找以.sh结尾的行，并只列出文件

```bash
[root@localhost ~]# grep -lR '.sh$'
.bash_history
test.sh
.viminfo
log.tar

```

​            （5）查找非x开关的行内容

```bash
[root@localhost ~]# grep '^[^x]' test.log
马士兵教育：www.mashibing.com

```

​            （6）显示包含 ed 或者 at 字符的内容行

```bash
[root@localhost ~]# grep -E 'ed|at' test.log
```



#### 12、ps命令

ps （显示系统中的进程信息）

​           **作用：**用于显示当前进程 (process) 的状态

​           **语法：**

```
ps [options] [--help]
```

```
linux上进程有5种状态:

1. 运行(正在运行或在运行队列中等待)
2. 中断(休眠中, 受阻, 在等待某个条件的形成或接受到信号)
3. 不可中断(收到信号不唤醒和不可运行, 进程必须等待直到有中断发生)
4. 僵死(进程已终止, 但进程描述符存在, 直到父进程调用wait4()系统调用后释放)
5. 停止(进程收到SIGSTOP, SIGSTP, SIGTIN, SIGTOU信号后停止运行运行)
```

​      ps 工具标识进程的5种状态码:

```
D 不可中断 uninterruptible sleep (usually IO)
R 运行 runnable (on run queue)
S 中断 sleeping
T 停止 traced or stopped
Z 僵死 a defunct (”zombie”) process
```

​           **命令参数**

| 参数 | 参数描述                   |
| ---- | -------------------------- |
| -A   | 显示所有进程               |
| a    | 显示所有进程               |
| -a   | 显示同一终端下所有进程     |
| c    | 显示进程真实名称           |
| e    | 显示环境变量               |
| f    | 显示进程间的关系           |
| r    | 显示当前终端运行的进程     |
| -aux | 显示所有包含其它使用的进程 |

​            （1）显示当前所有进程环境变量及进程间关系

```bash
[root@localhost ~]# ps -ef
```

​            （2）显示当前所有进程

```bash
[root@localhost ~]# ps -A
```

​            （3）与grep联用查找某进程

```bash
[root@localhost ~]# ps -aux | grep apache
root      20112  0.0  0.0 112824   980 pts/0    S+   15:30   0:00 grep --color=auto apache
```

​            （4）找出与 cron 与 syslog 这两个服务有关的 PID 号码

```bash
[root@localhost ~]# ps aux | grep '(cron|syslog)'
root      20454  0.0  0.0 112824   984 pts/0    S+   15:30   0:00 grep --color=auto (cron|syslog)

```



#### 13、kill命令

kill

kill 命令用于删除执行中的程序或工作

```
kill 可将指定的信息送至程序。预设的信息为 SIGTERM(15)，可将指定程序终止。若仍无法终止该程序，可使用 SIGKILL(9) 信息尝试强制删除程序。程序或工作的编号可利用 ps 指令或 jobs 指令查看
```

​           **语法：**

```
kill [-s <信息名称或编号>][程序]　或　kill [-l <信息编号>]
```

​           **常用参数：**

| 参数 | 参数描述                                                     |
| ---- | ------------------------------------------------------------ |
| -l   | 信号，若果不加信号的编号参数，则使用“-l”参数会列出全部的信号名称 |
| -a   | 当处理当前进程时，不限制命令名和进程号的对应关系             |
| -p   | 指定kill 命令只打印相关进程的进程号，而不发送任何信号        |
| -s   | 指定发送信号                                                 |
| -u   | 指定用户                                                     |

​            （1）先使用ps查找进程pro1，然后用kill杀掉

```bash
[root@localhost ~]# kill -9 $(ps -ef | grep pro1)
-bash: kill: root: 参数必须是进程或任务 ID
-bash: kill: (27319) - 没有那个进程
-bash: kill: (27317) - 没有那个进程
```



```

```