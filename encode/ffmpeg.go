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

func ffmpeg(avsPath string, outputPath string, isDualMonoMode bool, option *EncodeOptions, logger func(string)) error {
	args := []string{}

	if option.hwaccel != "" {
		args = append(args, "-hwaccel", option.hwaccel)
	}
	if option.hwaccelOutputFormat != "" {
		args = append(args, "-hwaccel_output_format", option.hwaccelOutputFormat)
	}
	if option.vaapiDevice != "" {
		args = append(args, "-vaapi_device", option.vaapiDevice)
	}

	args = append(args, []string{
		"-f", "avisynth",
		"-i", avsPath,
		"-c:v", option.vcodec,
		"-c:a", option.acodec,
		"-progress", "-",
	}...)
	if option.vf != "" {
		args = append(args, "-vf", option.vf)
	}
	if isDualMonoMode {
		args = append(args, "-dual_mono_mode", "main")
	}
	args = append(args, outputPath)

	return execute("ffmpeg", args, logger)
}
