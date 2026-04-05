// Package liglicko2 implements Lichess-style glicko2 ratings
//
// This implementation is based on https://github.com/niklasf/liglicko2
// and adapted for this codebase.
package liglicko2

import (
	"errors"
	"math"
	"time"
)

const (
	InternalRatingScale = 173.7178

	MinRating     = 400.0
	MaxRating     = 4000.0
	DefaultRating = 1500.0

	MinVolatility     = 0.01
	MaxVolatility     = 0.1
	DefaultVolatility = 0.09

	MinDeviation     = 45.0
	MaxDeviation     = 500.0
	DefaultDeviation = 500.0

	FirstAdvantage      = 0.0
	Tau                 = 0.75
	ConvergenceTol      = 1e-6
	MaxConvergenceIters = 1000
	MaxRatingDelta      = 700.0
	RegulatorFactor     = 1.02
)

type Rating struct {
	Rating     float64
	Deviation  float64
	Volatility float64
	At         float64
}

func InstantFromTime(t time.Time) float64 {
	return float64(t.UnixNano()) / float64(24*time.Hour)
}

func Clamp(v, minV, maxV float64) float64 {
	if v < minV {
		return minV
	}
	if v > maxV {
		return maxV
	}
	return v
}

func clampRating(r Rating) Rating {
	return Rating{
		Rating:     Clamp(r.Rating, MinRating, MaxRating),
		Deviation:  Clamp(r.Deviation, MinDeviation, MaxDeviation),
		Volatility: Clamp(r.Volatility, MinVolatility, MaxVolatility),
		At:         r.At,
	}
}

func toInternal(v float64) float64 {
	return v / InternalRatingScale
}

func toExternal(v float64) float64 {
	return v * InternalRatingScale
}

func impact(phi float64) float64 {
	return 1.0 / math.Sqrt(1.0+3.0*phi*phi/math.Pi/math.Pi)
}

func expected(diff, g float64) float64 {
	return 1.0 / (1.0 + math.Exp(-g*diff))
}

func newDeviation(phi, sigma, elapsed float64) float64 {
	elapsed = math.Max(elapsed, 0.0)
	return math.Sqrt(phi*phi + elapsed*sigma*sigma)
}

func previewDeviation(r Rating, at float64) float64 {
	r = clampRating(r)
	phi := newDeviation(toInternal(r.Deviation), r.Volatility, at-r.At)
	return Clamp(toExternal(phi), MinDeviation, MaxDeviation)
}

func regulate(rating, delta float64) float64 {
	factor := 1.0
	if delta > 0.0 && rating < DefaultRating+MaxDeviation {
		factor = RegulatorFactor
	}
	return rating + Clamp(factor*delta, -MaxRatingDelta, MaxRatingDelta)
}

func UpdateSingle(player, opponent Rating, score, now, advantage float64) (Rating, error) {
	player = clampRating(player)
	opponent = clampRating(opponent)

	phi := toInternal(previewDeviation(player, now-1.0))
	oppImpact := impact(toInternal(previewDeviation(opponent, now-1.0)))
	e := expected(toInternal(player.Rating-opponent.Rating+advantage), oppImpact)
	v := 1.0 / (oppImpact * oppImpact * e * (1.0 - e))
	delta := v * oppImpact * (score - e)

	a := math.Log(player.Volatility * player.Volatility)
	f := func(x float64) float64 {
		ex := math.Exp(x)
		return ex*(delta*delta-phi*phi-v-ex)/(2.0*math.Pow(phi*phi+v+ex, 2.0)) - (x-a)/(Tau*Tau)
	}

	A := a
	var B float64
	if delta*delta > phi*phi+v {
		B = math.Log(delta*delta - phi*phi - v)
	} else {
		k := 1.0
		for f(a-k*Tau) < 0.0 {
			k += 1.0
		}
		B = a - k*Tau
	}

	fA := f(A)
	fB := f(B)

	for i := 0; math.Abs(B-A) > ConvergenceTol; i++ {
		if i > MaxConvergenceIters {
			return Rating{}, errors.New("liglicko2 update failed to converge")
		}

		C := A + (A-B)*fA/(fB-fA)
		fC := f(C)

		if fC*fB <= 0.0 {
			A = B
			fA = fB
		} else {
			fA /= 2.0
		}

		B = C
		fB = fC
	}

	sigmaPrime := math.Exp(A / 2.0)
	phiStar := newDeviation(phi, sigmaPrime, math.Min(now-player.At, 1.0))
	phiPrime := 1.0 / math.Sqrt(1.0/(phiStar*phiStar)+1.0/v)
	muDelta := phiPrime * phiPrime * oppImpact * (score - e)

	return clampRating(Rating{
		Rating:     regulate(player.Rating, toExternal(muDelta)),
		Deviation:  toExternal(phiPrime),
		Volatility: sigmaPrime,
		At:         now,
	}), nil
}
