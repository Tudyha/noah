package main

import (
	"noah/internal/server"
)

func main() {
	printBanner()

	// new server
	s := server.NewServer()

	//start server
	s.Run()
}

func printBanner() {
	b := `
#    #  ####    ##   #    #
##   # #    #  #  #  #    #
# #  # #    # #    # ######
#  # # #    # ###### #    #
#   ## #    # #    # #    #
#    #  ####  #    # #    #

`
	println(b)
}
