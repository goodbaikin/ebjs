package encode

import "path/filepath"

func getJLFileName(channelID uint64) string {
	for _, id := range []uint64{
		3208043008,
		3272102056,
		400101,
		400103,
	} {
		if channelID == id {
			return "JL_NHK.txt"
		}
	}

	for _, id := range []uint64{
		3272202064,
		400161,
	} {
		if channelID == id {
			return "JL_MBS.txt"
		}
	}

	for _, id := range []uint64{
		700333,
	} {
		if channelID == id {
			return "JL_ATX.txt"
		}
	}

	return "JL_標準.txt"
}

func joinLogoScp(tempdir, jldir, logoframeTxtPath, chapterExePath string, channelID uint64, logger func(string)) (string, error) {
	jlPath := filepath.Join(jldir, getJLFileName(channelID))
	joinLogoScpAvsPath := filepath.Join(tempdir, "join_logo_scp.avs")

	if err := execute("join_logo_scp", []string{
		"-inlogo",
		logoframeTxtPath,
		"-inscp",
		chapterExePath,
		"-incmd",
		jlPath,
		"-o",
		joinLogoScpAvsPath,
	}, logger); err != nil {
		return "", nil
	}

	return joinLogoScpAvsPath, nil
}
