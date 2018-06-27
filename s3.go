package main
//
//func run(encrypted_arr *[]byte) {
//	c := get_cwd()
//	d := *c
//	content := *encrypted_arr
//	tmpfile, err := ioutil.TempFile(d, "bltempfile")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	defer os.Remove(tmpfile.Name()) // clean up
//
//	if _, err := tmpfile.Write(content); err != nil {
//		log.Fatal(err)
//		fmt.Println("adsad")
//	}
//
//
//
//
//	if err := tmpfile.Close(); err != nil {
//		log.Fatal(err)
//		fmt.Println("hahah")
//	}
//
//}
//
//func four_upload(fields *[]byte){
//	x := make(map[string]string)
//	x["key"],_ = jsonparser.GetString(*fields,"key")
//	x["acl"],_ = jsonparser.GetString(*fields,"acl")
//	x["x-amz-algorithm"],_ = jsonparser.GetString(*fields,"x-amz-algorithm")
//	x["x-amz-credential"],_ = jsonparser.GetString(*fields,"x-amz-credential")
//	x["x-amz-date"],_ = jsonparser.GetString(*fields,"x-amz-date")
//	x["policy"],_ = jsonparser.GetString(*fields,"policy")
//	x["x-amz-signature"],_ = jsonparser.GetString(*fields,"x-amz-signature")
//	file, _ := os.Open("c.py")
//
//	req.Post("http://bashlist-78.s3.amazonaws.com", req.FileUpload{
//		//"Key" : x["key"],
//
//		File:      file,
//		FieldName: "file",       // FieldName is form field name
//		FileName:  "avatar.png", //Filename is the name of the file that you wish to upload. We use this to guess the mimetype as well as pass it onto the server
//	})
//
//}
//
//func three_upload(fields *[]byte){
//	c := *fields
//	fmt.Println(string(c))
//	f, _ := filepath.Abs("config.go")
//	bytesOfFile, _ := ioutil.ReadFile(f)
//	resp,body,err:=gorequest.New().Post("http://bashlist-78.s3.amazonaws.com").
//		Type("multipart").
//		Send(c).
//		SendFile(bytesOfFile).
//		End(printStatus)
//	fmt.Println(resp)
//	fmt.Println(body)
//	fmt.Println(err)
//
//}
//func printStatus(resp gorequest.Response, body string, errs []error){
//	fmt.Println(resp.Status)
//}
//func two_upload(fields *[]byte){
//
//	x := make(map[string]string)
//	x["key"],_ = jsonparser.GetString(*fields,"key")
//	x["acl"],_ = jsonparser.GetString(*fields,"acl")
//	x["x-amz-algorithm"],_ = jsonparser.GetString(*fields,"x-amz-algorithm")
//	x["x-amz-credential"],_ = jsonparser.GetString(*fields,"x-amz-credential")
//	x["x-amz-date"],_ = jsonparser.GetString(*fields,"x-amz-date")
//	x["policy"],_ = jsonparser.GetString(*fields,"policy")
//	x["x-amz-signature"],_ = jsonparser.GetString(*fields,"x-amz-signature")
//	fmt.Println(x)
//	c := new(http.Client)
//	req := request.NewRequest(c)
//	f, err := os.Open("filewriter.go")
//	req.Files = []request.FileField{
//		request.FileField{"file", "filewriter.go", f},
//	}
//	//req.Data = x
//	//f, err := os.Open("filewriter.go")
//	if err!=nil{
//		fmt.Println("fucked")
//	}
//
//	req.Files = []request.FileField{
//		request.FileField{"file", "filewriter.go", f},
//	}
//	//x["key"],_ = jsonparser.GetString(*fields,"key")
//	req.Data = x
//	resp, err := req.Post("http://bashlist-78.s3.amazonaws.com")
//	fmt.Println((resp.Text()))
//
//
//
//}
//
//
//
//
//
//func aws_fields_to_map(fields *[]byte)(map[string]string){
//
//	fieldStruct := *fields
//
//	x := make(map[string]string)
//	        // Just to demonstrate how this works.
//	//x["monkey"] = "George"
//	err := json.Unmarshal([]byte(fieldStruct), &x)
//	if err != nil {
//	    fmt.Printf("Error: %s\n", err)
//	    // return
//	}
//	// for key, value := range x {
//	//    fmt.Printf("%s -> %s\n", key, value)
//	// }
//	fmt.Println(x)
//	return x
//}
