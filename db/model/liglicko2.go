package model

import (
	"errors"
	"math"
	"time"
)

const (
	liglicko2InternalRatingScale = 173.7178
	liglicko2MinRating           = 400.0
	liglicko2MaxRating           = 4000.0
	liglicko2DefaultRating       = 1500.0
	liglicko2MinVolatility       = 0.01
	liglicko2MaxVolatility       = 0.1
	liglicko2DefaultVolatility   = 0.09
	liglicko2MinDeviation        = 45.0
	liglicko2MaxDeviation        = 500.0
	liglicko2FirstAdvantage      = 0.0
	liglicko2Tau                 = 0.75
	liglicko2ConvergenceTol      = 1e-6
	liglicko2MaxConvergenceIters = 1000
	liglicko2MaxRatingDelta      = 700.0
	liglicko2RegulatorFactor     = 1.02
)

type Liglicko2Rating struct {
	Rating     float64
	Deviation  float64
	Volatility float64
	At         float64
}

func liglicko2InstantFromTime(t time.Time) float64 {
	return float64(t.UnixNano()) / float64(24*time.Hour)
}

func liglicko2Clamp(v, minV, maxV float64) float64 {
	if v < minV {
		return minV
	}
	if v > maxV {
		return maxV
	}
	return v
}

func liglicko2ToInternal(diff float64) float64 {
	return diff / liglicko2InternalRatingScale
}

func liglicko2ToExternal(internal float64) float64 {
	return internal * liglicko2InternalRatingScale
}

func liglicko2G(deviationInternal float64) float64 {
	return 1.0 / math.Sqrt(1.0+3.0*deviationInternal*deviationInternal/math.Pi/math.Pi)
}

func liglicko2Expectation(diffInternal, g float64) float64 {
	return 1.0 / (1.0 + math.Exp(-g*diffInternal))
}

func liglicko2NewDeviation(deviationInternal, volatility, elapsedPeriods float64) float64 {
	elapsed := math.Max(elapsedPeriods, 0.0)
	return math.Sqrt(deviationInternal*deviationInternal + elapsed*volatility*volatility)
}

func clampLiglicko2Rating(r Liglicko2Rating) Liglicko2Rating {
	return Liglicko2Rating{
		Rating:     liglicko2Clamp(r.Rating, liglicko2MinRating, liglicko2MaxRating),
		Deviation:  liglicko2Clamp(r.Deviation, liglicko2MinDeviation, liglicko2MaxDeviation),
		Volatility: liglicko2Clamp(r.Volatility, liglicko2MinVolatility, liglicko2MaxVolatility),
		At:         r.At,
	}
}

func liglicko2PreviewDeviation(r Liglicko2Rating, at float64) float64 {
	cr := clampLiglicko2Rating(r)
	dev := liglicko2NewDeviation(liglicko2ToInternal(cr.Deviation), cr.Volatility, at-cr.At)
	return liglicko2Clamp(liglicko2ToExternal(dev), liglicko2MinDeviation, liglicko2MaxDeviation)
}

func liglicko2Regulate(rating, delta float64) float64 {
	factor := 1.0
	if delta > 0.0 && rating < liglicko2DefaultRating+liglicko2MaxDeviation {
		factor = liglicko2RegulatorFactor
	}

	return rating + liglicko2Clamp(factor*delta, -liglicko2MaxRatingDelta, liglicko2MaxRatingDelta)
}

func updateLiglicko2Single(us, them Liglicko2Rating, score, now, advantage float64) (Liglicko2Rating, error) {
	us = clampLiglicko2Rating(us)
	them = clampLiglicko2Rating(them)

	phi := liglicko2ToInternal(liglicko2PreviewDeviation(us, now-1.0))
	theirG := liglicko2G(liglicko2ToInternal(liglicko2PreviewDeviation(them, now-1.0)))
	expected := liglicko2Expectation(liglicko2ToInternal(us.Rating-them.Rating+advantage), theirG)
	v := 1.0 / (theirG * theirG * expected * (1.0 - expected))
	delta := v * theirG * (score - expected)

	a := math.Log(us.Volatility * us.Volatility)
	f := func(x float64) float64 {
		ex := math.Exp(x)
		return ex*(delta*delta-phi*phi-v-ex)/(2.0*math.Pow(phi*phi+v+ex, 2.0)) - (x-a)/(liglicko2Tau*liglicko2Tau)
	}

	bigA := a
	bigB := 0.0
	if delta*delta > phi*phi+v {
		bigB = math.Log(delta*delta - phi*phi - v)
	} else {
		k := 1.0
		for f(a-k*liglicko2Tau) < 0.0 {
			k += 1.0
		}
		bigB = a - k*liglicko2Tau
	}

	fA := f(bigA)
	fB := f(bigB)

	iterations := 0
	for math.Abs(bigB-bigA) > liglicko2ConvergenceTol {
		iterations++
		if iterations > liglicko2MaxConvergenceIters {
			return Liglicko2Rating{}, errors.New("liglicko2 update failed to converge")
		}

		bigC := bigA + (bigA-bigB)*fA/(fB-fA)
		fC := f(bigC)

		if fC*fB <= 0.0 {
			bigA = bigB
			fA = fB
		} else {
			fA /= 2.0
		}

		bigB = bigC
		fB = fC
	}

	sigmaPrime := math.Exp(bigA / 2.0)
	phiStar := liglicko2NewDeviation(phi, sigmaPrime, math.Min(now-us.At, 1.0))
	phiPrime := 1.0 / math.Sqrt(1.0/(phiStar*phiStar)+1.0/v)
	muPrimeDiff := phiPrime * phiPrime * theirG * (score - expected)

	return clampLiglicko2Rating(Liglicko2Rating{
		Rating:     liglicko2Regulate(us.Rating, liglicko2ToExternal(muPrimeDiff)),
		Deviation:  liglicko2ToExternal(phiPrime),
		Volatility: sigmaPrime,
		At:         now,
	}), nil
}

func playerLiglicko2Rating(p *Player) Liglicko2Rating {
	return Liglicko2Rating{
		Rating:     p.Liglicko2Rating,
		Deviation:  p.Liglicko2Deviation,
		Volatility: p.Liglicko2Volatility,
		At:         p.Liglicko2At,
	}
}

func setPlayerLiglicko2Rating(p *Player, r Liglicko2Rating) {
	p.Liglicko2Rating = r.Rating
	p.Liglicko2Deviation = r.Deviation
	p.Liglicko2Volatility = r.Volatility
	p.Liglicko2At = r.At
}

func liglicko2ScoreFromOutcome(outcome Outcome) float64 {
	return liglicko2Clamp(float64(outcome), 0.0, 1.0)
}

func CalculateLiglicko2(a, b *Player, outcome Outcome, playedAt time.Time) (float64, float64, error) {
	now := liglicko2InstantFromTime(playedAt)
	first := playerLiglicko2Rating(a)
	second := playerLiglicko2Rating(b)
	score := liglicko2ScoreFromOutcome(outcome)

	updatedFirst, err := updateLiglicko2Single(first, second, score, now, liglicko2FirstAdvantage)
	if err != nil {
		return 0, 0, err
	}

	updatedSecond, err := updateLiglicko2Single(second, first, 1.0-score, now, -liglicko2FirstAdvantage)
	if err != nil {
		return 0, 0, err
	}

	deltaA := updatedFirst.Rating - first.Rating
	deltaB := updatedSecond.Rating - second.Rating

	setPlayerLiglicko2Rating(a, updatedFirst)
	setPlayerLiglicko2Rating(b, updatedSecond)

	return deltaA, deltaB, nil
}
