package main 

import "github.com/fatih/color"
import "fmt"
// import "io/ioutil"
// import "github.com/skratchdot/open-golang/open"
import "path/filepath"
import "os"
import "encoding/json"
import "github.com/howeyc/gopass"
import "github.com/imroc/req"
import "github.com/buger/jsonparser"

import "github.com/docker/docker-credential-helpers/credentials"
import "github.com/docker/docker-credential-helpers/osxkeychain"


//FOR MAC

var URL string = "http://127.0.0.1:5000"

var nativeStore = osxkeychain.Osxkeychain{}

type User struct{
	Email string
	Password string
}

func up_prompt() (string,string) {
	//Ask for Username
	color.Set(color.FgGreen)
	fmt.Print(("Bashlist Email: "))
	color.Unset() 
    var input string
    fmt.Scanln(&input)
	
	//Ask for Password
	color.Set(color.FgGreen)
	fmt.Print(("Bashlist Password: "))
	color.Unset() 
	pass, _ := gopass.GetPasswd()
	string_pass := string(pass)
	return input,string_pass

}


func get_account_url() string {

	endpoint:=URL+"/account"
	token:=authentication_handler(false)
	val := "JWT " + token

	header := req.Header{
		"Content-Type":"application/json",
		"Authorization": val,
	}
	r, err := req.Get(endpoint, header)
	if err != nil {
		fmt.Println("da")
	}
	c := r.String()
	return c	
}


func incorrect_auth_loop() string {

	for {
		u1,p1 := up_prompt()
		t1,s1 := get_token(u1,p1)
		if s1==true{
			save_secret("Bashlist-Credentials/pass",u1,p1)
			save_secret("Bashlist-Credentials/token",u1,t1)
			return t1
			break
		}
		//Invalid Credentials
		if t1=="XRQ23"{
			fmt.Println("Incorrect Username or Password")
		//Some other BS
		} else if t1=="SSD12"{
			fmt.Println("Couldn't establish connection. Please check your network and try again")
		} else {
			fmt.Println("Authentication Error. Please try again.")
		}
		fmt.Println("Forgot Password? run   bls account to reset your password")
	}
	return "NEVER-RETURNED"
}





func authentication_handler(exp bool) string{

	//token has expired, method is asking for a fresh token
	if exp{

		//get password from secure storage
		u,p:=retreive_secret("Bashlist-Credentials/pass")

		//password has been set before:this is not the first time password is being fetched
		if p!="X"{
			token,success := get_token(u,p)
			//Auth succeeded
			if success{
				save_secret("Bashlist-Credentials/token",u,token)
				return token
			} else {
			//Auth Failed
				res:= incorrect_auth_loop()
				return res
			}
		//Password isn't set. This should never happen, since method is asking token for the second time.
		} else {
			res := incorrect_auth_loop()
			return res
		}
	//method is asking for the token the first time
	} else {

		u,t:=retreive_secret("Bashlist-Credentials/token")
		//not the first time token will be called
		if t!="X"{
			return t

		// token is not set
		} else {
			u1,p1 := retreive_secret("Bashlist-Credentials/pass")
			//token is not set but password is set
			if p1!="X" {
				token,success := get_token(u1,p1)
				//Auth succeeded
				if success{
					save_secret("Bashlist-Credentials/token",u,token)
					return token
				} else {
				//Auth Failed
					res:= incorrect_auth_loop()
					return res
				}
			//token is not set, password is not set
			} else{
				color.Cyan("Welcome to Bashlist!")
				res:=incorrect_auth_loop()
				color.Green("Bashlist is Ready")
				return res
			}
		}
	}
}

func open_url(url string) {
	// opens url in browser
	_:=open.Run(url)

}

func upload_file(s string) {
	// """Uploads a file"""
	current_dir := get_current_dir()
	file_path := current_dir+"/"+s
	if objectExists(file_path) == false{
		fmt.Println(s+": No such file or directory")
		return
	}

	endpoint:=URL+"/filesync"
	token:=authentication_handler(false)
	val := "JWT " + token
	header := req.Header{
		"Content-Type":"application/json",
		"Authorization": val,
	}
	file, _ := os.Open(s)
	c,err_conn:=req.Post(endpoint, req.FileUpload{
		File:      file,
		FieldName: "file",       
		FileName:  s, 
	},header)
	if err_conn!=nil{
		fmt.Println("Error Connecting. Please Check Your Connection & Try Again.")
		return
	}
	byte_resp,_:=c.ToBytes()
	
	response,_:=jsonparser.GetString(byte_resp, "BLCODE")
	if response=="LMV23"{
		fmt.Println("Successfully uploaded "+ s)
		return
	}
	token=authentication_handler(true)
	val = "JWT " + token
	header = req.Header{
		"Content-Type":"application/json",
		"Authorization": val,
	}
	c,_=req.Post(endpoint, req.FileUpload{
		File:      file,
		FieldName: "file",       
		FileName:  s, 
	},header)
	byte_resp,_=c.ToBytes()
	response,_=jsonparser.GetString(byte_resp, "BLCODE")
	if response=="LMV23"{
		fmt.Println("Successfully uploaded "+ s)
		return
	}
	fmt.Println("Error Uploading File")
	return
}

func get_current_dir() string{
    dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
    return dir
}


func objectSize(path string) (int64, error) {
    var size int64
    err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
        if !info.IsDir() {
            size += info.Size()
        }
        return err
    })
    return size, err
}

func objectExists(name string) bool {
    if _, err := os.Stat(name); err != nil {
    if os.IsNotExist(err) {
            return false
        }
    }
    return true
}


func download_file(s string) {
	// """Downloads a file"""

	current_dir := get_current_dir()
	file_path := current_dir+"/"+s 
	ex:=objectExists(file_path)
	if ex==true{
		fmt.Print(s+" already Exists. Do you want to continue?[Y/n]")
		var response string
    	fmt.Scanln(&response)
    	if response=="n"||response=="N"||response=="No"||response=="no"{
    		return
    	}
	}

	endpoint:=URL+"/filedown/"+s
	token:=authentication_handler(false)
	val := "JWT " + token
	header := req.Header{
		"Content-Type":"application/json",
		"Authorization": val,
	}
	r, conn_error := req.Get(endpoint,header)
	if conn_error!=nil{
		fmt.Println("Error Connecting. Please check your connection and try again")
		return
	}

	byte_resp,_:=r.ToBytes()
	response,JSONerr:=jsonparser.GetString(byte_resp, "BLCODE")
	if JSONerr!=nil{
		r.ToFile(s)
		return
	}else if response=="NE235"{
		fmt.Println(s+": No such file exists in your Bashlist.")
		return
	}else if response=="INC23"{
		token=authentication_handler(true)
		val = "JWT " + token
		header = req.Header{
			"Content-Type":"application/json",
			"Authorization": val,
		}
		r, conn_error = req.Get(endpoint,header)
		if conn_error!=nil{
			fmt.Println("Error Connecting. Please check your connection and try again")
			return
		}
		byte_resp,_:=r.ToBytes()
		_,JSONerr:=jsonparser.GetString(byte_resp, "BLCODE")
		if JSONerr!=nil{
			r.ToFile(s)
			return
		}else{
			fmt.Println("Error downloading file.")
			return
		}
	}
}


func get_token(u string, p string) (string,bool) {

	b := User{
		Email:u,
		Password:p,
	}
	mmd, _ := json.Marshal(b)
	
	loc:="/bashlistauth"
	endpoint:=URL+loc

	r, err := req.Post(endpoint, req.BodyJSON(mmd))
	if err != nil {
		return "SSD12",false
	}
	var data map[string]interface{}
	err_ := r.ToJSON(&data)
	if err_ != nil {
		fmt.Println("PRS23")
	}
	if str, ok := data["BLCODE"].(string); ok {
   		if str=="CTR23" {
   			if tok,ok1 := data["access_token"].(string);ok1{
   				return tok,true
   			}else{
   				return "JSE52",false
   			}
   		}
   		return str,false
   	}
   	return "JSE54",false
	
}



func save_secret(url string, u string, t string) {
	c := &credentials.Credentials{
	    ServerURL: url,
	    Username: u,
	    Secret: t,
	}
	nativeStore.Add(c)
}

func retreive_secret(url string)(string,string) {
	username, tok,err := nativeStore.Get(url)
	if err!=nil {
		return "MMW43","X"
	}
	return username,tok
}	


func get_storage_list() {
	// """ Gets a list of stored objects for the user"""
	return
}

func display_files(){
	// """ Displays list of objects for the user in a pretty format"""
	return
}

func show_help() {
	// """ Shows the help page"""
}


func main() {
	// m,d:=up_prompt()
	// l,_:=get_token(m,d)
	// fmt.Println(l)
	// z,_:=get_token(m,d)
	// fmt.Println(z)
	// setup()
	// s:="harsh"
	// p:="Sairam"
	// j:=get_token(s,p)
	// fmt.Println(j)
	// fmt.Println(USERNAME)
	// USERNAME = "ddr"
	// fmt.Println(USERNAME)
	// save_secret("dsdsdsdsdsds","harsh","Sairam")
	// d,c:=retreive_secret("dsdsdsdsdsds")
	// fmt.Println(d)
	// fmt.Println(c)
	// authentication_handler(false)
	// c:=get_account_url()
	// fmt.Println(c)
	// save_token(j)
	// _, tt,_ := nativeStore.Get("Bashlist Credentials")
	// fmt.Println(secret)
	// fmt.Println(tt)
	// fmt.Println(username)
	// t := upload_file("cli_mac.go")
	// fmt.Println(t)
	// test_upload()
	// upload_file("asasdad.go")
	// download_file("pdasaspr.go")
	// current_dir := get_current_dir()
	// file_path := current_dir+"/"+"ppr.go"
	// deleteFile(file_path)
	// fmt.Println(c)
	// m:=get_current_dir()
	// fmt.Println(m)
}