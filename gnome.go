package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func setBackground(path string) error {
	cmd := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "color-scheme")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	colorSchemeArg := "picture-uri"
	if strings.TrimSpace(string(output)) == "'prefer-dark'" {
		colorSchemeArg = "picture-uri-dark"
	}

	cmd = exec.Command("gsettings", "set", "org.gnome.desktop.background", colorSchemeArg, "file://"+path)
	if _, err := cmd.Output(); err != nil {
		return err
	}

	return nil
}

type Resolution struct {
	Width  int
	Height int
}

func highestMonitorResolution() (Resolution, error) {
	cmd := exec.Command("xrandr", "--listmonitors")
	out, err := cmd.Output()
	if err != nil {
		return Resolution{}, fmt.Errorf("xrandr failed %w", err)
	}

	lines := strings.Split(string(out), "\n")

	best := Resolution{}
	bestArea := 0

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		resPart := fields[2]
		resPart = strings.Split(resPart, "+")[0]

		parts := strings.Split(resPart, "x")
		if len(parts) != 2 {
			continue
		}

		wStr := strings.Split(parts[0], "/")[0]
		hStr := strings.Split(parts[1], "/")[0]

		w, err1 := strconv.Atoi(wStr)
		h, err2 := strconv.Atoi(hStr)
		if err1 != nil || err2 != nil {
			continue
		}

		area := w * h
		if area > bestArea {
			bestArea = area
			best = Resolution{Width: w, Height: h}
		}
	}

	if bestArea == 0 {
		return Resolution{}, fmt.Errorf("no monitors detected")
	}

	return best, nil
}
