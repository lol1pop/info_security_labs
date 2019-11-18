package voting

import (
	"github.com/lol1pop/info_security_labs/lab5/client"
	"github.com/lol1pop/info_security_labs/lab5/server"
)

func StartAnonVote() {

	key := server.InitServer()
	vote_Alisa := client.InitUser("Alisa", 0, key)
	vS := server.ReadVote(vote_Alisa)

}
