package main

import (
	"github.com/imroc/req"
	"time"
	"fmt"
	"github.com/buger/jsonparser"
	"encoding/json"
	"github.com/fatih/color"
	"os"
	"bufio"
)

type PushURLRequester struct{
	Size int
	Name string
}

func get_upload_url(){
	/* Gets upload url*/
}

func encrypt_bucket(){
	/* Encrypts the bucket*/
}

func upload_bucket(){
	/*Uploads the bucket*/
}



func upload_handler(dirname string) {
	/* Method to manage upload*/

	usernamePtr, passwordPtr,err:=authHandler(0)
	var comp_bytes *[]byte
	endpoint := URL + PUSH_BUCKET_REQ
	ex := directory_exists(dirname)
	if ex ==false{
		return
	}
	color.Cyan("Initiating Push:")
	fmt.Println("  - "+ dirname)
	conf_comp := make(chan *[]byte)
	go dir_to_compressed_bytes(dirname,conf_comp)

	size,err := get_size(dirname)
	if err!=nil{
		color.Red("Bashlist encountered an unexpected error. Please try again later.")
		os.Exit(1)
	}

	username := *usernamePtr
	hashedPassword := *passwordPtr
	r := req.New()
	authHeader := req.Header{
		"Email":        username,
		"Password": hashedPassword,
	}
	//param := req.Param{
	//	"Content-Type":  "application/json",
	//}

	values := PushURLRequester{Name:dirname,Size:int(size)}
	jsonvalues,_ := json.Marshal(values)
	r.SetTimeout(25 * time.Second)
	c, err := r.Post(endpoint, authHeader,jsonvalues)

	if err!=nil{
		<-conf_comp
		color.Red("Could not connect to server. Please check your connection and try again later.")
		return
	}
	d := c.Response().StatusCode
	if d==403 {
		comp_bytes = <-conf_comp
		usernamePtr, passwordPtr, _ = authHandler(1)
		username = *usernamePtr
		hashedPassword = *passwordPtr
		authHeader = req.Header{
			"Email":    username,
			"Password": hashedPassword,
		}
		c, err = r.Post(endpoint, authHeader, param, jsonvalues)
		d = c.Response().StatusCode
		if d == 403 {
			color.Red("Authentication Error!. try again later.")
			return
		}
	} else if d==423 {
		<-conf_comp
		color.Red("Insufficient remaining space to add %s ", dirname)
		return
	} else if d==424{
		<-conf_comp
		color.Red("%s is a shared directory. The owner has insufficient remaining space for this push",dirname)
		return
	} else if d==399 {
		<-conf_comp
		color.Red("An unexpected error occurred while pushing %s. Please try again later", dirname)
		return
	} else{
		var desc string = "NU"
		byteResp, _ := c.ToBytes()
		exists, err := jsonparser.GetString(byteResp, "Exist")
		if err!=nil{
			<-conf_comp
			color.Red("An unexpected error occurred while pushing %s. Please try again later", dirname)
			return
		}
		_, pass, err := retreive_secret("Bashlist-Credentials/Safe-Credentials")
		if err != nil {
			<-conf_comp
			color.Red("An unexpected error occurred while pushing %s. Please try again later", dirname)
			return
		}

		shared,err := jsonparser.GetString(byteResp,"Shared")
		var file_key *[]byte
		if shared=="Y" {

			key, errkey := jsonparser.GetString(byteResp, "keyval")
			enc_privKey, errprivkey := jsonparser.GetString(byteResp, "PrivKey")
			if errkey != nil || errprivkey != nil {
				<-conf_comp
				color.Red("An unexpected error occurred while pushing %s. Please try again later", dirname)
				return
			}

			privKey, _ := decrypt_secret(&enc_privKey, *pass)
			privKeyObj, err := ParseRsaPrivateKeyFromPemStr(*privKey)
			if err != nil {
				<-conf_comp
				color.Red("An unexpected error occurred while pushing %s. Please try again later", dirname)
				return
			}
			byte_key := []byte(key)
			file_key, _ = DecryptWithPrivKey(privKeyObj, &byte_key)
		} else if shared=="N"{
			key, errkey := jsonparser.GetString(byteResp, "Key")
			if errkey!=nil{
				<-conf_comp
				color.Red("An unexpected error occurred while pushing %s. Please try again later", dirname)
				return
			}

			file_key_string_ptr, _ := decrypt_secret(&key, *pass)
			file_key_str := *file_key_string_ptr
			file_key_byte := []byte(file_key_str)
			file_key = &file_key_byte
		}

		comp_bytes = <-conf_comp
		conf_encryption := make(chan *[]byte)
		go EncryptObject(comp_bytes,file_key,conf_encryption)

		if exists=="Y"{
			fmt.Println()
			fmt.Print("A Directory "+dirname+" already exists in your bashlist. Pushing will update its contents. " +
				"Do you want to continue?[Y/n]")
			var response string
			fmt.Scanln(&response)
			if response=="n"||response=="N"||response=="No"||response=="no"{
				return
			}
		} else {
			fmt.Print("Description (Press Enter to Leave Blank): ")
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				desc = scanner.Text()
				break
			}
		}

		encrypted_bytes :=<-conf_encryption
		if encrypted_bytes!=nil{
			color.Red("An unexpected error occurred while pushing %s. Please try again later", dirname)
			return
		}
		fmt.Print(desc)


	}


}

