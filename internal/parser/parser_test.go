package parser

import (
	"testing"

	"github.com/matryer/is"
)

func TestInputsAll(t *testing.T) {

	tests := []struct {
		input    string
		expected Consumo
	}{
		{
			input: `07/02/23                              Pago interes por saldo en cuenta                                        $ 5,90                              $ 631.288,84
                                      Del 01/01/23 al 31/01/23`,
			expected: Consumo{
				Date:        "07/02/23",
				Code:        "",
				Description: "Pago interes por saldo en cuenta Del 01/01/23 al 31/01/23",
				Amount:      590,
				Balance:     63128884,
			},
		},
		{
			input: `05/07/21               10280171       Una Compra con tarjeta de debito                               -$ 650,00                                $ 104.095,74
    Mercadopago*recargatuenti - tarj nro. 1866`,
			expected: Consumo{

				Date:        "05/07/21",
				Code:        "10280171",
				Description: "Una Compra con tarjeta de debito Mercadopago*recargatuenti - tarj nro. 1866",
				Amount:      -65000,
				Balance:     10409574,
			},
		},
		{
			input: `25/08/21           25593863      Transferencia realizada                                             -$ 6.000,00                                $ 96.424,39
    A ganga carlos ignacio / varios - var / 201645877712`,
			expected: Consumo{
				Date:        "25/08/21",
				Code:        "25593863",
				Description: "Transferencia realizada A ganga carlos ignacio / varios - var / 201645877712",
				Amount:      -600000,
				Balance:     9642439,
			},
		},
		{
			input: `02/12/22                 Compra con tarjeta de debito                                       -$ 548,00                                $ 166.696,92
    Autoservicio santa ana - tarj nro. 1866`,
			expected: Consumo{
				Date:        "02/12/22",
				Code:        "",
				Description: "Compra con tarjeta de debito Autoservicio santa ana - tarj nro. 1866",
				Amount:      -54800,
				Balance:     16669692,
			},
		},
		{
			input: `16/01/23 1899579                 Compra con tarjeta en el exterior                                                     -U$S 3,49          U$S 1.594,74
    Google wm max llc - tarj nro. 1866`,
			expected: Consumo{
				Date:        "16/01/23",
				Code:        "1899579",
				Description: "Compra con tarjeta en el exterior Google wm max llc - tarj nro. 1866",
				Amount:      -349,
				Balance:     159474,
			},
		},
	}

	for i, test := range tests {
		is := is.New(t)

		p := New(test.input)

		consumo, err := p.Parse()

		if err != nil {
			t.Fatalf("Test %d: %s", i, err)
		}

		is.Equal(consumo, test.expected)
	}
}

