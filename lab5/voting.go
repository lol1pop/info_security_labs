package voting

import (
	"github.com/lol1pop/info_security_labs/lab5/server"
)

func StartAnonVote() {
	vote_server := new(server.Server).InitServer()
	client := vote_server.NewClient()
	vote_Alisa := client.InitVote("Alisa", 0)
	vS := vote_server.ReadVote(vote_Alisa.VH)
	n, s := vote_Alisa.Signature(vS)
	result := vote_server.CheckVote(n, s)
	println(result)
}
