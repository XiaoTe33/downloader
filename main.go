package main

import "downloader/router"

func main() {
	go router.Accept()
	router.InitRouters()
}
