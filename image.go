package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	"runtime"
	"sort"
	"sync"

	xdraw "golang.org/x/image/draw"
)

const sampleSize = 512

func decode(data []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode data %w", err)
	}

	dst := image.NewRGBA(image.Rect(0, 0, sampleSize, sampleSize))
	xdraw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, img.Bounds(), xdraw.Over, nil)

	return dst, nil
}

func luminanceDistribution(img image.Image) []float64 {
	b := img.Bounds()
	width := b.Dx()
	height := b.Dy()

	workers := min(runtime.NumCPU()/2, height)

	rowsPerWorker := height / workers

	results := make([][]float64, workers)
	var wg sync.WaitGroup

	for i := range workers {
		startY := i * rowsPerWorker
		endY := startY + rowsPerWorker
		if i == workers-1 {
			endY = height
		}

		wg.Add(1)
		go func(idx, y0, y1 int) {
			defer wg.Done()

			local := make([]float64, 0, (y1-y0)*width)

			for y := y0; y < y1; y++ {
				for x := range width {
					r, g, b, _ := img.At(x, y).RGBA()

					R := float64(r) / 65535.0
					G := float64(g) / 65535.0
					B := float64(b) / 65535.0

					Y := 0.2126*R + 0.7152*G + 0.0722*B
					local = append(local, Y)
				}
			}

			results[idx] = local
		}(i, startY, endY)
	}

	wg.Wait()

	// Merge
	total := 0
	for _, r := range results {
		total += len(r)
	}

	lums := make([]float64, 0, total)
	for _, r := range results {
		lums = append(lums, r...)
	}

	return lums
}

func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}

	idx := int(float64(len(sorted)-1) * p)
	return sorted[idx]
}

func isDark(data []byte) (bool, error) {
	if !c.darkEnabled {
		return true, nil
	}

	img, err := decode(data)
	if err != nil {
		return false, err
	}

	lums := luminanceDistribution(img)
	sort.Float64s(lums)

	median := percentile(lums, 0.5)
	p90 := percentile(lums, 0.9)

	return median <= c.darkMedian && p90 <= c.darkP90, nil
}
