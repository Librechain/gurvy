package bw761

import (
	"testing"

	"github.com/consensys/gurvy/bw761/fp"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/commands"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// ------------------------------------------------------------
// tests

func TestE2ReceiverIsOperand(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100

	properties := gopter.NewProperties(parameters)

	genA := GenE2()
	genB := GenE2()
	genfp := GenFp()

	properties.Property("Having the receiver as operand (addition) should output the same result", prop.ForAll(
		func(a, b *E2) bool {
			var c, d E2
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
		func(a, b *E2) bool {
			var c, d E2
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
		func(a, b *E2) bool {
			var c, d E2
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
		func(a *E2) bool {
			var b E2
			b.Square(a)
			a.Square(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("Having the receiver as operand (neg) should output the same result", prop.ForAll(
		func(a *E2) bool {
			var b E2
			b.Neg(a)
			a.Neg(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("Having the receiver as operand (double) should output the same result", prop.ForAll(
		func(a *E2) bool {
			var b E2
			b.Double(a)
			a.Double(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("Having the receiver as operand (mul by non residue) should output the same result", prop.ForAll(
		func(a *E2) bool {
			var b E2
			b.MulByNonResidue(a)
			a.MulByNonResidue(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("Having the receiver as operand (Inverse) should output the same result", prop.ForAll(
		func(a *E2) bool {
			var b E2
			b.Inverse(a)
			a.Inverse(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("Having the receiver as operand (Conjugate) should output the same result", prop.ForAll(
		func(a *E2) bool {
			var b E2
			b.Conjugate(a)
			a.Conjugate(a)
			return a.Equal(&b)
		},
		genA,
	))

	properties.Property("Having the receiver as operand (mul by element) should output the same result", prop.ForAll(
		func(a *E2, b fp.Element) bool {
			var c E2
			c.MulByElement(a, &b)
			a.MulByElement(a, &b)
			return a.Equal(&c)
		},
		genA,
		genfp,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func TestE2State(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100

	subadd := &commands.ProtoCommand{
		Name: "SUBADD",
		RunFunc: func(systemUnderTest commands.SystemUnderTest) commands.Result {
			var a, b E2
			b.SetRandom()
			a.Add(systemUnderTest.(*E2), &b).Sub(&a, &b)
			return systemUnderTest.(*E2).Equal(&a)
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
			var a, b E2
			b.SetRandom()
			a.Mul(systemUnderTest.(*E2), &b)
			b.Inverse(&b)
			a.Mul(&a, &b)
			return systemUnderTest.(*E2).Equal(&a)
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
			var a E2
			a.Inverse(systemUnderTest.(*E2)).Inverse(&a)
			return systemUnderTest.(*E2).Equal(&a)
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
			var a E2
			a.Neg(systemUnderTest.(*E2)).Neg(&a)
			return systemUnderTest.(*E2).Equal(&a)
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
			var a, b, c E2
			c.Set(systemUnderTest.(*E2))
			a.Square(systemUnderTest.(*E2))
			b.Mul(systemUnderTest.(*E2), systemUnderTest.(*E2))
			return a.Equal(&b) && c.Equal(systemUnderTest.(*E2)) // check that the system hasn't changed
		},
		PostConditionFunc: func(state commands.State, result commands.Result) *gopter.PropResult {
			if result.(bool) {
				return &gopter.PropResult{Status: gopter.PropTrue}
			}
			return &gopter.PropResult{Status: gopter.PropFalse}
		},
	}

	mulbyelmtinverse := &commands.ProtoCommand{
		Name: "MULBYELMT",
		RunFunc: func(systemUnderTest commands.SystemUnderTest) commands.Result {
			var a E2
			var b fp.Element
			a.Set(systemUnderTest.(*E2))
			b.SetRandom()
			a.MulByElement(&a, &b)
			b.Inverse(&b)
			a.MulByElement(&a, &b)
			return a.Equal(systemUnderTest.(*E2))
		},
		PostConditionFunc: func(state commands.State, result commands.Result) *gopter.PropResult {
			if result.(bool) {
				return &gopter.PropResult{Status: gopter.PropTrue}
			}
			return &gopter.PropResult{Status: gopter.PropFalse}
		},
	}

	doublemul := &commands.ProtoCommand{
		Name: "DOUBLEMUL",
		RunFunc: func(systemUnderTest commands.SystemUnderTest) commands.Result {
			var a, b, c E2
			var d fp.Element
			c.Set(systemUnderTest.(*E2))
			d.SetUint64(2)
			a.MulByElement(systemUnderTest.(*E2), &d)
			b.Double(systemUnderTest.(*E2))
			return a.Equal(&b) && c.Equal(systemUnderTest.(*E2)) // check that the system hasn't changed
		},
		PostConditionFunc: func(state commands.State, result commands.Result) *gopter.PropResult {
			if result.(bool) {
				return &gopter.PropResult{Status: gopter.PropTrue}
			}
			return &gopter.PropResult{Status: gopter.PropFalse}
		},
	}

	mulbynonres := &commands.ProtoCommand{
		Name: "MULBYNONRESIDUE",
		RunFunc: func(systemUnderTest commands.SystemUnderTest) commands.Result {
			var a, b, c, w E2
			w.Set(systemUnderTest.(*E2))
			a.MulByNonResidue(systemUnderTest.(*E2))
			b.A1.SetOne()
			c.Mul(systemUnderTest.(*E2), &b)
			return w.Equal(systemUnderTest.(*E2)) && a.Equal(&c)
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
			var a, d, w E2
			var b, c fp.Element
			w.Set(systemUnderTest.(*E2))
			b.Double(&systemUnderTest.(*E2).A0)
			c.Double(&systemUnderTest.(*E2).A1)
			a.Conjugate(systemUnderTest.(*E2))
			d.Add(systemUnderTest.(*E2), &a)
			a.Sub(systemUnderTest.(*E2), &a)
			return d.A1.IsZero() && a.A0.IsZero() && d.A0.Equal(&b) && a.A1.Equal(&c) && w.Equal(systemUnderTest.(*E2))
		},
		PostConditionFunc: func(state commands.State, result commands.Result) *gopter.PropResult {
			if result.(bool) {
				return &gopter.PropResult{Status: gopter.PropTrue}
			}
			return &gopter.PropResult{Status: gopter.PropFalse}
		},
	}

	e2commands := &commands.ProtoCommands{
		NewSystemUnderTestFunc: func(_ commands.State) commands.SystemUnderTest {
			var a E2
			a.SetRandom()
			return &a
		},
		InitialStateGen: gen.Const(false),
		GenCommandFunc: func(state commands.State) gopter.Gen {
			return gen.OneConstOf(subadd, mulinverse, inversetwice, negtwice, squaremul, mulbyelmtinverse, doublemul, mulbynonres, conjugate)
		},
	}

	properties := gopter.NewProperties(parameters)
	properties.Property("E2 state", commands.Prop(e2commands))
	properties.TestingRun(t, gopter.ConsoleReporter(false))

}

// ------------------------------------------------------------
// benches

func BenchmarkE2Add(b *testing.B) {
	var a, c E2
	a.SetRandom()
	c.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Add(&a, &c)
	}
}

func BenchmarkE2Sub(b *testing.B) {
	var a, c E2
	a.SetRandom()
	c.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Sub(&a, &c)
	}
}

func BenchmarkE2Mul(b *testing.B) {
	var a, c E2
	a.SetRandom()
	c.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Mul(&a, &c)
	}
}

func BenchmarkE2MulByElement(b *testing.B) {
	var a E2
	var c fp.Element
	c.SetRandom()
	a.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.MulByElement(&a, &c)
	}
}

func BenchmarkE2Square(b *testing.B) {
	var a E2
	a.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Square(&a)
	}
}

func BenchmarkE2Inverse(b *testing.B) {
	var a E2
	a.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Inverse(&a)
	}
}

func BenchmarkE2MulNonRes(b *testing.B) {
	var a E2
	a.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.MulByNonResidue(&a)
	}
}

func BenchmarkE2Conjugate(b *testing.B) {
	var a E2
	a.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Conjugate(&a)
	}
}
