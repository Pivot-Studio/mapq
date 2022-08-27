package mapq

import (
	"fmt"
)

const (
	TYPE_PLUS      = iota // "+"
	TYPE_SUB              // "-"
	TYPE_MUL              // "*"
	TYPE_DIV              // "/"
	TYPE_LP               // "("
	TYPE_RP               // ")"
	TYPE_VAR              // "([a-z]|[A-Z])([a-z]|[A-Z]|[0-9])*"
	TYPE_RES_TRUE         // "true"
	TYPE_RES_FALSE        // "false"
	TYPE_AND              // "&&"
	TYPE_OR               // "||"
	TYPE_EQ               // "=="
	TYPE_LG               // ">"
	TYPE_SM               // "<"
	TYPE_LEQ              // ">="
	TYPE_SEQ              // "<="
	TYPE_NEQ              // "!="
	TYPE_STR              // a quoted string(单引号)
	TYPE_INT              // an integer
	TYPE_FLOAT            // 小数，x.y这种
	TYPE_UNKNOWN          // 未知的类型
	TYPE_NOT              // "!"
	TYPE_DOT              // "."
	TYPE_RES_NULL         // "null"
)

var (
	reserved = map[string]int{
		"true":  TYPE_RES_TRUE,
		"false": TYPE_RES_FALSE,
		"null":  TYPE_RES_NULL,
	}
	ErrEOS  = fmt.Errorf("eos error")
	ErrTYPE = fmt.Errorf("the next token doesn't match the expected type")
)

// Lexer 词法分析器
type Lexer struct {
	input string
	pos   int
	runes []rune
}

// SetInput 设置输入
func (l *Lexer) SetInput(s string) {
	l.pos = 0
	l.input = s
	l.runes = []rune(s)
}

func (l *Lexer) getCh() (ch rune, end bool) {
	defer func() {
		l.pos++
	}()
	return l.Peek()
}

// Peek 看下一个字符
func (l *Lexer) Peek() (ch rune, end bool) {
	le := len(l.runes)
	if l.pos >= le {
		return ch, true
	}
	ch = l.runes[l.pos]
	return ch, false
}

func (l *Lexer) getChSkipEmpty() (ch rune, end bool) {
	ch, end = l.getCh()
	if end {
		return
	}
	if ch == ' ' || ch == '\t' {
		return l.getChSkipEmpty()
	}
	return
}
func isLetter(ch rune) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}
func isLetterOrUnderscore(ch rune) bool {
	return isLetter(ch) || ch == '_'
}
func isNum(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

// Checkpoint 检查点
type Checkpoint struct {
	pos int
}

// SetCheckpoint 设置检查点
func (l *Lexer) SetCheckpoint() Checkpoint {
	return Checkpoint{
		pos: l.pos,
	}
}

// GobackTo 回到一个检查点
func (l *Lexer) GobackTo(c Checkpoint) {
	l.pos = c.pos
}

// GetPos 获取现在位置
func (l *Lexer) GetPos() int {
	return l.pos
}

// Currpos 获取现在位置
func (l *Lexer) Currpos(pos int) (line, off int) {
	line = 1
	last := 0
	for i, v := range l.runes[:pos] {
		if v == '\n' {
			line++
			last = i
		}
	}
	for i := 0; true; i++ {
		if l.runes[pos+i] != ' ' {
			return line, pos + i - last
		}
	}
	return

}

// ScanType 扫描一个特定Token，下一个token不是这个类型则自动回退，返回err
func (l *Lexer) ScanType(code int) (token string, err error) {
	if code == TYPE_LG {
		chp := l.SetCheckpoint()
		ch, end := l.getChSkipEmpty()
		if end {
			return "", ErrEOS
		}
		if ch == '>' {
			return ">", nil
		}
		l.GobackTo(chp)
	}
	ch := l.SetCheckpoint()
	c, t, e := l.Scan()
	if c == code {
		return t, nil
	} else if e {
		return "", ErrEOS
	}
	// fmt.Println(pos, t)
	l.GobackTo(ch)
	return "", ErrTYPE
}

func (l *Lexer) scanStr(ch rune, end bool) (code int, token string, eos bool) {

	if ch == '\'' {
		i := []rune{}
		for {
			c, end := l.getCh()
			if end {
				if end {
					return TYPE_UNKNOWN, "", true
				}
			}
			if c == '\\' {
				c, end := l.getCh()
				if end {
					return TYPE_UNKNOWN, "", true
				}
				switch c {
				case '\'':
					i = append(i, c)
				case 't':
					i = append(i, '\t')
				case 'n':
					i = append(i, '\n')
				case 'r':
					i = append(i, '\r')
				case '\\':
					i = append(i, '\\')
				case '0':
					i = append(i, '\x00')
				default:
					return TYPE_UNKNOWN, "", false
				}
				continue

			}
			if c == '\'' {
				break
			}
			i = append(i, c)
		}
		return TYPE_STR, string(i), end
	}
	return -1, "", false
}

func (l *Lexer) scanVar(ch rune, end bool) (code int, token string, eos bool) {
	i := []rune{ch}
	for {
		c, end := l.getCh()
		if end {
			break
		}
		if !isLetterOrUnderscore(c) && !isNum(c) {
			l.pos--
			break
		}
		i = append(i, c)
	}
	token = string(i)
	if tp, ok := reserved[token]; ok {
		return tp, token, end
	}
	return TYPE_VAR, string(i), end
}

func (l *Lexer) scanNum(ch rune, end bool) (code int, token string, eos bool) {
	i := []rune{ch}
	t := TYPE_INT
	next, _ := l.Peek()
	if ch == '0' && (next == 'b' ||
		next == 'o' ||
		next == 'x') {
		l.getCh()
		i = append(i, next)
	}
	for {
		c, end := l.getCh()
		if end {
			break
		}
		if c == '.' {
			i = append(i, c)
			t = TYPE_FLOAT
			continue
		}
		if !isNum(c) {
			l.pos--
			break
		}
		i = append(i, c)
	}
	return t, string(i), end
}

func (l *Lexer) scanNumericOp(ch rune, end bool) (code int, token string, eos bool) {
	switch ch {
	case '+':
		return TYPE_PLUS, "+", end
	case '-':
		return TYPE_SUB, "-", end
	case '*':
		return TYPE_MUL, "*", end
	case '.':
		return TYPE_DOT, ".", end
	case '/':
		return TYPE_DIV, "/", end
	case '(':
		return TYPE_LP, "(", end
	case ')':
		return TYPE_RP, ")", end
	}
	return -1, "", false
}

func (l *Lexer) scanLogicOp(ch rune, end bool) (code int, token string, eos bool) {
	switch ch {
	case '=':
		if ne, _ := l.Peek(); ne == '=' {
			l.getCh()
			return TYPE_EQ, "==", end
		}
		return TYPE_UNKNOWN, "=", end
	case '&':
		if ne, _ := l.Peek(); ne == '&' {
			l.getCh()
			return TYPE_AND, "&&", end
		}
		return TYPE_UNKNOWN, "&", end
	case '|':
		if ne, _ := l.Peek(); ne == '|' {
			l.getCh()
			return TYPE_OR, "||", end
		}
		return TYPE_OR, "|", end
	case '>':
		ne, _ := l.Peek()
		if ne == '=' {
			l.getCh()
			return TYPE_LEQ, ">=", end
		}
		return TYPE_LG, ">", end
	case '<':
		ne, _ := l.Peek()
		if ne == '=' {
			l.getCh()
			return TYPE_SEQ, "<=", end
		}
		return TYPE_SM, "<", end
	case '!':
		if ne, _ := l.Peek(); ne == '=' {
			l.getCh()
			return TYPE_NEQ, "!=", end
		}
		return TYPE_NOT, "!", end
	}
	return -1, "", false
}

// Scan scan a token
func (l *Lexer) Scan() (code int, token string, eos bool) {
	ch, end := l.getChSkipEmpty()
	if end {
		eos = end
		return
	}
	code, token, eos = l.scanStr(ch, end)
	if code != -1 {
		return
	}
	if isLetterOrUnderscore(ch) {
		return l.scanVar(ch, end)
	}
	if isNum(ch) {
		return l.scanNum(ch, end)
	}
	code, token, eos = l.scanNumericOp(ch, end)
	if code != -1 {
		return
	}
	code, token, eos = l.scanLogicOp(ch, end)
	if code != -1 {
		return
	}
	return TYPE_UNKNOWN, "", true

}
