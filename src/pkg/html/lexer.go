package html

import (
	"fmt"
	"io"
)

type TokenizerFlag uint8

type TokenType uint8

const (
	TokenTypeDefault TokenType = iota
	TokenTypeInlineText
	TokenTypeDoctype
	TokenTypeComment
	TokenTypeOpening
	TokenTypeClosing
	TokenTypeSelfClosing
)

func (tokenType TokenType) String() string {
	switch tokenType {
	case TokenTypeDefault:
		return "Default"
	case TokenTypeInlineText:
		return "InlineText"
	case TokenTypeDoctype:
		return "Doctype"
	case TokenTypeComment:
		return "Comment"
	case TokenTypeOpening:
		return "Opening"
	case TokenTypeClosing:
		return "Closing"
	case TokenTypeSelfClosing:
		return "SelfClosing"
	default:
		return fmt.Sprintf("%d", tokenType)
	}
}

const (
	TokenizerFlagDefault TokenizerFlag = iota
	TokenizerFlagInlineText
	TokenizerFlagDoctype
	TokenizerFlagComment
	TokenizerFlagOpening
	TokenizerFlagClosing
	TokenizerFlagSelfClosing
)

type Tokenizer struct {
	stream          io.Reader
	currentTagBlock []byte
	tokens          []Token
	flag            TokenizerFlag
	lastChar        byte
}

type Token struct {
	tokenBlock []byte
	tokenType  TokenType
}

func NewLexer(stream io.Reader) *Tokenizer {
	return &Tokenizer{
		stream: stream,
	}
}

func getTokenTypeFromFlag(flag TokenizerFlag) TokenType {
	switch flag {
	case TokenizerFlagInlineText:
		return TokenTypeInlineText
	case TokenizerFlagDoctype:
		return TokenTypeDoctype
	case TokenizerFlagComment:
		return TokenTypeComment
	case TokenizerFlagOpening:
		return TokenTypeOpening
	case TokenizerFlagClosing:
		return TokenTypeClosing
	case TokenizerFlagSelfClosing:
		return TokenTypeSelfClosing
	default:
		return TokenTypeDefault
	}
}

func (t *Tokenizer) Read() {
	for {
		var char [1]byte
		byteCount, err := t.stream.Read(char[:])
		if err != nil && byteCount == 0 {
			// eof
			break
		}
		switch char[0] {
		case '<':
			if t.flag == TokenizerFlagInlineText {
				t.tokens = append(t.tokens, Token{
					tokenBlock: t.currentTagBlock,
					tokenType:  getTokenTypeFromFlag(t.flag),
				})
				t.currentTagBlock = []byte{}
				t.flag = TokenizerFlagDefault
			}
			if t.flag == TokenizerFlagDefault {
				t.flag = TokenizerFlagOpening
			}
			t.currentTagBlock = append(t.currentTagBlock, char[0])
		case '!':
			if t.flag == TokenizerFlagOpening && t.lastChar == '<' {
				t.flag = TokenizerFlagDoctype
			}
			t.currentTagBlock = append(t.currentTagBlock, char[0])
		case '-':
			if t.flag == TokenizerFlagDoctype && t.lastChar == '!' {
				t.flag = TokenizerFlagComment
			}
			t.currentTagBlock = append(t.currentTagBlock, char[0])
		case '/':
			if t.flag == TokenizerFlagOpening && t.lastChar == '<' {
				t.flag = TokenizerFlagClosing
			}
			t.currentTagBlock = append(t.currentTagBlock, char[0])
		case '>':
			if t.flag == TokenizerFlagOpening && t.lastChar == '/' {
				t.flag = TokenizerFlagSelfClosing
			}
			t.currentTagBlock = append(t.currentTagBlock, char[0])
			t.tokens = append(t.tokens, Token{
				tokenBlock: t.currentTagBlock,
				tokenType:  getTokenTypeFromFlag(t.flag),
			})
			t.currentTagBlock = []byte{}
			t.flag = TokenizerFlagDefault
		default:
			if t.flag == TokenizerFlagDefault && t.lastChar == '>' && char[0] != '\n' {
				t.flag = TokenizerFlagInlineText
				t.currentTagBlock = append(t.currentTagBlock, char[0])
			} else if t.flag != TokenizerFlagDefault {
				t.currentTagBlock = append(t.currentTagBlock, char[0])
			}
		}
		t.lastChar = char[0]
	}
	for i, token := range t.tokens {
		fmt.Printf("%5d - (%s): %s\n", i, token.tokenType, string(token.tokenBlock))
	}
}
