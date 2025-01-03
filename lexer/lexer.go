package lexer

import (
	"log"
	"monkey-lang-interpreted/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current character (ch variable))
	readPosition int  // current reading position in input (after current character)
	ch           byte // current character under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// gives us the next character and advance our position in the input string
func (l *Lexer) readChar() {
	/*
		Checks if we've reached the input of the input. If so, assign 0 to ch which in
		ASCII is the "NUL" character. The "NUL" / 0 will signal we finished reading and lexing
		the input (end of file) or we didn't start yet.
	*/
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else { // if not at the end of our input, set the ch value to the next character
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()

	switch l.ch {
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LEFT_PAREN, l.ch)
	case ')':
		tok = newToken(token.RIGHT_PAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '=':
		/**
		if the character after current char is "="
		then it is an equality operator and not assignment
		*/
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar() // advancing the current char so the next char will be used. "=="
			tok = token.Token{Type: token.EQUAL, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQUAL, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '>':
		tok = newToken(token.GREATER_THAN, l.ch)
	case '<':
		tok = newToken(token.LESS_THAN, l.ch)
	case '{':
		tok = newToken(token.LEFT_BRACE, l.ch)
	case '}':
		tok = newToken(token.RIGHT_BRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			// log.Printf("tok: %v", tok)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			log.Printf("illegal l.ch: %v, string: %v", l.ch, string(l.ch))
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	/*
		keeps read chars until we've reached an empty space or non alpha or "_"
		value. Then return the joined string of the starting position to the last known
		position and that will be the full identifier value
	*/
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// iterates through text until we get to a non empty space character
// to parse
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

/*
returns the character after the current character
to decide if the next character is impactful in determining
what the token / identifier is. ex) if we're at "=" and the next
character is "==" then we know it's an equality operator and not
just an assignment operator. It doesn't increment the lexer position
value so we stay there at current position.
*/
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// supported first chars in an identifier
func isLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

// check that the character byte code is either between '0' and '9'
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
