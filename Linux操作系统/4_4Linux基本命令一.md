# Python 全栈文档

## 第四章  Linux基本命令一

#### 1、mkdir 命令

**mkdir**

​           **作用：**命令用来创建指定的名称的目录，要求创建目录的用户在当前目录中具有写权限，并且指定的目录名不能是当前目录中已有的目录

​           **语法：**mkdir   [选项]   目录

​           **命令功能：**通过 mkdir 命令可以实现在指定位置创建以 DirName(指定的文件名)命名的文件夹或目录。要创建文件夹或目录的用户必须对所创建的文件夹的父文件夹具有写权限。并且，所创建的文件夹(目录)不能与其父目录(即父文件夹)中的文件名重名，即同一个目录下不能有同名的(区分大小写)

​           **命令参数：**

| 选项参数 | 完整参数                             | 功能描述                                                     |
| -------- | ------------------------------------ | ------------------------------------------------------------ |
| -m       | --mode=模式                          | 设定权限<模式> (类似 chmod)，而不是 rwxrwxrwx 减 umask       |
| -p       | --parents                            | 可以是一个路径名称。<br/>此时若路径中的某些目录尚不存在,加上此选项后,<br/>系统将自动建立好那些尚不存在的目录,即一次可以建立多个目录; |
| -v       | --verbose<br />--help<br />--version | --verbose 每次创建新目录都显示信息<br/>--help显示此帮助信息并退出<br/>--version输出版本信息并退出 |

​           **示例如下：**

```bash
# 创建目录                        正常使用无参数
[root@localhost text.txt]# mkdir mydir
[root@localhost text.txt]# ll
总用量 0
drwxr-xr-x. 2 root root 6 9月   7 13:06 mydir

# 切换至mydir目录
[root@localhost text.txt]# cd mydir/

# 参数-p 进行递归创建目录          使用参数 -p
[root@localhost mydir]# mkdir -p text1/text2

# 切换至text2目录
[root@localhost mydir]# cd text1/text2/

[root@localhost text2]# pwd
/root/mydir/text1/text2

总结：# 使用参数 -p 可以创建多层路径文件目录





# 分别创建三个目录，并设置权限                  使用参数 -m
[root@localhost text2]# mkdir -m 777  exercise1
[root@localhost text2]# mkdir -m 765  exercise2
[root@localhost text2]# mkdir -m 654  exercise3

[root@localhost text2]# ll
总用量 0
drwxrwxrwx. 2 root root 6 9月   7 13:22 exercise1
drwxrw-r-x. 2 root root 6 9月   7 13:23 exercise2
drw-r-xr--. 2 root root 6 9月   7 13:23 exercise3

总结：# 使用参数 -m 可以创建文件并设置文件的权限  
注意(777，765，654 其中每一个数字，分别表示User、Group、及Other的权限。r=4，w=2，x=1)




# 创建一个目录，并设置权限                  使用参数 -v   -m
[root@localhost text2]# mkdir -v -m 654  exercise4
mkdir: 已创建目录 "exercise4"

```

#### 2、touch 命令

**touch**

主要用于 **创建空文件** 或 **更新文件的时间戳**。

​           **作用：**用于修改文件或者目录的时间属性，包括存取时间和更改时间。若文件不存在，系统会建立一个新的文件。

​                      ls -l 可以显示档案的时间记录。

​           **语法：**touch \[-acfm] \[-d<日期时间>] \[-r<参考文件或目录>]  \[-t<日期时间>]  \[--help]  \[--version]  \[文件或目录…]

​           **命令参数：**

| 参数        | 参数描述                                                  |
| ----------- | --------------------------------------------------------- |
| -a          | 只更新访问时间，不改变修改时间                            |
| -m          | 改变修改时间记录                                          |
| -c          | 不创建不存在的文件                                        |
| -f          | 不使用，是为了与其他 unix 系统的相容性而保留。            |
| -m          | 只更新修改时间，不改变访问时间                            |
| -r file     | 使用文件file的时间更新文件的时间                          |
| -t          | 将时间修改为参数指定的日期,如：07081556代表7月8号15点56分 |
| --no-create | 不会建立文件                                              |
| --help      | 列出指令格式                                              |
| --version   | 列出版本讯息                                              |

​            **示例如下：**

​            （1）使用指令"touch"创建"text1.txt"文件

```bash
[root@localhost text2]# cd exercise1
[root@localhost exercise1]# touch text1.txt
[root@localhost exercise1]# ls -l
总用量 0
-rw-r--r--. 1 root root 0 9月   7 13:44 text1.txt
# 当然touch可以一次创建多个文件，例如：touch text1.txt  text2.txt text3.txt ...
[root@localhost exercise1]# touch text2.txt text3.txt
[root@localhost exercise1]# ls -l
总用量 0
-rw-r--r--. 1 root root 0 9月   7 13:44 text1.txt
-rw-r--r--. 1 root root 0 9月   7 13:45 text2.txt
-rw-r--r--. 1 root root 0 9月   7 13:45 text3.txt



```

​            （2）使用指令"touch"修改文件"text1.txt"的时间属性为当前系统时间，发现如果文件存在，这里直接是修改时间了

```bash
[root@localhost exercise1]# touch text1.txt
[root@localhost exercise1]# ls -l
总用量 0
-rw-r--r--. 1 root root 0 9月   7 13:46 text1.txt
# 发现时间是已经修改了
```

​            （3）强制避免使用touch命令创建新文件   使用参数 -c

​                     有时，如果新文件不存在，则需要避免创建新文件。 在这种情况下，可以使用touch命令使用'-c'选项

```bash
[root@localhost exercise1]# ls
text1.txt  text2.txt  text3.txt
[root@localhost exercise1]# touch -c text4.txt
[root@localhost exercise1]# ls
text1.txt  text2.txt  text3.txt

```

​            （4）更改文件的访问和修改时间         使用参数 -a

​                     touch命令的另一个用途是更改文件的访问时间和修改时间

```bash
[root@localhost exercise1]# stat text3.txt 
  文件："text3.txt"
  大小：0               块：0          IO 块：4096   普通空文件
设备：fd00h/64768d      Inode：16626       硬链接：1
权限：(0644/-rw-r--r--)  Uid：(    0/    root)   Gid：(    0/    root)
环境：unconfined_u:object_r:admin_home_t:s0
最近访问：2020-09-07 14:04:07.539848285 +0800
最近更改：2020-09-07 14:04:07.539848285 +0800
最近改动：2020-09-07 14:04:07.539848285 +0800
创建时间：-

[root@localhost exercise1]# touch -a text3.txt 

[root@localhost exercise1]# stat text3.txt 
  文件："text3.txt"
  大小：0               块：0          IO 块：4096   普通空文件
设备：fd00h/64768d      Inode：16626       硬链接：1
权限：(0644/-rw-r--r--)  Uid：(    0/    root)   Gid：(    0/    root)
环境：unconfined_u:object_r:admin_home_t:s0
最近访问：2020-09-07 14:08:33.788865586 +0800
最近更改：2020-09-07 14:04:07.539848285 +0800
最近改动：2020-09-07 14:08:33.788865586 +0800
创建时间：-

我们看到先是使用stat查看文件详细信息  最近访问和最近更改时间
当我们使用命令使用参数-a 对文件 text3.txt 做了一些操作 touch -a text3.txt 
再查看时间时，发现文件已经修改好了
```

​            （4）仅更改此文件的修改时间          使用参数 -m

```bash

[root@localhost exercise1]# touch -m text3.txt 
[root@localhost exercise1]# stat text3.txt 
  文件："text3.txt"
  大小：0               块：0          IO 块：4096   普通空文件
设备：fd00h/64768d      Inode：16626       硬链接：1
权限：(0644/-rw-r--r--)  Uid：(    0/    root)   Gid：(    0/    root)
环境：unconfined_u:object_r:admin_home_t:s0
最近访问：2020-09-07 14:08:33.788865586 +0800
最近更改：2020-09-07 14:15:59.782894567 +0800
最近改动：2020-09-07 14:15:59.782894567 +0800
创建时间：-

```

​     （5）将访问和修改时间从一个文件复制到另一个文件   使用参数 -r

```bash
[root@localhost exercise1]# touch text3.txt  -r  text1.txt 

[root@localhost exercise1]# stat text1.txt 
  文件："text1.txt"
  大小：0               块：0          IO 块：4096   普通空文件
设备：fd00h/64768d      Inode：16624       硬链接：1
权限：(0644/-rw-r--r--)  Uid：(    0/    root)   Gid：(    0/    root)
环境：unconfined_u:object_r:admin_home_t:s0
最近访问：2020-09-07 13:57:38.308822993 +0800
最近更改：2020-09-07 13:57:38.308822993 +0800
最近改动：2020-09-07 13:57:38.308822993 +0800
创建时间：-

[root@localhost exercise1]# stat text3.txt 
  文件："text3.txt"
  大小：0               块：0          IO 块：4096   普通空文件
设备：fd00h/64768d      Inode：16626       硬链接：1
权限：(0644/-rw-r--r--)  Uid：(    0/    root)   Gid：(    0/    root)
环境：unconfined_u:object_r:admin_home_t:s0
最近访问：2020-09-07 13:57:38.308822993 +0800
最近更改：2020-09-07 13:57:38.308822993 +0800
最近改动：2020-09-07 14:19:35.732908600 +0800
创建时间：-

# 输出显示text3.txt现在具有与text1.txt相同的访问和修改值
```

​            （6）使用指定的时间戳创建新文件         使用参数-t

```bash
[root@localhost exercise1]# touch -t 2001011314.52 time.log

[root@localhost exercise1]# stat time.log 
  文件："time.log"
  大小：0               块：0          IO 块：4096   普通空文件
设备：fd00h/64768d      Inode：16627       硬链接：1
权限：(0644/-rw-r--r--)  Uid：(    0/    root)   Gid：(    0/    root)
环境：unconfined_u:object_r:admin_home_t:s0
最近访问：2020-01-01 13:14:52.000000000 +0800
最近更改：2020-01-01 13:14:52.000000000 +0800
最近改动：2020-09-07 14:26:23.916935124 +0800
创建时间：-
# 最近访问与最近更该时间为设定的时间 2020-01-01 13:14:52.000000000
```

​            （7）将文件的时间戳更改为其他时间     使用参数

```bash
[root@localhost exercise1]# touch -c -t 1801011314.52 time.log

[root@localhost exercise1]# stat time.log 
  文件："time.log"
  大小：0               块：0          IO 块：4096   普通空文件
设备：fd00h/64768d      Inode：16627       硬链接：1
权限：(0644/-rw-r--r--)  Uid：(    0/    root)   Gid：(    0/    root)
环境：unconfined_u:object_r:admin_home_t:s0
最近访问：2018-01-01 13:14:52.000000000 +0800
最近更改：2018-01-01 13:14:52.000000000 +0800
最近改动：2020-09-07 14:30:45.188952101 +0800
创建时间：-
# 使用参数修改为指定的时间戳
```

#### 3、rm 命令

**rm**

​           **作用：**用于删除一个文件或者目录。

​           **语法：**rm [选项] 文件…

​           **命令参数：**

| 参数 | 参数描述                                         |
| ---- | ------------------------------------------------ |
| -i   | 删除前逐一询问确认                               |
| -f   | 即使原档案属性设为唯读，亦直接删除，无需逐一确认 |
| -r   | 将目录及以下之档案亦逐一删除  (递归删除)         |

​            （1）删除文件或者目录前提示  使用参数  -i

```bash
[root@localhost exercise1]# ls
text1.txt  text2.txt  time.log

[root@localhost exercise1]# rm -i text2.txt
rm：是否删除普通空文件 "text2.txt"？y
[root@localhost exercise1]# ls
text1.txt  time.log

# 这里提示是否删除，输入 y确认
```

​            （2）删除 test 子目录及子目录中所有档案删除，并且不用一一确认   -rf

​                    我们先切换到上一级目录

```bash
[root@localhost exercise1]# cd ../

[root@localhost text2]# ls
exercise1  exercise2  exercise3  exercise4
[root@localhost text2]# rm -rf exercise1
[root@localhost text2]# ls
exercise2  exercise3  exercise4
[root@localhost text2]# 

```

#### 4、rmdir 命令

rmdir

​           **作用：**用于删除空的目录。从一个目录中删除一个或多个子目录项，删除某目录时也必须具有对其父目录的写权限。

​           **语法：**rmdir [-p] dirName

​           **命令参数：**

| 参数 | 参数描述                                               |
| ---- | ------------------------------------------------------ |
| -p   | 当子目录被删除后使它也成为空目录的话，则顺便一并删除。 |

​           注意（不能删除非空目录）

​            （1）正常删除目录

```bash
[root@localhost text2]# ls
exercise2  exercise3  exercise4

[root@localhost text2]# rmdir exercise2
[root@localhost text2]# ls
exercise3  exercise4

[root@localhost text2]# rmdir exercise3
[root@localhost text2]# ls
exercise4

# 这里删除空目录
```

​            （1）工作目录下的 text1 目录中，删除名为 text2 的子目录。若 text2 删除后，text1 目录成为空目录，则 text1 亦予删除，类推

```bash
[root@localhost ~]# ls -R mydir/
mydir/:
text1

mydir/text1:
text2

mydir/text1/text2:


[root@localhost ~]# rmdir -p mydir/text1/text2/
[root@localhost ~]# ls
anaconda-ks.cfg

```

#### 5、mv 命令

mv

​           **作用：**用来为文件或目录改名、或将文件或目录移入其它位置。

​           **语法：**

```
mv [options] source dest
mv [options] source... directory
```

​           **命令参数：**

| 参数 | 参数描述                                        |
| ---- | ----------------------------------------------- |
| -i   | 若指定目录已有同名文件，则先询问是否覆盖旧文件; |
| -f   | 在 mv 操作要覆盖某已有的目标文件时不给任何指示  |

​           mv参数设置与运行结果

| 命令格式         | 运行结果                                                     |
| :--------------- | :----------------------------------------------------------- |
| mv 文件名 文件名 | 将源文件名改为目标文件名                                     |
| mv 文件名 目录名 | 将文件移动到目标目录                                         |
| mv 目录名 目录名 | 目标目录已存在，将源目录移动到目标目录；目标目录不存在则改名 |
| mv 目录名 文件名 | 出错                                                         |

​            先给个环境

```bash
[root@localhost ~]# ls -R mydir/
mydir/:
test1  test2

mydir/test1:
text1.log  text1.txt  text2.log

mydir/test2:

```

​            （1）将文件 text1.log 重命名为 text2.txt

```bash
[root@localhost test1]# ls
text1.log  text1.txt  text2.log

[root@localhost test1]# mv text1.log text2.txt
[root@localhost test1]# ls
text1.txt  text2.log  text2.txt


```

​            （2）将文件 text1.txt    text2.log    text2.txt   移动到mydir的 test2 目录中

```bash
[root@localhost test1]# mv text1.txt  text2.log  text2.txt  ../test2
[root@localhost test1]# ls
[root@localhost test1]# cd ../test2
[root@localhost test2]# ls
text1.txt  text2.log  text2.txt

```

​            （3）将文件 file1 改名为 file2，如果 file2 已经存在，则询问是否覆盖  使用参数 -i 询问

```bash
[root@localhost test2]# ls
text1.txt  text2.log  text2.txt

[root@localhost test2]# mv -i text2.txt text1.txt 
mv：是否覆盖"text1.txt"？ y
[root@localhost test2]# ls
text1.txt  text2.log
[root@localhost test2]# 

```

​            （4）移动当前文件夹下的所有文件到上一级目录

```bash
[root@localhost test2]# mv * ../
[root@localhost test2]# ls ../
test1  test2  text1.txt  text2.log

```

#### 6、cp 命令

cp

​           **作用：**用于复制文件或目录。

​           **语法：**

```bash
cp [options] source dest

cp [options] source... directory
```

​           **命令参数：**

| 参数 | 参数功能描述                                                 |
| ---- | ------------------------------------------------------------ |
| -a   | 此选项通常在复制目录时使用，它保留链接、文件属性，并复制目录下的所有内容。其作用等于dpR参数组合 |
| -d   | 复制时保留链接。这里所说的链接相当于Windows系统中的快捷方式。 |
| -f   | 覆盖已经存在的目标文件而不给出提示。                         |
| -i   | 与-f选项相反，在覆盖目标文件之前给出提示，要求用户确认是否覆盖，回答"y"时目标文件将被覆盖。 |
| -p   | 除复制文件的内容外，还把修改时间和访问权限也复制到新文件中。 |
| -r   | 若给出的源文件是一个目录文件，此时将复制该目录下所有的子目录和文件。 |
| -l   | 不复制文件，只是生成链接文件。                               |

```bash
# 常用的一些有以下参数
-i 提示
-r 复制目录及目录内所有项目
-a 复制的文件与原文件时间一样
```

```bash
[root@localhost ~]# ls -R mydir/
mydir/:
test1  test2  text1.txt  text2.log

mydir/test1:

mydir/test2:

# 目录文件环境
```

​            （1）复制 text1.txt 到 test1 目录下，保持原文件时间，如果原文件存在提示是否覆盖。使用参数 -ai

```bash
[root@localhost ~]# cp -ai mydir/text1.txt  mydir/test1
[root@localhost ~]# ls mydir/test1
text1.txt

```

​            （2）为 text1.txt 建议一个链接（快捷方式） 使用参数 -s

```bash
[root@localhost ~]# cp -s mydir/text2.log  link_text2
[root@localhost ~]# ls
anaconda-ks.cfg  link_text2  mydir

# 注意（只能于当前目录中创建相对的符号链接）
```

​            （3）将当前或者制定的目录下的所有文件复制到新目录

```bash
[root@localhost ~]# cp -r mydir/ test

[root@localhost ~]# ls -R test
test:
test1  test2  text1.txt  text2.log

test/test1:
text1.txt

test/test2:
test1  test2

test/test2/test1:
text1.txt

test/test2/test2:


[root@localhost ~]# ls -R mydir/
mydir/:
test1  test2  text1.txt  text2.log

mydir/test1:
text1.txt

mydir/test2:
test1  test2

mydir/test2/test1:
text1.txt

mydir/test2/test2:

# 两个文件一模一样的内容

```

#### 7、cat命令

cat

​           **作用：**用于连接文件并打印到标准输出设备上

​           **语法：**cat [-AbeEnstTuv] [--help] [--version] fileName

​           **命令参数：**

|        |                        |                                                    |
| ------ | ---------------------- | -------------------------------------------------- |
| **-n** | **--number**           | 由 1 开始对所有输出的行数编号                      |
| **-b** | **--number-nonblank**  | 和 -n 相似，只不过对于空白行不编号                 |
| **-s** | **--squeeze-blank**    | 当遇到有连续两行以上的空白行，就代换为一行的空白行 |
| **-v** | **--show-nonprinting** | 使用 ^ 和 M- 符号，除了 LFD 和 TAB 之外            |
| **-E** | **--show-ends**        | 在每行结束处显示 $                                 |
| **-T** | **--show-tabs**        | 将 TAB 字符显示为 ^I                               |
| **-A** | **--show-all**         | 等价于 -vET                                        |
| **-e** |                        | 等价于"-vE"选项                                    |
| **-t** |                        | 等价于"-vT"选项                                    |

```bash
参考环境，事先编辑好的文本

[root@localhost ~]# vi mydir/text1.txt 

# 编辑文本，使用vi打开文本
# 按i键进入编辑界面
# 输入对应的内容后
# 按ESC退出编辑模式
# 按住shift+；进入命令行界面
# 输入wq回车接口保存退出


```

​            （1）**显示整个文件内容 使用 cat 正常进查看**

```bash
[root@localhost ~]# cat mydir/text1.txt
# 编辑文本，使用vi打开文本
# 按i键进入编辑界面
# 输入对应的内容后
# 按ESC退出编辑模式
# 按住shift+；进入命令行界面
# 输入wq回车接口保存退出
```

​            （2）把   text1.txt   的文档内容加上行号后输入 text2.log  这个文档里：

```bash
[root@localhost mydir]# cat -n text1.txt > text2.log 
[root@localhost mydir]# cat text2.log 
     1  # 编辑文本，使用vi打开文本
     2  # 按i键进入编辑界面
     3  # 输入对应的内容后
     4  # 按ESC退出编辑模式
     5  # 按住shift+；进入命令行界面
     6  # 输入wq回车接口保存退出


```

​            （3）把 text1.txt 和 text2.txt 的文档内容加上行号（空白行不加）之后将内容附加到 text3.txt 文档里

```bash
[root@localhost mydir]# cat -b text1.txt  text2.log >> text3.txt
[root@localhost mydir]# cat text3.txt
     1  # 编辑文本，使用vi打开文本
     2  # 按i键进入编辑界面
     3  # 输入对应的内容后
     4  # 按ESC退出编辑模式
     5  # 按住shift+；进入命令行界面
     6  # 输入wq回车接口保存退出
     7       1  # 编辑文本，使用vi打开文本
     8       2  # 按i键进入编辑界面
     9       3  # 输入对应的内容后
    10       4  # 按ESC退出编辑模式
    11       5  # 按住shift+；进入命令行界面
    12       6  # 输入wq回车接口保存退出
```

​            （3）清空 /mydir/text1.txt 文档内容

```bash
[root@localhost mydir]# cat /dev/null > text1.txt
[root@localhost mydir]# cat text1.txt
[root@localhost mydir]# ls
test1  test2  text1.txt  text2.log  text3.txt

```

#### 8、more 命令

more

​           **作用：**类似 cat ，**不过会以一页一页的形式显示**，更方便使用者逐页阅读，而最基本的指令就是按空白键（space）就往下一页显示，按 b 键就会往回（back）一页显示，而且还有搜寻字串的功能（与 vi 相似），使用中的说明文件，请按 h 

​           **语法：**more [-dlfpcsu] [-num] [+/pattern] [+linenum] [fileNames..]

​           **命令参数：**

| 参数      | 参数功能描述                                                 |
| --------- | ------------------------------------------------------------ |
| -num      | 一次显示的行数                                               |
| -d        | 提示使用者，在画面下方显示 [Press space to continue, 'q' to quit.] ，<br />如果使用者按错键，则会显示 [Press 'h' for instructions.] 而不是 '哔' 声 |
| -l        | 取消遇见特殊字元 ^L（送纸字元）时会暂停的功能                |
| -f        | 计算行数时，以实际上的行数，而非自动换行过后的行数（有些单行字数太长的会被扩展为两行或两行以上） |
| -p        | 不以卷动的方式显示每一页，而是先清除萤幕后再显示内容         |
| -c        | 跟 -p 相似，不同的是先显示内容再清除其他旧资料               |
| -s        | 当遇到有连续两行以上的空白行，就代换为一行的空白行           |
| -u        | 不显示下引号 （根据环境变数 TERM 指定的 terminal 而有所不同） |
| +/pattern | 在每个文档显示前搜寻该字串（pattern），然后从该字串之后开始显示 |
| +num      | 从第 num 行开始显示                                          |
| fileNames | 需要显示内容的文档，可为复数个数                             |

​           常用的操作命令

| 按键   | 按键功能描述                     |
| ------ | -------------------------------- |
| Enter  | 向下 n 行，需要定义。默认为 1 行 |
| Ctrl+F | 向下滚动一屏                     |
| 空格键 | 向下滚动一屏                     |
| Ctrl+B | 返回上一屏                       |
| =      | 输出当前行的行号                 |
| :f     | 输出文件名和当前行的行号         |
| V      | 调用vi编辑器                     |
| !命令  | 调用Shell，并执行命令            |
| q      | 退出more                         |

​            （1）显示文件中从第3行起的内容

```bash
[root@localhost mydir]# more +3 text3.txt
     3  # 输入对应的内容后
     4  # 按ESC退出编辑模式
     5  # 按住shift+；进入命令行界面
     6  # 输入wq回车接口保存退出
     7       1  # 编辑文本，使用vi打开文本
     8       2  # 按i键进入编辑界面
     9       3  # 输入对应的内容后
    10       4  # 按ESC退出编辑模式
    11       5  # 按住shift+；进入命令行界面
    12       6  # 输入wq回车接口保存退出

```

​            （2）在所列出文件目录详细信息，借助管道使每次显示 5 行

```bash
[root@localhost mydir]# ls -l / | more -5
总用量 20
lrwxrwxrwx.   1 root root    7 8月  31 15:48 bin -> usr/bin
dr-xr-xr-x.   5 root root 4096 8月  31 15:58 boot
drwxr-xr-x.  20 root root 3240 9月   5 13:07 dev
drwxr-xr-x.  75 root root 8192 9月   7 10:30 etc
--More--

# 空格会显示下5行
# 回车会显示下1行
```

#### 9、less 命令

less

​           **作用：**less 与 more 类似，但使用 **less 可以随意浏览文件，而 more 仅能向前移动，却不能向后移动**，而且 less 在查看之前不会加载整个文件。

​           **语法：**less [参数] 文件 

​           **命令参数：**

| 参数       | 参数功能描述                                   |
| ---------- | ---------------------------------------------- |
| -i         | 忽略搜索时的大小写                             |
| -N         | 显示每行的行号                                 |
| -o         | <文件名> 将less 输出的内容在指定文件中保存起来 |
| -s         | 显示连续空行为一行                             |
| /字符串：  | 向下搜索“字符串”的功能                         |
| ?字符串：  | 向上搜索“字符串”的功能                         |
| n          | 重复前一个搜索（与 / 或 ? 有关）               |
| N          | 反向重复前一个搜索（与 / 或 ? 有关）           |
| -x <数字>  | 将“tab”键显示为规定的数字空格                  |
| b          | 向后翻一页                                     |
| d          | 向后翻半页                                     |
| h          | 显示帮助界面                                   |
| Q          | 退出less 命令                                  |
| u          | 向前滚动半页                                   |
| y          | 向前滚动一行                                   |
| 空格键     | 滚动一行                                       |
| 回车键     | 滚动一页                                       |
| [pagedown] | 向下翻动一页                                   |
| [pageup]   | 向上翻动一页                                   |

​            （1）查看文件

```bash
[root@localhost mydir]# less text3.txt 
     1  # 编辑文本，使用vi打开文本
     2  # 按i键进入编辑界面
     ...省略...
     
     
# 进入后查看，Q键退出界面
```

​            （2）ps查看进程信息并通过less分页显示

```bash
[root@localhost mydir]# less text3.txt 
UID         PID   PPID  C STIME TTY          TIME CMD
root          1      0  0 05:06 ?        00:00:02 /usr/lib/systemd/systemd --switched-root --system --deserialize 22
root          2      0  0 05:06 ?        00:00:00 [kthreadd]
...省略...
```

​            （3）查看命令历史使用记录并通过less分页显示

```bash
[root@localhost test]# history | less
    1  exit
    2  reboot
    3  shutdowm -r now
    4  poweroff
    5  cd 

...省略...
```

​            （4）查看多个文件

```bash
[root@localhost mydir]# less text3.txt 

# 此时如果需要查看多个文件可以使用 可以输入shift+；进入命令行模式
# 使用 p 和 n 进行上下页面翻页查看
```

```bash
附加备注
1.全屏导航

ctrl + F - 向前移动一屏
ctrl + B - 向后移动一屏
ctrl + D - 向前移动半屏
ctrl + U - 向后移动半屏
2.单行导航

j - 向前移动一行
k - 向后移动一行
3.其它导航

G - 移动到最后一行
g - 移动到第一行
q / ZZ - 退出 less 命令
4.其它有用的命令

v - 使用配置的编辑器编辑当前文件
h - 显示 less 的帮助文档
&pattern - 仅显示匹配模式的行，而不是整个文件
5.标记导航

当使用 less 查看大文件时，可以在任何一个位置作标记，可以通过命令导航到标有特定标记的文本位置：

ma - 使用 a 标记文本的当前位置
'a - 导航到标记 a 处
Linux 命令大全 Linux 命令大全
```

#### cat、more、less比较

`cat` 命令用于将文件内容一次性显示到终端，适合查看小型文件的全部内容。

`more` 命令**按页显示**文件内容，适合查看较长文件。

`less` 命令是功能更强大的文件查看工具，**支持向上和向下滚动。**

