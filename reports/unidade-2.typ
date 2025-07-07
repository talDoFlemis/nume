#import "@preview/lovelace:0.3.0": *
#import "@preview/algorithmic:1.0.0"
#import algorithmic: algorithm, algorithm-figure, style-algorithm
#show: style-algorithm;
#set text(lang: "br")
#set page(numbering: "1")

#set heading(numbering: "1.")


= Integrais com Newton-Cotes

As integrais com Newton-Cotes são uma forma de aproximar a integral de uma função usando polinômios interpoladores. A ideia é usar pontos igualmente espaçados para avaliar a função e construir um polinômio que passe por esses pontos usando a seguinte função base.

#algorithm-figure("Newton-cotes", {
  import algorithmic: *
  Function("Newton-Cotes", ("strategy", "left", "right", "partitions"), {
    Assign[area][$0$]
    Assign[delta][($"right" - "left"$) / $"partitions"$]
    For([$l$ := left until $"right"$ with step delta ], {
      Assign[area][area + strategy.interpolate($l$, $l$ + delta)]
    })
    Return($"area"$)
  })
}) <newton-cotes-base>

Onde a função `strategy.interpolate` é definida com base em algum polinômio interpolador, `left` e `right` são os limites da integral, e `partitions` é o número de partições que você deseja usar para a aproximação.

Dentre as várias estratégias de interpolação, dividimos em duas categorias: as chamadas abertas, que usam os limites de integração e as fechadas que são usadas quando a função vai ao infinito ou não é definida nos limites de integração, mas possui uma integral no intervalo completo.

Segue a baixo um quadro com as estratégias de interpolação mais comuns, quantidade de pontos usados, classificação em aberta ou fechada e a fórmula de interpolação usada.

#figure(
  table(
    columns: 5,
    "Strategy", "Pontos", "Tipo", "Delta", "Fórmula",
    table.hline(),
    "Regra do trapézio", "2", "Fechada", $(b - a)$, $h/2 (f(a) + f(b))$,
    "Regra de Simpson 1/3", "3", "Fechada", $(b - a)/2$, $h/3 (f(a) + 4f(a + h) + f(b))$,
    "Regra de Simpson 3/8", "4", "Fechada", $(b - a)/3$, $3h/8 (f(a) + 3f(a + h) + 3f(a + 2h) + f(b))$,
    table.hline(),
    "Regra do trapézio aberta", "2", "Aberta", $(b - a)/3$, $3h/2 (f(a + h) + f(a + 2h))$,
    "Regra de Milne", "3", "Aberta", $(b - a)/4$, $4h/3 (2f(a + h) - f(a + 2h) + 2f(a + 3h))$,
    "Fórmula aberta 3º grau", "4", "Aberta", $(b - a)/5$, $5h/24 (11f(a + h) + f(a + 2h) + f(a + 3h) + 11f(a + 4h))$,
    table.hline(),
  ),
  caption: [Tabela de estratégias de Newton-Cotes, com pontos, tipo, cálculo de delta e fórmula.],
) <table-newton-cotes>


Para comprovar a implementação dessas estratégias foi usado o seguinte código:

```go
func TestNewtonCotes(t *testing.T) {
	// Arrange
	t.Parallel()

	strategies := []NewtonCotesStrategy{
		&OpenTrapezoidalRule{},
		&MilneRule{},
		&ThirdDegreeOpenNewtonCotesStrategy{},
		&TrapezoidalRule{},
		&SimpsonsOneThirdRule{},
		&SimpsonsThreeEighthsRule{},
	}

	testCases := []newtonCotesTestCase{
		{
			name:          "sin(x)",
			leftInterval:  0,
			rightInterval: math.Pi / 2,
			expectedValue: 1,
			tolerance:     10e-3,
			simpleExpr: func(x float64) float64 {
				return math.Sin(x)
			},
			amountOfPartitions: 1000,
		},
    // Outros casos de teste
	}

	for _, strategy := range strategies {
		for _, testCase := range testCases {
			testName := fmt.Sprintf("%s - %s from %.2f to %.2f using %d partitions",
				strategy.Description(), testCase.name,
				testCase.leftInterval, testCase.rightInterval, testCase.amountOfPartitions)

			t.Run(testName, func(t *testing.T) {
				// Act
				useCase := NewNewtonCotesUseCase(strategy)

				actualArea, err := useCase.Calculate(
					t.Context(),
					testCase.simpleExpr,
					testCase.leftInterval,
					testCase.rightInterval,
					testCase.amountOfPartitions,
				)

				// Assert
				assert.NoError(t, err, "Expected no error during integration")
				assert.InDelta(t, testCase.expectedValue, actualArea, testCase.tolerance)
			})
		}
	}
}
```

Onde `NewtonCotesStrategy` é uma interface que define o método `interpolate`, e cada estratégia de Newton-Cotes implementa essa interface. O caso de teste `newtonCotesTestCase` contém os parâmetros necessários para testar cada estratégia, como o intervalo de integração, a função a ser integrada e o número de partições.

= Quadraturas de Gauss

A quadratura de Gauss é uma técnica de integração numérica que utiliza pontos específicos (chamados de pontos de Gauss) e pesos associados para calcular a integral de uma função. Esses pontos e pesos são escolhidos de forma a maximizar a precisão da aproximação da integral.

A ideia central é aproximar a integral por uma soma ponderada de valores da função em pontos específicos. A implementação é dividida em duas funções principais:

A função *principal* (`Gauss-Quadrature-Calculate`) que gerencia o processo de integração e decide se deve usar particionamento baseado na estratégia escolhida.

A função *de cálculo de partição* (`Calculate-Partition`) que realiza o cálculo efetivo da quadratura para um intervalo específico usando os nós e pesos da estratégia.

O fluxo de execução verifica primeiro se a estratégia permite particionamento. Se não permitir (como no caso das quadraturas especiais), chama diretamente a função de cálculo de partição. Se permitir, divide o intervalo em partições e calcula cada uma separadamente.

#algorithm-figure("Gauss-quadrature-main", {
  import algorithmic: *
  Function("Gauss-Quadrature-Calculate", ("strategy", "left", "right", "partitions"), {
    If(FnInline[strategy.AllowPartitioning][], {
      Assign([partitionArea], FnInline[Calculate-Partition][strategy, left, right])
      Return("partitionArea")
    })
    Assign[area][$0$]
    Assign[delta][($"right" - "left"$) / $"partitions"$]
    For([$l$ := left until $"right"$ with step delta ], {
      Assign([partitionArea], FnInline[Calculate-Partition][strategy, $l$, $l +$delta$$])

      Assign[area][area + partitionArea]
    })
    Return($"area"$)
  })
}) <gauss-quadrature-main>

#algorithm-figure("Gauss-quadrature-partition", {
  import algorithmic: *
  Function("Calculate-Partition", ("strategy", "left", "right"), {
    Assign[nodes][strategy.GetNodes()]
    Assign[weights][strategy.GetWeights()]
    Assign[area][$0$]
    Assign[offset][strategy.GetOffset($"left"$, $"right"$)]
    Assign[scale][strategy.GetScalingFactor($"left"$, $"right"$)]
    For([$i$ := $0$ until $"len(nodes)"$ ], {
      Assign[x][offset + scale × nodes[$i$]]
      Assign[area][area + weights[$i$] × f(x)]
    })
    Assign[area][area × scale]
    Return($"area"$)
  })
}) <gauss-quadrature-partition>

Onde `strategy` define o tipo de quadratura (Legendre, Chebyshev, Hermite ou Laguerre), `nodes` são os pontos de Gauss, `weights` são os pesos correspondentes, e as transformações de escala e offset permitem aplicar a quadratura em intervalos diferentes do intervalo padrão.

A função principal verifica se a estratégia permite particionamento através do método `AllowPartitioning()`. Estratégias como Gauss-Legendre permitem particionamento para melhorar a precisão, enquanto estratégias especiais (Chebyshev, Hermite, Laguerre) têm restrições específicas de intervalo e não permitem particionamento.

== Tipos de Quadraturas de Gauss

As quadraturas de Gauss podem ser classificadas em diferentes tipos, cada uma adequada para diferentes tipos de funções e intervalos:

*Gauss-Legendre:* A quadratura mais comum, adequada para funções suaves em intervalos finitos. Utiliza os polinômios de Legendre como base.

*Quadraturas Especiais:* As demais quadraturas (Chebyshev, Hermite e Laguerre) são especializadas para casos específicos e requerem funções com características particulares ou intervalos específicos.

#figure(
  table(
    columns: 6,
    "Quadratura", "Pontos", "Intervalo", "Função de Peso", "Aplicação", "Precisão",
    table.hline(),
    "Gauss-Legendre", "2-4", "[-1, 1]", "1", "Funções suaves", "Polinômios grau 2n-1",
    table.hline(),
    "Gauss-Chebyshev", "2-4", "[-1, 1]", $1/sqrt(1-x^2)$, "Funções com singularidades", "Integrais com peso",
    "Gauss-Hermite", "2-4", "(-∞, +∞)", $e^(-x^2)$, "Funções com decaimento gaussiano", "Integrais infinitas",
    "Gauss-Laguerre", "2-4", "[0, +∞)", $e^(-x)$, "Funções com decaimento exponencial", "Integrais semi-infinitas",
    table.hline(),
  ),
  caption: [Tabela de quadraturas de Gauss, com pontos, intervalos, funções de peso e aplicações.],
) <table-gauss-quadratures>

=== Detalhamento dos Nós e Pesos

Para a quadratura de Gauss-Legendre, que é a mais utilizada, os nós e pesos para diferentes ordens são:

#figure(
  table(
    columns: 4,
    "Ordem", "Nós", "Pesos", "Precisão",
    table.hline(),
    "2", $±1/sqrt(3) ≈ ±0.577$, $1.0, 1.0$, "Polinômios grau 3",
    "3", $0, ±sqrt(3/5) ≈ ±0.775$, $8/9, 5/9, 5/9$, "Polinômios grau 5",
    "4", $±0.339, ±0.861$, $0.652, 0.348, 0.348, 0.652$, "Polinômios grau 7",
    table.hline(),
  ),
  caption: [Nós e pesos para quadratura de Gauss-Legendre (valores aproximados).],
) <table-gauss-legendre-nodes>

=== Quadraturas Especiais

As quadraturas de Chebyshev, Hermite e Laguerre são especializadas para casos específicos:

*Gauss-Chebyshev:* Utiliza uma função de peso $w(x) = 1/√(1-x²)$ e é adequada para integrais da forma $∫_{-1}^{1} f(x)/√(1-x²) d x$. Requer o intervalo [-1, 1] e é especialmente útil para funções com singularidades nos extremos.

*Gauss-Hermite:* Utiliza uma função de peso $w(x) = e^(-x^2)$ e é adequada para integrais da forma $∫_(-infinity)^(+infinity) f(x)e^(-x^2) d x$. Requer intervalos infinitos e é ideal para funções com decaimento gaussiano.

*Gauss-Laguerre:* Utiliza uma função de peso $w(x) = e^{-x}$ e é adequada para integrais da forma $∫_(0)^(+infinity) f(x)e^(-x) d x$. Requer o intervalo [0, +∞) e é ideal para funções com decaimento exponencial.

Para comprovar a implementação dessas estratégias foi usado o seguinte código:

```go
func TestGaussianQuadratures(t *testing.T) {
	// Arrange
	t.Parallel()

	// Create strategies for different Gaussian quadratures
	strategies := []GaussianQuadrature{}

	// Gauss-Legendre strategies (orders 2, 3, 4)
	for order := 2; order <= 4; order++ {
		strategy, err := NewGaussLegendre(order)
		assert.NoError(t, err)
		strategies = append(strategies, strategy)
	}

  // Fazer a mesma coisa para Gauss-Chebyshev, Gauss-Hermite e Gauss-Laguerre

	testCases := []gaussQuadratureTestCase{
		{
			name:          "x²",
			leftInterval:  0,
			rightInterval: 1,
			expectedArea:  1.0 / 3.0, // ∫₀¹ x² dx = 1/3
			tolerance:     1e-10,
			expr: func(x float64) float64 {
				return x * x
			},
		},
		// Outros casos de teste...
	}

	for _, strategy := range strategies {
		for _, testCase := range testCases {
			testName := fmt.Sprintf("%s Order %d - %s from %.2f to %.2f",
				strategy.Describe(), strategy.Order(), testCase.name,
				testCase.leftInterval, testCase.rightInterval)

			t.Run(testName, func(t *testing.T) {
				// Act
				useCase := NewGaussCalculatorUseCase(strategy)

				actualArea, err := useCase.Calculate(
					context.Background(),
					testCase.expr,
					testCase.leftInterval,
					testCase.rightInterval,
					1, // Gauss quadrature typically uses 1 partition
				)

        // Assert
        assert.NoError(t, err)
        assert.InDelta(t, testCase.expectedArea, actualArea, testCase.tolerance)
			})
		}
	}
}
```

Onde `GaussianQuadrature` é uma interface que define os métodos necessários para cada tipo de quadratura, incluindo `GetNodes()`, `GetWeights()`, `Integrate()`, e métodos de validação específicos para cada tipo. O caso de teste `gaussQuadratureTestCase` contém os parâmetros necessários para testar cada estratégia, incluindo tolerâncias específicas para cada tipo de quadratura.

= Integrais duplas

As integrais duplas são uma extensão das integrais simples para funções de duas variáveis, permitindo calcular volumes sob superfícies tridimensionais e áreas de regiões complexas no plano. Elas são particularmente úteis em aplicações de engenharia, física e matemática aplicada para resolver problemas que envolvem grandezas distribuídas em duas dimensões.

A integral dupla de uma função $f(x,y)$ sobre uma região $R$ é definida como:

$ integral.double_R f(x,y) d x d y $

Numericamente, aproximamos esta integral usando a *regra do ponto médio* aplicada em duas dimensões, dividindo a região retangular em uma grade de sub-retângulos e avaliando a função no centro de cada sub-retângulo.

#algorithm-figure("Double-integral", {
  import algorithmic: *
  Function("Calculate-Double-Integral", ("expr", "leftX", "rightX", "leftY", "rightY", "partitions"), {
    If($"leftX" = "rightX"$, {
      Return("Error: Zero width interval")
    })
    If($"leftY" = "rightY"$, {
      Return("Error: Zero width interval")
    })
    If($"partitions" = 0$, {
      Assign[partitions][$1$]
    })
    Assign[deltaX][($"rightX" - "leftX"$) / $"partitions"$]
    Assign[deltaY][($"rightY" - "leftY"$) / $"partitions"$]
    Assign[area][$0$]
    For([$i$ := $0$ until $"partitions"$ ], {
      For([$j$ := $0$ until $"partitions"$ ], {
        Assign[midX][$"leftX" + (i + 0.5) times "deltaX"$]
        Assign[midY][$"leftY" + (j + 0.5) times "deltaY"$]
        Assign[functionValue][expr(midX, midY)]
        Assign[area][area + functionValue × deltaX × deltaY]
      })
    })
    Return($"area"$)
  })
}) <double-integral-algorithm>

Onde `expr` é uma função de duas variáveis $f(x,y)$, `leftX` e `rightX` definem os limites de integração em x, `leftY` e `rightY` definem os limites em y, e `partitions` define o número de subdivisões em cada dimensão.

== Aplicações Práticas

As integrais duplas têm diversas aplicações práticas:

*Cálculo de Áreas:* Para calcular a área de regiões complexas, definimos uma função característica que vale 1 dentro da região e 0 fora dela.

*Cálculo de Volumes:* Para calcular o volume sob uma superfície $z = f(x,y)$, integramos a função sobre a região de interesse.

*Centros de Massa:* Para encontrar o centro de massa de uma lâmina com densidade variável.

*Momentos de Inércia:* Para calcular momentos de inércia de objetos bidimensionais.

== Exemplo: Cálculo da Área de um Círculo

Um exemplo prático interessante é o cálculo da área de um círculo usando integrais duplas. Para um círculo de raio $r = 1$ centrado na origem, definimos uma função característica:

$
  f(x,y) = cases(
    1 "se" x^2 + y^2 <= 1,
    0 "caso contrário"
  )
$

A área do círculo é então:

$ "Área" = integral.double_(-1)^1 integral.double_(-1)^1 f(x,y) d x d y = pi $

Este método é especialmente útil para formas complexas onde não há fórmulas analíticas simples.

== Exemplo: Cálculo da Área de uma Elipse

Similarmente, para uma elipse com semi-eixos $a = 3$ e $b = 2$, definimos:

$
  f(x,y) = cases(
    1 "se" (x/3)^2 + (y/2)^2 <= 1,
    0 "caso contrário"
  )
$

A área da elipse é:

$ "Área" = integral.double_(-3)^3 integral.double_(-2)^2 f(x,y) d x d y = pi a b = 6pi $

Para comprovar a implementação desta técnica foi usado o seguinte código:

```go
func TestDoubleIntegralCalculateArea(t *testing.T) {
	// Arrange
	t.Parallel()

	tests := []doubleIntegralTestCase{
		{
			name: "Circle Approximation with radius 1 and center = 0",
			expr: func(x, y float64) float64 {
				radius := 1.0
				distanceSquared := x*x + y*y

				if distanceSquared <= radius*radius {
					return 1.0 // Dentro do círculo
				}
				return 0.0 // Fora do círculo
			},
			leftIntervalX:      -1,
			rightIntervalX:     1,
			leftIntervalY:      -1,
			rightIntervalY:     1,
			numberOfPartitions: 1000,
			expectedArea:       math.Pi, // Área do círculo unitário é π
			tolerance:          0.01,
			description:        "Área de um círculo usando função característica",
		},
		// Outros testes...
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			useCase := NewDoubleIntegralUseCase()

			// Act
			result, err := useCase.CalculateArea(
				context.Background(),
				tc.expr,
				tc.leftIntervalX,
				tc.rightIntervalX,
				tc.leftIntervalY,
				tc.rightIntervalY,
				tc.numberOfPartitions,
			)

			// Assert
			assert.NoError(t, err)
			assert.InDelta(t, tc.expectedArea, result, tc.tolerance,
				"Expected area %v but got %v for %s",
				tc.expectedArea, result, tc.name)
		})
	}
}
```
