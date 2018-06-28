package main

import (
	"github.com/imroc/req"
	"fmt"
	"encoding/json"
	"github.com/fatih/color"
	"os"
	"strconv"
	"github.com/buger/jsonparser"
	"bufio"


)

type PushURLRequester struct{
	Size string
	Name string
}

type PushConfirmer struct{
	Name string
	Shared string
	Description string
}

func get_upload_url(){
	/* Gets upload url*/
}

func encrypt_bucket(){
	/* Encrypts the bucket*/
}


func upload_handler(dirname string) {
	/* Method to manage upload*/


	//Fetch username and password
	usernamePtr, passwordPtr,pass:=authHandler(0)

	//Allocate byte array for compressed directory
	var comp_bytes *[]byte

	//Endpoint to get upload URL
	endpoint := URL + PUSH_BUCKET_REQ

	//Check if directory Exists
	ex := directory_exists(dirname)
	//No Directory Exists Return
	if ex==false{
		return
	}

	//Initiate operation
	color.Cyan("Initiating Push:")
	fmt.Println("  - "+ dirname)

	//Channel to receive compressed directory
	conf_comp := make(chan *[]byte)

	//Initiate goroutine
	go dir_to_compressed_bytes(dirname,conf_comp)

	//Get Size
	size,err := get_size(dirname)
	size = size/1000 //Convert from bytes to KB
	if err!=nil{
		color.Red("Bashlist encountered an unexpected error. Please try again later.")
		os.Exit(1)
	}

	//Retrieve Username & Password
	username := *usernamePtr
	hashedPassword := *passwordPtr

	//Header for getting push URL
	header := req.Header{
		"Content-Type":"application/json",
		"Email": username,
		"Password": hashedPassword,
	}

	//JSON Representation of directory name and size
	vals := PushURLRequester{Name:dirname,Size:strconv.Itoa(int(size))}
	jsonvals,_ := json.Marshal(vals)

	//Perform Post and receive URL
	resp, err := req.Post(endpoint, req.BodyJSON(jsonvals),header)
	//Error Performing Post
	if err != nil {
		//Do it once More
		resp, err = req.Post(endpoint, req.BodyJSON(jsonvals),header)
	}

	//Second Time Error
	if err!=nil{
		<-conf_comp
		color.Red("Could not connect to server. Please check your connection and try again later.")
		return
	}

	//Get Response Code
	respCode := resp.Response().StatusCode

	//Authentication Error
	if respCode==403 {
		//Finish Compression
		comp_bytes = <-conf_comp
		//Request Fresh Username & Password From AuthHandler
		usernamePtr, passwordPtr, pass = authHandler(1)

		//Retrieve u/p values
		username = *usernamePtr
		hashedPassword = *passwordPtr
		header = req.Header{
			"Content-Type":"application/json",
			"Email":    username,
			"Password": hashedPassword,
		}
		resp, err1 := req.Post(endpoint, header, req.BodyJSON(jsonvals))
		if err1!=nil{
			unexpected_event()
		}
		respCode = resp.Response().StatusCode
		//Technically should never happen,as U/P values are saved after auth_check with server.
		if respCode == 403 {
			color.Red("Authentication Error!. try again later.")
			return
		}
	//Insufficient Space Error
	} else if respCode==423 {
		<-conf_comp
		color.Red("Insufficient remaining space to add %s to your bashlist.", dirname)
		return
	//Insufficient Space - Shared Directory
	} else if respCode==424{
		<-conf_comp
		color.Red("%s is a shared directory. The owner has insufficient remaining space for this push",dirname)
		return
	//Other Server Error
	} else if respCode==399 {
		<-conf_comp
		color.Red("An unexpected error occurred while pushing %s. Please try again later", dirname)
		return
	//Received a valid response from server
	} else{

		//Set Description to Null Value
		var desc string = "~N////V~"

		//Get Byte Resp
		byteResp, _ := resp.ToBytes()

		//Check if the directory being pushed already exists. Need it to show confirmation messages.
		exists, err := jsonparser.GetString(byteResp, "Exist")
		if err!=nil{
			//If exist variable isn't there. Unexpected response
			<-conf_comp
			color.Red("An unexpected error occurred while pushing %s. Please try again later", dirname)
			return
		}


		//Get whether directory is shared or not
		shared,err := jsonparser.GetString(byteResp,"Shared")
		//Allocate byte array for file key
		var file_key *[]byte


		if shared=="Y" {

			//Get filekey encrypted with private key
			key, errkey := jsonparser.GetString(byteResp, "keyval")
			//Get private key encrypted with password
			enc_privKey, errprivkey := jsonparser.GetString(byteResp, "PrivKey")
			//If either of those values are not present, return
			if errkey != nil || errprivkey != nil {
				<-conf_comp
				color.Red("An unexpected error occurred while pushing %s. Please try again later", dirname)
				return
			}
			//Decrypt Private key
			privKey, _ := decrypt_secret(&enc_privKey, *pass)
			//Build Private Key object
			privKeyObj, err := ParseRsaPrivateKeyFromPemStr(*privKey)
			//If error Parsing privatekey, return
			if err != nil {
				<-conf_comp
				color.Red("An unexpected error occurred while pushing %s. Please try again later", dirname)
				return
			}
			//Convert key to byte array
			byte_key := []byte(key)
			//Decrypt with private key and get file key
			file_key, _ = DecryptWithPrivKey(privKeyObj, &byte_key)

		} else if shared=="N"{

			//Get filedecryption key encrypted by password
			key, errkey := jsonparser.GetString(byteResp, "Key")
			//If error, return
			if errkey!=nil{
				<-conf_comp
				color.Red("An unexpected error occurred while pushing %s. Please try again later", dirname)
				return
			}
			//Retreive filekey by decrypting from password
			file_key_string_ptr, _ := decrypt_secret(&key, *pass)
			file_key_str := *file_key_string_ptr
			file_key_byte := []byte(file_key_str)
			file_key = &file_key_byte
		}
		//make sure compression is done, to start encryption
		comp_bytes = <-conf_comp
		//make channel for encryption
		conf_encryption := make(chan *[]byte)
		//Start encryption
		go EncryptObject(comp_bytes,file_key,conf_encryption)

		//Some lines to seperate from processing stdout
		fmt.Println()
		fmt.Println()

		//Description Messages Printer
		if exists=="Y"{
			if shared=="N" {
				fmt.Print("A Directory " + dirname + " already exists in your bashlist. Pushing will update its contents. " +
					"Do you want to continue?[Y/n] ")
				var response string
				fmt.Scanln(&response)
				if response == "n" || response == "N" || response == "No" || response == "no" {
					return
				}
			}else{
				fmt.Print("Directory " + dirname + " already exists and is a shared directory. Pushing will update its contents for all members. " +
					"Do you want to continue?[Y/n] ")
				var response string
				fmt.Scanln(&response)
				if response == "n" || response == "N" || response == "No" || response == "no" {
					return
				}
			}
		} else {
			fmt.Print("Description (Press Enter to Leave Blank): ")
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				desc = scanner.Text()
				break
			}
		}
		fmt.Println("Encrypting contents")
		encrypted_bytes :=<-conf_encryption
		if encrypted_bytes==nil{
			color.Red("An unexpected error occurred while pushing %s. Please try again later", dirname)
			return
		}

		d,_,_,_ := jsonparser.Get(byteResp,"URL","fields")
	
		uurl ,_:=jsonparser.GetString(byteResp,"URL","url")
		upload_helper(&d,encrypted_bytes,uurl)


		//send confirmation
		header = req.Header{
			"Content-Type":"application/json",
			"Email": username,
			"Password": hashedPassword,
		}

		valsConfirmer := PushConfirmer{Name:dirname,Shared:shared,Description:desc}
		jsonvals,_ = json.Marshal(valsConfirmer)
		endpoint = URL + PUSH_BUCKET_CONF


		resp, err = req.Post(endpoint, req.BodyJSON(jsonvals),header)
		respCode = resp.Response().StatusCode
		if err!=nil{
			resp, err = req.Post(endpoint, req.BodyJSON(jsonvals),header)
			if err!=nil{
				unexpected_event()
			}
		}
		if respCode!=225{
			fmt.Println(respCode)
			resp, err = req.Post(endpoint, req.BodyJSON(jsonvals),header)
			if err!=nil{
				resp, err = req.Post(endpoint, req.BodyJSON(jsonvals),header)
				if err!=nil{
					unexpected_event()
				}
			}

		} else {
			fmt.Println(dirname+" uploaded successfully.")

		}

	}


}

