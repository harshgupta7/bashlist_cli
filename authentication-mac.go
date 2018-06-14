package main

import "os"
import "fmt"
import "github.com/fatih/color"
import "github.com/howeyc/gopass"
import "github.com/docker/docker-credential-helpers/credentials"
import "github.com/docker/docker-credential-helpers/osxkeychain"

var nativeStore = osxkeychain.Osxkeychain{}


func get_username_password()(*string,*string){
	/*Shows a prompt to enter username and password and returns them*/

	//Ask for Username
	color.Set(color.FgGreen)
	fmt.Print(("Bashlist Email: "))
	color.Unset() 
	var username string
    fmt.Scanln(&username)

	//Ask for Password
	color.Set(color.FgGreen)
	fmt.Print(("Bashlist Password: "))
	color.Unset() 
	pass, _ := gopass.GetPasswd()
	string_pass := string(pass)
	return &username,&string_pass
}


func save_secret(url *string, username *string, secret *string) {

	/* Saves secret in users credential manager*/

	c := &credentials.Credentials{
	    ServerURL: *url,
	    Username: *username,
	    Secret: *secret,
	}
	nativeStore.Add(c)
}

func retreive_secret(url string)(*string,*string) {

	/* Retreives secret from users credentials manager*/

	username, tok,err := nativeStore.Get(url)
	if err!=nil {
		fmt.Println("An error occured while retreiving your credentials.Please try again later.")
		os.Exit(1)
	}
	return &username,&tok
}	



func incorrect_auth_loop(){
	/* Infinite Loop that runs till users enters a wrong username/password combination*/
}

func get_token(){
	/*gets a token from server*/
}

func refresh_token(){
	/*refresh an expired token*/
}

func change_password(){
	/*password change handler*/
}


