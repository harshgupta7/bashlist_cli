package main

import (
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/fatih/color"
	"github.com/imroc/req"
	"github.com/olekukonko/tablewriter"
	"os"
	"time"
)

type BLObject struct {
	Name        string
	Size        string
	Updated     string
	Description string
	Status      string
}

func print_list() {
	endpoint := URL + BASHLIST_LIST_URL
	usernamePtr, passwordPtr, _ := authHandler(0)

	username := *usernamePtr
	hashedPassword := *passwordPtr
	r := req.New()
	authHeader := req.Header{
		"Email":    username,
		"Password": hashedPassword,
	}
	r.SetTimeout(25 * time.Second)
	c, err := r.Get(endpoint, authHeader)
	if err != nil {
		c, err = r.Get(endpoint, authHeader)
		if err != nil {
			color.Red("Error contacting Server. Please check you connection and try again")
			os.Exit(1)
		}
	}
	d := c.Response().StatusCode
	if d == 403 {
		authHandler(1)
		print_list()
		return
	} else if d == 399 {
		color.Red("Bashlist encountered an unexpected error! Please try again later.")
		return
	} else {
		byteResp, _ := c.ToBytes()
		empty, err := jsonparser.GetString(byteResp, "Empty")
		if err != nil {
			color.Red("An unexpected error occurred! Aborting operation.")
			return
		}
		if empty == "T" {
			cyan := color.New(color.FgCyan).SprintFunc()
			d := color.New(color.FgGreen, color.Bold).SprintFunc()
			color.Cyan("Your Bashlist Storage is Empty!")
			fmt.Printf("%s %s \n", cyan("Upload your first directory using "), d("bls push"))
			fmt.Println()
			fmt.Printf("%s %s \n", cyan("View help page using "), d("bls help"))
			os.Exit(1)
			upload_handler("swcli")
			return
		} else {

			response, _, _, err := jsonparser.Get(byteResp, "Data")
			//fmt.Print(string(response))
			if err != nil {
				color.Red("An unexpected error occurred! Aborting operation.")
				return

			}
			keys := make([]BLObject, 0)
			json.Unmarshal(response, &keys)
			table := tablewriter.NewWriter(os.Stdout)
			table.SetRowLine(false)
			table.SetBorder(false)
			table.SetColumnSeparator(" ")
			var s []string
			s = append(s, "Name", "Updated", "Description", "Size", "Status")
			table.Append(s)

			for _, obj := range keys {
				var s []string
				var c string
				if obj.Description == "~N////V~" {
					c = "None"
				} else {
					c = obj.Description
				}
				s = append(s, obj.Name, obj.Updated, c, string(obj.Size), obj.Status)
				table.Append(s)
			}
			table.Render()
			return

		}
	}
}
