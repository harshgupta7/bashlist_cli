package main

import (
	"github.com/imroc/req"
	"time"
	"github.com/buger/jsonparser"
)

func get_existing_creds()(*string,*string){


	endpoint := URL + PULL_BUCKET_ENDPOINT
	usernamePtr, passwordPtr,_ := authHandler(0)

	r := req.New()
	authHeader := req.Header{
		"Email":        *usernamePtr,
		"Password": *passwordPtr,
	}
	r.SetTimeout(25 * time.Second)
	c, err := r.Get(endpoint, authHeader)
	if err!=nil{
		c,err = r.Get(endpoint,authHeader)
		if err!=nil{
			return nil,nil
		}
	}
	respCode := c.Response().StatusCode
	if respCode==403{
		usernamePtr, passwordPtr,_ := authHandler(1)
		r = req.New()
		authHeader = req.Header{
			"Email":        *usernamePtr,
			"Password": *passwordPtr,
		}
		r.SetTimeout(25 * time.Second)
		c, err = r.Get(endpoint, authHeader)
		if err!=nil{
			unexpected_event()
		}
		if c.Response().StatusCode==403 {
			unexpected_event()
		}
	}

	if c.Response().StatusCode==399{
		unexpected_event()
	} else if c.Response().StatusCode==398{
		byteVal,err := c.ToBytes()
		if err!=nil{
			unexpected_event()
		}
		fileKey,err := jsonparser.GetString(byteVal,"fKey")
		privKey,err1 := jsonparser.GetString(byteVal,"pKey")
		if err!=nil || err1!=nil{
			unexpected_event()
		}
		return &fileKey,&privKey
	} else{
		unexpected_event()
	}
	return nil,nil
}

func reencrypt(fileKey *string, privKey *string, oldPassword *string, newPassword *string)(*string,*string){

	realFK,err := decrypt_secret(fileKey,*oldPassword)
	if err!=nil{
		unexpected_event()
	}
	realPrivKey,err := decrypt_secret(fileKey,*oldPassword)
	if err!=nil{
		unexpected_event()
	}

	encFK,err := encrypt_secret(realFK,*newPassword)
	if err!=nil{
		unexpected_event()
	}
	encPK,err := encrypt_secret(realPrivKey, *newPassword)
	if err!=nil{
		unexpected_event()
	}
	return encFK,encPK

}

func postUpdated(){

}

func changePassManager(){

}

