package encode

type EncodeOptions struct {
	hwaccel             string
	hwaccelOutputFormat string
	vaapiDevice         string
	vcodec              string
	acodec              string
	vf                  string
}

type EncodeOption func(*EncodeOptions)

func NewEncodeOptions() *EncodeOptions {
	return &EncodeOptions{
		vcodec: "libx264",
		acodec: "libfdk_aac",
	}
}

func WithVCodec(vcodec string) EncodeOption {
	return func(eo *EncodeOptions) {
		eo.vcodec = vcodec
	}
}

func WithACodec(acodec string) EncodeOption {
	return func(eo *EncodeOptions) {
		eo.acodec = acodec
	}
}

func WithVF(vf string) EncodeOption {
	return func(eo *EncodeOptions) {
		eo.vf = vf
	}
}

func WithHWAccel(hw string) EncodeOption {
	return func(eo *EncodeOptions) {
		eo.hwaccel = hw
	}
}

func WithHWAccelOutputFormat(format string) EncodeOption {
	return func(eo *EncodeOptions) {
		eo.hwaccelOutputFormat = format
	}
}

func WithVAAPIDevice(device string) EncodeOption {
	return func(eo *EncodeOptions) {
		eo.vaapiDevice = device
	}
}
