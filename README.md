# Mapq

用类似js的语法查询特定的map


```go
// 返回true
mapq.QueryMap(map[string]interface{}{"a": 1, "b": 2, "c": map[string]interface{}{"d": 3}}, "a == 1 && b == 2 && c.d*(a+b) == 9")
```

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

小括号：改变运算优先级
