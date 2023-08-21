package lexer

import (
	"fmt"
	"regexp"
	"testing"
)

func TestNextToken(t *testing.T) {

	input := `05/07/21               10280171       Una Compra con tarjeta de debito                               -$ 650,00                                $ 104.095,74
    Mercadopago*recargatuenti - tarj nro. 1866`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{INT, "05"},
		{SLASH, "/"},
		{INT, "07"},
		{SLASH, "/"},
		{INT, "21"},
		{INT, "10280171"},
		{DESC, "Una Compra con tarjeta de debito"},
		{MINUS, "-"},
		{DOLLAR, "$"},
		{INT, "650"},
		{COMMA, ","},
		{INT, "00"},
		{DOLLAR, "$"},
		{INT, "104"},
		{DOT, "."},
		{INT, "095"},
		{COMMA, ","},
		{INT, "74"},
		{DESC, "Mercadopago*recargatuenti - tarj nro. 1866"},
		{EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		fmt.Printf("tok=%v\n", tok)

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] [position: %d] - tokentype wrong. expected=%q, got=%q", i, l.position, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}

}

func TestSinComprobante(t *testing.T) {

	input := `07/02/23                              Pago interes por saldo en cuenta                                        $ 5,90                              $ 631.288,84
                                      Del 01/01/23 al 31/01/23`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{INT, "07"},
		{SLASH, "/"},
		{INT, "02"},
		{SLASH, "/"},
		{INT, "23"},
		{DESC, "Pago interes por saldo en cuenta"},
		{DOLLAR, "$"},
		{INT, "5"},
		{COMMA, ","},
		{INT, "90"},
		{DOLLAR, "$"},
		{INT, "631"},
		{DOT, "."},
		{INT, "288"},
		{COMMA, ","},
		{INT, "84"},
		{DESC, "Del 01/01/23 al 31/01/23"},
		{EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {

		tok := l.NextToken()

		fmt.Printf("tok=%v\n", tok)

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] [position: %d] - tokentype wrong. expected=%q, got=%q", i, l.position, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}

}

func TestConUSD(t *testing.T) {

	input := `16/01/23 1899579                 Compra con tarjeta en el exterior                                                     -U$S 3,49          U$S 1.594,74
                                 Google wm max llc - tarj nro. 1866`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{INT, "16"},
		{SLASH, "/"},
		{INT, "01"},
		{SLASH, "/"},
		{INT, "23"},
		{INT, "1899579"},
		{DESC, "Compra con tarjeta en el exterior"},
		{MINUS, "-"},
		{USD, "U$S"},
		{INT, "3"},
		{COMMA, ","},
		{INT, "49"},
		{USD, "U$S"},
		{INT, "1"},
		{DOT, "."},
		{INT, "594"},
		{COMMA, ","},
		{INT, "74"},
		{DESC, "Google wm max llc - tarj nro. 1866"},
		{EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {

		tok := l.NextToken()

		fmt.Printf("tok=%v\n", tok)

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] [position: %d] - tokentype wrong. expected=%q, got=%q", i, l.position, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}

}

func TestPeekCharAt(t *testing.T) {

	input := "01234     5678"

	l := New(input)

	tok := l.NextToken()

	if tok.Type != INT {
		t.Fatalf("tokentype wrong. expected=%q, got=%q", INT, tok.Type)
	}

	b := l.peekCharAt(5)

	if b != '5' {
		t.Fatalf("expected '5' found=%q", b)
	}

}

func TestReadNChar(t *testing.T) {

	input := "012345"

	l := New(input)

	l.readNChar(4)

	if l.position != 4 {
		t.Fatalf("expected position to be 4, got=%d", l.position)
	}
}

func TestRegex1(t *testing.T) {

	input := `Descubierto                                               Traspaso entre cuentas
Límite: $ 31.000,00             Hasta: 31/05/23           Desactivado




Movimientos en pesos
                                                                                                    Caja de Ahorro en   Cuenta Corriente en
           Fecha      Comprobante     Movimiento                                                                                                  Saldo en cuenta
                                                                                                               pesos                 pesos


27/01/23                              Saldo Inicial                                                  $ 407.575,95                   $ 0,00        $ 407.575,95

30/01/23 88945528                     Debito debin                                                     -$ 3.660,00                                $ 403.915,95
                                      Id debin z0kv8794vp8wml42peydx4 cuit 30703088534

30/01/23              2236282         Compra con tarjeta de debito                                     -$ 1.509,00                                $ 402.406,95
                                      Pedidosya - tarj nro. 1866

30/01/23              42238514        Compra con tarjeta de debito                                          -$ 50,00                              $ 402.356,95
    Pedidos ya - propinas - tarj nro. 1866`

	re := regexp.MustCompile(`(?m)^([0-9]{2}/[0-9]{2}/[0-9]{2})\s+([0-9]+)\s+(.*?)\s{2,}(.*?)\s{2,}(.*)\n(.*)`)
	// re, err := regexp.Compile(`^[0-9]{2}/[0-9]{2}/[0-9]{2}`)
	// re := regexp.MustCompile(`(?m)^30.*`)

	res := re.FindAllString(input, -1)

	if len(res) == 0 {

		t.Fatalf("len 0")

	}

	fmt.Printf("len=%v\n", len(res))

	for _, m := range res {

		fmt.Printf("m=%v\n", m)
	}
}
