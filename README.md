# Mapq

用类似js的语法查询特定的map


```go
// 返回true
mapq.QueryMap(map[string]interface{}{"a": 1, "b": 2, "c": map[string]interface{}{"d": 3}}, "a == 1 && b == 2 && c.d*(a+b) == 9")
```

## 相关知识和参考资料（仅供参考）

### 知识点
- 递归下降分析法
- 编译原理
- 抽象语法树（AST）（完成此任务不一定需要抽象语法树）


### 参考资料
[如何制作一个小解释器](https://ruslanspivak.com/lsbasi-part1/)

## 支持的数据类型

- bool
- string(单引号)
- number(float64)
- null

## 支持的单目运算符

- \-
- \+
- ！

## 支持的双目运算符

- ==
- !=
- \>
- <
- \>=
- <=
- \+
- \-
- \*
- /
- ||
- &&

## 特殊支持

- 小括号：改变运算优先级
- 嵌套map查询