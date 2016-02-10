package ascii

const (
	keyEscape                    = 27
	ascii_palette                = "   ...',;:clodxkO0KXNWM"
	ascii_palette_length         = len(ascii_palette) - 1
	threshold_low        float64 = 0.1
	threshold_high       float64 = 1 - threshold_low
)

var (
	colorBlack   = []byte{keyEscape, '[', '3', '0', 'm'}
	colorRed     = []byte{keyEscape, '[', '3', '1', 'm'}
	colorGreen   = []byte{keyEscape, '[', '3', '2', 'm'}
	colorYellow  = []byte{keyEscape, '[', '3', '3', 'm'}
	colorBlue    = []byte{keyEscape, '[', '3', '4', 'm'}
	colorMagenta = []byte{keyEscape, '[', '3', '5', 'm'}
	colorCyan    = []byte{keyEscape, '[', '3', '6', 'm'}
	colorWhite   = []byte{keyEscape, '[', '3', '7', 'm'}
	colorReset   = []byte{keyEscape, '[', '0', 'm'}
)
