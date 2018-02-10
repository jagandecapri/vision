package color

import "github.com/lucasb-eyer/go-colorful"

type ColorHelperInterface interface{
	GetRandomColors(num_of_colors int) []string
}

type ColorHelper struct{

}

func (ch *ColorHelper) GetRandomColors(num_of_colors int) []string{
	colors := colorful.FastHappyPalette(num_of_colors)
	hex_colors := []string{}
	for _, color := range colors{
		hex_color := color.Hex()
		hex_colors = append(hex_colors, hex_color)
	}
	return hex_colors
}


