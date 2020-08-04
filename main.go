package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
)

var Literals []string       // The tokens representing literal strings
var Keywords []string       // The keyword tokens
var Tokens []string         // All of the tokens (including literals and keywords)
var TokenIds map[string]int // A map from the token names to their int ids
var Lexer *lexmachine.Lexer // The lexer object. Use this to construct a Scanner

func initTokens() {
	Literals = []string{
		"\"\"",
		"'",
		"`",
		"(",
		")",
	}

	Keywords = []string{
		"SELECT",
		"FROM",
		"AS",
	}

	Tokens = []string{
		"COMMENT",
		"ID",
		"FIELD",
	}

	Tokens = append(Tokens, Keywords...)
	Tokens = append(Tokens, Literals...)
	TokenIds = make(map[string]int)

	for i, tok := range Tokens {
		TokenIds[tok] = i
	}
}

// a lex.Action function with constructs a Token of the given token type by
// the token type's name.
func token(name string) lexmachine.Action {
	return func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
		return s.Token(TokenIds[name], string(m.Bytes), m), nil
	}
}

// a lexmachine.Action function which skips the match.
func skip(*lexmachine.Scanner, *machines.Match) (interface{}, error) {
	return nil, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	lexer := lexmachine.NewLexer()

	for _, lit := range Literals {
		r := "\\" + strings.Join(strings.Split(lit, ""), "\\")
		lexer.Add([]byte(r), token(lit))
	}

	for _, name := range Keywords {
		lexer.Add([]byte(strings.ToLower(name)), token(name))
	}

	lexer.Add([]byte(`//[^\n]*\n?`), token("COMMENT"))
	lexer.Add([]byte(`/\*([^*]|\r|\n|(\*+([^*/]|\r|\n)))*\*+/`), token("COMMENT"))
	lexer.Add([]byte(`([a-z]|[A-Z]|[0-9]|_)+`), token("ID"))
	lexer.Add([]byte(`([a-z]|[A-Z]|[0-9]|_)+\.([a-z]|[A-Z]|[0-9]|_)+`), token("FIELD"))

	lexer.Add([]byte("( |\t|\n|\r)+"), skip)

	err := lexer.Compile()
	if err != nil {
		log.Fatal(err)
	}

	s, err := lexer.Scanner([]byte(`select t.a
     // haha cc
	/* comm */ from (select a from tb) as t`))
	if err != nil {
		log.Fatal(err)
	}

	for tok, err, eof := s.Next(); !eof; tok, err, eof = s.Next() {
		if ui, is := err.(*machines.UnconsumedInput); is {
			// to skip bad token do:
			// s.TC = ui.FailTC
			log.Println(ui)
			log.Fatal(err) // however, we will just fail the program
		} else if err != nil {
			log.Fatal(err)
		}

		token := tok.(*lexmachine.Token)
		fmt.Printf("%+v\n", token)
		fmt.Println("token:", token.String())
		fmt.Println("type:", token.Type)
	}

	fmt.Println("vim-go")
}
