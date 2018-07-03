package main
import "os"
import (
	"fmt"
	"github.com/buger/jsonparser"
)

func writerPy(filname string, fields *[]byte,encfilePath string, uurl string)int{

	fieldVals := *fields

	acl, err := jsonparser.GetString(fieldVals, "acl")
	if err != nil {
		unexpected_event()
	}
	key, err := jsonparser.GetString(fieldVals, "key")
	if err != nil {
		unexpected_event()
	}
	xaa, err := jsonparser.GetString(fieldVals, "x-amz-algorithm")
	if err != nil {
		unexpected_event()
	}
	xac, err := jsonparser.GetString(fieldVals, "x-amz-credential")
	if err != nil {
		unexpected_event()
	}
	xad, err := jsonparser.GetString(fieldVals, "x-amz-date")
	if err != nil {
		unexpected_event()
	}
	policy, err := jsonparser.GetString(fieldVals, "policy")
	if err != nil {
		unexpected_event()
	}
	xas, err := jsonparser.GetString(fieldVals, "x-amz-signature")
	if err != nil {
		unexpected_event()
	}

	aclVal := acl + "\n"
	keyVal := key + "\n"
	xaaVal := xaa + "\n"
	xacVal := xac + "\n"
	xadVal := xad + "\n"
	policyVal := policy + "\n"
	xasVal := xas + "\n"

	urlVal := uurl


	file, err := os.Create(filname)
	if err != nil {
		fmt.Println("hh")
		return 0
	}
	defer file.Close()

	fmt.Fprintf(file, "import sys\n")
	fmt.Fprintf(file, "import os\n")
	fmt.Fprintf(file, "import requests\n")
	fmt.Fprintf(file, "if sys.version_info[0]<3\n")
	fmt.Fprintf(file, "    import ConfigParser\n")
	fmt.Fprintf(file, "else:\n")
	fmt.Fprintf(file, "    import configparser as ConfigParser\n")
	fmt.Fprintf(file, "from collections import OrderedDict\n")
	fmt.Fprintf(file, "config = ConfigParser.ConfigParser()\n")
	fmt.Fprintf(file, "inputDict = OrderedDict()\n")
	fmt.Fprintf(file, "inputDict['acl'] = %s",aclVal)
	fmt.Fprintf(file, "inputDict['key'] = %s",keyVal)
	fmt.Fprintf(file, "inputDict['x-amz-algorithm']= %s",xaaVal)
	fmt.Fprintf(file, "inputDict['x-amz-credential']=%s",xacVal)
	fmt.Fprintf(file, "inputDict['x-amz-date']=%s",xadVal)
	fmt.Fprintf(file, "inputDict['policy']=%s",policyVal)
	fmt.Fprintf(file, "inputDict['x-amz-signature']=%s",xasVal)
	fmt.Fprintf(file, "files = {'file':open(%s,'rb')}\n",encfilePath)
	fmt.Fprintf(file, "resp = requests.post(url=%s,data=inputDict,files=files\n)",urlVal)
	fmt.Fprintf(file, "print(resp.status_code)\n")

	return 1

}
