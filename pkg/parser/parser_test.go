package parser

import (
	"fmt"
	"github.com/francoganga/pagoda_bun/pkg/internal/lexer"
	"testing"
)

func TestParseDate(t *testing.T) {

	someDate := "05/08/22"

	input := fmt.Sprintf("%s           99999999", someDate)

	l := lexer.New(input)

	p := New(l)

	d := p.parseDate()

	checkParserErrors(t, p)

	if d != someDate {
		t.Fatalf("Error parsing date, expected=%s, got=%s", someDate, d)
	}
}

func TestParseAmount(t *testing.T) {
	input := "$ 104.095,74"

	expectedAmount := 10409574

	l := lexer.New(input)

	p := New(l)

	a := p.parseAmount()

	checkParserErrors(t, p)

	if a != expectedAmount {
		t.Fatalf("Error parsing amount, expected=%d, got=%d", expectedAmount, a)
	}

}

func TestParseConsumo(t *testing.T) {

	input := `05/07/21               10280171       Compra con tarjeta de debito                               -$ 650,00                                $ 104.095,74
    Mercadopago*recargatuenti - tarj nro. 1866`

	l := lexer.New(input)

	p := New(l)

	c := p.ParseConsumo()

	checkParserErrors(t, p)

	if c.Date != "05/07/21" {
		t.Fatalf("expected string to be=%s, got=%s", "05/07/2021", c.Date)
	}

	if c.Code != "10280171" {
		t.Fatalf("expected code to be=%s, got=%s", "10280171", c.Code)
	}

	if c.Description != "Compra con tarjeta de debito: Mercadopago*recargatuenti - tarj nro. 1866" {
		t.Fatalf("expected Description to be=%s, got=%s", "Compra con tarjeta de debito Mercadopago*recargatuenti - tarj nro. 1866", c.Description)
	}

	if c.Amount != -65000 {
		t.Fatalf("expected Amount to be=%d, got=%d", 6500, c.Amount)
	}

	if c.Balance != 10409574 {
		t.Fatalf("expected Balance to be=%d, got=%d", 10409574, c.Balance)
	}

}

func TestParseConsumo2(t *testing.T) {

	input := `25/08/21           25593863      Transferencia realizada                                             -$ 6.000,00                                $ 96.424,39
    A ganga carlos ignacio / varios - var / 201645877712`

	l := lexer.New(input)

	p := New(l)

	c := p.ParseConsumo()

	checkParserErrors(t, p)

	fmt.Printf("c=%v\n", c)

}

func TestParseConsumo3(t *testing.T) {
	input := `06/08/21               79378436       Acreditacion de haberes                                        $ 105.319,15                                 $ 185.136,86
    307113401661 210805007universidad nacional a jauretc`

	l := lexer.New(input)

	p := New(l)

	c := p.ParseConsumo()

	checkParserErrors(t, p)

	fmt.Printf("c=%v\n", c)
}

func TestParseConsumo4(t *testing.T) {

	input := `02/12/22 1714314                      Compra con tarjeta de debito                                       -$ 548,00                                $ 166.696,92
    Autoservicio santa ana - tarj nro. 1866`

	p := FromInput(input)

	c := p.ParseConsumo()

	checkParserErrors(t, p)

	fmt.Printf("c=%v\n", c)
}

func TestParseConsumoNoCode(t *testing.T) {

	input := `02/12/22                 Compra con tarjeta de debito                                       -$ 548,00                                $ 166.696,92
    Autoservicio santa ana - tarj nro. 1866`

	p := FromInput(input)

	c := p.ParseConsumo()

	checkParserErrors(t, p)

	fmt.Printf("c=%v\n", c)
}

func TestParseConsumoUSD(t *testing.T) {

	input := `16/01/23 1899579                 Compra con tarjeta en el exterior                                                     -U$S 3,49          U$S 1.594,74
    Google wm max llc - tarj nro. 1866`

	p := FromInput(input)

	c := p.ParseConsumo()

	checkParserErrors(t, p)

	fmt.Printf("c=%v\n", c)
}

// TODO: Maybe refactor
// parseAmount starts by doing nexToken so for testing purposes
// i need to have a token before the amount
func TestParseAmountUSD(t *testing.T) {
	input := "9 -U$S 3,49"

	p := FromInput(input)

	amount := p.parseAmount()

	checkParserErrors(t, p)

	fmt.Printf("amount=%v\n", amount)
}

func TestParseAmountNegative(t *testing.T) {
	input := "9 -$ 3,49"

	p := FromInput(input)

	amount := p.parseAmount()

	checkParserErrors(t, p)

	if amount != -349 {
		t.Fatalf("expected amount to be %d, got=%d", -349, amount)
	}

	fmt.Printf("amount=%v\n", amount)
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}
