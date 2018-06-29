package main

import (
	"github.com/imroc/req"
	"time"
	"github.com/fatih/color"
"github.com/pierrre/archivefile/zip"
	"fmt"
	"bytes"
	"github.com/buger/jsonparser"
)

func get_download_url(bucketname string)[]byte{
	/* Gets download url for bucket*/

	endpoint := URL + PULL_BUCKET_ENDPOINT+bucketname
	usernamePtr, passwordPtr,_ := authHandler(0)

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
		}
	}
	d := c.Response().StatusCode
	if d==403{
		usernamePtr, passwordPtr,_ = authHandler(1)
		authHeader = req.Header{
			"Email":        *usernamePtr,
			"Password": *passwordPtr,
		}
		c, err = r.Get(endpoint, authHeader)
		d = c.Response().StatusCode
		if d==403{
			color.Cyan("Bashlist could not authenticate you at the moment. Please try again.")
			return nil
		}

	} else if d==284 {
		msg := bucketname + ": Does not exist. View your available directories using bashls."
		color.Cyan(msg)
		return nil
	} else if d==285{
		byteResp, err := c.ToBytes()
		if err!=nil{
			unexpected_event()
		}
		return byteResp
	} else{
		unexpected_event()
	}
	return nil
}

func download_bucket_to_bytes(url string,done chan *[]byte){
	/* Downloads bucket*/
	r := req.New()
	r.SetTimeout(100 * time.Second)
	resp,err := r.Get(url)
	if err!=nil{
		resp,err = r.Get(url)
		if err!=nil{
			done<-nil
		}
	}
	res,err := resp.ToBytes()
	if err!=nil{
		done<-nil
	}
	done<-&res
	close(done)
}


func decrypt_private_file_key(enc_key *string)(*[]byte,error){
	_,_,passwordPtr := authHandler(0)
	password := *passwordPtr
	fileKey,err := decrypt_secret(enc_key,password)
	if err!=nil{
		return nil,err
	}
	fileKeyval := *fileKey
	fileKeyBytes := []byte(fileKeyval)
	return &fileKeyBytes,nil
}




func decrypt_shared_file_key(encrypted_private_key *string,enc_file_key *string)(*[]byte,error){

	_,_,pass := authHandler(0)

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

func UnencryptContents(enc_contents *[]byte, keyval *[]byte)(*[]byte,error){
	contents,err := DecryptObject(enc_contents,keyval)
	if err!=nil{
		return nil,err
	}
	return contents,nil

}

func unzipContents(zippedContents *[]byte,outpath string){
	donesig := color.New(color.FgGreen).SprintFunc()
	progress := func(archivePath string) {
		fmt.Printf("Receiving: %s....%s\n", archivePath, donesig("OK"))
	}
	contents := *zippedContents
	c := int64(len(contents))
	r := bytes.NewReader(contents)
	zip.Unarchive(r,c,outpath,progress)
}


func download_manager(bucketname string){
	/* Download Manager*/

	resp:= get_download_url(bucketname)
	url,err := jsonparser.GetString(resp,"url")
	if err!=nil{
		unexpected_event()
	}
	downloadDone := make(chan *[]byte)
	go download_bucket_to_bytes(url,downloadDone)
	exists := directory_exists(bucketname,"pull")
	if exists==true{
		msg := "A directory " + bucketname + " already exists. Pulling will overwrite its contents." +
			"Do you want to continue?[Y/n] "
		color.Cyan(msg)
		//color.Cyan("Do you want to continue?[Y/n] ")
		var response string
		fmt.Scanln(&response)
		if response == "n" || response == "N" || response == "No" || response == "no" {
			return
		}
	}
	sharedVal,err := jsonparser.GetString(resp,"shared")
	if err!=nil{
		unexpected_event()
	}
	var fileKey *[]byte
	if sharedVal=="False"{
		decKey,err := jsonparser.GetString(resp,"key")
		if err!=nil{
			unexpected_event()
		}
		fileKey,err = decrypt_private_file_key(&decKey)
		if err!=nil{
			unexpected_event()
		}
	} else{
		privKey,err := jsonparser.GetString(resp,"unlock_key")
		if err!=nil{
			unexpected_event()
		}
		decKey, err := jsonparser.GetString(resp,"key")
		if err!=nil{
			unexpected_event()
		}
		fileKey,err = decrypt_shared_file_key(&privKey,&decKey)
		if err!=nil{
			unexpected_event()
		}
	}
	cwd := get_cwd()
	cwdStr := *cwd
	outpath := cwdStr

	encContents := <-downloadDone
	contents,err := UnencryptContents(encContents,fileKey)
	if err!=nil{
		unexpected_event()
	}
	unzipContents(contents,outpath)

}


