package main

import (
	"fmt"
	//"os/exec"
	//"bytes"
	//"os"
	//"github.com/buger/jsonparser"
	"io/ioutil"

	"bytes"
	"os/exec"
	"strings"
)

func upload_helper(fields *[]byte, encBytes *[]byte, uurl string) int {

	codePath := *get_cwd()

	encfilepth := codePath + "/" + ".encontents.bls"
	content := *encBytes
	err := ioutil.WriteFile(encfilepth, content, 0777)
	if err != nil {
		fmt.Println("somewhere else")
	}
	pypath := codePath + "/" + ".uploader.py"
	k := writerPy(pypath, fields, encfilepth, uurl)
	if k != 1 {
		fmt.Println(k)
		unexpected_event()
	}
	cmd := exec.Command("python", pypath)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err = cmd.Run()
	if err != nil {
		fmt.Println("haha")
		fmt.Println(err)
		unexpected_event()
	}
	statusCode := strings.TrimSpace(outb.String())
	if statusCode != "204" {
		fmt.Println(statusCode)
		unexpected_event()
	}
	//err2 := os.Remove(pypath)
	//err1 := os.Remove(encfilepth)
	//if err2 != nil || err1 != nil {
	//	return 0
	//}
	return 1
}
