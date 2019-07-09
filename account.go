package main

import (
	"github.com/buger/jsonparser"
	"github.com/fatih/color"
	"github.com/imroc/req"
	"github.com/skratchdot/open-golang/open"
	"os"
	"time"
)

func get_account_url() string {

	/* Fetches unique account url*/

	endpoint := URL + GET_ACCOUNT_URL_ENDPOINT

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
		}
	}

	respCode := c.Response().StatusCode
	if respCode == 403 {
		usernamePtr, passwordPtr, _ = authHandler(1)
		username = *usernamePtr
		hashedPassword = *passwordPtr
		authHeader = req.Header{
			"Email":    username,
			"Password": hashedPassword,
		}
		c, err = r.Get(endpoint, authHeader)
		if err != nil {
			unexpected_event()
		}
		respCode = c.Response().StatusCode
		if respCode == 403 {
			unexpected_event()
		}
	} else {
		byteResp, _ := c.ToBytes()
		response, err := jsonparser.GetString(byteResp, "Url")
		if err != nil {
			unexpected_event()
		}
		return response
	}
	return "NONE"
}

func open_account_url(url string) {
	/* Open account URL in browser*/
	err := open.Run(url)
	if err != nil {
		color.Red("Bashlist encountered an unexpected error. Exiting application.")
		os.Exit(1)
	}
}

func open_account_handler() {
	url := get_account_url()
	if url == "NONE" {
		unexpected_event()
	}
	open_account_url(url)
	return
}
