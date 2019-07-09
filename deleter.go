package main

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/fatih/color"
	"github.com/imroc/req"
	"os"
	"time"
)

func requestDeletetionDetails(dir_name string, usernamePtr *string, hashedPasswordPtr *string) *[]byte {

	endpoint := URL + DELETIONDETAIL + "/" + dir_name

	r := req.New()
	authHeader := req.Header{
		"Email":    *usernamePtr,
		"Password": *hashedPasswordPtr,
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
	respCode := c.Response().StatusCode

	//auth error

	if respCode == 403 {
		usernamePtr, passwordPtr, _ := authHandler(1)
		authHeader = req.Header{
			"Email":    *usernamePtr,
			"Password": *passwordPtr,
		}
		c, err = r.Get(endpoint, authHeader)
		respCode = c.Response().StatusCode
		if respCode == 403 {
			color.Red("Bashlist could not authenticate you. Please try again.")
			os.Exit(1)
		}
		//no bucket available
	} else if respCode == 367 {
		msg := dir_name + ": No such directory exists in your bashlist"
		color.Red(msg)
		cyan := color.New(color.FgCyan).SprintFunc()
		d := color.New(color.FgGreen, color.Bold).SprintFunc()
		fmt.Printf("%s %s \n", cyan("View your available directories using"), d("bashls"))
		os.Exit(1)
		//found bucket
	} else {
		resp, err := c.ToBytes()
		if err != nil {
			unexpected_event()
		}
		return &resp
	}
	return nil

}

func sendDeleteConfirmation(dir_name string, owner string, shared string, usernamePtr *string, hashedPasswordPtr *string) int {

	r_, err := generate_random_string(16)
	if err != nil {
		r_ = "lsakndasndas"
	}
	endpoint := URL + "api/v01/dbuckconf" + "/" + dir_name + "/" + shared + "/" + owner + "/" + r_

	r := req.New()
	authHeader := req.Header{
		"Email":    *usernamePtr,
		"Password": *hashedPasswordPtr,
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
	respCode := c.Response().StatusCode

	//auth error

	if respCode == 403 {
		usernamePtr, passwordPtr, _ := authHandler(1)
		authHeader = req.Header{
			"Email":    *usernamePtr,
			"Password": *passwordPtr,
		}
		c, err = r.Get(endpoint, authHeader)
		respCode = c.Response().StatusCode
		if respCode == 403 {
			color.Cyan("Bashlist could not authenticate you. Please try again.")
			os.Exit(1)
		}
	} else if respCode != 200 {
		unexpected_event()
	} else {
		return 1
	}
	return 0

}

func deletionHandler(dirname string) {

	usernamePtr, passwordPtr, _ := authHandler(0)
	resp := requestDeletetionDetails(dirname, usernamePtr, passwordPtr)
	if resp == nil {
		unexpected_event()
	}
	exist, err := jsonparser.GetString(*resp, "exist")
	if err != nil {
		unexpected_event()
	}
	if exist != "True" {
		msg := dirname + ": No such directory exists in your bashlist"
		color.Red(msg)
		cyan := color.New(color.FgCyan).SprintFunc()
		d := color.New(color.FgGreen, color.Bold).SprintFunc()
		fmt.Printf("%s %s \n", cyan("View your available directories using"), d("bashls"))
		os.Exit(1)

	}

	owner, errO := jsonparser.GetString(*resp, "owner")
	shared, errS := jsonparser.GetString(*resp, "shared")
	if errO != nil || errS != nil {
		unexpected_event()
	}
	ownerVal := owner == "True"
	sharedVal := shared == "True"
	if (ownerVal && !sharedVal) || (!ownerVal && sharedVal) {
		r := sendDeleteConfirmation(dirname, owner, shared, usernamePtr, passwordPtr)
		if r != 1 {
			unexpected_event()
		} else {
			color.Green("Successfully deleted " + dirname)
			return
		}
	} else if ownerVal && sharedVal {
		color.Set(color.FgCyan)
		fmt.Print(dirname + " is a shared directory. Deleting will remove access for all participants." +
			"Do you want to continue?[Y/n] ")
		color.Unset()
		var response string
		fmt.Scan(&response)
		if response == "n" || response == "N" || response == "No" || response == "no" {
			return
		}
		r := sendDeleteConfirmation(dirname, owner, shared, usernamePtr, passwordPtr)
		if r != 1 {
			unexpected_event()
		} else {
			color.Green("Successfully deleted " + dirname)
			return
		}

	} else {
		unexpected_event()
	}

}
