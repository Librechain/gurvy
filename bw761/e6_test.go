package bw761

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/commands"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// ------------------------------------------------------------
// tests

func TestE6ReceiverIsOperand(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100

	properties := gopter.NewProperties(parameters)

	genA := GenE6()
	genB := GenE6()

	properties.Property("Having the receiver as operand (addition) should output the same result", prop.ForAll(
		func(a, b *E6) bool {
			var c, d E6
			d.Set(a)
			c.Add(a, b)
			a.Add(a, b)
			b.Add(&d, b)
			return a.Equal(b) && a.Equal(&c) && b.Equal(&c)
		},
		genA,
		genB,
	))

	properties.Property("Having the receiver as operand (sub) should output the same result", prop.ForAll(
		func(a, b *E6) bool {
			var c, d E6
			d.Set(a)
			c.Sub(a, b)
			a.Sub(a, b)
			b.Sub(&d, b)
			return a.Equal(b) && a.Equal(&c) && b.Equal(&c)
		},
		genA,
		genB,
	))

	properties.Property("Having the receiver as operand (mul) should output the same result", prop.ForAll(
		func(a, b *E6) bool {
			var c, d E6
			d.Set(a)
			c.Mul(a, b)
			a.Mul(a, b)
			b.Mul(&d, b)
			return a.Equal(b) && a.Equal(&c) && b.Equal(&c)
		},
		genA,
		genB,
	))

	properties.Property("Having the receiver as operand (square) should output the same result", prop.ForAll(
		func(a *E6) bool {
			var b E6
			b.Square(a)
			a.Square(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("Having the receiver as operand (neg) should output the same result", prop.ForAll(
		func(a *E6) bool {
			var b E6
			b.Neg(a)
			a.Neg(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("Having the receiver as operand (double) should output the same result", prop.ForAll(
		func(a *E6) bool {
			var b E6
			b.Double(a)
			a.Double(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("Having the receiver as operand (Inverse) should output the same result", prop.ForAll(
		func(a *E6) bool {
			var b E6
			b.Inverse(a)
			a.Inverse(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("Having the receiver as operand (mul by E2) should output the same result", prop.ForAll(
		func(a *E6) bool {
			var b E6
			var c E2
			c.SetRandom()
			b.MulByE2(a, &c)
			a.MulByE2(a, &c)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("Having the receiver as operand (cyclotomic square) should output the same result", prop.ForAll(
		func(a *E6) bool {
			var b E6
			b.CyclotomicSquare(a)
			a.CyclotomicSquare(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("Having the receiver as operand (frobenius) should output the same result", prop.ForAll(
		func(a *E6) bool {
			var b E6
			b.Frobenius(a)
			a.Frobenius(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("Having the receiver as operand (frobenius square) should output the same result", prop.ForAll(
		func(a *E6) bool {
			var b E6
			b.FrobeniusSquare(a)
			a.FrobeniusSquare(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("Having the receiver as operand (frobenius cube) should output the same result", prop.ForAll(
		func(a *E6) bool {
			var b E6
			b.FrobeniusSquare(a)
			a.FrobeniusSquare(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func TestE6State(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100

	subadd := &commands.ProtoCommand{
		Name: "SUBADD",
		RunFunc: func(systemUnderTest commands.SystemUnderTest) commands.Result {
			var a, b E6
			b.SetRandom()
			a.Add(systemUnderTest.(*E6), &b).Sub(&a, &b)
			return systemUnderTest.(*E6).Equal(&a)
		},
		PostConditionFunc: func(state commands.State, result commands.Result) *gopter.PropResult {
			if result.(bool) {
				return &gopter.PropResult{Status: gopter.PropTrue}
			}
			return &gopter.PropResult{Status: gopter.PropFalse}
		},
	}

	mulinverse := &commands.ProtoCommand{
		Name: "MULINVERSE",
		RunFunc: func(systemUnderTest commands.SystemUnderTest) commands.Result {
			var a, b E6
			b.SetRandom()
			a.Mul(systemUnderTest.(*E6), &b)
			b.Inverse(&b)
			a.Mul(&a, &b)
			return systemUnderTest.(*E6).Equal(&a)
		},
		PostConditionFunc: func(state commands.State, result commands.Result) *gopter.PropResult {
			if result.(bool) {
				return &gopter.PropResult{Status: gopter.PropTrue}
			}
			return &gopter.PropResult{Status: gopter.PropFalse}
		},
	}

	inversetwice := &commands.ProtoCommand{
		Name: "INVERSETWICE",
		RunFunc: func(systemUnderTest commands.SystemUnderTest) commands.Result {
			var a E6
			a.Inverse(systemUnderTest.(*E6)).Inverse(&a)
			return systemUnderTest.(*E6).Equal(&a)
		},
		PostConditionFunc: func(state commands.State, result commands.Result) *gopter.PropResult {
			if result.(bool) {
				return &gopter.PropResult{Status: gopter.PropTrue}
			}
			return &gopter.PropResult{Status: gopter.PropFalse}
		},
	}

	negtwice := &commands.ProtoCommand{
		Name: "NEGTWICE",
		RunFunc: func(systemUnderTest commands.SystemUnderTest) commands.Result {
			var a E6
			a.Neg(systemUnderTest.(*E6)).Neg(&a)
			return systemUnderTest.(*E6).Equal(&a)
		},
		PostConditionFunc: func(state commands.State, result commands.Result) *gopter.PropResult {
			if result.(bool) {
				return &gopter.PropResult{Status: gopter.PropTrue}
			}
			return &gopter.PropResult{Status: gopter.PropFalse}
		},
	}

	squaremul := &commands.ProtoCommand{
		Name: "SQUAREMUL",
		RunFunc: func(systemUnderTest commands.SystemUnderTest) commands.Result {
			var a, b, c E6
			c.Set(systemUnderTest.(*E6))
			a.Square(systemUnderTest.(*E6))
			b.Mul(systemUnderTest.(*E6), systemUnderTest.(*E6))
			return a.Equal(&b) && c.Equal(systemUnderTest.(*E6)) // check that the system hasn't changed
		},
		PostConditionFunc: func(state commands.State, result commands.Result) *gopter.PropResult {
			if result.(bool) {
				return &gopter.PropResult{Status: gopter.PropTrue}
			}
			return &gopter.PropResult{Status: gopter.PropFalse}
		},
	}

	doubleadd := &commands.ProtoCommand{
		Name: "DOUBLEADD",
		RunFunc: func(systemUnderTest commands.SystemUnderTest) commands.Result {
			var a, b, c E6
			c.Set(systemUnderTest.(*E6))
			a.Add(systemUnderTest.(*E6), systemUnderTest.(*E6))
			b.Double(systemUnderTest.(*E6))
			return a.Equal(&b) && c.Equal(systemUnderTest.(*E6)) // check that the system hasn't changed
		},
		PostConditionFunc: func(state commands.State, result commands.Result) *gopter.PropResult {
			if result.(bool) {
				return &gopter.PropResult{Status: gopter.PropTrue}
			}
			return &gopter.PropResult{Status: gopter.PropFalse}
		},
	}

	frobenius := &commands.ProtoCommand{
		Name: "FROBENIUS",
		RunFunc: func(systemUnderTest commands.SystemUnderTest) commands.Result {
			var a E6
			a.Frobenius(systemUnderTest.(*E6)).
				Frobenius(&a).
				Frobenius(&a).
				Frobenius(&a).
				Frobenius(&a).
				Frobenius(&a)
			return a.Equal(systemUnderTest.(*E6))
		},
		PostConditionFunc: func(state commands.State, result commands.Result) *gopter.PropResult {
			if result.(bool) {
				return &gopter.PropResult{Status: gopter.PropTrue}
			}
			return &gopter.PropResult{Status: gopter.PropFalse}
		},
	}

	frobeniussquare := &commands.ProtoCommand{
		Name: "FROBENIUSSQUARE",
		RunFunc: func(systemUnderTest commands.SystemUnderTest) commands.Result {
			var a E6
			a.FrobeniusSquare(systemUnderTest.(*E6)).
				FrobeniusSquare(&a).
				FrobeniusSquare(&a)
			return a.Equal(systemUnderTest.(*E6))
		},
		PostConditionFunc: func(state commands.State, result commands.Result) *gopter.PropResult {
			if result.(bool) {
				return &gopter.PropResult{Status: gopter.PropTrue}
			}
			return &gopter.PropResult{Status: gopter.PropFalse}
		},
	}

	frobeniuscube := &commands.ProtoCommand{
		Name: "FROBENIUSCUBE",
		RunFunc: func(systemUnderTest commands.SystemUnderTest) commands.Result {
			var a E6
			a.FrobeniusCube(systemUnderTest.(*E6)).
				FrobeniusCube(&a)
			return a.Equal(systemUnderTest.(*E6))
		},
		PostConditionFunc: func(state commands.State, result commands.Result) *gopter.PropResult {
			if result.(bool) {
				return &gopter.PropResult{Status: gopter.PropTrue}
			}
			return &gopter.PropResult{Status: gopter.PropFalse}
		},
	}

	cyclotomicsquare := &commands.ProtoCommand{
		Name: "CYCLOTOMICSQUARE",
		RunFunc: func(systemUnderTest commands.SystemUnderTest) commands.Result {
			var a, w, b, s, sc E6
			w.Set(systemUnderTest.(*E6))
			a.FrobeniusCube(systemUnderTest.(*E6))
			b.Inverse(systemUnderTest.(*E6))
			a.Mul(&a, &b)
			b.Set(&a)
			a.Frobenius(&a).Mul(&a, &b)
			s.Square(&a)
			sc.CyclotomicSquare(&a)
			return w.Equal(systemUnderTest.(*E6)) && s.Equal(&sc)
		},
		PostConditionFunc: func(state commands.State, result commands.Result) *gopter.PropResult {
			if result.(bool) {
				return &gopter.PropResult{Status: gopter.PropTrue}
			}
			return &gopter.PropResult{Status: gopter.PropFalse}
		},
	}

	e6commands := &commands.ProtoCommands{
		NewSystemUnderTestFunc: func(_ commands.State) commands.SystemUnderTest {
			var a E6
			a.SetRandom()
			return &a
		},
		InitialStateGen: gen.Const(false),
		GenCommandFunc: func(state commands.State) gopter.Gen {
			return gen.OneConstOf(subadd, mulinverse, inversetwice, negtwice, squaremul, doubleadd, frobenius, frobeniussquare, frobeniuscube, cyclotomicsquare)
		},
	}

	properties := gopter.NewProperties(parameters)
	properties.Property("E6 state", commands.Prop(e6commands))
	properties.TestingRun(t, gopter.ConsoleReporter(false))

}

// ------------------------------------------------------------
// benches

func BenchmarkE6Add(b *testing.B) {
	var a, c E6
	a.SetRandom()
	c.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Add(&a, &c)
	}
}

func BenchmarkE6Sub(b *testing.B) {
	var a, c E6
	a.SetRandom()
	c.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Sub(&a, &c)
	}
}

func BenchmarkE6Mul(b *testing.B) {
	var a, c E6
	a.SetRandom()
	c.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Mul(&a, &c)
	}
}

func BenchmarkE6Square(b *testing.B) {
	var a E6
	a.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Square(&a)
	}
}

func BenchmarkE6Inverse(b *testing.B) {
	var a E6
	a.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Inverse(&a)
	}
}
