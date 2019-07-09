package main

import (
	"fmt"
	"github.com/buger/jsonparser"
	"os"
)

func writerPy(filname string, fields *[]byte, encfilePath string, uurl string) int {

	fieldVals := *fields

	aclVal, err := jsonparser.GetString(fieldVals, "acl")
	if err != nil {
		unexpected_event()
	}
	keyVal, err := jsonparser.GetString(fieldVals, "key")
	if err != nil {
		unexpected_event()
	}
	xaaVal, err := jsonparser.GetString(fieldVals, "x-amz-algorithm")
	if err != nil {
		unexpected_event()
	}
	xacVal, err := jsonparser.GetString(fieldVals, "x-amz-credential")
	if err != nil {
		unexpected_event()
	}
	xadVal, err := jsonparser.GetString(fieldVals, "x-amz-date")
	if err != nil {
		unexpected_event()
	}
	policyVal, err := jsonparser.GetString(fieldVals, "policy")
	if err != nil {
		unexpected_event()
	}
	xasVal, err := jsonparser.GetString(fieldVals, "x-amz-signature")
	if err != nil {
		unexpected_event()
	}
	file, err := os.Create(filname)
	if err != nil {
		fmt.Println("hh")
		return 0
	}
	defer file.Close()

	fmt.Fprintf(file, "import sys\n")
	fmt.Fprintf(file, "import os\n")
	fmt.Fprintf(file, "import requests\n")
	fmt.Fprintf(file, "if sys.version_info[0]<3:\n")
	fmt.Fprintf(file, "    import ConfigParser\n")
	fmt.Fprintf(file, "else:\n")
	fmt.Fprintf(file, "    import configparser as ConfigParser\n")
	fmt.Fprintf(file, "from collections import OrderedDict\n")
	fmt.Fprintf(file, "config = ConfigParser.ConfigParser()\n")
	fmt.Fprintf(file, "inputDict = OrderedDict()\n")
	fmt.Fprintf(file, "inputDict['acl'] = '%s'\n", aclVal)
	fmt.Fprintf(file, "inputDict['key'] = '%s'\n", keyVal)
	fmt.Fprintf(file, "inputDict['x-amz-algorithm']= '%s'\n", xaaVal)
	fmt.Fprintf(file, "inputDict['x-amz-credential']='%s'\n", xacVal)
	fmt.Fprintf(file, "inputDict['x-amz-date']='%s'\n", xadVal)
	fmt.Fprintf(file, "inputDict['policy']='%s'\n", policyVal)
	fmt.Fprintf(file, "inputDict['x-amz-signature']='%s'\n", xasVal)
	fmt.Fprintf(file, "files = {'file':open('%s','rb')}\n", encfilePath)
	fmt.Fprintf(file, "resp = requests.post(url='%s',data=inputDict,files=files)\n", uurl)
	fmt.Fprintf(file, "print(resp.status_code)\n")
	return 1

}
