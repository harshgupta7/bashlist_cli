package main

import (
	"bytes"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/fatih/color"
	"github.com/imroc/req"
	"github.com/pierrre/archivefile/zip"
	"time"
	"os"
)

func get_download_url(bucketname string, usernamePtr *string, passwordPtr *string)[]byte{

	/* Gets download url for bucket*/

	endpoint := URL + PULL_BUCKET_ENDPOINT + bucketname

	r := req.New()
	authHeader := req.Header{
		"Email":        *usernamePtr,
		"Password": *passwordPtr,
	}
	r.SetTimeout(25 * time.Second)
	c, err := r.Get(endpoint, authHeader)
	if err!=nil{
		c, err = r.Get(endpoint, authHeader)
		if err!=nil {
			color.Red("Error contacting Server. Please check you connection and try again")
			os.Exit(1)
		}
	}

	respCode := c.Response().StatusCode
	//auth error
	if respCode==403{
		usernamePtr, passwordPtr,_ = authHandler(1)
		authHeader = req.Header{
			"Email":        *usernamePtr,
			"Password": *passwordPtr,
		}
		c, err = r.Get(endpoint, authHeader)
		respCode = c.Response().StatusCode
		if respCode==403{
			color.Cyan("Bashlist could not authenticate you. Please try again.")
			os.Exit(1)
		}
	//no bucket available
	} else if respCode==284 {
		cyan := color.New(color.FgCyan).SprintFunc()
		d := color.New(color.FgGreen, color.Bold).SprintFunc()
		fmt.Printf("%s %s.\n", cyan(": Does not exist. View your available directories using"), d("bashls"))
		os.Exit(1)
	//found bucket
	} else if respCode==285{
		byteResp, err := c.ToBytes()
		if err!=nil{
			unexpected_event()
		}
		return byteResp
	//unexpected
	} else{
		unexpected_event()
	}
	return nil
}

func download_bucket_to_bytes(url string,done chan *[]byte){
	/* Downloads bucket from s3 to bytes*/

	r := req.New()
	r.SetTimeout(100 * time.Second)
	resp,err := r.Get(url)
	if err!=nil{
		resp,err = r.Get(url)
		if err!=nil{
			done<-nil
		}
	}

	if resp.Response().StatusCode!=200{
		done<-nil
		close(done)
		return
	}
	//convert response to bytes
	res,err := resp.ToBytes()
	if err!=nil{
		done<-nil
	}
	//ping channel
	done<-&res
	//close channel
	close(done)
}

func decrypt_private_file_key(enc_key *string,pass *string)(*[]byte,error){

	//Decrypts bucket encryption key using password

	password := *pass
	fileKey,err := decrypt_secret(enc_key,password)
	if err!=nil{
		return nil,err
	}
	fileKeyval := *fileKey
	fileKeyBytes := []byte(fileKeyval)
	return &fileKeyBytes,nil
}




func decrypt_shared_file_key(encrypted_private_key *string,enc_file_key *string, pass *string)(*[]byte,error){

	//decrypts shared encryption key using private key and password

	privKeyVal,err := decrypt_secret(encrypted_private_key,*pass)
	if err!=nil{
		return nil,err
	}

	privKey,err := ParseRsaPrivateKeyFromPemStr(*privKeyVal)
	if err!=nil{
		return nil,err
	}
	encKeyStr := *enc_file_key
	encKeyByte := []byte(encKeyStr)
	keyVal,err := DecryptWithPrivKey(privKey,&encKeyByte)
	if err!=nil{
		return nil,err
	}
	return keyVal,nil

}

func DecryptContents(enc_contents *[]byte, keyval *[]byte)(*[]byte,error){

	//Decrypt downloaded contents using file key

	contents,err := DecryptObject(enc_contents,keyval)
	if err!=nil{
		return nil,err
	}
	return contents,nil

}

func unzipContents(zippedContents *[]byte,outpath string){

	//Unzip contents

	donesig := color.New(color.FgGreen).SprintFunc()
	progress := func(archivePath string) {
		fmt.Printf("Receiving: %s....%s\n", archivePath, donesig("OK"))
	}
	contents := *zippedContents
	c := int64(len(contents))
	r := bytes.NewReader(contents)
	zip.Unarchive(r,c,outpath,progress)
}


func download_manager(bucketname string) {

	/* Download Manager*/
	usernamePtr, passwordPtr,pass := authHandler(0)

	resp := get_download_url(bucketname, usernamePtr,passwordPtr)
	//should never happen
	if resp == nil {
		unexpected_event()
	}
	//retreive the url value from response
	url, err := jsonparser.GetString(resp, "url")
	if err != nil {
		unexpected_event()
	}

	//channel for downloaded contents
	downloadDone := make(chan *[]byte)
	//start downloading
	go download_bucket_to_bytes(url, downloadDone)

	//check if directory already exists in the current location
	exists := directory_exists(bucketname, "pull")
	if exists == true {
		color.Set(color.FgCyan)
		fmt.Print("A directory " + bucketname + " already exists. Pulling will overwrite its contents." +
			"Do you want to continue?[Y/n] ")
		color.Unset()
		//color.Cyan("Do you want to continue?[Y/n] ")
		var response string
		fmt.Scan(&response)
		if response == "n" || response == "N" || response == "No" || response == "no" {
			return
		}
	}
	//check if the downloaded directory is shared
	sharedVal, err := jsonparser.GetString(resp, "shared")
	if err != nil {
		unexpected_event()
	}

	//create filekey slice
	var fileKey *[]byte
	//Directory is private
	if sharedVal == "False" {
		decKey, err := jsonparser.GetString(resp, "key")
		if err != nil {
			unexpected_event()
		}
		fileKey, err = decrypt_private_file_key(&decKey,pass)
		if err != nil {
			unexpected_event()
		}
	//Directory is shared
	} else {
		privKey, err := jsonparser.GetString(resp, "unlock_key")
		if err != nil {
			unexpected_event()
		}
		decKey, err := jsonparser.GetString(resp, "key")
		if err != nil {
			unexpected_event()
		}
		fileKey, err = decrypt_shared_file_key(&privKey, &decKey,pass)
		if err != nil {
			unexpected_event()
		}
	}

	cwd := get_cwd()
	outPath := *cwd

	encContents := <-downloadDone

	//Unencrypt data
	contents, err := DecryptContents(encContents, fileKey)
	if err != nil {
		unexpected_event()
	}

	//Unzip and save data
	unzipContents(contents, outPath)
	return
}


