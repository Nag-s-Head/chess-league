package theme 

type ColourHex string // hex code (with hash) i.e: #ffffff, #a000ff

type AppIconType int

const (
	AppIconType_Png AppIconType = iota
	AppIconType_Jpg
	AppIconType_Svg
)

type Theme struct {
	AppName            string
	VenueName          string
	PrimaryColour      ColourHex
	SecondaryColour    ColourHex
	TitleBarTextColour string
	AppIcon            []byte
	AppIconType        AppIconType
}

func DefaultTheme() Theme {
	appIcon, appIconType := DefaultIcon()

	return Theme{
		PrimaryColour:      "#7000f0",
		SecondaryColour:    "#300080",
		TitleBarTextColour: "#ffffff",
		AppIcon:            appIcon,
		AppIconType:        appIconType,
		AppName:            "Chess League",
		VenueName:          "our club",
	}
}

func DefaultIcon() ([]byte, AppIconType) {
	return []byte(`<svg
  xmlns:dc="http://purl.org/dc/elements/1.1/"
  xmlns:cc="http://creativecommons.org/ns#"
  xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
  xmlns:svg="http://www.w3.org/2000/svg"
  xmlns="http://www.w3.org/2000/svg"
  xmlns:sodipodi="http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd"
  xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape"
  viewBox="0 0 12.7 15.875"
  version="1.1"
  x="0px"
  y="0px"
  height="80"
>
  <g transform="translate(0,-284.29998)">
    <path
      style=""
      d="m 5.3133,286.32059 a 4.0915775,4.0915775 0 0 0 -2.4829,3.2639 8.6577146,8.6577146 0 0 0 0.9277,4.77439 H 2.8514 a 0.27711209,0.27711209 0 0 0 -0.277,0.27711 v 1.55749 a 0.27711209,0.27711209 0 0 0 0.277,0.27711 H 9.259 a 0.27711209,0.27711209 0 0 0 0.2771,-0.27711 v -1.55749 A 0.27711209,0.27711209 0 0 0 9.259,294.35888 H 8.4909 v -0.67878 a 0.55422429,0.55422429 0 0 0 -0.045,-0.21988 9.4437562,9.4437562 0 0 0 -2.1058,-3.04384 h 0.4677 l 0.5031,-0.24373 a 1.1084485,1.1084485 0 0 1 0.9697,0.002 l 0.6909,0.33744 a 0.80203069,0.80203069 0 0 0 1.0669,-1.08422 l -2.4608,-2.40942 a 3.0709507,3.0709507 0 0 0 -1.1636,-0.71448 8.3202992,8.3202992 0 0 1 -0.118,-1.4746 2.578302,2.578302 0 0 0 -0.9822,1.49168 z"
      fill="#ffffff"
      stroke="none"
    />
  </g>
</svg>`), AppIconType_Svg
}
