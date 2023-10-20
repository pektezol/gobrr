package html

import (
	"fmt"
	"io"
	"strings"
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
	tokenBlock      []byte
	tokenType       TokenType
	tokenAttributes []string
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
				t.insertToken()
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
			t.insertToken()
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
		fmt.Printf("%5d - (%s): %s\t%v\n", i, token.tokenType, string(token.tokenBlock), token.tokenAttributes)
	}
}

func (t *Tokenizer) insertToken() {
	token := Token{
		tokenBlock: t.currentTagBlock,
		tokenType:  getTokenTypeFromFlag(t.flag),
	}
	if t.flag != TokenizerFlagInlineText {
		attributes := strings.Split(string(t.currentTagBlock), " ")[1:]
		if len(attributes) != 0 {
			if attributes[len(attributes)-1][len(attributes[len(attributes)-1])-1] == '>' {
				attributes[len(attributes)-1] = attributes[len(attributes)-1][:len(attributes[len(attributes)-1])-1]
			}
			token.tokenAttributes = attributes
		}
	}
	t.tokens = append(t.tokens, token)
	t.currentTagBlock = []byte{}
	t.flag = TokenizerFlagDefault
}
