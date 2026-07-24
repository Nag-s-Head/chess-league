package ele_charts

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
)

type EloChange struct {
	Delta      int
	currentElo int
}

type Params struct {
	// Must be in ascending time order
	Changes []EloChange
	EndElo  int
}

type eloStats struct {
	min, max int
	mean     int
}

const ContentType = "image/png"
const MaxEloChanges = 10

const (
	imageHeight                        = 70
	pixelsPerEntry                     = 30
	lineThickness                      = 3
	backgroundColourMultiplier float32 = 0.8
)

func lineInterp(start, end, x int) int {
	t := float64(x) / float64(pixelsPerEntry)
	smoothT := (1.0 - math.Cos(t*math.Pi)) / 2.0
	return int(float64(start) + smoothT*float64(end-start))
}

func shadedInterp(low, high, x, width int) int {
	return low + int((float32(x)/float32(width))*float32(high-low))
}

func clamp(y int) int {
	if y >= imageHeight {
		return imageHeight - 1
	}
	return y
}

func renderPoint(img draw.Image, xStart, lastElo, currentElo int, stats eloStats) {
	eloGainBaseColour := color.RGBA{G: 255, A: 255}
	eloLossBaseColour := color.RGBA{R: 255, A: 255}

	diff := stats.max - stats.min

	for i := range pixelsPerEntry {
		x := xStart + i
		unscaledY := lineInterp(lastElo, currentElo, i)

		y := 0
		if diff != 0 {
			normalised := (unscaledY - stats.min) * imageHeight / diff
			y = imageHeight - normalised
		}

		y = imageHeight - y

		c := eloGainBaseColour
		if lastElo < currentElo {
			c = eloLossBaseColour
		}

		src := image.NewUniform(c)
		draw.Draw(img,
			image.Rect(x, y, x+1, clamp(y+lineThickness)),
			src,
			image.Point{},
			draw.Over)

		fillStart := y + lineThickness
		fillHeight := imageHeight - fillStart
		for i := imageHeight - 1; i >= y+lineThickness; i-- {
			currentStep := i - fillStart
			alpha := shadedInterp(int(255.0*backgroundColourMultiplier), 50, currentStep, fillHeight)
			shadedColor := color.NRGBA{
				R: c.R,
				G: c.G,
				B: c.B,
				A: uint8(alpha),
			}

			src := image.NewUniform(shadedColor)
			draw.Draw(img,
				image.Rect(x, i, x+1, i+1),
				src,
				image.Point{},
				draw.Over)
		}
	}
}

func Render(params Params) ([]byte, error) {
	width := pixelsPerEntry * len(params.Changes)
	if width == 0 {
		width = pixelsPerEntry
	}

	img := image.NewRGBA(image.Rect(0, 0, width, imageHeight+lineThickness))

	stats := eloStats{
		min: params.EndElo,
		max: params.EndElo,
	}

	if len(params.Changes) > 0 {
		total := 0
		currentElo := params.EndElo
		for i := len(params.Changes) - 1; i >= 0; i-- {
			entry := &params.Changes[i]
			currentElo += entry.Delta
			entry.currentElo = currentElo

			total += currentElo

			if currentElo > stats.max {
				stats.max = currentElo
			} else if currentElo < stats.min {
				stats.min = currentElo
			}
		}

		stats.mean = total / len(params.Changes)
	}

	lastElo := 0
	for i, entry := range params.Changes {
		if i == 0 {
			lastElo = entry.currentElo
			continue
		}

		xStart := pixelsPerEntry * (i - 1)
		renderPoint(img, xStart, lastElo, entry.currentElo, stats)
		lastElo = entry.currentElo
	}

	if len(params.Changes) == 0 {
		lastElo = params.EndElo
	}

	xStart := width - pixelsPerEntry
	renderPoint(img, xStart, lastElo, params.EndElo, stats)

	buf := bytes.NewBuffer(nil)
	err := png.Encode(buf, img)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot encode chart"), err)
	}

	return buf.Bytes(), nil
}
