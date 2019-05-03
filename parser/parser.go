package parser

import (
	"fmt"
	"strconv"

	"github.com/skatsuta/monkey-compiler/ast"
	"github.com/skatsuta/monkey-compiler/lexer"
	"github.com/skatsuta/monkey-compiler/token"
)

const (
	_ int = iota
	// LOWEST represents the lowest precedence.
	LOWEST
	// OR represents precedence of logical OR.
	OR
	// AND represents precedence of logical AND.
	AND
	// EQUALS represents precedence of equals.
	EQUALS // ==
	// LESSGREATER represents precedence of less than or greater than.
	LESSGREATER // > or <
	// SUM represents precedence of sum.
	SUM // +
	// PRODUCT represents precedence of product.
	PRODUCT // *
	// PREFIX represents precedence of prefix operator.
	PREFIX // -X or !X
	// CALL represents precedence of function call.
	CALL // myFunc(X)
	// INDEX represents precedence of array index operator.
	INDEX // array[index]
)

var precedences = map[token.Type]int{
	token.Or:       OR,
	token.And:      AND,
	token.Eq:       EQUALS,
	token.NEq:      EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.LE:       LESSGREATER,
	token.GE:       LESSGREATER,
	token.Plus:     SUM,
	token.Minus:    SUM,
	token.Slash:    PRODUCT,
	token.Astarisk: PRODUCT,
	token.LParen:   CALL,
	token.LBracket: INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser is a parser of Monkey programming language.
type Parser struct {
	l      lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

// New returns a new Parser.
func New(l lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = map[token.Type]prefixParseFn{
		token.Ident:    p.parseIdent,
		token.Int:      p.parseIntegerLiteral,
		token.Float:    p.parseFloatLiteral,
		token.Bang:     p.parsePrefixExpression,
		token.Minus:    p.parsePrefixExpression,
		token.True:     p.parseBoolean,
		token.False:    p.parseBoolean,
		token.Nil:      p.parseNil,
		token.LParen:   p.parseGroupedExpression,
		token.If:       p.parseIfExpression,
		token.Function: p.parseFunctionLiteral,
		token.String:   p.parseStringLiteral,
		token.LBracket: p.parseArrayLiteral,
		token.LBrace:   p.parseHashLiteral,
		token.Macro:    p.parseMacroLiteral,
	}

	p.infixParseFns = map[token.Type]infixParseFn{
		token.Plus:     p.parseInfixExpression,
		token.Minus:    p.parseInfixExpression,
		token.Astarisk: p.parseInfixExpression,
		token.Slash:    p.parseInfixExpression,
		token.Eq:       p.parseInfixExpression,
		token.NEq:      p.parseInfixExpression,
		token.LT:       p.parseInfixExpression,
		token.GT:       p.parseInfixExpression,
		token.LE:       p.parseInfixExpression,
		token.GE:       p.parseInfixExpression,
		token.And:      p.parseInfixExpression,
		token.Or:       p.parseInfixExpression,
		token.LParen:   p.parseCallExpression,
		token.LBracket: p.parseIndexExpression,
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Errors returns error messages.
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(typ token.Type) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", typ, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) curTokenIs(typ token.Type) bool {
	return p.curToken.Type == typ
}

func (p *Parser) peekTokenIs(typ token.Type) bool {
	return p.peekToken.Type == typ
}

func (p *Parser) expectPeek(typ token.Type) bool {
	if p.peekTokenIs(typ) {
		p.nextToken()
		return true
	}

	p.peekError(typ)
	return false
}

// ParseProgram parses a program and returns a new Program AST node.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.Let:
		return p.parseLetStatement()
	case token.Ident, token.Int, token.Float, token.String, token.Function, token.LParen,
		token.LBracket, token.Minus, token.Bang:
		return p.parseSimpleStatement()
	case token.Return:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.Ident) {
		return nil
	}

	stmt.Name = &ast.Ident{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectPeek(token.Assign) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if fl, ok := stmt.Value.(*ast.FunctionLiteral); ok {
		fl.Name = stmt.Name.Value
	}

	for p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseSimpleStatement() (stmt ast.Statement) {
	lhs := p.parseExpression(LOWEST)

	switch p.peekToken.Type {
	case token.Assign, token.AddAssign, token.SubAssign, token.MulAssign, token.DivAssign:
		p.nextToken()

		tok := p.curToken

		p.nextToken()

		rhs := p.parseExpression(LOWEST)

		// Give an anonymous closure a variable name
		if fl, ok := rhs.(*ast.FunctionLiteral); ok {
			if ident, ok := lhs.(*ast.Ident); ok {
				fl.Name = ident.Value
			}
		}

		stmt = &ast.AssignStatement{Token: tok, LHS: lhs, RHS: rhs}

	default:
		// Expression
		stmt = &ast.ExpressionStatement{Token: p.curToken, Expression: lhs}
	}

	for p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{
		Token: p.curToken,
	}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	for p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Token:      p.curToken,
		Expression: p.parseExpression(LOWEST),
	}

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		msg := fmt.Sprintf("no prefix parse function for %s found", p.curToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	expr := prefix()

	for !p.curTokenIs(token.Semicolon) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return expr
		}

		p.nextToken()

		expr = infix(expr)
	}

	return expr
}

func (p *Parser) parseIdent() ast.Expression {
	return &ast.Ident{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	tok := p.curToken

	val, err := strconv.ParseInt(tok.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", tok.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	return &ast.IntegerLiteral{Token: tok, Value: val}
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	tok := p.curToken

	val, err := strconv.ParseFloat(tok.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", tok.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	return &ast.FloatLiteral{Token: tok, Value: val}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	tok := p.curToken

	p.nextToken()

	return &ast.PrefixExpression{
		Token:    tok,
		Operator: tok.Literal,
		Right:    p.parseExpression(PREFIX),
	}

}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	tok := p.curToken
	prec := p.curPrecedence()

	p.nextToken()

	return &ast.InfixExpression{
		Token:    tok,
		Operator: tok.Literal,
		Left:     left,
		Right:    p.parseExpression(prec),
	}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.curToken,
		Value: p.curTokenIs(token.True),
	}
}

func (p *Parser) parseNil() ast.Expression {
	return &ast.Nil{Token: p.curToken}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	expr := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RParen) {
		return nil
	}

	return expr
}

func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LParen) {
		return nil
	}

	p.nextToken()

	expr.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RParen) {
		return nil
	}

	if !p.expectPeek(token.LBrace) {
		return nil
	}

	expr.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.Else) {
		p.nextToken()

		if !p.expectPeek(token.LBrace) {
			return nil
		}

		expr.Alternative = p.parseBlockStatement()
	}

	return expr
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token:      p.curToken,
		Statements: []ast.Statement{},
	}

	p.nextToken()

	for !p.curTokenIs(token.RBrace) && !p.curTokenIs(token.EOF) {
		if stmt := p.parseStatement(); stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LParen) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBrace) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Ident {
	idents := []*ast.Ident{}

	if p.peekTokenIs(token.RParen) {
		p.nextToken()
		return idents
	}

	p.nextToken()

	ident := &ast.Ident{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	idents = append(idents, ident)

	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Ident{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
		idents = append(idents, ident)
	}

	if !p.expectPeek(token.RParen) {
		return nil
	}

	return idents
}

func (p *Parser) parseExpressionList(end token.Type) []ast.Expression {
	list := make([]ast.Expression, 0)

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	return &ast.CallExpression{
		Token:     p.curToken,
		Function:  function,
		Arguments: p.parseExpressionList(token.RParen),
	}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	return &ast.ArrayLiteral{
		Token:    p.curToken,
		Elements: p.parseExpressionList(token.RBracket),
	}
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expr := &ast.IndexExpression{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken()
	expr.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBracket) {
		return nil
	}

	return expr
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{
		Token: p.curToken,
		Pairs: make(map[ast.Expression]ast.Expression),
	}

	for !p.peekTokenIs(token.RBrace) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.Colon) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)
		hash.Pairs[key] = value

		if !p.peekTokenIs(token.RBrace) && !p.expectPeek(token.Comma) {
			return nil
		}
	}

	if !p.expectPeek(token.RBrace) {
		return nil
	}

	return hash
}

func (p *Parser) parseMacroLiteral() ast.Expression {
	tok := p.curToken

	if !p.expectPeek(token.LParen) {
		return nil
	}

	params := p.parseFunctionParameters()

	if !p.expectPeek(token.LBrace) {
		return nil
	}

	body := p.parseBlockStatement()

	return &ast.MacroLiteral{
		Token:      tok,
		Parameters: params,
		Body:       body,
	}
}
