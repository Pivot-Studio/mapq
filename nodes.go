package mapq

import "reflect"

// Node 节点
type Node interface {
	Eval(data map[string]interface{}) interface{}
}

// BinNode 表达式节点
type BinNode struct {
	Left, Right Node
	Op          int
}

func toF64(i interface{}) float64 {
	switch v := i.(type) {
	case int:
		return float64(v)
	case float64:
		return v
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	}
	return 0
}
func trytoF64(i interface{}) interface{} {
	switch v := i.(type) {
	case int, float64, int32, int64, uint32, uint64, float32:
		return toF64(v)

	}
	return i
}

// Eval 查询
func (n *BinNode) Eval(data map[string]interface{}) interface{} {
	left := n.Left.Eval(data)
	// early stop logic
	switch n.Op {
	case TYPE_AND:
		if !left.(bool) {
			return false
		}
	case TYPE_OR:
		if left.(bool) {
			return true
		}
	}
	right := n.Right.Eval(data)
	fleft := toF64(left)
	fright := toF64(right)
	switch n.Op {
	case TYPE_EQ:
		return reflect.DeepEqual(trytoF64(left), trytoF64(right))
	case TYPE_NEQ:
		return !reflect.DeepEqual(trytoF64(left), trytoF64(right))
	case TYPE_LG:
		return fleft > fright
	case TYPE_LEQ:
		return fleft >= fright
	case TYPE_SM:
		return fleft < fright
	case TYPE_SEQ:
		return fleft <= fright
	case TYPE_AND:
		return left.(bool) && right.(bool)
	case TYPE_OR:
		return left.(bool) || right.(bool)
	case TYPE_PLUS:
		return fleft + fright
	case TYPE_SUB:
		return fleft - fright
	case TYPE_MUL:
		return fleft * fright
	case TYPE_DIV:
		return fleft / fright
	}
	return false
}

// BoolConstNode 布尔常量节点
type BoolConstNode struct {
	Val bool
}

// Eval 查询
func (n *BoolConstNode) Eval(data map[string]interface{}) interface{} {
	return n.Val
}

// NotNode 非节点
type NotNode struct {
	Bool Node
}

// Eval 查询
func (n *NotNode) Eval(data map[string]interface{}) interface{} {
	return !n.Bool.Eval(data).(bool)
}

// NumNode 数字节点
type NumNode struct {
	Num float64
}

// Eval 查询
func (n *NumNode) Eval(data map[string]interface{}) interface{} {
	return n.Num
}

// VarNode 变量节点
type VarNode struct {
	Token string
	Next  *VarNode
}

// Eval 查询
func (n *VarNode) Eval(data map[string]interface{}) interface{} {
	if n.Next != nil {
		return n.Next.Eval(data[n.Token].(map[string]interface{}))
	}
	return data[n.Token]
}

// StringNode 字符串节点
type StringNode struct {
	Str string
}

// Eval 查询
func (n *StringNode) Eval(data map[string]interface{}) interface{} {
	return n.Str
}

// NullNode null节点
type NullNode struct {
}

// Eval 查询
func (n *NullNode) Eval(data map[string]interface{}) interface{} {
	return nil
}

// UnaryNode 单目运算节点
type UnaryNode struct {
	Val Node
	Op  int
}

// Eval 查询
func (n *UnaryNode) Eval(data map[string]interface{}) interface{} {
	switch n.Op {
	case TYPE_SUB:
		return -toF64(n.Val.Eval(data))
	case TYPE_PLUS:
		return toF64(n.Val.Eval(data))
	}
	return nil
}
