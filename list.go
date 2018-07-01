package main

import (
	"github.com/fatih/color"
	"github.com/imroc/req"
	"time"
	"github.com/buger/jsonparser"
	"os"
	"encoding/json"
	"github.com/olekukonko/tablewriter"
	"fmt"
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
			color.Cyan("Your Bashlist Storage is Empty!")
			color.Cyan("Upload your first directory using bashls push")
			return
		} else {

			response, _, _, err := jsonparser.Get(byteResp, "Data")
			//fmt.Print(string(response))
			if err != nil {
				color.Red("An unexpected error occurred! Aborting operation.")
				return

			}
			var keys []BLObject
			json.Unmarshal(response, &keys)
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Description", "Size", "Updated", "Status"})
			fmt.Println(keys)
			for _, obj := range keys {
				var s []string
				s = append(s, obj.Name, obj.Description, obj.Size, obj.Updated, obj.Status)
				table.Append(s)
			}
			table.Render()
			return

		}
	}
}
