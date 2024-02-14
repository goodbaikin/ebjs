package encode

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func mergeAvsFiles(tempdir, inputAvsPath, logoframeAvsPath, joinLogoScpAvsPath string) (string, error) {
	mergedAvsPath := filepath.Join(tempdir, "input_cut.avs")
	file, err := os.Create(mergedAvsPath)
	if err != nil {
		return "", fmt.Errorf("failed to create a file: %v", err)
	}

	inputAvs, err := os.Open(inputAvsPath)
	if err != nil {
		return "", fmt.Errorf("failed to open a input avs: %v", err)
	}
	defer inputAvs.Close()
	if _, err := io.Copy(file, inputAvs); err != nil {
		return "", fmt.Errorf("failed to write a input avs: %v", err)
	}

	logoframeAvs, err := os.Open(logoframeAvsPath)
	if err != nil {
		return "", fmt.Errorf("failed to open a logoframe avs: %v", err)
	}
	defer logoframeAvs.Close()
	if _, err := io.Copy(file, logoframeAvs); err != nil {
		return "", fmt.Errorf("failed to write a logoframe avs: %v", err)
	}

	joinLocoScpAvs, err := os.Open(joinLogoScpAvsPath)
	if err != nil {
		return "", fmt.Errorf("failed to open a join_logo_scp avs: %v", err)
	}
	defer joinLocoScpAvs.Close()
	if _, err := io.Copy(file, joinLocoScpAvs); err != nil {
		return "", fmt.Errorf("failed to write a join_logo_scp avs: %v", err)
	}

	return mergedAvsPath, nil
}

func ffprobe(avsPath string, logger func(string)) error {
	return execute("ffprobe", []string{
		"-f", "avisynth",
		"-i", avsPath,
		"-v", "0",
		"-show_entries", "format=duration",
	}, logger)
}

func ffmpeg(avsPath string, outputPath string, isDualMonoMode bool, logger func(string)) error {
	args := []string{
		"-f", "avisynth",
		"-i", avsPath,
		"-vf", "yadif,scale=1920:1080",
		"-c:v", "libx264",
		"-c:a", "aac",
		"-preset", "veryfast",
		"-progress", "-",
		outputPath,
	}
	if isDualMonoMode {
		args = append(args, "-dual_mono_mode", "main")
	}

	return execute("ffmpeg", args, logger)
}
