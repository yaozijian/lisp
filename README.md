
一个非常简单的Lisp解析器，学习Go语言后练手的。使用方法如下。

* 支持整数的四则运算
* 用set设置变量;用(set 变量名 nil)清除变量
* 用(print 变量名)输出变量信息


```
E:\GoProjects\MyGithub\src\github.com\yaozijian\lisp>lisp
>(+ 2 3)
5
>(+ 2 (* 4 5))
22
>(set (+ 2 (* 4 5)))
set函数只接受两个操作数
>(set a (+ 2 (* 4 5)))
a: 22
>(print a)
a: 22
>(set a nil)
>(print a)
a
>(print)
-: 0x401db0
*: 0x402180
/: 0x402450
print: 0x401000
set: 0x401380
+: 0x401b00
>exit
```
