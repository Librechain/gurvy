package bn256

import (
	"math/big"

	"github.com/consensys/gurvy"
	"github.com/consensys/gurvy/bn256/fp"
	"github.com/consensys/gurvy/utils"
)

// E: y**2=x**3+1
// Etwist: y**2 = x**3+u**-1
// Tower: Fp->Fp2, u**2=-1 -> Fp12, v**6=9+u
// Generator (BN family): x=4965661367192848881
// optimal Ate loop: 6x+2
// Fp: p=21888242871839275222246405745257275088696311157297823662689037894645226208583
// Fr: r=21888242871839275222246405745257275088548364400416034343698204186575808495617

// ID bn256 ID
var ID = gurvy.BN256

// B b coeff of the curve
var B fp.Element

// generators of the r-torsion group, resp. in ker(pi-id), ker(Tr)
var g1Gen G1Jac
var g2Gen G2Jac

// point at infinity
var g1Infinity G1Jac
var g2Infinity G2Jac

// optimal Ate loop counter (=trace-1 = x in BLS family)
var loopCounter [66]int8

// Parameters useful for the GLV scalar multiplication. The third roots define the
//  endomorphisms phi1 and phi2 for <G1> and <G2>. lambda is such that <r, phi-lambda> lies above
// <r> in the ring Z[phi]. More concretely it's the associated eigenvalue
// of phi1 (resp phi2) restricted to <G1> (resp <G2>)
// cf https://www.cosic.esat.kuleuven.be/nessie/reports/phase2/GLV.pdf
var thirdRootOneG1 fp.Element
var thirdRootOneG2 fp.Element
var lambdaGLV big.Int

// parameters for pippenger ScalarMulByGen
// TODO get rid of this, keep only double and add, and the multi exp
const sGen = 4
const bGen = sGen

var tGenG1 [((1 << bGen) - 1)]G1Jac
var tGenG2 [((1 << bGen) - 1)]G2Jac

func init() {

	B.SetUint64(1)

	g1Gen.X.SetString("20567171726433170376993012834626974355708098753738075953327671604980729474588")
	g1Gen.Y.SetString("14259118686601658563517637559143782061303537174604067025175876803301021346267")
	g1Gen.Z.SetString("1")

	g2Gen.X.SetString("14433365730775072582213482468844163390964025019096075555058505630999708262443",
		"3683446723006852480794963570030936618743148392137679437247363531986401769417")
	g2Gen.Y.SetString("21253271987667943455369004300257637004831224612428754877033343975009216128128",
		"12495620673937637012904672587588023149812491484245871073230980321212840773339")
	g2Gen.Z.SetString("1",
		"0")

	g1Infinity.X.SetOne()
	g1Infinity.Y.SetOne()
	g2Infinity.X.SetOne()
	g2Infinity.Y.SetOne()

	thirdRootOneG1.SetString("2203960485148121921418603742825762020974279258880205651966")
	thirdRootOneG2.Square(&thirdRootOneG1)
	lambdaGLV.SetString("4407920970296243842393367215006156084916469457145843978461", 10)

	// binary decomposition of 15132376222941642752 little endian
	optimaAteLoop, _ := new(big.Int).SetString("29793968203157093288", 10)
	utils.NafDecomposition(optimaAteLoop, loopCounter[:])

	tGenG1[0].Set(&g1Gen)
	for j := 1; j < len(tGenG1)-1; j = j + 2 {
		tGenG1[j].Set(&tGenG1[j/2]).DoubleAssign()
		tGenG1[j+1].Set(&tGenG1[(j+1)/2]).AddAssign(&tGenG1[j/2])
	}
	tGenG2[0].Set(&g2Gen)
	for j := 1; j < len(tGenG2)-1; j = j + 2 {
		tGenG2[j].Set(&tGenG2[j/2]).DoubleAssign()
		tGenG2[j+1].Set(&tGenG2[(j+1)/2]).AddAssign(&tGenG2[j/2])
	}
}
