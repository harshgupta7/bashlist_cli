package main

import (
	"fmt"
	"github.com/zalando/go-keyring"
	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	"github.com/imroc/req"
	"github.com/buger/jsonparser"
	"time"
)

var service string = "Bashlist-Credentials"

func get_username_password() (*string, *string, *string) {
	/*Shows a prompt to enter username and password and returns them*/

	//Ask for Username
	color.Set(color.FgGreen)
	fmt.Print("Bashlist Email: ")
	color.Unset()
	var username string
	fmt.Scanln(&username)

	//Ask for Password
	color.Set(color.FgGreen)
	fmt.Print("Bashlist Password: ")
	color.Unset()
	pass, _ := gopass.GetPasswdMasked()
	stringPass := string(pass)
	hashedPass := AuthPassFromPassword(stringPass)
	return &username, &hashedPass, &stringPass
}

func save_secret(url string, secret *string) {
	/* Saves secret in users credential manager*/
	err := keyring.Set(service, url, *secret)
	if err != nil {
		unexpected_event()
	}

}

func retreive_secret(url string) (*string, error) {

	/* Retreives secret from users credentials manager*/

	secret, err := keyring.Get(service, url)
	if err != nil {
		return nil, err
	}
	return &secret, nil
}

func delete_secret(url string) {
	err := keyring.Delete(service, url)
	if err != nil {
		unexpected_event()
	}
}

func incorrect_auth_loop() (*string, *string, *string) {
	/* Infinite Loop that runs till users enters a wrong username/password combination*/
	endpoint := URL + TEST_AUTH_ENDPOINT
	for {
		usernamePtr, passwordPtr, realPassPtr := get_username_password()
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
			color.Red("Authentication Error! Please try again")
			continue
		}
		d := c.Response().StatusCode
		if d == 403 {
			fmt.Println("Incorrect Username or Password. Please try again.")
		} else {
			byteResp, _ := c.ToBytes()
			response, err := jsonparser.GetString(byteResp, "Valid")
			if err != nil {
				continue
			}
			if response == "T" {
				color.Set(color.FgCyan)
				fmt.Print("Do you wish to save your credentials on this machine? [Y/n] ")
				color.Unset()
				var response string
				fmt.Scanln(&response)
				if response == "n" || response == "N" || response == "No" || response == "no" {
					return usernamePtr, passwordPtr, realPassPtr
				}
				save_secret("Bashlist-Credentials/Username", usernamePtr)
				save_secret("Bashlist-Credentials/HashedPassword", passwordPtr)
				save_secret("Bashlist-Credentials/RealPassword", realPassPtr)
				return usernamePtr, passwordPtr, realPassPtr
			} else {
				continue
			}
		}
	}
}

func authHandler(incorrect int) (*string, *string, *string) {
	if incorrect == 0 {
		usernamePtr, err := retreive_secret("Bashlist-Credentials/Username")
		passwordPtr, err1 := retreive_secret("Bashlist-Credentials/HashedPassword")
		realPassPtr, err2 := retreive_secret("Bashlist-Credentials/RealPassword")
		if err != nil || err1 != nil || err2 != nil {
			usernamePtr, passwordPtr, realPassPtr := incorrect_auth_loop()
			return usernamePtr, passwordPtr, realPassPtr
		}
		return usernamePtr, passwordPtr, realPassPtr
	} else {
		usernamePtr, passwordPtr, realPassPtr := incorrect_auth_loop()
		return usernamePtr, passwordPtr, realPassPtr
	}
}
