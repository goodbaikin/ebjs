package encode

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Encoder struct {
	basedir       string
	recordedDir   string
	logodir       string
	jldir         string
	encodeOptions *EncodeOptions
}

func NewEncoder(basedir, recordedDir string, encodeOptions *EncodeOptions) Encoder {
	return Encoder{
		basedir:       basedir,
		recordedDir:   recordedDir,
		logodir:       filepath.Join(basedir, "logo"),
		jldir:         filepath.Join(basedir, "JL"),
		encodeOptions: encodeOptions,
	}
}

func aviSynth(basedir string, tsPath string) (string, error) {
	temp, err := os.MkdirTemp(basedir, "tmp")
	if err != nil {
		return "", err
	}

	inputAvs := strings.Join([]string{
		"TSFilePath=\"" + tsPath + "\"",
		"LWLibavVideoSource(TSFilePath, fpsnum=30000, fpsden=1001)",
		"AudioDub(last,LWLibavAudioSource(TSFilePath, av_sync=true))",
	}, "\n")
	inputAvsPath := filepath.Join(temp, "input.avs")
	if os.WriteFile(inputAvsPath, []byte(inputAvs), 0600) != nil {
		return "", err
	}

	return inputAvsPath, nil
}

func chapterExe(tempdir string, avsPath string, logger func(string)) (string, error) {
	chapterExePath := filepath.Join(tempdir, "chapter_exe.txt")
	if err := execute("chapter_exe", []string{
		"-v",
		avsPath,
		"-o",
		chapterExePath,
	}, logger); err != nil {
		if _, err := os.Create(chapterExePath); err != nil {
			return "", err
		}
	}

	return chapterExePath, nil
}

func logoframe(tempdir, logodir, avsPath string, channelID uint64, logger func(string)) (string, string, error) {
	logo := fmt.Sprintf("SID%d-1.lgd", channelID%100000)
	logoPath := filepath.Join(logodir, logo)
	if _, err := os.Lstat(logoPath); err != nil {
		fmt.Printf("channelID %dのロゴが見つかりません\n", channelID)
		for {
			time.Sleep(time.Second * 5)
			if _, err := os.Lstat(logoPath); err == nil {
				fmt.Printf("ロゴファイル%sを確認\n", logoPath)
				break
			}
		}
	}

	logoframeTxtPath := filepath.Join(tempdir, "logoframe.txt")
	logoframeAvsPath := filepath.Join(tempdir, "logoframe.avs")
	args := []string{
		avsPath,
		"-oa",
		logoframeTxtPath,
		"-o",
		logoframeAvsPath,
		"-logo",
		logoPath,
	}

	if err := execute("logoframe", args, logger); err != nil {
		return "", "", err
	}
	return logoframeTxtPath, logoframeAvsPath, nil
}

func (e *Encoder) Encode(input string, output string, channelID uint64, isDualMonoMode bool, logger func(string)) error {
	tsPath := filepath.Join(e.recordedDir, input)
	inputAvsPath, err := aviSynth(e.basedir, tsPath)
	if err != nil {
		return err
	}
	tempdir := filepath.Dir(inputAvsPath)
	defer func() {
		os.RemoveAll(tempdir)
		os.Remove(tempdir)
	}()

	chapterExePath, err := chapterExe(tempdir, inputAvsPath, logger)
	if err != nil {
		return fmt.Errorf("chapter_exe failed: %v", err)
	}
	defer os.Remove(filepath.Join(e.recordedDir, input+".lwi"))

	logoframeTxtPath, logoframeAvsPath, err := logoframe(tempdir, e.logodir, inputAvsPath, channelID, logger)
	if err != nil {
		return fmt.Errorf("logoframe failed: %v", err)
	}

	joinLogoScpAvsPath, err := joinLogoScp(tempdir, e.jldir, logoframeTxtPath, chapterExePath, channelID, logger)
	if err != nil {
		return fmt.Errorf("join_logo_scp failed: %v", err)
	}

	mergedAvsPath, err := mergeAvsFiles(tempdir, inputAvsPath, logoframeAvsPath, joinLogoScpAvsPath)
	if err != nil {
		return fmt.Errorf("merge avs files failed: %v", err)
	}

	if err = ffprobe(mergedAvsPath, logger); err != nil {
		return fmt.Errorf("ffprobe failed: %v", err)
	}

	if err = ffmpeg(mergedAvsPath, filepath.Join(e.recordedDir, output), isDualMonoMode, e.encodeOptions, logger); err != nil {
		return fmt.Errorf("ffmpeg failed: %v", err)
	}

	return nil
}
