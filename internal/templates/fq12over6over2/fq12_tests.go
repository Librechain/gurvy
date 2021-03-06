package fq12over6over2

const Fq12Tests = `

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/commands"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// ------------------------------------------------------------
// tests

func TestE12ReceiverIsOperand(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100

	properties := gopter.NewProperties(parameters)

	genA := GenE12()
	genB := GenE12()

	properties.Property("Having the receiver as operand (addition) should output the same result", prop.ForAll(
		func(a, b *E12) bool {
			var c, d E12
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
		func(a, b *E12) bool {
			var c, d E12
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
		func(a, b *E12) bool {
			var c, d E12
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
		func(a *E12) bool {
			var b E12
			b.Square(a)
			a.Square(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("Having the receiver as operand (double) should output the same result", prop.ForAll(
		func(a *E12) bool {
			var b E12
			b.Double(a)
			a.Double(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("Having the receiver as operand (Inverse) should output the same result", prop.ForAll(
		func(a *E12) bool {
			var b E12
			b.Inverse(a)
			a.Inverse(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("Having the receiver as operand (Cyclotomic square) should output the same result", prop.ForAll(
		func(a *E12) bool {
			var b E12
			b.CyclotomicSquare(a)
			a.CyclotomicSquare(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("Having the receiver as operand (Conjugate) should output the same result", prop.ForAll(
		func(a *E12) bool {
			var b E12
			b.Conjugate(a)
			a.Conjugate(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("Having the receiver as operand (Frobenius) should output the same result", prop.ForAll(
		func(a *E12) bool {
			var b E12
			b.Frobenius(a)
			a.Frobenius(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("Having the receiver as operand (FrobeniusSquare) should output the same result", prop.ForAll(
		func(a *E12) bool {
			var b E12
			b.FrobeniusSquare(a)
			a.FrobeniusSquare(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("Having the receiver as operand (FrobeniusCube) should output the same result", prop.ForAll(
		func(a *E12) bool {
			var b E12
			b.FrobeniusCube(a)
			a.FrobeniusCube(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func TestE12State(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100

	subadd := &commands.ProtoCommand{
		Name: "SUBADD",
		RunFunc: func(systemUnderTest commands.SystemUnderTest) commands.Result {
			var a, b E12
			b.SetRandom()
			a.Add(systemUnderTest.(*E12), &b).Sub(&a, &b)
			return systemUnderTest.(*E12).Equal(&a)
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
			var a, b E12
			b.SetRandom()
			a.Mul(systemUnderTest.(*E12), &b)
			b.Inverse(&b)
			a.Mul(&a, &b)
			return systemUnderTest.(*E12).Equal(&a)
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
			var a E12
			a.Inverse(systemUnderTest.(*E12)).Inverse(&a)
			return systemUnderTest.(*E12).Equal(&a)
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
			var a, b, c E12
			c.Set(systemUnderTest.(*E12))
			a.Square(systemUnderTest.(*E12))
			b.Mul(systemUnderTest.(*E12), systemUnderTest.(*E12))
			return a.Equal(&b) && c.Equal(systemUnderTest.(*E12)) // check that the system hasn't changed
		},
		PostConditionFunc: func(state commands.State, result commands.Result) *gopter.PropResult {
			if result.(bool) {
				return &gopter.PropResult{Status: gopter.PropTrue}
			}
			return &gopter.PropResult{Status: gopter.PropFalse}
		},
	}

	conjugate := &commands.ProtoCommand{
		Name: "CONJUGATE",
		RunFunc: func(systemUnderTest commands.SystemUnderTest) commands.Result {
			var a, w, sum, diff E12
			var real, im, zero E6
			w.Set(systemUnderTest.(*E12))
			real.Double(&systemUnderTest.(*E12).C0)
			im.Double(&systemUnderTest.(*E12).C1)
			a.Conjugate(systemUnderTest.(*E12))
			sum.Add(systemUnderTest.(*E12), &a)
			diff.Sub(systemUnderTest.(*E12), &a)
			return w.Equal(systemUnderTest.(*E12)) && sum.C1.Equal(&zero) && sum.C0.Equal(&real) && diff.C0.Equal(&zero) && diff.C1.Equal(&im)
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
			var a, w E12
			w.Set(systemUnderTest.(*E12))
			a.Frobenius(systemUnderTest.(*E12)).
				Frobenius(&a).
				Frobenius(&a).
				Frobenius(&a).
				Frobenius(&a).
				Frobenius(&a).
				Frobenius(&a).
				Frobenius(&a).
				Frobenius(&a).
				Frobenius(&a).
				Frobenius(&a).
				Frobenius(&a)
			return w.Equal(systemUnderTest.(*E12)) && a.Equal(&w)
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
			var a, w E12
			w.Set(systemUnderTest.(*E12))
			a.FrobeniusSquare(systemUnderTest.(*E12)).
				FrobeniusSquare(&a).
				FrobeniusSquare(&a).
				FrobeniusSquare(&a).
				FrobeniusSquare(&a).
				FrobeniusSquare(&a)
			return w.Equal(systemUnderTest.(*E12)) && a.Equal(&w)
		},
		PostConditionFunc: func(state commands.State, result commands.Result) *gopter.PropResult {
			if result.(bool) {
				return &gopter.PropResult{Status: gopter.PropTrue}
			}
			return &gopter.PropResult{Status: gopter.PropFalse}
		},
	}

	frobeniusscube := &commands.ProtoCommand{
		Name: "FROBENIUSCUBE",
		RunFunc: func(systemUnderTest commands.SystemUnderTest) commands.Result {
			var a, w E12
			w.Set(systemUnderTest.(*E12))
			a.FrobeniusCube(systemUnderTest.(*E12)).
				FrobeniusCube(&a).
				FrobeniusCube(&a).
				FrobeniusCube(&a)
			return w.Equal(systemUnderTest.(*E12)) && a.Equal(&w)
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
			var a, b, w, s, sc E12
			w.Set(systemUnderTest.(*E12))
			a.FrobeniusCube(systemUnderTest.(*E12)).
				FrobeniusCube(&a)
			b.Inverse(&w)
			a.Mul(&a, &b)
			b.Set(&a)
			a.FrobeniusSquare(&a).Mul(&a, &b) // a is now in the cyclotomic subgroup
			s.Square(&a)
			sc.CyclotomicSquare(&a)
			return w.Equal(systemUnderTest.(*E12)) && s.Equal(&sc)
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
			var a E12
			a.SetRandom()
			return &a
		},
		InitialStateGen: gen.Const(false),
		GenCommandFunc: func(state commands.State) gopter.Gen {
			return gen.OneConstOf(subadd, mulinverse, inversetwice, squaremul, conjugate, frobenius, frobeniussquare, frobeniusscube, cyclotomicsquare)
		},
	}

	properties := gopter.NewProperties(parameters)
	properties.Property("E12 state", commands.Prop(e6commands))
	properties.TestingRun(t, gopter.ConsoleReporter(false))

}

// ------------------------------------------------------------
// benches

func BenchmarkE12Add(b *testing.B) {
	var a, c E12
	a.SetRandom()
	c.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Add(&a, &c)
	}
}

func BenchmarkE12Sub(b *testing.B) {
	var a, c E12
	a.SetRandom()
	c.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Sub(&a, &c)
	}
}

func BenchmarkE12Mul(b *testing.B) {
	var a, c E12
	a.SetRandom()
	c.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Mul(&a, &c)
	}
}

func BenchmarkE12Cyclosquare(b *testing.B) {
	var a E12
	a.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.CyclotomicSquare(&a)
	}
}

func BenchmarkE12Square(b *testing.B) {
	var a E12
	a.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Square(&a)
	}
}

func BenchmarkE12Inverse(b *testing.B) {
	var a E12
	a.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Inverse(&a)
	}
}

func BenchmarkE12Conjugate(b *testing.B) {
	var a E12
	a.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Conjugate(&a)
	}
}

func BenchmarkE12Frobenius(b *testing.B) {
	var a E12
	a.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Frobenius(&a)
	}
}

func BenchmarkE12FrobeniusSquare(b *testing.B) {
	var a E12
	a.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.FrobeniusSquare(&a)
	}
}

func BenchmarkE12FrobeniusCube(b *testing.B) {
	var a E12
	a.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.FrobeniusCube(&a)
	}
}

func BenchmarkE12Expt(b *testing.B) {
	var a E12
	a.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Expt(&a)
	}
}

func BenchmarkE12FinalExponentiation(b *testing.B) {
	var a E12
	a.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.FinalExponentiation(&a)
	}
}

`
