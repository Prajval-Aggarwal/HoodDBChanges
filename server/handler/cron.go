package handler

import (
	"fmt"
	"main/server/services/socket"
	"time"

	socketio "github.com/googollee/go-socket.io"
	"github.com/robfig/cron/v3"
)

func StartCron(server *socketio.Server) {
	c := cron.New()
	fmt.Println("Current time is:", time.Now())

	//Concept

	//Run cron at every minute
	//check the next reward time is same as current time. If yes give the rewards
	c.AddFunc("*/1 * * * *", func() {
		socket.GiveArenaPerks2(server)
	})

	// c.AddFunc("*/6 * * * *", func() {

	// 	fmt.Println("Giving medium perks")
	// 	fmt.Println("medium arena time is:", time.Now())

	// 	socket.GiveArenaPerks2(int64(utils.MEDIUM), server)
	// })

	// c.AddFunc("*/8 * * * *", func() {
	// 	fmt.Println("Giving hard perks")
	// 	fmt.Println("hard arena time is:", time.Now())

	// 	socket.GiveArenaPerks2(int64(utils.HARD), server)
	// })
	c.Start()

}
