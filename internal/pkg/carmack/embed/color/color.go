// Package color
// source of color values: https://sashamaps.net/docs/resources/20-colors/
package color

const Purple = 0xa970ff
const Blue = 0x0c2aac
const Red = 0xe6194B
const Green = 0x3cb44b
const Yellow = 0xffe119
const Orange = 0xf58231
const Cyan = 0x42d4f4
const Magenta = 0xf032e6
const Lime = 0xbfef45
const Pink = 0xfabed4
const Teal = 0x469990
const Lavender = 0xdcbeff
const Brown = 0x9A6324
const Beige = 0xfffac8
const Maroon = 0x800000
const Mint = 0xaaffc3
const Olive = 0x808000
const Apricot = 0xffd8b1
const Navy = 0x000075
const Grey = 0xa9a9a9

func all() []int {
	return []int{
		Red, Green, Yellow, Blue, Orange, Purple, Cyan, Magenta, Lime, Pink, Teal, Lavender, Brown, Beige, Maroon, Mint, Olive, Apricot, Navy, Grey,
	}
}

func FromIndex(index int) int {
	colors := all()
	return colors[index%len(colors)]
}
