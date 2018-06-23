package main

import "fmt"
import "github.com/fatih/color"
import "github.com/howeyc/gopass"
import "github.com/imroc/req"
import "github.com/docker/docker-credential-helpers/credentials"

import (
	"github.com/docker/docker-credential-helpers/osxkeychain"
    "github.com/buger/jsonparser"
	"time"
)

var nativeStore = osxkeychain.Osxkeychain{}


func get_username_password()(*string,*string){
	/*Shows a prompt to enter username and password and returns them*/

	// TODO: VALIDATE STRINGS

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
	hashedPass:= AuthPassFromPassword(stringPass)
	return &username,&hashedPass
}

func save_secret(url string, username *string, secret *string) {

	/* Saves secret in users credential manager*/

	c := &credentials.Credentials{
	    ServerURL: url,
	    Username: *username,
	    Secret: *secret,
	}
	nativeStore.Add(c)
}

func retreive_secret(url string)(*string,*string,error) {

	/* Retreives secret from users credentials manager*/

	username, secret,err := nativeStore.Get(url)
	if err!=nil {
		t := ""
		return &t,&t,err
	}
	return &username,&secret,nil
}

func delete_secret(url string){
	nativeStore.Delete(url)
}



func incorrect_auth_loop() {
	/* Infinite Loop that runs till users enters a wrong username/password combination*/
	endpoint := URL + TEST_AUTH_ENDPOINT
	for {
		usernamePtr, passwordPtr:=get_username_password()
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
			fmt.Println("Error contacting server! Please check your connection.")
		}
		d := c.Response().StatusCode
		if d==403{
			fmt.Println("Incorrect Username or Password. Please try again.")
		} else{
			byteResp, _ := c.ToBytes()
			response, err := jsonparser.GetString(byteResp, "Valid")
			if err!=nil{
				continue
			}
			if response=="T"{
				save_secret("Bashlist-Credentials/Credentials",usernamePtr,passwordPtr)
				return
			}else{
				continue
			}
		}
	}
}



func change_password(){
	/*password change handler*/

}

func authHandler(incorrect int)(*string,*string,error){
	if incorrect==0 {
		usernamePtr, passwordPtr, err := retreive_secret("Bashlist-Credentials/Credentials")
		if err != nil {
			incorrect_auth_loop()
			usernamePtr, passwordPtr, err = retreive_secret("Bashlist-Credentials/Credentials")
			return usernamePtr, passwordPtr, err
		}
		return usernamePtr, passwordPtr, err
	} else{
		incorrect_auth_loop()
		usernamePtr, passwordPtr, err := retreive_secret("Bashlist-Credentials/Credentials")
		return usernamePtr, passwordPtr, err

	}
}
