package main

import (
	"fmt"
	"github.com/fatih/color"
)

func show_description() {

	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgHiCyan).SprintFunc()
	fmt.Println()
	fmt.Println("Bashlist :: Available Commands")
	fmt.Println()
	fmt.Printf("%s -- List all current directories in bashlist\n", green("bashls"))
	fmt.Printf("  - %s:  bashls \n", cyan("Example"))
	fmt.Println()

	fmt.Printf("%s -- Uploads a directory to bashlist\n", green("push"))
	fmt.Printf("  - %s:  bashls push photos\n", cyan("Example"))
	fmt.Println()

	fmt.Printf("%s -- Downloads a directory from bashlist\n", green("pull"))
	fmt.Printf("  - %s:  bashls pull photos\n", cyan("Example"))
	fmt.Println()

	fmt.Printf("%s -- Deletes a directory from bashlist\n", green("del"))
	fmt.Printf("  - %s:  bashls del photos\n", cyan("Example"))
	fmt.Println()

	fmt.Printf("%s -- Opens bashlist account page in browser\n", green("account"))
	fmt.Printf("  - %s:  bashls account\n", cyan("Example"))
	fmt.Println()

}
