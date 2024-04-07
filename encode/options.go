package encode

type EncodeOptions struct {
	vcodec string
	acodec string
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
