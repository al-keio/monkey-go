package ast

import (
	"bytes"
	"strings"

	"github.com/al-keio/monkey-go/token"
)

type Node interface {
	TokenLiteral() string
	String() string
	Copy() Node
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
func (p *Program) Copy() Node {
	statements := []Statement{}
	for _, stmt := range p.Statements {
		statements = append(statements, stmt.Copy().(Statement))
	}
	return &Program{Statements: statements}
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}
func (ls *LetStatement) Copy() Node {
	return &LetStatement{Token: ls.Token, Name: ls.Name.Copy().(*Identifier), Value: ls.Value.Copy().(Expression)}
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}
func (rs *ReturnStatement) Copy() Node {
	return &ReturnStatement{Token: rs.Token, ReturnValue: rs.ReturnValue.Copy().(Expression)}
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}
func (es *ExpressionStatement) Copy() Node {
	return &ExpressionStatement{Token: es.Token, Expression: es.Expression.Copy().(Expression)}
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string {
	return i.Value
}
func (i *Identifier) Copy() Node {
	return &Identifier{Token: i.Token, Value: i.Value}
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }
func (il *IntegerLiteral) Copy() Node {
	return &IntegerLiteral{Token: il.Token, Value: il.Value}
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }
func (sl *StringLiteral) Copy() Node {
	return &StringLiteral{Token: sl.Token, Value: sl.Value}
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }
func (b *Boolean) Copy() Node {
	return &Boolean{Token: b.Token, Value: b.Value}
}

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
func (al *ArrayLiteral) Copy() Node {
	elements := []Expression{}
	for _, el := range al.Elements {
		elements = append(elements, el.Copy().(Expression))
	}
	return &ArrayLiteral{Token: al.Token, Elements: elements}
}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for key, value := range hl.Pairs {
		elements = append(elements, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("}")

	return out.String()
}
func (hl *HashLiteral) Copy() Node {
	pairs := make(map[Expression]Expression)
	for key, value := range hl.Pairs {
		pairs[key.Copy().(Expression)] = value.Copy().(Expression)
	}
	return &HashLiteral{Token: hl.Token, Pairs: pairs}
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("]")
	out.WriteString(")")

	return out.String()
}
func (ie *IndexExpression) Copy() Node {
	return &IndexExpression{Token: ie.Token, Left: ie.Left.Copy().(Expression), Index: ie.Index.Copy().(Expression)}
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}
func (pe *PrefixExpression) Copy() Node {
	return &PrefixExpression{Token: pe.Token, Operator: pe.Operator, Right: pe.Right.Copy().(Expression)}
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}
func (ie *InfixExpression) Copy() Node {
	return &InfixExpression{Token: ie.Token, Left: ie.Left.Copy().(Expression), Operator: ie.Operator, Right: ie.Right.Copy().(Expression)}
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}
func (ie *IfExpression) Copy() Node {
	return &IfExpression{Token: ie.Token, Condition: ie.Condition.Copy().(Expression), Consequence: ie.Consequence.Copy().(*BlockStatement), Alternative: ie.Alternative.Copy().(*BlockStatement)}
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}
func (bs *BlockStatement) Copy() Node {
	statements := []Statement{}
	for _, stmt := range bs.Statements {
		statements = append(statements, stmt.Copy().(Statement))
	}
	return &BlockStatement{Token: bs.Token, Statements: statements}
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}
func (fl *FunctionLiteral) Copy() Node {
	identifiers := []*Identifier{}
	for _, identifier := range fl.Parameters {
		identifiers = append(identifiers, identifier.Copy().(*Identifier))
	}
	return &FunctionLiteral{Token: fl.Token, Parameters: identifiers, Body: fl.Body.Copy().(*BlockStatement)}
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
func (ce *CallExpression) Copy() Node {
	args := []Expression{}
	for _, arg := range ce.Arguments {
		args = append(args, arg.Copy().(Expression))
	}
	return &CallExpression{Token: ce.Token, Function: ce.Function.Copy().(Expression), Arguments: args}
}

type MacroLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (ml *MacroLiteral) expressionNode()      {}
func (ml *MacroLiteral) TokenLiteral() string { return ml.Token.Literal }
func (ml *MacroLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range ml.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(ml.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(ml.Body.String())

	return out.String()
}
func (ml *MacroLiteral) Copy() Node {
	identifiers := []*Identifier{}
	for _, identifier := range ml.Parameters {
		identifiers = append(identifiers, identifier.Copy().(*Identifier))
	}
	return &MacroLiteral{Token: ml.Token, Parameters: identifiers, Body: ml.Body.Copy().(*BlockStatement)}
}
