package main

import (
	"github.com/common-nighthawk/go-figure"
	voting "github.com/lol1pop/info_security_labs/lab5"
)

func main() {
	lab5 := figure.NewFigure("lab 5 AnonVote", "puffy", true)
	lab5.Print()
	/** lab 5 **/
	voting.StartAnonVote()

	//myFigure := figure.NewFigure("lab 4 poker", "invita", true)
	//myFigure.Print()
	///** lab 4 **/
	//poker.StartPoker()

	/** lab 3 **/
	//signature_el_gam.StartElGamale()
	//signature_rsa.StartRsaSign()
	//signature_gost.StartGost()

	/** lab 2 **/
	//shamir.StartShamir()
	//el_gamalya.StartElGamalya()
	//rsa.StartRSA()
	//vernam.VernamExample()
}
func ShowAllFonts() {
	var fonts = []string{"3-d", "3x5", "5lineoblique", "acrobatic", "alligator", "alligator2", "alphabet", "avatar", "banner", "banner3-D", "banner3", "banner4", "barbwire", "basic", "bell", "big", "bigchief", "binary", "block", "bubble", "bulbhead", "calgphy2", "caligraphy", "catwalk", "chunky", "coinstak", "colossal", "computer", "contessa", "contrast", "cosmic", "cosmike", "cricket", "cursive", "cyberlarge", "cybermedium", "cybersmall", "diamond", "digital", "doh", "doom", "dotmatrix", "drpepper", "eftichess", "eftifont", "eftipiti", "eftirobot", "eftitalic", "eftiwall", "eftiwater", "epic", "fender", "fourtops", "fuzzy", "goofy", "gothic", "graffiti", "hollywood", "invita", "isometric1", "isometric2", "isometric3", "isometric4", "italic", "ivrit", "jazmine", "jerusalem", "katakana", "kban", "larry3d", "lcd", "lean", "letters", "linux", "lockergnome", "madrid", "marquee", "maxfour", "mike", "mini", "mirror", "mnemonic", "morse", "moscow", "nancyj-fancy", "nancyj-underlined", "nancyj", "nipples", "ntgreek", "o8", "ogre", "pawp", "peaks", "pebbles", "pepper", "poison", "puffy", "pyramid", "rectangles", "relief", "relief2", "rev", "roman", "rot13", "rounded", "rowancap", "rozzo", "runic", "runyc", "sblood", "script", "serifcap", "shadow", "short", "slant", "slide", "slscript", "small", "smisome1", "smkeyboard", "smscript", "smshadow", "smslant", "smtengwar", "speed", "stampatello", "standard", "starwars", "stellar", "stop", "straight", "tanja", "tengwar", "term", "thick", "thin", "threepoint", "ticks", "ticksslant", "tinker-toy", "tombstone", "trek", "tsalagi", "twopoint", "univers", "usaflag", "wavy", "weird"}
	for _, font := range fonts {
		myFigure := figure.NewFigure("lab 4 poker", font, true)
		myFigure.Print()
		println(font)
	}
}
