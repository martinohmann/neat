package style

import "strconv"

// Attribute is an attribute of an ANSI escape sequence.
type Attribute interface {
	sequence() string
}

// SimpleAttribute is an attribute that is just one uint8 value.
type SimpleAttribute uint8

// sequence implements Attribute.
func (a SimpleAttribute) sequence() string {
	return strconv.Itoa(int(a))
}

// Base style attributes
const (
	Reset SimpleAttribute = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

// Font attributes
const (
	DefaultFont SimpleAttribute = iota + 10
	AltFont1
	AltFont2
	AltFont3
	AltFont4
	AltFont5
	AltFont6
	AltFont7
	AltFont8
	AltFont9
)

// Style reset attributes
const (
	Fraktur SimpleAttribute = iota + 20
	DoubleUnderline
	Normal
	NoFraktur
	NoUnderline
	NoBlink
	ProportionalSpacing
	NoReverse
	Reveal
	NoCrossedOut
)

// Foreground text colors
const (
	FgBlack SimpleAttribute = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
	FgColor
	FgDefault
)

// Background text colors
const (
	BgBlack SimpleAttribute = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
	BgColor
	BgDefault
)

// Foreground Hi-Intensity text colors
const (
	FgHiBlack SimpleAttribute = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

// Background Hi-Intensity text colors
const (
	BgHiBlack SimpleAttribute = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)

const (
	colorModeRGB SimpleAttribute = 2
	colorMode256 SimpleAttribute = 5
)

var (
	// AttributeMap contains a mapping between names and attributes. This is
	// used by StyleString to search an replace attribute names with their ANSI
	// color escape sequences. It is exported so that users can add more
	// mappings if desired.
	AttributeMap = map[string]Attribute{
		// 0-9
		"reset":        Reset,
		"bold":         Bold,
		"faint":        Faint,
		"italic":       Italic,
		"underline":    Underline,
		"blinkslow":    BlinkSlow,
		"blinkrapid":   BlinkRapid,
		"reversevideo": ReverseVideo,
		"concealed":    Concealed,
		"crossedout":   CrossedOut,

		// 20-29
		"fraktur":             Fraktur,
		"doubleunderline":     DoubleUnderline,
		"normal":              Normal,
		"nofraktur":           NoFraktur,
		"nounderline":         NoUnderline,
		"noblink":             NoBlink,
		"proportionalspacing": ProportionalSpacing,
		"noreverse":           NoReverse,
		"reveal":              Reveal,
		"nocrossedout":        NoCrossedOut,

		// 30-39
		"fgblack":   FgBlack,
		"fgred":     FgRed,
		"fggreen":   FgGreen,
		"fgyellow":  FgYellow,
		"fgblue":    FgBlue,
		"fgmagenta": FgMagenta,
		"fgcyan":    FgCyan,
		"fgwhite":   FgWhite,
		"fgcolor":   FgColor,
		"fgdefault": FgDefault,

		// 40-49
		"bgblack":   BgBlack,
		"bgred":     BgRed,
		"bggreen":   BgGreen,
		"bgyellow":  BgYellow,
		"bgblue":    BgBlue,
		"bgmagenta": BgMagenta,
		"bgcyan":    BgCyan,
		"bgwhite":   BgWhite,
		"bcolor":    BgColor,
		"bgdefault": BgDefault,

		// 90-97
		"fghiblack":   FgHiBlack,
		"fghired":     FgHiRed,
		"fghigreen":   FgHiGreen,
		"fghiyellow":  FgHiYellow,
		"fghiblue":    FgHiBlue,
		"fghimagenta": FgHiMagenta,
		"fghicyan":    FgHiCyan,
		"fghiwhite":   FgHiWhite,

		// 100-107
		"bghiblack":   BgHiBlack,
		"bghired":     BgHiRed,
		"bghigreen":   BgHiGreen,
		"bghiyellow":  BgHiYellow,
		"bghiblue":    BgHiBlue,
		"bghimagenta": BgHiMagenta,
		"bghicyan":    BgHiCyan,
		"bghiwhite":   BgHiWhite,

		// fg convenience variants
		"black":   FgBlack,
		"red":     FgRed,
		"green":   FgGreen,
		"yellow":  FgYellow,
		"blue":    FgBlue,
		"magenta": FgMagenta,
		"cyan":    FgCyan,
		"white":   FgWhite,
		"default": FgDefault,

		"rst":           Reset,
		"blink":         BlinkSlow,
		"strike":        CrossedOut,
		"strikethrough": CrossedOut,
	}
)
