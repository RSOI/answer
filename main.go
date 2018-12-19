// QUESTION SERVICE

package main

import (
	"fmt"
	"os"

	"github.com/RSOI/answer/controller"
	"github.com/RSOI/answer/database"
	"github.com/RSOI/answer/utils"
	"github.com/valyala/fasthttp"
)

// PORT application port
const PORT = 8081

func main() {
	if len(os.Args) > 1 {
		utils.DEBUG = os.Args[1] == "debug"
	}
	utils.LOG("Launched in debug mode...")
	utils.LOG(fmt.Sprintf("Answer service is starting on localhost: %d", PORT))

	controller.Init(database.Connect())
	fasthttp.ListenAndServe(fmt.Sprintf(":%d", PORT), initRoutes().Handler)
}
