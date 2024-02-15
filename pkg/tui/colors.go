package tui

import (
	"github.com/charmbracelet/lipgloss"
)

type ColorPair struct {
	Name  string
	Value lipgloss.Color
}

const (
	White = lipgloss.Color("#fff")
	Black = lipgloss.Color("#000")

	Blue      = lipgloss.Color("#0D6EFD")
	Blue100   = lipgloss.Color("#CFE2FF")
	Blue200   = lipgloss.Color("#9EC5FE")
	Blue300   = lipgloss.Color("#6EA8FE")
	Blue400   = lipgloss.Color("#3D8BFD")
	Blue500   = lipgloss.Color("#0D6EFD")
	Blue600   = lipgloss.Color("#0A58CA")
	Blue700   = lipgloss.Color("#084298")
	Blue800   = lipgloss.Color("#052C65")
	Blue900   = lipgloss.Color("#031633")
	Indigo    = lipgloss.Color("#6610F2")
	Indigo100 = lipgloss.Color("#E0CFFC")
	Indigo200 = lipgloss.Color("#C29FFA")
	Indigo300 = lipgloss.Color("#A370F7")
	Indigo400 = lipgloss.Color("#8540F5")
	Indigo500 = lipgloss.Color("#6610F2")
	Indigo600 = lipgloss.Color("#520DC2")
	Indigo700 = lipgloss.Color("#3D0A91")
	Indigo800 = lipgloss.Color("#290661")
	Indigo900 = lipgloss.Color("#140330")
	Purple    = lipgloss.Color("#6F42C1")
	Purple100 = lipgloss.Color("#E2D9F3")
	Purple200 = lipgloss.Color("#C5B3E6")
	Purple300 = lipgloss.Color("#A98EDA")
	Purple400 = lipgloss.Color("#8C68CD")
	Purple500 = lipgloss.Color("#6F42C1")
	Purple600 = lipgloss.Color("#59359A")
	Purple700 = lipgloss.Color("#432874")
	Purple800 = lipgloss.Color("#2C1A4D")
	Purple900 = lipgloss.Color("#160D27")
	Pink      = lipgloss.Color("#D63384")
	Pink100   = lipgloss.Color("#F7D6E6")
	Pink200   = lipgloss.Color("#EFADCE")
	Pink300   = lipgloss.Color("#E685B5")
	Pink400   = lipgloss.Color("#DE5C9D")
	Pink500   = lipgloss.Color("#D63384")
	Pink600   = lipgloss.Color("#AB296A")
	Pink700   = lipgloss.Color("#801F4F")
	Pink800   = lipgloss.Color("#561435")
	Pink900   = lipgloss.Color("#2B0A1A")
	Red       = lipgloss.Color("#DC3545")
	Red100    = lipgloss.Color("#F8D7DA")
	Red200    = lipgloss.Color("#F1AEB5")
	Red300    = lipgloss.Color("#EA868F")
	Red400    = lipgloss.Color("#E35D6A")
	Red500    = lipgloss.Color("#DC3545")
	Red600    = lipgloss.Color("#B02A37")
	Red700    = lipgloss.Color("#842029")
	Red800    = lipgloss.Color("#58151C")
	Red900    = lipgloss.Color("#2C0B0E")
	Orange    = lipgloss.Color("#FD7E14")
	Orange100 = lipgloss.Color("#FFE5D0")
	Orange200 = lipgloss.Color("#FECBA1")
	Orange300 = lipgloss.Color("#FEB272")
	Orange400 = lipgloss.Color("#FD9843")
	Orange500 = lipgloss.Color("#FD7E14")
	Orange600 = lipgloss.Color("#CA6510")
	Orange700 = lipgloss.Color("#984C0C")
	Orange800 = lipgloss.Color("#653208")
	Orange900 = lipgloss.Color("#331904")
	Yellow    = lipgloss.Color("#FFC107")
	Yellow100 = lipgloss.Color("#FFF3CD")
	Yellow200 = lipgloss.Color("#FFE69C")
	Yellow300 = lipgloss.Color("#FFDA6A")
	Yellow400 = lipgloss.Color("#FFCD39")
	Yellow500 = lipgloss.Color("#FFC107")
	Yellow600 = lipgloss.Color("#CC9A06")
	Yellow700 = lipgloss.Color("#997404")
	Yellow800 = lipgloss.Color("#664D03")
	Yellow900 = lipgloss.Color("#332701")
	Green     = lipgloss.Color("#198754")
	Green100  = lipgloss.Color("#D1E7DD")
	Green200  = lipgloss.Color("#A3CFBB")
	Green300  = lipgloss.Color("#75B798")
	Green400  = lipgloss.Color("#479F76")
	Green500  = lipgloss.Color("#198754")
	Green600  = lipgloss.Color("#146C43")
	Green700  = lipgloss.Color("#0F5132")
	Green800  = lipgloss.Color("#0A3622")
	Green900  = lipgloss.Color("#051B11")
	Teal      = lipgloss.Color("#20C997")
	Teal100   = lipgloss.Color("#D2F4EA")
	Teal200   = lipgloss.Color("#A6E9D5")
	Teal300   = lipgloss.Color("#79DFC1")
	Teal400   = lipgloss.Color("#4DD4AC")
	Teal500   = lipgloss.Color("#20C997")
	Teal600   = lipgloss.Color("#1AA179")
	Teal700   = lipgloss.Color("#13795B")
	Teal800   = lipgloss.Color("#0D503C")
	Teal900   = lipgloss.Color("#06281E")
	Cyan      = lipgloss.Color("#0DCAF0")
	Cyan100   = lipgloss.Color("#CFF4FC")
	Cyan200   = lipgloss.Color("#9EEAF9")
	Cyan300   = lipgloss.Color("#6EDFF6")
	Cyan400   = lipgloss.Color("#3DD5F3")
	Cyan500   = lipgloss.Color("#0DCAF0")
	Cyan600   = lipgloss.Color("#0AA2C0")
	Cyan700   = lipgloss.Color("#087990")
	Cyan800   = lipgloss.Color("#055160")
	Cyan900   = lipgloss.Color("#032830")
	Gray      = lipgloss.Color("#ADB5BD")
	Gray100   = lipgloss.Color("#EFF0F2")
	Gray200   = lipgloss.Color("#DEE1E5")
	Gray300   = lipgloss.Color("#CED3D7")
	Gray400   = lipgloss.Color("#BDC4CA")
	Gray500   = lipgloss.Color("#ADB5BD")
	Gray600   = lipgloss.Color("#8A9197")
	Gray700   = lipgloss.Color("#686D71")
	Gray800   = lipgloss.Color("#45484C")
	Gray900   = lipgloss.Color("#232426")
)

var (
	BlueFamily = []ColorPair{
		{
			Name:  "Blue100",
			Value: Blue100,
		},
		{
			Name:  "Blue200",
			Value: Blue200,
		},
		{
			Name:  "Blue300",
			Value: Blue300,
		},
		{
			Name:  "Blue400",
			Value: Blue400,
		},
		{
			Name:  "Blue500",
			Value: Blue500,
		},
		{
			Name:  "Blue600",
			Value: Blue600,
		},
		{
			Name:  "Blue700",
			Value: Blue700,
		},
		{
			Name:  "Blue800",
			Value: Blue800,
		},
		{
			Name:  "Blue900",
			Value: Blue900,
		},
	}
	IndigoFamily = []ColorPair{
		{
			Name:  "Indigo100",
			Value: Indigo100,
		},
		{
			Name:  "Indigo200",
			Value: Indigo200,
		},
		{
			Name:  "Indigo300",
			Value: Indigo300,
		},
		{
			Name:  "Indigo400",
			Value: Indigo400,
		},
		{
			Name:  "Indigo500",
			Value: Indigo500,
		},
		{
			Name:  "Indigo600",
			Value: Indigo600,
		},
		{
			Name:  "Indigo700",
			Value: Indigo700,
		},
		{
			Name:  "Indigo800",
			Value: Indigo800,
		},
		{
			Name:  "Indigo900",
			Value: Indigo900,
		},
	}
	PurpleFamily = []ColorPair{
		{
			Name:  "Purple100",
			Value: Purple100,
		},
		{
			Name:  "Purple200",
			Value: Purple200,
		},
		{
			Name:  "Purple300",
			Value: Purple300,
		},
		{
			Name:  "Purple400",
			Value: Purple400,
		},
		{
			Name:  "Purple500",
			Value: Purple500,
		},
		{
			Name:  "Purple600",
			Value: Purple600,
		},
		{
			Name:  "Purple700",
			Value: Purple700,
		},
		{
			Name:  "Purple800",
			Value: Purple800,
		},
		{
			Name:  "Purple900",
			Value: Purple900,
		},
	}
	PinkFamily = []ColorPair{
		{
			Name:  "Pink100",
			Value: Pink100,
		},
		{
			Name:  "Pink200",
			Value: Pink200,
		},
		{
			Name:  "Pink300",
			Value: Pink300,
		},
		{
			Name:  "Pink400",
			Value: Pink400,
		},
		{
			Name:  "Pink500",
			Value: Pink500,
		},
		{
			Name:  "Pink600",
			Value: Pink600,
		},
		{
			Name:  "Pink700",
			Value: Pink700,
		},
		{
			Name:  "Pink800",
			Value: Pink800,
		},
		{
			Name:  "Pink900",
			Value: Pink900,
		},
	}
	RedFamily = []ColorPair{
		{
			Name:  "Red100",
			Value: Red100,
		},
		{
			Name:  "Red200",
			Value: Red200,
		},
		{
			Name:  "Red300",
			Value: Red300,
		},
		{
			Name:  "Red400",
			Value: Red400,
		},
		{
			Name:  "Red500",
			Value: Red500,
		},
		{
			Name:  "Red600",
			Value: Red600,
		},
		{
			Name:  "Red700",
			Value: Red700,
		},
		{
			Name:  "Red800",
			Value: Red800,
		},
		{
			Name:  "Red900",
			Value: Red900,
		},
	}
	OrangeFamily = []ColorPair{
		{
			Name:  "Orange100",
			Value: Orange100,
		},
		{
			Name:  "Orange200",
			Value: Orange200,
		},
		{
			Name:  "Orange300",
			Value: Orange300,
		},
		{
			Name:  "Orange400",
			Value: Orange400,
		},
		{
			Name:  "Orange500",
			Value: Orange500,
		},
		{
			Name:  "Orange600",
			Value: Orange600,
		},
		{
			Name:  "Orange700",
			Value: Orange700,
		},
		{
			Name:  "Orange800",
			Value: Orange800,
		},
		{
			Name:  "Orange900",
			Value: Orange900,
		},
	}
	YellowFamily = []ColorPair{
		{
			Name:  "Yellow100",
			Value: Yellow100,
		},
		{
			Name:  "Yellow200",
			Value: Yellow200,
		},
		{
			Name:  "Yellow300",
			Value: Yellow300,
		},
		{
			Name:  "Yellow400",
			Value: Yellow400,
		},
		{
			Name:  "Yellow500",
			Value: Yellow500,
		},
		{
			Name:  "Yellow600",
			Value: Yellow600,
		},
		{
			Name:  "Yellow700",
			Value: Yellow700,
		},
		{
			Name:  "Yellow800",
			Value: Yellow800,
		},
		{
			Name:  "Yellow900",
			Value: Yellow900,
		},
	}
	GreenFamily = []ColorPair{
		{
			Name:  "Green100",
			Value: Green100,
		},
		{
			Name:  "Green200",
			Value: Green200,
		},
		{
			Name:  "Green300",
			Value: Green300,
		},
		{
			Name:  "Green400",
			Value: Green400,
		},
		{
			Name:  "Green500",
			Value: Green500,
		},
		{
			Name:  "Green600",
			Value: Green600,
		},
		{
			Name:  "Green700",
			Value: Green700,
		},
		{
			Name:  "Green800",
			Value: Green800,
		},
		{
			Name:  "Green900",
			Value: Green900,
		},
	}
	TealFamily = []ColorPair{
		{
			Name:  "Teal100",
			Value: Teal100,
		},
		{
			Name:  "Teal200",
			Value: Teal200,
		},
		{
			Name:  "Teal300",
			Value: Teal300,
		},
		{
			Name:  "Teal400",
			Value: Teal400,
		},
		{
			Name:  "Teal500",
			Value: Teal500,
		},
		{
			Name:  "Teal600",
			Value: Teal600,
		},
		{
			Name:  "Teal700",
			Value: Teal700,
		},
		{
			Name:  "Teal800",
			Value: Teal800,
		},
		{
			Name:  "Teal900",
			Value: Teal900,
		},
	}
	CyanFamily = []ColorPair{
		{
			Name:  "Cyan100",
			Value: Cyan100,
		},
		{
			Name:  "Cyan200",
			Value: Cyan200,
		},
		{
			Name:  "Cyan300",
			Value: Cyan300,
		},
		{
			Name:  "Cyan400",
			Value: Cyan400,
		},
		{
			Name:  "Cyan500",
			Value: Cyan500,
		},
		{
			Name:  "Cyan600",
			Value: Cyan600,
		},
		{
			Name:  "Cyan700",
			Value: Cyan700,
		},
		{
			Name:  "Cyan800",
			Value: Cyan800,
		},
		{
			Name:  "Cyan900",
			Value: Cyan900,
		},
	}
	GrayFamily = []ColorPair{
		{
			Name:  "Gray100",
			Value: Gray100,
		},
		{
			Name:  "Gray200",
			Value: Gray200,
		},
		{
			Name:  "Gray300",
			Value: Gray300,
		},
		{
			Name:  "Gray400",
			Value: Gray400,
		},
		{
			Name:  "Gray500",
			Value: Gray500,
		},
		{
			Name:  "Gray600",
			Value: Gray600,
		},
		{
			Name:  "Gray700",
			Value: Gray700,
		},
		{
			Name:  "Gray800",
			Value: Gray800,
		},
		{
			Name:  "Gray900",
			Value: Gray900,
		},
	}

	ColorsFamilies = []string{
		"Blue",
		"Indigo",
		"Purple",
		"Pink",
		"Red",
		"Orange",
		"Yellow",
		"Green",
		"Teal",
		"Cyan",
		"Gray",
	}

	// All colors, grouped by their family
	ColorsByFamily = map[string][]ColorPair{
		"Blue":   BlueFamily,
		"Indigo": IndigoFamily,
		"Purple": PurpleFamily,
		"Pink":   PinkFamily,
		"Red":    RedFamily,
		"Orange": OrangeFamily,
		"Yellow": YellowFamily,
		"Green":  GreenFamily,
		"Teal":   TealFamily,
		"Cyan":   CyanFamily,
		"Gray":   GrayFamily,
	}

	// All known colors in a map to easily look up their name to value
	AllColors = map[string]lipgloss.Color{
		"Blue100":   Blue100,
		"Blue200":   Blue200,
		"Blue300":   Blue300,
		"Blue400":   Blue400,
		"Blue500":   Blue500,
		"Blue600":   Blue600,
		"Blue700":   Blue700,
		"Blue800":   Blue800,
		"Blue900":   Blue900,
		"Indigo100": Indigo100,
		"Indigo200": Indigo200,
		"Indigo300": Indigo300,
		"Indigo400": Indigo400,
		"Indigo500": Indigo500,
		"Indigo600": Indigo600,
		"Indigo700": Indigo700,
		"Indigo800": Indigo800,
		"Indigo900": Indigo900,
		"Purple100": Purple100,
		"Purple200": Purple200,
		"Purple300": Purple300,
		"Purple400": Purple400,
		"Purple500": Purple500,
		"Purple600": Purple600,
		"Purple700": Purple700,
		"Purple800": Purple800,
		"Purple900": Purple900,
		"Pink100":   Pink100,
		"Pink200":   Pink200,
		"Pink300":   Pink300,
		"Pink400":   Pink400,
		"Pink500":   Pink500,
		"Pink600":   Pink600,
		"Pink700":   Pink700,
		"Pink800":   Pink800,
		"Pink900":   Pink900,
		"Red100":    Red100,
		"Red200":    Red200,
		"Red300":    Red300,
		"Red400":    Red400,
		"Red500":    Red500,
		"Red600":    Red600,
		"Red700":    Red700,
		"Red800":    Red800,
		"Red900":    Red900,
		"Orange100": Orange100,
		"Orange200": Orange200,
		"Orange300": Orange300,
		"Orange400": Orange400,
		"Orange500": Orange500,
		"Orange600": Orange600,
		"Orange700": Orange700,
		"Orange800": Orange800,
		"Orange900": Orange900,
		"Yellow100": Yellow100,
		"Yellow200": Yellow200,
		"Yellow300": Yellow300,
		"Yellow400": Yellow400,
		"Yellow500": Yellow500,
		"Yellow600": Yellow600,
		"Yellow700": Yellow700,
		"Yellow800": Yellow800,
		"Yellow900": Yellow900,
		"Green100":  Green100,
		"Green200":  Green200,
		"Green300":  Green300,
		"Green400":  Green400,
		"Green500":  Green500,
		"Green600":  Green600,
		"Green700":  Green700,
		"Green800":  Green800,
		"Green900":  Green900,
		"Teal100":   Teal100,
		"Teal200":   Teal200,
		"Teal300":   Teal300,
		"Teal400":   Teal400,
		"Teal500":   Teal500,
		"Teal600":   Teal600,
		"Teal700":   Teal700,
		"Teal800":   Teal800,
		"Teal900":   Teal900,
		"Cyan100":   Cyan100,
		"Cyan200":   Cyan200,
		"Cyan300":   Cyan300,
		"Cyan400":   Cyan400,
		"Cyan500":   Cyan500,
		"Cyan600":   Cyan600,
		"Cyan700":   Cyan700,
		"Cyan800":   Cyan800,
		"Cyan900":   Cyan900,
		"Gray100":   Gray100,
		"Gray200":   Gray200,
		"Gray300":   Gray300,
		"Gray400":   Gray400,
		"Gray500":   Gray500,
		"Gray600":   Gray600,
		"Gray700":   Gray700,
		"Gray800":   Gray800,
		"Gray900":   Gray900,
	}
)

var AlphaMap = map[int]string{
	100: "FF",
	99:  "FC",
	98:  "FA",
	97:  "F7",
	96:  "F5",
	95:  "F2",
	94:  "F0",
	93:  "ED",
	92:  "EB",
	91:  "E8",
	90:  "E6",
	89:  "E3",
	88:  "E0",
	87:  "DE",
	86:  "DB",
	85:  "D9",
	84:  "D6",
	83:  "D4",
	82:  "D1",
	81:  "CF",
	80:  "CC",
	79:  "C9",
	78:  "C7",
	77:  "C4",
	76:  "C2",
	75:  "BF",
	74:  "BD",
	73:  "BA",
	72:  "B8",
	71:  "B5",
	70:  "B3",
	69:  "B0",
	68:  "AD",
	67:  "AB",
	66:  "A8",
	65:  "A6",
	64:  "A3",
	63:  "A1",
	62:  "9E",
	61:  "9C",
	60:  "99",
	59:  "96",
	58:  "94",
	57:  "91",
	56:  "8F",
	55:  "8C",
	54:  "8A",
	53:  "87",
	52:  "85",
	51:  "82",
	50:  "80",
	49:  "7D",
	48:  "7A",
	47:  "78",
	46:  "75",
	45:  "73",
	44:  "70",
	43:  "6E",
	42:  "6B",
	41:  "69",
	40:  "66",
	39:  "63",
	38:  "61",
	37:  "5E",
	36:  "5C",
	35:  "59",
	34:  "57",
	33:  "54",
	32:  "52",
	31:  "4F",
	30:  "4D",
	29:  "4A",
	28:  "47",
	27:  "45",
	26:  "42",
	25:  "40",
	24:  "3D",
	23:  "3B",
	22:  "38",
	21:  "36",
	20:  "33",
	19:  "30",
	18:  "2E",
	17:  "2B",
	16:  "29",
	15:  "26",
	14:  "24",
	13:  "21",
	12:  "1F",
	11:  "1C",
	10:  "1A",
	9:   "17",
	8:   "14",
	7:   "12",
	6:   "0F",
	5:   "0D",
	4:   "0A",
	3:   "08",
	2:   "05",
	1:   "03",
	0:   "00",
}
