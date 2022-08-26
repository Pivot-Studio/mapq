# Mapq

用类似js的语法查询特定的map


```go
// 返回true
mapq.QueryMap(map[string]interface{}{"a": 1, "b": 2, "c": map[string]interface{}{"d": 3}}, "a == 1 && b == 2 && c.d*(a+b) == 9")
```


