package main

import (
	"github.com/imroc/req"
	"time"
	"github.com/buger/jsonparser"
	"github.com/fatih/color"
	"github.com/skratchdot/open-golang/open"
	"os"
)

func get_account_url()(string){
	/* Fetches unique account url*/
	endpoint := URL + GET_ACCOUNT_URL_ENDPOINT
	usernamePtr, passwordPtr, err:=authHandler(0)
	if err!=nil {
		color.Red("Authentication Error! Please try again later")
	}
	username := *usernamePtr
	hashedPassword := *passwordPtr
	r := req.New()
	authHeader := req.Header{
		"Email":        username,
		"Password": hashedPassword,
	}
	r.SetTimeout(25 * time.Second)
	c, err := r.Get(endpoint, authHeader)
	if err!=nil{
		c, err = r.Get(endpoint, authHeader)
		if err!=nil {
			color.Red("Error contacting Server. Please check you connection and try again")
		}
	}
	d := c.Response().StatusCode
	if d==403{
		authHandler(1)
		e := get_account_url()
		return e
	} else{
		byteResp, _ := c.ToBytes()
		response, err := jsonparser.GetString(byteResp, "Url")
		if err!=nil{
			color.Red("An unexpected error occurred! Aborting operation.")
		}
		return response
	}
	return "NONE"
}

func open_account_url(url string ){
	/* Open account URL in browser*/
	err := open.Run(url)
	if err!=nil{
		color.Red("Bashlist encountered an unexpected error. Exiting application.")
		os.Exit(1)
	}
}

func open_account_handler(){
	url := get_account_url()
	open_account_url(url)
	return
}