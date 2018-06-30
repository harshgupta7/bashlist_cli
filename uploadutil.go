package main 

import (
    "fmt"
    "os/exec"
    "bytes"
	"os"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"strings"
)


func upload_helper(fields *[]byte, encBytes *[]byte,uurl string)int{

	cwdPtr := get_cwd()
	cwd := *cwdPtr
	filepth := cwd+"/"+".encontents.bls"
	content := *encBytes
	err := ioutil.WriteFile(filepth,content,0777)
	if err!=nil{
		fmt.Println("somewhere else")
	}

	fieldVals := *fields

	acl,err := jsonparser.GetString(fieldVals,"acl")
	if err!=nil{
		unexpected_event()
	}
	key,err := jsonparser.GetString(fieldVals,"key")
	if err!=nil{
		unexpected_event()
	}
	xaa,err := jsonparser.GetString(fieldVals,"x-amz-algorithm")
	if err!=nil{
		unexpected_event()
	}
	xac,err := jsonparser.GetString(fieldVals,"x-amz-credential")
	if err!=nil{
		unexpected_event()
	}
	xad,err := jsonparser.GetString(fieldVals,"x-amz-date")
	if err!=nil{
		unexpected_event()
	}
	policy,err := jsonparser.GetString(fieldVals,"policy")
	if err!=nil{
		unexpected_event()
	}
	xas,err := jsonparser.GetString(fieldVals,"x-amz-signature")
	if err!=nil{
		unexpected_event()
	}

	aclVal := "acl="+acl+"\n"
	keyVal := "key="+key+"\n"
	xaaVal := "x-amz-algorithm="+xaa+"\n"
	xacVal := "x-amz-credential="+xac+"\n"
	xadVal := "x-amz-date="+xad+"\n"
	policyVal := "policy="+policy+"\n"
	xasVal := "x-amz-signature="+xas+"\n"

	urlVal := "url="+uurl+"\n"

	blconfigfilepath := cwd+"/"+".bashlistuploadconfig.txt"
	encfilePathVal := "encfile="+filepth+"\n"

	file,err := os.Create(blconfigfilepath)
	if err!=nil{
		//try in different directory
		//fmt.Println("hooha")

	}
	defer file.Close()
	fmt.Fprintf(file,"[default]\n")
	fmt.Fprintf(file,aclVal)
	fmt.Fprint(file,keyVal)
	fmt.Fprint(file,xaaVal)
	fmt.Fprint(file,xacVal)
	fmt.Fprint(file,xadVal)
	fmt.Fprint(file,policyVal)
	fmt.Fprint(file,xasVal)
	fmt.Fprint(file,encfilePathVal)
	fmt.Fprint(file,urlVal)
	pypath := get_code_path()+"/"+"uploadutil.py"
	cmd := exec.Command("python", pypath,string(blconfigfilepath))
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err = cmd.Run()
	if err != nil {
		unexpected_event()
	}
	statusCode := strings.TrimSpace(outb.String())
	if statusCode!="204"{
		unexpected_event()
	}
	err2:=os.Remove(blconfigfilepath)
	err1:=os.Remove(filepth)
	if err2 !=nil || err1 !=nil{
		return 0
	}
	return 1


}

