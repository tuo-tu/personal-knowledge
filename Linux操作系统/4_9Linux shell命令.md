# Python 全栈文档

## 第九章  Linux shell命令

#### 1、Shell简介

Shell 脚本是一种用来自动化完成系统管理任务、批量处理操作或运行复杂程序的脚本语言。**它将一系列命令组合到一个文件中**，用户可以通过执行脚本来自动完成这些任务。

**`bash` 是 `sh` 的增强版**，支持更强大的功能，适合现代脚本编写。

​          Shell 是一个用 C 语言编写的程序，它是用户使用 Linux 的桥梁。Shell 既是一种命令语言，又是一种程序设计语言。

​          Shell 是指一种应用程序，这个应用程序提供了一个界面，用户通过这个界面访问操作系统内核的服务。

​          Ken Thompson 的 sh 是第一种 Unix Shell，Windows Explorer 是一个典型的图形界面 Shell。

shell 脚本

​          是一种为 shell 编写的脚本程序。通常说的 shell 通常都是指 shell 脚本，但是注意 shell 和 shell script 是两个不同的概念。

​          Shell 编程跟 JavaScript、php 编程一样，只要有一个能编写代码的文本编辑器和一个能解释执行的脚本解释器就可以了。

​          Linux 的 Shell 种类众多，常见的有：

```bash
- Bourne Shell（/usr/bin/sh或/bin/sh）
- Bourne Again Shell（/bin/bash）


- C Shell（/usr/bin/csh）
- K Shell（/usr/bin/ksh）
- Shell for Root（/sbin/sh）
```

```bash
由于易用和免费，Bash 在日常工作中被广泛使用。
同时，Bash 也是大多数Linux 系统默认的 Shell。

在一般情况下，并不区分 Bourne Shell 和 Bourne Again Shell，所以，像 #!/bin/sh，它同样也可以改为 #!/bin/bash。

#! 告诉系统其后路径所指定的程序即是解释此脚本文件的 Shell 程序。
```

​	**那此时就可以模拟一段编程看看基本的熟悉方法**

​          新建一个文件 test.sh，扩展名为 sh（sh代表shell），扩展名并不影响脚本执行，见名知意就好，

​          如果你用 python 写 shell 脚本，扩展名就用 py 好了。

```bash
[root@localhost ~]# touch test.sh
[root@localhost ~]# vim test.sh


#!/bin/bash
echo 'hello word!'

# 运行方法1：
[root@localhost ~]# ./test.sh
-bash: ./test.sh: 权限不够
[root@localhost ~]# chmod +x ./test.sh
[root@localhost ~]# ./test.sh
hello word!


# 运行方法2：
[root@localhost ~]# /bin/sh test.sh
hello word!

# 但是注意第二种方法的话在，就不需要在第一行指定解释器信息，没有作用了
```

#### 2、Shell变量

##### Shell 脚本中的变量 `$n`

在 Shell 脚本中，`$n` 用于表示传递给脚本或函数的 **位置参数**。位置参数是从命令行传递给脚本或函数的输入参数，`$n` 是其中的第 `n` 个参数。

### **位置参数的具体含义**

| **变量** | **含义**                                                     |
| -------- | :----------------------------------------------------------- |
| `$0`     | 脚本或程序的名称（包含路径，若有）。                         |
| `$1`     | 第一个参数。                                                 |
| `$2`     | 第二个参数。                                                 |
| ...      | 依次类推，表示第 `n` 个参数。                                |
| `$#`     | 参数的个数。                                                 |
| `$*`     | 以一个字符串形式，表示所有参数，参数之间用空格分隔。         |
| `"$*"`   | 以一个字符串形式，表示所有参数，参数之间用空格分隔（整体视为一个参数）。 |
| `$@`     | 以独立的形式，表示所有参数，每个参数都保留为单独的字符串。   |
| `"$@"`   | 保留每个参数的独立性，视为多个独立参数，而非整体视为一个。   |
| `$$`     | 当前脚本或进程的 PID（进程 ID）。                            |
| `$?`     | 上一个命令或脚本的退出状态码（`0` 表示成功，其他表示失败）。 |
| `$_`     | 上一个命令的最后一个参数。                                   |

**定义变量**

```bash
your_name="mashibing"
```

​          注意，变量名和等号之间不能有空格，这和其他的编程语言不太一样。同时，变量名的命名须遵循如下规则：

```bash
- 命名只能使用英文字母，数字和下划线，首个字符不能以数字开头。
- 中间不能有空格，可以使用下划线（_）。
- 不能使用标点符号。
- 不能使用bash里的关键字（可用help命令查看保留关键字）。
```

​          例如下方有效的变量名：

```bash
ZHANGSAN
mashibing
_let
let2
```

​          无效的如下：

```bash
?let=123
user*name=mashibing
```

**使用变量**

​          使用一个定义过的变量，只要在变量名前面加美元符号即可，如：

```bash
your_name="shell学习"
echo $your_name
echo ${your_name}
```

​          变量名外面的花括号是可选的，加不加都行，加花括号是为了帮助解释器识别变量的边界

​          推荐给所有变量加上花括号，这是个好的编程习惯。

​          已定义的变量，可以被重新定义或者重新赋值，如：

```bash
your_name="Tony"
echo $your_name
your_name="alibaba"
echo $your_name
```

​          这样写是合法的，但注意，第二次赋值的时候不能写$your_name="alibaba"，使用变量的时候才加美元符(​\$)。

**只读变量（静态变量？）**

在 Shell 脚本中，所谓的“只读变量”通常指的是 **在脚本运行过程中其值保持不变的变量**。这些变量在脚本的生命周期内是固定的，它们常被用于定义一些固定的配置值或常量，便于维护和复用。

​          使用 readonly 命令可以将变量定义为只读变量，只读变量的值不能被改变。



​          下面的例子尝试更改只读变量，结果报错：

```bash
#!/bin/bash
myUrl="https://www.baidu.com"
readonly myUrl
myUrl="https://www.baidu.com"

[root@localhost ~]#  /bin/sh test.sh
test.sh:行4: myUrl: 只读变量

```

**删除变量**

```bash
语法：unset variable_name  删除变量名

[root@localhost ~]# vim test.sh

#!/bin/sh
myUrl="https://www.baidu.com"
unset myUrl
echo $myUrl

[root@localhost ~]# /bin/sh test.sh

# 注意：以上实例执行将没有任何输出。
```

​          注意（变量被删除后不能再次使用。unset 命令不能删除只读变量。）

​          变量类型

​          运行shell时，会同时存在三种变量：

```
- 1) 局部变量    局部变量在脚本或命令中定义，仅在当前shell实例中有效，其他shell启动的程序不能访问局部变量。
- 2) 环境变量    所有的程序，包括shell启动的程序，都能访问环境变量，有些程序需要环境变量来保证其正常运行。必要的时候shell脚本  也可以定义环境变量。
- 3) shell变量   shell变量是由shell程序设置的特殊变量。shell变量中有一部分是环境变量，有一部分是局部变量，这些变量保证了shell的正常运行
```

`export` 命令用于将一个 Shell 变量导出为环境变量，从而使该变量在当前 Shell 进程及其子进程中都可用。

#### 3、Shell字符串

​          字符串可以用单引号，也可以用双引号，也可以不用引号

​          **单引号**

```bash
str='this is a string'
```

​          单引号字符串的限制：

```
单引号里的任何字符都会原样输出，单引号字符串中的变量是无效的；
单引号字串中不能出现单独一个的单引号（对单引号使用转义符后也不行），但可成对出现，作为字符串拼接使用。
```

​          **双引号**

```bash
your_name='mashibing'
str="Hello, I know you are \"$your_name\"! \n"
echo -e $str
```

​          输出结果为：

```bash
Hello, I know you are "mashibing"! 
```

​          双引号的优点：

```
- 双引号里可以有变量
- 双引号里可以出现转义字符
```

​          **拼接字符串**

```bash
your_name="Tony"
# 使用双引号拼接
greeting="hello, "$your_name" !"
greeting_1="hello, ${your_name} !"
echo $greeting  $greeting_1
# 使用单引号拼接
greeting_2='hello, '$your_name' !'
greeting_3='hello, ${your_name} !'
echo $greeting_2  $greeting_3
```

​          输出结果为：

```bash
hello, Tony ! hello, Tony !
hello, Tony ! hello, ${Tony} !
```

​          **获取字符串长度**

```bash
string="abcd"
echo ${#string} 

#输出 4
```

​          **提取子字符串**

​          以下实例从字符串第 **2** 个字符开始截取 **4** 个字符：

```bash
string="This is a book"
echo ${string:1:4}

# 输出 his 
# 空格也是一个字符
```

​          **注意**：第一个字符的索引值为 **0**。

​          **查找子字符串**

​          查找字符 **i** 或 **o** 的位置(哪个字母先出现就计算哪个)：

```bash
string="This is a book"
echo `expr index "$string" io`  

# 输出 2
```

​          **注意：** 以上脚本中 **`** 是反引号，而不是单引号 **'**，不要看错了哦。

------

#### 4、Shell 数组

​          bash支持一维数组（不支持多维数组），并且没有限定数组的大小。

​          数组元素的下标由 0 开始编号。获取数组中的元素要利用下标，下标可以是整数或算术表达式，其值应大于或等于 0。

​          **定义数组**

​          在 Shell 中，用括号来表示数组，数组元素用"空格"符号分割开。定义数组的一般形式为：

```bash
数组名=(值1 值2 ... 值n)
```

​          例如：

```bash
array_name=(value0 value1 value2 value3)
```

​          或者

```bash
array_name=(
value0
value1
value2
value3
)
```

​          还可以单独定义数组的各个索引位置值：

```bash
array_name[0]=value0
array_name[1]=value1
array_name[n]=valuen
```

​          可以不使用连续的下标，而且下标的范围没有限制。

​          **读取数组**

​          读取数组元素值的一般格式是：

```bash
${数组名[下标]}
```

​          例如：

```bash
valuen=${array_name[n]}
```

​          使用 **@** 符号可以获取数组中的所有元素，例如：

```bash
echo ${array_name[@]}
```

​          **获取数组的长度**

​          获取数组长度的方法与获取字符串长度的方法相同，例如：

```bash
# 取得数组元素的个数
length=${#array_name[@]}
# 或者
length=${#array_name[*]}
# 取得数组单个元素的长度
lengthn=${#array_name[n]}
```

------

#### 5、Shell 注释

​          以 **#** 开头的行就是注释，会被解释器忽略。

​          通过每一行加一个 **#** 号设置多行注释，像这样：

```bash
#--------------------------------------------
# 这是一个注释
# author：马士兵教育
# site：www.mashibing.com
# slogan：学的不仅是技术，更是梦想！
#--------------------------------------------
##### 用户配置区 开始 #####
#
#
# 这里可以添加脚本描述信息
# 
#
##### 用户配置区 结束  #####
```

​          如果在开发过程中，遇到大段的代码需要临时注释起来，过一会儿又取消注释，怎么办呢？

​          每一行加个#符号太费力了，可以把这一段要注释的代码用一对花括号括起来，定义成一个函数，没有地方调用这个函数，这块代码就不会执行，达到了和注释一样的效果。

​          **多行注释**

​          多行注释还可以使用以下格式：

```bash
:<<EOF
注释内容...
注释内容...
注释内容...
EOF
```

​          EOF 也可以使用其他符号:

```bash
:<<'
注释内容...
注释内容...
注释内容...
'

:<<!
注释内容...
注释内容...
注释内容...
!
```

#### 6、Shell 传递参数

​          我们可以在执行 Shell 脚本时，向脚本传递参数，脚本内获取参数的格式为：**$n**。**n** 代表一个数字，1 为执行脚本的第一个参数，2 为执行脚本的第二个参数，以此类推……

​          以下实例我们向脚本传递三个参数，并分别输出，其中 **$0** 为执行的文件名（包含文件路径）：

```bash
#!/bin/bash

echo "Shell 传递参数实例！";
echo "执行的文件名：$0";
echo "第一个参数为：$1";
echo "第二个参数为：$2";
echo "第三个参数为：$3";
```

​          为脚本设置可执行权限，并执行脚本，输出结果如下所示：

```shell
[root@localhost ~]# ./test.sh 1 2 3
Shell 传递参数实例！
执行的文件名：./test.sh
第一个参数为：1
第二个参数为：2
第三个参数为：3

```

​          另外，还有几个特殊字符用来处理参数：

| 参数处理 | 说明                                                         |
| :------- | :----------------------------------------------------------- |
| $#       | 传递到脚本的参数个数                                         |
| $*       | 以一个单字符串显示所有向脚本传递的参数。 如"$*"用「"」括起来的情况、以"$1 $2 … $n"的形式输出所有参数。 |
| $$       | 脚本运行的当前进程ID号                                       |
| $!       | 后台运行的最后一个进程的ID号                                 |
| $@       | 与$*相同，但是使用时加引号，并在引号中返回每个参数。 如"$@"用「"」括起来的情况、以"$1" "$2" … "$n" 的形式输出所有参数。 |
| $-       | 显示Shell使用的当前选项，与[set命令](https://www.runoob.com/linux/linux-comm-set.html)功能相同。 |
| $?       | 显示最后命令的退出状态。0表示没有错误，其他任何值表明有错误。 |

```shell
#!/bin/bash

echo "Shell 传递参数实例！";
echo "第一个参数为：$1";

echo "参数个数为：$#";
echo "传递的参数作为一个字符串显示：$*";
```

​          执行脚本，输出结果如下所示：

```shell
[root@localhost ~]# ./test.sh 1 2 3
Shell 传递参数实例！
第一个参数为：1
参数个数为：3
传递的参数作为一个字符串显示：1 2 3

```

​          \$* 与 $@ 区别：

```
- 相同点：都是引用所有参数。
- 不同点：只有在双引号中体现出来。假设在脚本运行时写了三个参数 1、2、3，，则 " * " 等价于 "1 2 3"（传递了一个参数），而 "@" 等价于 "1" "2" "3"（传递了三个参数）。
```



```shell
#!/bin/bash

echo "-- \$* 演示 ---"
for i in "$*"; do
    echo $i
done

echo "-- \$@ 演示 ---"
for i in "$@"; do
    echo $i
done
```

​          执行脚本，输出结果如下所示：

```shell
[root@localhost ~]# ./test.sh 1 2 3
-- $* 演示 ---
1 2 3
-- $@ 演示 ---
1
2
3

```

#### 7、运算符

​          Shell 和其他编程语言一样，支持多种运算符：

```
算数运算符
关系运算符
布尔运算符
字符串运算符
文件测试运算符
```

​          注意：原生bash是不支持简单的数学运算，但是可以通过其他命令来实现，例如 awk 和 expr，expr 最常用

​                    expr 是一款表达式计算工具，使用它能完成表达式的求值操作

​                    例如，两个数相加(**注意使用的是反引号 ` 而不是单引号 '**)

```shell
[root@localhost ~]# vim test.sh

#!/bin/bash

val=`expr 2 + 2`
echo "两数之和为 : $val"


[root@localhost ~]# ./test.sh
两数之和为 : 4

```

​          注意点：

```shell
表达式和运算符之间要有空格，例如 2+2 是不对的，必须写成 2 + 2，这与我们熟悉的大多数编程语言不一样。

完整的表达式要被 ` ` 包含，注意这个字符不是常用的单引号，在 Esc 键下边。
```

​          **算术运算符**

```
环境如下：变量 a 为 10，变量 b 为 20
```

| 运算符 | 说明                                          | 举例                          |
| :----- | :-------------------------------------------- | :---------------------------- |
| +      | 加法                                          | `expr $a + $b` 结果为 30。    |
| -      | 减法                                          | `expr $a - $b` 结果为 -10。   |
| *      | 乘法                                          | `expr $a \* $b` 结果为  200。 |
| /      | 除法                                          | `expr $b / $a` 结果为 2。     |
| %      | 取余                                          | `expr $b % $a` 结果为 0。     |
| =      | 赋值                                          | a=$b 将把变量 b 的值赋给 a。  |
| ==     | 相等。用于比较两个数字，相同则返回 true。     | [ $a == $b ] 返回 false。     |
| !=     | 不相等。用于比较两个数字，不相同则返回 true。 | [ $a != $b ] 返回 true。      |

**（注意：**条件表达式要放在方括号之间，并且要有空格，例如: **[$a==$b]** 是错误的，必须写成 **[ $a == $b ]）**

​          **注意：乘号(*)前边必须加反斜杠(\)才能实现乘法运算，也是转义字符**

```shell
#!/bin/bash

a=10
b=20

val=`expr $a + $b`
echo "a + b : $val"
# 执行结果：a + b : 30



val=`expr $a - $b`
echo "a - b : $val"
# 执行结果：a - b : -10



val=`expr $a \* $b`
echo "a * b : $val"
# 执行结果：a * b : 200



val=`expr $b / $a`
echo "b / a : $val"
# 执行结果：b / a : 2



val=`expr $b % $a`
echo "b % a : $val"
# 执行结果：b % a : 0



if [ $a == $b ]
then
   echo "a 等于 b"
fi
if [ $a != $b ]
then
   echo "a 不等于 b"
fi
# 执行结果：a 不等于 b
```

​          注意事项：

```
if...then...fi 是条件语句，后续将会讲解。

在 MAC 中 shell 的 expr 语法是：$((表达式))，此处表达式中的 "*" 不需要转义符号 "\" 。
```

​          **关系型运算符**

​          关系运算符只支持数字，不支持字符串，除非字符串的值是数字。

```
环境如下：变量 a 为 10，变量 b 为 20
```

| 运算符 | 说明                                                  | 举例                       |
| :----- | :---------------------------------------------------- | :------------------------- |
| -eq    | 检测两个数是否相等，相等返回 true。                   | [ $a -eq $b ] 返回 false。 |
| -ne    | 检测两个数是否不相等，不相等返回 true。               | [ $a -ne $b ] 返回 true。  |
| -gt    | 检测左边的数是否大于右边的，如果是，则返回 true。     | [ $a -gt $b ] 返回 false。 |
| -lt    | 检测左边的数是否小于右边的，如果是，则返回 true。     | [ $a -lt $b ] 返回 true。  |
| -ge    | 检测左边的数是否大于等于右边的，如果是，则返回 true。 | [ $a -ge $b ] 返回 false。 |
| -le    | 检测左边的数是否小于等于右边的，如果是，则返回 true。 | [ $a -le $b ] 返回 true。  |

```shell
#!/bin/bash

a=10
b=20

if [ $a -eq $b ]
then
   echo "$a -eq $b : a 等于 b"
else
   echo "$a -eq $b: a 不等于 b"
fi
# 执行结果：10 -eq 20: a 不等于 b




if [ $a -ne $b ]
then
   echo "$a -ne $b: a 不等于 b"
else
   echo "$a -ne $b : a 等于 b"
fi
# 执行结果：10 -ne 20: a 不等于 b




if [ $a -gt $b ]
then
   echo "$a -gt $b: a 大于 b"
else
   echo "$a -gt $b: a 不大于 b"
fi
# 执行结果：10 -gt 20: a 不大于 b




if [ $a -lt $b ]
then
   echo "$a -lt $b: a 小于 b"
else
   echo "$a -lt $b: a 不小于 b"
fi
# 执行结果：10 -lt 20: a 小于 b




if [ $a -ge $b ]
then
   echo "$a -ge $b: a 大于或等于 b"
else
   echo "$a -ge $b: a 小于 b"
fi
# 执行结果：10 -ge 20: a 小于 b





if [ $a -le $b ]
then
   echo "$a -le $b: a 小于或等于 b"
else
   echo "$a -le $b: a 大于 b"
fi
# 执行结果：10 -le 20: a 小于或等于 b
```

​          **布尔运算符**

```
环境如下：变量 a 为 10，变量 b 为 20
```

| 运算符 | 说明                                                | 举例                                     |
| :----- | :-------------------------------------------------- | :--------------------------------------- |
| !      | 非运算，表达式为 true 则返回 false，否则返回 true。 | [ ! false ] 返回 true。                  |
| -o     | 或运算，有一个表达式为 true 则返回 true。           | [ $a -lt 20 -o $b -gt 100 ] 返回 true。  |
| -a     | 与运算，两个表达式都为 true 才返回 true。           | [ $a -lt 20 -a $b -gt 100 ] 返回 false。 |

```shell
#!/bin/bash

a=10
b=20

if [ $a != $b ]
then
   echo "$a != $b : a 不等于 b"
else
   echo "$a == $b: a 等于 b"
fi
# 执行结果：10 != 20 : a 不等于 b




if [ $a -lt 100 -a $b -gt 15 ]
then
   echo "$a 小于 100 且 $b 大于 15 : 返回 true"
else
   echo "$a 小于 100 且 $b 大于 15 : 返回 false"
fi
# 执行结果：10 小于 100 且 20 大于 15 : 返回 true





if [ $a -lt 100 -o $b -gt 100 ]
then
   echo "$a 小于 100 或 $b 大于 100 : 返回 true"
else
   echo "$a 小于 100 或 $b 大于 100 : 返回 false"
fi
# 执行结果：10 小于 100 或 20 大于 100 : 返回 true




if [ $a -lt 5 -o $b -gt 100 ]
then
   echo "$a 小于 5 或 $b 大于 100 : 返回 true"
else
   echo "$a 小于 5 或 $b 大于 100 : 返回 false"
fi
# 执行结果：10 小于 5 或 20 大于 100 : 返回 false
```

​          **逻辑运算符**

```
环境如下：变量 a 为 10，变量 b 为 20
```

| 运算符 | 说明       | 举例                                       |
| :----- | :--------- | :----------------------------------------- |
| &&     | 逻辑的 AND | [[ \$a -lt 100 && $b -gt 100 ]] 返回 false |
| \|\|   | 逻辑的 OR  | [[ $a -lt 100 || $b -gt 100 ]] 返回 true   |



```shell
#!/bin/bash

a=10
b=20


if [[ $a -lt 100 && $b -gt 100 ]]
then
   echo "返回 true"
else
   echo "返回 false"
fi
# 执行结果：返回 false




if [[ $a -lt 100 || $b -gt 100 ]]
then
   echo "返回 true"
else
   echo "返回 false"
fi
# 执行结果：返回 true
```

​          **字符串运算符**

```
环境如下：变量 a 为 "abc"，变量 b 为 "efg"
```

| 运算符 | 说明                                         | 举例                     |
| :----- | :------------------------------------------- | :----------------------- |
| =      | 检测两个字符串是否相等，相等返回 true。      | [ $a = $b ] 返回 false。 |
| !=     | 检测两个字符串是否相等，不相等返回 true。    | [ $a != $b ] 返回 true。 |
| -z     | 检测字符串长度是否为0，为0返回 true。        | [ -z $a ] 返回 false。   |
| -n     | 检测字符串长度是否不为 0，不为 0 返回 true。 | [ -n "$a" ] 返回 true。  |
| $      | 检测字符串是否为空，不为空返回 true。        | [ $a ] 返回 true。       |

​          **文件测试运算符**

文件测试运算符用于检测 Unix 文件的各种属性。

| 操作符  | 说明                                                         | 举例                      |
| :------ | :----------------------------------------------------------- | :------------------------ |
| -b file | 检测文件是否是块设备文件，如果是，则返回 true。              | [ -b $file ] 返回 false。 |
| -c file | 检测文件是否是字符设备文件，如果是，则返回 true。            | [ -c $file ] 返回 false。 |
| -d file | 检测文件是否是目录，如果是，则返回 true。                    | [ -d $file ] 返回 false。 |
| -f file | 检测文件是否是普通文件（既不是目录，也不是设备文件），如果是，则返回 true。 | [ -f $file ] 返回 true。  |
| -g file | 检测文件是否设置了 SGID 位，如果是，则返回 true。            | [ -g $file ] 返回 false。 |
| -k file | 检测文件是否设置了粘着位(Sticky Bit)，如果是，则返回 true。  | [ -k $file ] 返回 false。 |
| -p file | 检测文件是否是有名管道，如果是，则返回 true。                | [ -p $file ] 返回 false。 |
| -u file | 检测文件是否设置了 SUID 位，如果是，则返回 true。            | [ -u $file ] 返回 false。 |
| -r file | 检测文件是否可读，如果是，则返回 true。                      | [ -r $file ] 返回 true。  |
| -w file | 检测文件是否可写，如果是，则返回 true。                      | [ -w $file ] 返回 true。  |
| -x file | 检测文件是否可执行，如果是，则返回 true。                    | [ -x $file ] 返回 true。  |
| -s file | 检测文件是否为空（文件大小是否大于0），不为空返回 true。     | [ -s $file ] 返回 true。  |
| -e file | 检测文件（包括目录）是否存在，如果是，则返回 true。          | [ -e $file ] 返回 true。  |

```
-S: 判断某文件是否 socket。
-L: 检测文件是否存在并且是一个符号链接。
```

```shell
[root@localhost ~]# vim test.sh


#!/bin/bash

file="/root/test.sh"

if [ -r $file ]
then
   echo "文件可读"
else
   echo "文件不可读"
fi
# 执行结果：文件可读




if [ -w $file ]
then
   echo "文件可写"
else
   echo "文件不可写"
fi
# 执行结果：文件可写




if [ -x $file ]
then
   echo "文件可执行"
else
   echo "文件不可执行"
fi
# 执行结果：文件可执行




if [ -f $file ]
then
   echo "文件为普通文件"
else
   echo "文件为特殊文件"
fi
# 执行结果：文件为普通文件





if [ -d $file ]
then
   echo "文件是个目录"
else
   echo "文件不是个目录"
fi
# 执行结果：文件不是个目录




if [ -s $file ]
then
   echo "文件不为空"
else
   echo "文件为空"
fi
# 执行结果：文件不为空




if [ -e $file ]
then
   echo "文件存在"
else
   echo "文件不存在"
fi
# 执行结果：文件存在


[root@localhost ~]# ./test.sh
文件可读
文件可写
文件可执行
文件为普通文件
文件不是个目录
文件不为空
文件存在

```



#### 8、Shell命令

（1）**echo命令：用于字符串的输出**

​          语法：

```shell
echo string
```

​          **显示普通字符串**

```shell
echo "It is a test"
或者
echo It is a test

# 执行结果：It is a test
```

​          **显示转义字符**

```shell
echo "\"It is a test\""

# 执行结果："It is a test"
```

​          **显示变量**

​          read 命令从标准输入中读取一行,并把输入行的每个字段的值指定给 shell 变量

```shell
#!/bin/sh
read name 
echo "$name It is a test"
```

​          运行代码：

```shell
[root@localhost ~]# ./test.sh
mashibing
mashibing It is a test
```

​          **显示换行与不换行**

```shell
#!/bin/sh
echo -e "OK! \n" # -e 开启转义   \n 开始换行
# echo -e "OK! \c"   同理 \c 没有换行了
echo "It is a test"
```

​          运行代码：

```shell
[root@localhost ~]# ./test.sh
OK! 

It is a test

```

​          **显示结果定向至文件**

```shell
[root@localhost ~]# vim test.sh

#!/bin/sh
echo "It is a test" > myfile

[root@localhost ~]# ./test.sh
[root@localhost ~]# ls
anaconda-ks.cfg  link_text2  mydir  myfile  test.sh
[root@localhost ~]# cat myfile 
It is a test

```

​          **原样输出字符串，不进行转义或者取变量**

```shell
#!/bin/sh
echo '$name\"'

[root@localhost ~]# ./test.sh 
$name\"
```

​          **显示命令执行结果**

```shell
#!/bin/sh
echo `date`

# 注意： 这里使用的是反引号 `, 而不是单引号 '。

[root@localhost ~]# ./test.sh 
2020年 09月 05日 星期二 16:40:22 CST
# 结果将显示当前日期
```

（2）**printf命令**

​                   printf 命令模仿 C 程序库（library）里的 printf() 程序。

​                   printf 由 POSIX 标准所定义，因此使用 printf 的脚本比使用 echo 移植性好。

​                   printf 使用引用文本或空格分隔的参数，可以在 printf 中使用格式化字符串，还可以制定字符串的宽度、左右对齐方式等。

​                   默认 printf 不会像 echo 自动添加换行符，我们可以手动添加 \n

​                   语法：

```shell
printf  format-string  [arguments...]
```

​                   参数说明：

```
format-string: 为格式控制字符串

arguments: 为参数列表。
```

​                   例如换行这里可以如下方法：

```shell
#!/bin/sh

$ echo "Hello, Shell"
# 执行结果：Hello, Shell


$ printf "Hello, Shell\n"
# 执行结果： Hello, Shell
            $
```

​                   来看一个实际的脚本

```shell
#!/bin/bash

printf "%-10s %-8s %-4s\n" 姓名 性别 体重kg
printf "%-10s %-8s %-4.2f\n" 大锤 男 66.1234  
printf "%-10s %-8s %-4.2f\n" 小虎 男 60.6543
printf "%-10s %-8s %-4.2f\n" 小美 女 47.9876

```

​                   运行结果如下：

```shell
[root@localhost ~]# ./test.sh 
姓名     性别   体重kg
大锤     男      66.12
小虎     男      60.65
小美     女      47.99
```

扩展：

```
%s %c %d %f都是格式替代符

%-10s 指一个宽度为10个字符（-表示左对齐，没有则表示右对齐），任何字符都会被显示在10个字符宽的字符内，如果不足则自动以空格填充，超过也会将内容全部显示出来。

%-4.2f 指格式化为小数，其中.2指保留2位小数。
```



```shell
#!/bin/bash

# format-string为双引号
printf "%d %s\n" 1 "abc"

# 单引号与双引号效果一样 
printf '%d %s\n' 1 "abc" 

# 运行结果：
[root@localhost ~]# ./test.sh 
1 abc
1 abc

```

​                   printf转义

| 序列  | 说明                                                         |
| :---- | :----------------------------------------------------------- |
| \a    | 警告字符，通常为ASCII的BEL字符                               |
| \b    | 后退                                                         |
| \c    | 抑制（不显示）输出结果中任何结尾的换行字符（只在%b格式指示符控制下的参数字符串中有效），而且，任何留在参数里的字符、任何接下来的参数以及任何留在格式字符串中的字符，都被忽略 |
| \f    | 换页（formfeed）                                             |
| \n    | 换行                                                         |
| \r    | 回车（Carriage return）                                      |
| \t    | 水平制表符                                                   |
| \v    | 垂直制表符                                                   |
| \\    | 一个字面上的反斜杠字符                                       |
| \ddd  | 表示1到3位数八进制值的字符。仅在格式字符串中有效             |
| \0ddd | 表示1到3位的八进制值字符                                     |

```shell
#!/bin/bash	
printf "a string, no processing:<%s>\n" "A\nB"


运行结果：
[root@localhost ~]# ./test.sh 
a string, no processing:<A\nB>

```

（3）**test 命令：用于检查某个条件是否成立，它可以进行数值、字符和文件三个方面的测试**

​                   数值测试

| 参数 | 说明           |
| :--- | :------------- |
| -eq  | 等于则为真     |
| -ne  | 不等于则为真   |
| -gt  | 大于则为真     |
| -ge  | 大于等于则为真 |
| -lt  | 小于则为真     |
| -le  | 小于等于则为真 |

​        例如：

```shell
#!/bin/bash	

num1=50
num2=50
if test $[num1] -eq $[num2]
then
    echo '两个数相等！'
else
    echo '两个数不相等！'
fi
```

​        代码中的 **[  ]** 执行基本的算数运算        

​        运行代码：

```shell
[root@localhost ~]# ./test.sh 
两个数相等！
```

​        字符串测试

| 参数      | 说明                     |
| :-------- | :----------------------- |
| =         | 等于则为真               |
| !=        | 不相等则为真             |
| -z 字符串 | 字符串的长度为零则为真   |
| -n 字符串 | 字符串的长度不为零则为真 |

​        例如：

```shell
#!/bin/bash	

num1="mashibing"
num2="mashibin"
if test $num1 = $num2
then
    echo '两个字符串相等!'
else
    echo '两个字符串不相等!'
fi
```

​        运行代码：

```shell
[root@localhost ~]# ./test.sh 
两个字符串不相等!

```

​        文件测试

| 参数      | 说明                                 |
| :-------- | :----------------------------------- |
| -e 文件名 | 如果文件存在则为真                   |
| -r 文件名 | 如果文件存在且可读则为真             |
| -w 文件名 | 如果文件存在且可写则为真             |
| -x 文件名 | 如果文件存在且可执行则为真           |
| -s 文件名 | 如果文件存在且至少有一个字符则为真   |
| -d 文件名 | 如果文件存在且为目录则为真           |
| -f 文件名 | 如果文件存在且为普通文件则为真       |
| -c 文件名 | 如果文件存在且为字符型特殊文件则为真 |
| -b 文件名 | 如果文件存在且为块特殊文件则为真     |

​        例如：

```shell
#!/bin/bash

cd /bin
if test -e ./bash
then
    echo '文件已存在!'
else
    echo '文件不存在!'
fi
```

​        运行代码：

```shell
[root@localhost ~]# ./test.sh 
文件已存在!
```

（备注：还提供了与( -a )、或( -o )、非( ! )三个逻辑操作符用于将测试条件连接起来，其优先级为： **!** 最高， **-a** 次之， **-o** 最低）

#### 9、流程控制

​        sh的流程控制不可为空，在sh/bash里如果else分支没有语句执行，就不要写这个else，不像Python可以写一个pass

（1）**if else**

​        if 的 语法：

```shell
#!/bin/bash

if condition
then
    command1 
    command2
    ...
    commandN 
fi

# 也可以一行（适用于终端命令提示符）
if [ $(ps -ef | grep -c "ssh") -gt 1 ]; then echo "true"; fi
```

​        if else 的 语法：

```shell
#!/bin/bash

if condition
then
    command1 
    command2
    ...
    commandN
else
    command
fi
```

​        if else-if else 语法格式：

```shell
#!/bin/bash

if condition1
then
    command1
elif condition2 
then 
    command2
else
    commandN
fi
```

​        示例如下：

```shell
#!/bin/bash

a=10
b=20
if [ $a == $b ]
then
   echo "a 等于 b"
elif [ $a -gt $b ]
then
   echo "a 大于 b"
elif [ $a -lt $b ]
then
   echo "a 小于 b"
else
   echo "没有符合的条件"
fi
```

​        运行结果：

```shell
[root@localhost ~]# ./test.sh 
a 小于 b

```

（2）**for 循环**

​        for 循环 的 语法：

```shell
#!/bin/bash

for var in item1 item2 ... itemN
do
    command1
    command2
    ...
    commandN
done

# 这里也可以写成一行
for var in item1 item2 ... itemN; do command1; command2… done;
```

​        当变量值在列表里，for循环即执行一次所有命令，使用变量名获取列表中的当前取值。

​        命令可为任何有效的shell命令和语句。in列表可以包含替换、字符串和文件名。

​        in列表是可选的，如果不用它，for循环使用命令行的位置参数。

​        例如，顺序输出当前列表中的数字：

```shell
#!/bin/bash

for loop in 1 2 3 4 5
do
    echo "The value is: $loop"
done
```

​        运行脚本输出

```shell
[root@localhost ~]# ./test.sh 
The value is: 1
The value is: 2
The value is: 3
The value is: 4
The value is: 5

```

（2）**While 循环**

```
while循环用于不断执行一系列命令，也用于从输入文件中读取数据；命令通常为测试条件
```

​        语法：

```shell
#!/bin/bash	

while condition
do
    command
done
```

​        实现：

```
基本的while循环，测试条件是：如果int小于等于5，那么条件返回真。
int从0开始，每次循环处理时，int加1。运行上述脚本，返回数字1到5，然后终止
```

​        如下：

```shell
#!/bin/bash

int=1
while(( $int<=5 ))
do
    echo $int
    let "int++"
done


[root@localhost ~]# ./test.sh 
1
2
3
4
5
```

当然还有无限循环，也称呼为死循环

```shell
#!/bin/bash



while :
do
    command
done
# 或者


while true
do
    command
done
# 或者


for (( ; ; ))
```

（3）case选择语句

```
case语句为多选择语句。可以用case语句匹配一个值与一个模式，如果匹配成功，执行相匹配的命令
```

​        语法：

```shell
#!/bin/bash


case 值 in
模式1)
    command1
    command2
    ...
    commandN
    ;;
模式2）
    command1
    command2
    ...
    commandN
    ;;
esac
```

```
case工作方式如上所示。取值后面必须为单词in，每一模式必须以右括号结束。取值可以为变量或常数。匹配发现取值符合某一模式后，其间所有命令开始执行直至 ;;。

取值将检测匹配的每一个模式。一旦模式匹配，则执行完匹配模式相应命令后不再继续其他模式。如果无一匹配模式，使用星号 * 捕获该值，再执行后面的命令。
```

​        示例如下：

```shell
#!/bin/bash


echo '输入 1 到 4 之间的数字:'
echo '你输入的数字为:'
read aNum
case $aNum in
    1)  echo '你选择了 1'
    ;;
    2)  echo '你选择了 2'
    ;;
    3)  echo '你选择了 3'
    ;;
    4)  echo '你选择了 4'
    ;;
    *)  echo '你没有输入 1 到 4 之间的数字'
    ;;
esac
```

​        运行如下：

```shell
[root@localhost ~]# ./test.sh 
输入 1 到 4 之间的数字:
你输入的数字为:
1
你选择了 1

```

（4）跳出循环

​        实际我们经常看到跳出循环功能：break和continue。在Python中也是一样的存在

​        **break命令：break命令允许跳出所有循环（终止执行后面的所有循环）**

​        例如下方：(脚本进入死循环直至用户输入数字大于5。要跳出这个循环，返回到shell提示符下，需要使用break命令)

```shell
#!/bin/bash

while :
do
    echo -n "输入 1 到 5 之间的数字:"
    read aNum
    case $aNum in
        1|2|3|4|5) echo "你输入的数字为 $aNum!"
        ;;
        *) echo "你输入的数字不是 1 到 5 之间的! 游戏结束"
            break
        ;;
    esac
done
```

​        运行代码

```shell
[root@localhost ~]# ./test.sh 
输入 1 到 5 之间的数字:4
你输入的数字为 4!
输入 1 到 5 之间的数字:7
你输入的数字不是 1 到 5 之间的! 游戏结束
```

​        **continue命令与break命令类似，只有一点差别，它不会跳出所有循环，仅仅跳出当前循环。**

​        对上面的例子进行修改：

```shell
#!/bin/bash

while :
do
    echo -n "输入 1 到 5 之间的数字: "
    read aNum
    case $aNum in
        1|2|3|4|5) echo "你输入的数字为 $aNum!"
        ;;
        *) echo "你输入的数字不是 1 到 5 之间的!"
            continue
            echo "游戏结束"
        ;;
    esac
done
```

其实可以发现语句已经陷入死循环了，但是有一点，只要输入大于5的数字时，循环不会结束，语句 **echo "游戏结束"** 永远不会被执行

（注意：强行结束可以 ctrl + c）

#### 10、函数

`basename` 函数的主要作用是从一个路径中提取文件名部分，或可同时删除文件路径中的指定后缀。

`dirname` 用于提取路径中的目录部分，与 `basename` 互补使用。

​        shell 可以用户定义函数，然后在shell脚本中可以随便调用。

​        **基本函数语法如下**：

```shell
#!/bin/bash


[ function ] funname [()]

{

    action;

    [return int;]

}
```

​        参数说明：

```
1、可以带function fun() 定义，也可以直接fun() 定义,不带任何参数。
2、参数返回，可以显示加：return 返回，如果不加，将以最后一条命令运行结果，作为返回值。 return后跟数值n(0-255)
```

​        示例如下：

```shell
#!/bin/bash

demoFun(){
    echo "这是我的第一个 shell 函数!"
}
echo "-----函数开始执行-----"
demoFun
echo "-----函数执行完毕-----"
```

​        运行如下：

```shell
[root@localhost ~]# ./test.sh 
-----函数开始执行-----
这是我的第一个 shell 函数!
-----函数执行完毕-----
```

​        带return语句

```shell
#!/bin/bash


funWithReturn(){
    echo "这个函数会对输入的两个数字进行相加运算..."
    echo "输入第一个数字: "
    read aNum
    echo "输入第二个数字: "
    read anotherNum
    echo "两个数字分别为 $aNum 和 $anotherNum !"
    return $(($aNum+$anotherNum))
}
funWithReturn
echo "输入的两个数字之和为 $? !"
```

​        运行如下：

```shell
[root@localhost ~]# ./test.sh 
这个函数会对输入的两个数字进行相加运算...
输入第一个数字: 
5
输入第二个数字: 
6
两个数字分别为 5 和 6 !
输入的两个数字之和为 11 !

```

​        **函数参数**

```
调用函数时可以向其传递参数。在函数体内部，通过 $n 的形式来获取参数的值，例如，$1表示第一个参数，$2表示第二个参数...
```

```shell
#!/bin/bash

funWithParam(){
    echo "第一个参数为 $1 !"
    echo "第二个参数为 $2 !"
    echo "第十个参数为 $10 !"
    echo "第十个参数为 ${10} !"
    echo "第十一个参数为 ${11} !"
    echo "参数总数有 $# 个!"
    echo "作为一个字符串输出所有参数 $* !"
}
funWithParam 1 2 3 4 5 6 7 8 9 34 73
```

​        运行如下：

```
[root@localhost ~]# ./test.sh 
第一个参数为 1 !
第二个参数为 2 !
第十个参数为 10 !
第十个参数为 34 !
第十一个参数为 73 !
参数总数有 11 个!
作为一个字符串输出所有参数 1 2 3 4 5 6 7 8 9 34 73 !
```

（注意，$10 不能获取第十个参数，获取第十个参数需要${10}。当n>=10时，需要使用${n}来获取参数。）

| 参数处理 | 说明                                                         |
| :------- | :----------------------------------------------------------- |
| $#       | 传递到脚本或函数的参数个数                                   |
| $*       | 以一个单字符串显示所有向脚本传递的参数                       |
| $$       | 脚本运行的当前进程ID号                                       |
| $!       | 后台运行的最后一个进程的ID号                                 |
| $@       | 与$*相同，但是使用时加引号，并在引号中返回每个参数。         |
| $-       | 显示Shell使用的当前选项，与set命令功能相同。                 |
| $?       | 显示最后命令的退出状态。0表示没有错误，其他任何值表明有错误。 |

### **`sed` 命令**详解

`sed` 是 Linux 系统中的流编辑器，用于**对文本进行快速的处理和编辑**。它可以在管道中使用，也可以直接操作文件，支持**查找、替换、插入、删除**等功能。`sed` 的强大之处在于支持正则表达式，使其非常灵活和高效。

------

#### **语法**

```bash
sed [OPTION]... 'SCRIPT' [FILE]...
```

- **SCRIPT**：指定要执行的操作（如替换、删除）。
- **FILE**：要处理的文件。

------

#### **常用选项**

- `-e SCRIPT`：指定要执行的脚本（可多次使用）。
- `-f FILE`：从指定文件中读取脚本。
- `-i[SUFFIX]`：直接修改文件内容（原地编辑），可选择备份。
- `-n`：禁止自动打印，配合 `p` 显式打印结果。
- `-r` 或 `-E`：启用扩展正则表达式（等同于 `--regexp-extended`）。

------

#### **基本功能**

1. **文本替换**
2. **删除指定行**
3. **插入和追加内容**
4. **显示部分内容**
5. **正则表达式匹配与处理**

------

#### **常见操作示例**

##### **1. 文本替换**

**将文本中的 "hello" 替换为 "world"：**

```bash
sed 's/hello/world/' file.txt
```

**全局替换（替换每一行中所有匹配的内容）：**

```bash
sed 's/hello/world/g' file.txt
```

**忽略大小写：**

```bash
sed 's/hello/world/gi' file.txt
```

**原地替换（直接修改文件内容）：**

```bash
sed -i 's/hello/world/g' file.txt
```

**使用备份文件（修改前创建 `.bak` 备份）：**

```bash
sed -i.bak 's/hello/world/g' file.txt
```

------

##### **2. 删除指定行**

**删除第 3 行：**

```bash
sed '3d' file.txt
```

**删除第 2 到第 4 行：**

```bash
sed '2,4d' file.txt
```

**删除匹配特定模式的行：**

```bash
sed '/pattern/d' file.txt
```

------

##### **3. 插入和追加**

**在第 2 行前插入内容：**

```bash
sed '2i\This is a new line' file.txt
```

**在第 2 行后追加内容：**

```bash
sed '2a\This is an appended line' file.txt
```

**匹配行后追加内容：**

```bash
sed '/pattern/a\This is appended after pattern' file.txt
```

------

##### **4. 显示部分内容**

**仅显示匹配的行：**

```bash
sed -n '/pattern/p' file.txt
```

**显示第 2 到第 4 行：**

```bash
sed -n '2,4p' file.txt
```

**显示匹配的行及其上下文：**

```bash
sed -n -e '/pattern/{=;p}' file.txt
```

------

##### **5. 替换带分隔符的内容**

`sed` 支持自定义分隔符，适用于包含 `/` 的字符串。

**替换路径字符串：**

```bash
sed 's|/usr/bin|/usr/local/bin|g' file.txt
```

------

##### **6. 使用正则表达式**

**替换数字：**

```bash
sed 's/[0-9]/#/g' file.txt
```

**替换开头匹配的内容：**

```bash
sed 's/^pattern/replacement/' file.txt
```

**替换行尾匹配的内容：**

```bash
sed 's/pattern$/replacement/' file.txt
```

------

##### **7. 多操作组合**

**替换 "hello" 为 "world"，同时删除第 3 行：**

```bash
sed -e 's/hello/world/g' -e '3d' file.txt
```

------

#### **高级功能**

##### **1. 替换文件中的所有内容**

将文件内容全部替换为 "new content"：

```bash
sed 'c\new content' file.txt
```

##### **2. 替换文件中两次出现的模式**

仅替换第 2 次匹配到的内容：

```bash
sed 's/hello/world/2' file.txt
```

------

##### **结合管道使用**

`sed` 可以直接处理其他命令的输出：

**查找并替换目录列表中的文件名：**

```bash
ls | sed 's/^/File: /'
```

**处理命令输出：**

```bash
ps aux | sed '/root/d'
```

------

#### **适用场景**

1. **文本文件批量替换**
   - 修改日志文件中的日期格式。
   - 批量更新配置文件的路径或参数。
2. **日志分析**
   - 删除特定行或内容。
   - 提取关键字相关的上下文。
3. **自动化脚本**
   - 动态修改配置文件。
   - 从输出结果中提取所需信息。

------

#### **注意事项**

1. **大小写敏感**
   - 默认区分大小写，使用 `i` 修饰符忽略大小写。
2. **分隔符冲突**
   - 当内容中包含 `/` 时，建议使用其他分隔符（如 `|`）。
3. **原地编辑**
   - 使用 `-i` 时需谨慎，避免误操作。建议提前备份文件。
4. **多字节字符**
   - 对于非 ASCII 编码（如中文），需确保文件编码和系统环境一致。

------

#### **总结**

`sed` 是一个功能强大的文本处理工具，适用于快速编辑和自动化脚本操作。无论是简单的文本替换，还是复杂的正则处理，`sed` 都是 Linux 用户不可或缺的利器。



### sed 命令**`-e` 选项的作用**

1. **允许多个操作脚本**
    使用 `-e` 时，可以在同一条 `sed` 命令中应用多个操作脚本，而不需要写成多个独立的 `sed` 命令。
2. **提高灵活性**
    通过 `-e`，可以在一次读取输入时完成多个处理步骤，避免重复操作，提高效率。

------

#### **示例用法**

##### **1. 单个脚本**

使用 `-e` 指定单条脚本：

```
bash


复制编辑
sed -e 's/hello/world/g' file.txt
```

功能：将文件 `file.txt` 中所有的 `hello` 替换为 `world`。

------

##### **2. 多个脚本**

使用 `-e` 指定多条脚本：

```bash
sed -e 's/hello/world/g' -e 's/foo/bar/g' file.txt
```

功能：

1. 将 `hello` 替换为 `world`。
2. 将 `foo` 替换为 `bar`。

------

##### **3. 多行脚本**

也可以直接用多个 `-e` 指定多行操作：

```bash
sed -e '1d' -e 's/error/warning/g' file.txt
```

功能：

1. 删除第一行。
2. 将 `error` 替换为 `warning`。

### **`awk` 命令详解**

`awk` 是一种强大的文本处理工具，**适用于格式化文本文件的分析和处理（包括简单的统计）**。它以**逐行扫描文件为基础**，根据模式匹配并对每一行执行指定的操作。`awk` 是一种编程语言，支持变量、条件判断、循环等功能。

------

#### **基本语法**

```bash
awk [options] 'script' file
```

- **`script`**：定义 `awk` 的执行逻辑，可以包含模式匹配和操作。
- **`file`**：要处理的文件。
- `options`：常见选项如 `-F` 指定分隔符，`-v` 定义变量。

------

#### **`awk` 的结构**

`awk` 的程序结构通常分为以下三个部分：

1. **BEGIN 块**：在处理输入之前执行一次的代码块。
2. **主处理块**：对文件的每一行执行。
3. **END 块**：在处理所有输入之后执行一次的代码块。

```bash
awk 'BEGIN { actions } { actions } END { actions }' file
```

- **`BEGIN`**：初始化操作，例如设置分隔符、打印表头等。
- **`{}`**：对每一行进行操作。
- **`END`**：完成后处理工作，例如打印统计结果。

------

#### **常用功能与示例**

##### **1. 打印文件内容**

**直接打印每一行：**

```bash
awk '{ print }' file.txt
```

等价于 `cat file.txt`。

------

##### **2. 打印指定列**

**打印文件的第 1 列和第 3 列：**

```bash
awk '{ print $1, $3 }' file.txt
```

- `$1`：表示第 1 列。
- `$2`：表示第 2 列，依此类推。

**指定分隔符（例如用逗号分隔的文件）：**

```bash
awk -F ',' '{ print $1, $2 }' file.csv
```

------

##### **3. 条件过滤**

**打印包含关键字 "error" 的行：**

```bash
awk '/error/' file.txt
```

**打印第 2 列大于 50 的行：**

```bash
awk '$2 > 50' file.txt
```

**条件 + 列操作：**

```bash
awk '$2 > 50 { print $1, $3 }' file.txt
```

------

##### **4. BEGIN 和 END 块**

**在文件开头和结尾添加消息：**

```bash
awk 'BEGIN { print "Start of File" } { print } END { print "End of File" }' file.txt
```

输出示例：

```
Start of File
Line 1 content
Line 2 content
End of File
```

------

##### **5. 统计行数和总和**

**统计行数：**

```bash
awk 'END { print "Total lines:", NR }' file.txt
```

- `NR`：当前行号，表示总行数。

**统计某一列的总和（如第 3 列）：**

```bash
awk '{ sum += $3 } END { print "Total:", sum }' file.txt
```

------

##### **6. 使用内置变量**

`awk` 提供了许多内置变量，常用的有：

- **`NR`**：当前行号。
- **`NF`**：当前行的字段数（列数）。
- **`$0`**：当前行的完整内容。
- **`FS`**：输入字段分隔符（默认为空格或制表符）。
- **`OFS`**：输出字段分隔符。
- **`RS`**：记录分隔符（默认为换行）。
- **`ORS`**：输出记录分隔符。

**示例：打印行号和行内容：**

```bash
awk '{ print NR, $0 }' file.txt
```

------

##### **7. 自定义分隔符**

**以逗号为分隔符处理文件：**

```bash
awk -F ',' '{ print $1, $2 }' file.csv
```

**更改输出分隔符：**

```bash
awk 'BEGIN { OFS=":" } { print $1, $2 }' file.txt
```

------

##### **8. 复杂模式匹配**

**打印包含 "error" 或 "failed" 的行：**

```bash
awk '/error|failed/' file.txt
```

**只处理第 3 列为数字的行：**

```bash
awk '$3 ~ /^[0-9]+$/' file.txt
```

------

##### **9. 多条件处理**

**基于条件执行不同操作：**

```bash
awk '{ if ($2 > 50) print $1, $2; else print "Low:", $1 }' file.txt
```

**使用 `else if`：**

```bash
awk '{ 
  if ($2 > 80) print $1, "High";
  else if ($2 > 50) print $1, "Medium";
  else print $1, "Low";
}' file.txt
```

------

##### **10. 调用外部变量**

使用 `-v` 传递变量：

```bash
awk -v threshold=50 '$2 > threshold { print $1, $2 }' file.txt
```

------

#### **结合管道**

`awk` 可处理其他命令的输出：

**统计 `ps` 命令中每个用户的进程数量：**

```bash
ps aux | awk '{ print $1 }' | sort | uniq -c
```

------

#### **综合示例**

##### **统计文本信息**

统计文件中每列的总和并输出：

```bash
awk '{ for (i = 1; i <= NF; i++) sum[i] += $i } END { for (i in sum) print "Column", i, "Sum:", sum[i] }' file.txt
```

##### **处理日志文件**

提取日志文件中某个时间段的错误信息：

```bash
awk '$1 ~ /2025-03-13/ && /ERROR/ { print $0 }' logfile.log
```

------

#### **注意事项**

1. **分隔符处理**
   - 默认以空格或制表符分隔，需用 `-F` 指定其他分隔符。
2. **模式优先级**
   - 模式匹配在操作之前，确保模式正确。
3. **环境兼容性**
   - 不同操作系统上的 `awk` 可能有差异，例如 `gawk`、`mawk` 等。

------

#### **总结**

`awk` 是一个强大的文本处理工具，结合正则表达式、内置变量和条件逻辑，可以高效完成各种文本处理任务。适用于文件解析、日志分析和格式化输出等场景。

### sed和awk比较

| 工具      | 功能和用途                                                   |
| --------- | ------------------------------------------------------------ |
| **`sed`** | 主要用于**流式编辑**，擅长对文本进行**逐行处理**，适合文本的替换、删除、插入和修改等操作。常用于简单的文本处理任务，例如行编辑和字符串替换。 |
| **`awk`** | 主要用于**数据处理和报告生成**，擅长按字段对文本进行分割，支持模式匹配、条件判断、数学运算等功能，适合处理结构化数据，如日志、CSV 文件等。**重点关注每一行的列?** |

**总结：**

- `sed` 更适合简单的**逐行**编辑任务，如替换和删除。
- `awk` 更适合复杂的逻辑处理和按字段操作。

#### **总结**

- **`sed`：文本编辑工具**，专注于行操作（替换、删除、插入等），简单高效，适合基本文本编辑任务。
- **`awk`：文本处理工具**，功能类似编程语言，擅长字段操作和数据分析，适合结构化数据的复杂处理。
