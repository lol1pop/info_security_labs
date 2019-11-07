package main

import (
	"github.com/lol1pop/info_security_labs/lab3/signature_gost"
	"github.com/lol1pop/info_security_labs/lab3/signature_rsa"
)

func main() {
	//** lab 3 **
	//signature_el_gam.StartElGamale()
	signature_rsa.StartRsaSign()
	signature_gost.StartGost()

	//** lab 2 **
	//shamir.StartShamir()
	//el_gamalya.StartElGamalya()
	//rsa.StartRSA()
	//vernam.VernamExample()
}
