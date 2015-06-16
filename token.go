package main

// Token represents a lexical token.
type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Literals
	IDENT // fields, test_name

	// Misc characters
	COLON // :
	LP    // (
	RP    // )

	// Keywords
	EQUALS // ===
	DASHES // ---
	RUN
	PASS
	FAIL
	SKIP
	OK      // ok
	SECONDS // seconds
)
