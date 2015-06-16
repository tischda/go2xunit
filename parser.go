package main

import (
	"fmt"
	"io"
	"log"
)

type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// Parse parses a test statement.
func (p *Parser) Parse() (*Test, error) {
	test := &Test{}

	// First token should be a "===" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != EQUALS {
		return nil, fmt.Errorf("found %q, expected ===", lit)
	}

	if tok, lit := p.scanIgnoreWhitespace(); tok != RUN {
		return nil, fmt.Errorf("found %q, expected RUN", lit)
	}

	// Read test name
	if tok, lit := p.scanIgnoreWhitespace(); tok != IDENT {
		return nil, fmt.Errorf("found %q, expected field", lit)
	} else {
		test.Name = lit
	}

	// Next token should be a "---" keyword.
	if tok, _ := p.scanIgnoreWhitespace(); tok != DASHES {
		log.Println("not implemented: this is log output or panic")
	}

	status, _ := p.scanIgnoreWhitespace()
	switch status {
	case PASS:
		test.Passed = true
	case FAIL:
		test.Failed = true
	case SKIP:
		test.Skipped = true
	}

	// ignore colon
	p.scanIgnoreWhitespace()

	// Read test name again
	if tok, lit := p.scanIgnoreWhitespace(); tok != IDENT {
		return nil, fmt.Errorf("found %q, expected field", lit)
	} else {
		if test.Name != lit {
			return nil, fmt.Errorf("found %q, expected %q", lit, test.Name)
		}
	}

	// Read time
	p.scanIgnoreWhitespace()
	if tok, lit := p.scanIgnoreWhitespace(); tok != IDENT {
		return nil, fmt.Errorf("found %q, expected field", lit)
	} else {
		test.Time = lit
	}

	// Return the successfully parsed statement.
	return test, nil
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS || tok == COLON {
		tok, lit = p.scan()
	}
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }
