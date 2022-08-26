package mapq

import (
	"fmt"
	"strconv"
)

// Parser 语法分析器
type Parser struct {
	lexer *Lexer
}

func (p *Parser) boolexp() (node Node, err error) {
	ch := p.lexer.SetCheckpoint()
	defer func() {
		if err != nil {
			p.lexer.GobackTo(ch)
		}
	}()
	node, err = p.compare()
	if err != nil {
		return nil, err
	}
	for {
		cp := p.lexer.SetCheckpoint()
		co, _, eos := p.lexer.Scan()
		if eos || !(co == TYPE_AND || co == TYPE_OR) {
			p.lexer.GobackTo(cp)
			break
		}
		n := &BinNode{}
		n.Left = node
		n.Op = co
		node, err = p.compare()
		if err != nil {
			return nil, err
		}
		n.Right = node
		node = n
	}
	return
}

func (p *Parser) boolean() (node Node, err error) {
	ch1 := p.lexer.SetCheckpoint()
	defer func() {
		if err != nil {
			p.lexer.GobackTo(ch1)
		}
	}()
	_, err = p.lexer.ScanType(TYPE_RES_TRUE)
	if err == nil {
		return &BoolConstNode{Val: true}, nil
	}
	_, err = p.lexer.ScanType(TYPE_RES_FALSE)
	if err == nil {
		return &BoolConstNode{Val: false}, nil
	}

	n, err := p.addedFactor()

	if err == nil {
		return n, nil
	}

	code, _, eos := p.lexer.Scan()
	if eos {
		return nil, ErrEOS
	}
	switch code {
	case TYPE_NOT:
		node, err = p.boolean()
		if err != nil {
			return nil, err
		}
		return &NotNode{Bool: node}, nil
	case TYPE_LP:
		node, err = p.boolexp()
		if err != nil {
			return nil, err
		}
		_, err = p.lexer.ScanType(TYPE_RP)
		if err != nil {
			return nil, err
		}
		return

	}

	return nil, fmt.Errorf("parse failed")
}

func (p *Parser) compare() (node Node, err error) {
	ch := p.lexer.SetCheckpoint()
	defer func() {
		if err != nil {
			p.lexer.GobackTo(ch)
		}
	}()

	bn, err := p.boolean()
	if err != nil {
		return nil, err
	}
	for {
		n := &BinNode{}
		check := p.lexer.SetCheckpoint()
		code, _, eos := p.lexer.Scan()
		if eos {
			return bn, nil
		}
		switch code {
		case TYPE_EQ, TYPE_NEQ,
			TYPE_LG, TYPE_SM,
			TYPE_LEQ, TYPE_SEQ:
			n.Op = code
			n.Left = bn
		default:
			p.lexer.GobackTo(check)
			return bn, nil
		}
		n.Right, err = p.boolean()
		if err != nil {
			return nil, err
		}
		bn = n
	}
}

func (p *Parser) factor() (n Node, err error) {
	check := p.lexer.SetCheckpoint()
	defer p.goback(check, err)
	a, err := p.symbol()
	if err != nil {
		return nil, err
	}
	ch := p.lexer.SetCheckpoint()
	code, _, eos := p.lexer.Scan()
	for !eos && code == TYPE_DIV ||
		code == TYPE_MUL {
		b, err := p.symbol()
		if err != nil {
			return nil, err
		}
		a = &BinNode{
			Op:    code,
			Left:  a,
			Right: b,
		}
		ch = p.lexer.SetCheckpoint()
		code, _, eos = p.lexer.Scan()
	}
	if !eos {
		p.lexer.GobackTo(ch)
	}
	return a, nil
}

func (p *Parser) addedFactor() (n Node, err error) {
	check := p.lexer.SetCheckpoint()
	defer p.goback(check, err)
	a, err := p.factor()
	if err != nil {
		return nil, err
	}
	ch := p.lexer.SetCheckpoint()
	code, _, eos := p.lexer.Scan()
	for !eos && (code == TYPE_PLUS || code == TYPE_SUB) {
		b, err := p.factor()
		if err != nil {
			return nil, err
		}
		a = &BinNode{
			Op:    code,
			Left:  a,
			Right: b,
		}
		ch = p.lexer.SetCheckpoint()
		code, _, eos = p.lexer.Scan()
	}
	if !eos {
		p.lexer.GobackTo(ch)
	}
	return a, nil
}

func (p *Parser) goback(ch Checkpoint, err error) {
	if err != nil {
		p.lexer.GobackTo(ch)
	}
}

func (p *Parser) symbol() (n Node, err error) {
	ch := p.lexer.SetCheckpoint()
	defer func() {
		if err != nil {
			p.lexer.GobackTo(ch)
		}
	}()
	code, _, eos := p.lexer.Scan()
	if eos {
		return nil, ErrEOS
	}
	if code == TYPE_PLUS || code == TYPE_SUB {
		num, err := p.number()
		if err != nil {
			return nil, err
		}
		return &UnaryNode{Op: code, Val: num}, nil
	}
	p.lexer.GobackTo(ch)
	return p.number()
}

func (p *Parser) number() (n Node, err error) {
	ch := p.lexer.SetCheckpoint()
	defer func() {
		if err != nil {
			p.lexer.GobackTo(ch)
		}
	}()
	str, err := p.str()
	if err == nil {
		return str, nil
	}

	v, err := p.varblock()
	if err == nil {
		return v, nil
	}
	_, err = p.lexer.ScanType(TYPE_RES_NULL)
	if err == nil {
		return &NullNode{}, nil
	}
	code, t1, eos := p.lexer.Scan()
	if eos {
		return nil, ErrEOS
	}
	switch code {
	case TYPE_FLOAT:
		i, err := strconv.ParseFloat(t1, 64)
		if err != nil {
			return nil, err
		}
		return &NumNode{Num: i}, nil
	case TYPE_INT:
		i, err := strconv.ParseInt(t1, 10, 64)
		if err != nil {
			return nil, err
		}
		return &NumNode{Num: float64(i)}, nil

	}
	p.lexer.GobackTo(ch)
	_, err = p.lexer.ScanType(TYPE_LP)
	if err != nil {
		return nil, err
	}
	i, err := p.boolexp()
	if err != nil {
		return nil, err
	}
	_, err = p.lexer.ScanType(TYPE_RP)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (p *Parser) varblock() (n *VarNode, err error) {
	ch := p.lexer.SetCheckpoint()
	defer func() {
		if err != nil {
			p.lexer.GobackTo(ch)
		}
	}()
	token, err := p.lexer.ScanType(TYPE_VAR)
	if err != nil {
		return nil, err
	}
	n = &VarNode{Token: token}
	_, err = p.lexer.ScanType(TYPE_DOT)
	if err != nil {
		return n, nil
	}
	n.Next, err = p.varblock()
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (p *Parser) str() (n Node, err error) {
	str, err := p.lexer.ScanType(TYPE_STR)
	if err != nil {
		return nil, err
	}

	return &StringNode{Str: str}, nil
}

// Parse 生成ast
func (p *Parser) Parse(str string) (n Node, err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	p.lexer = &Lexer{}
	p.lexer.SetInput(str)
	return p.boolexp()
}
