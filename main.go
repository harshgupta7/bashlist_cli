package main

import (
	"os"
)

func main() {

	numArgs := len(os.Args)

	if numArgs > 3 {
		show_description()
		return
	} else if numArgs == 3 {
		dir := os.Args[2]
		if os.Args[1] == "pull" {
			download_manager(dir)
		} else if os.Args[1] == "push" {
			upload_handler(dir)
			return
		} else if os.Args[1] == "del" {
			deletionHandler(dir)
		} else {
			show_description()
			return
		}
	} else if numArgs == 2 {
		if os.Args[1] == "account" {
			open_account_handler()
			return
		} else if os.Args[1] == "help" {
			show_description()
			return
		} else {
			show_description()
			return
		}
	} else if numArgs == 1 {
		if os.Args[0] == "bls" {
			print_list()
			return
		} else {
			show_description()
			return
		}
	} else {
		show_description()
		return

	}
}
