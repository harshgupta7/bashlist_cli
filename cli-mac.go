package main 

import "bufio"
import "encoding/json"
import "fmt"
import "github.com/buger/jsonparser"
import "github.com/docker/docker-credential-helpers/credentials"
import "github.com/docker/docker-credential-helpers/osxkeychain"
import "github.com/fatih/color"
import "github.com/howeyc/gopass"
import "github.com/imroc/req"
import "github.com/olekukonko/tablewriter"
import "github.com/skratchdot/open-golang/open"
import "path/filepath"
import "os"
import "strconv"




var URL string = "http://127.0.0.1:5000"

var nativeStore = osxkeychain.Osxkeychain{}

type User struct{
	Email string
	Password string
}

func up_prompt() (string,string) {

	/* Displays prompt to ask for username, password. 
	Returns Strings */

	//Ask for Username
	color.Set(color.FgGreen)
	fmt.Print(("Bashlist Email: "))
	color.Unset() 
	var desc string
    scanner := bufio.NewScanner(os.Stdin)
  	for scanner.Scan() {
  	    desc = scanner.Text()
  	    break
  	}
	//Ask for Password
	color.Set(color.FgGreen)
	fmt.Print(("Bashlist Password: "))
	color.Unset() 
	pass, _ := gopass.GetPasswd()
	string_pass := string(pass)
	return desc,string_pass

}


func get_account_url() string {

	/* Gets Personal Account URL */

	var endpoint string = URL+"/account"
	var token string = authentication_handler(false)
	val := "JWT " + token

	header := req.Header{
		"Content-Type":"application/json",
		"Authorization": val,
	}

	r, err := req.Get(endpoint, header)
	if err != nil {
		fmt.Println("Error Connecting. Please check your internet connection.")
		return "ERR"
	}

	byte_resp,_:=r.ToBytes()
	response,err :=jsonparser.GetString(byte_resp,"URL")
	if err==nil{
		return response
	}

	response1,err1:=jsonparser.GetString(byte_resp, "BLCODE")
	if err1==nil{
		if response1=="INC23"{
			token=authentication_handler(true)
			val = "JWT " + token

			header = req.Header{
				"Content-Type":"application/json",
				"Authorization": val,
			}
			r, err = req.Get(endpoint, header)
			if err != nil {
				fmt.Println("Error Connecting. Please check your internet connection.")
				return "ERR"
			}

			byte_resp,_:=r.ToBytes()
			response,err :=jsonparser.GetString(byte_resp,"URL")
			if err==nil{
				return response
			}

		}
	}
	return "Error Retreiving Account URL"

}


func incorrect_auth_loop() string {

	/* Keeps asking for username and password until correct combination is entered*/

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

	/* This is the method that each method will call.
	Handles retreiving token, and if token is expired, fetching a new one*/


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

	/* Opens a url in browser*/

	// opens url in browser
	open.Run(url)

}


func upload_file(s string) {

	/* Uploads a file to the URL*/

	// """Uploads a file"""

	token:=authentication_handler(false)

	current_dir := get_current_dir()
	file_path := current_dir+"/"+s
	if objectExists(file_path) == false{
		fmt.Println(s+": No such file or directory")
		return
	}

	fmt.Print("Description (Press Enter to Leave Blank): ")

	var desc string
  	scanner := bufio.NewScanner(os.Stdin)
  	for scanner.Scan() {
  	    desc = scanner.Text()
  	    break
  	}


	endpoint:=URL+"/filesync"
	val := "JWT " + token
	header := req.Header{
		"Content-Type":"application/json",
		"Authorization": val,
		"Description":desc,
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
	}else if response=="OVX23"{
		fmt.Println("Not Enough Space to Upload File.")
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
	}else if response=="OVX23"{
		fmt.Println("Not Enough Space to Upload File.")
		return
	}

	fmt.Println("Error Uploading File")
	return
}

func get_current_dir() string{

	/* Gets current directory*/

    dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
    return dir
}


type Mailer struct{
	Recs string 
	Filename string
}

func mail_file(r string,f string)string{
	/* Mails a file in users bashlist storage*/

	token := authentication_handler(false)
	val:="JWT "+ token

	b := Mailer{
		Recs:r,
		Filename:f,
	}
	mmd, _ := json.Marshal(b)

	header := req.Header{
		"Content-Type":"application/json",
		"Authorization": val,
	}
	
	loc:="/sendmail"
	endpoint:=URL+loc

	resp, err := req.Post(endpoint, req.BodyJSON(mmd),header)
	if err != nil {
		return "SSD12"
	}
	fmt.Println(resp)
	return resp.String()

}


func objectSize(path string) (int64, error) {

	/* Gets the size of a file or directory*/

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

	/* Checks whether a file or directory with name exists or not.*/

    if _, err := os.Stat(name); err != nil {
    if os.IsNotExist(err) {
            return false
        }
    }
    return true
}


func download_file(s string) {

	/* Downloads a file from users bashlist*/

	// """Downloads a file"""
	token:=authentication_handler(false)
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

	/* Gets a token from the server*/

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
		return "PRS23",false
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

	/* Saves secret in users credential manager*/

	c := &credentials.Credentials{
	    ServerURL: url,
	    Username: u,
	    Secret: t,
	}
	nativeStore.Add(c)
}

func retreive_secret(url string)(string,string) {

	/* Retreives secret from users credentials manager*/

	username, tok,err := nativeStore.Get(url)
	if err!=nil {
		return "MMW43","X"
	}
	return username,tok
}	


type Item struct {
	Name string
	Size int
	Updated_On string 
	Description string
}


func lister(jsonBytes []byte)[][]string{

	/* Method to parse JSON and return it as a matrix of string*/

	var f interface{}
	err := json.Unmarshal(jsonBytes, &f)
	if err != nil {
		fmt.Println("Error parsing JSON: ", err)
	}
	var retarr [][]string
	// JSON object parses into a map with string keys
	itemsMap := f.(map[string]interface{})

	// Loop through the Items; we're not interested in the key, just the values
	for _, v := range itemsMap {
		switch jsonObj := v.(type) {
		case interface{}:
			var item Item
			for itemKey, itemValue := range jsonObj.(map[string]interface{}) {
				switch itemKey {
				case "Name":
					switch itemValue := itemValue.(type) {
					case string:
						item.Name = itemValue
					default:
						fmt.Println("Incorrect type for", itemKey)
					}
				case "Size":
					switch itemValue := itemValue.(type) {
					case float64:
						item.Size = int(itemValue)
					default:
						fmt.Println("Incorrect type for", itemKey)
					}

				case "Updated_On":
					switch itemValue := itemValue.(type) {
					case string:
						item.Updated_On = itemValue
					default:
						fmt.Println("Incorrect type for", itemKey)
					}
				case "Description":
					switch itemValue := itemValue.(type) {
					case string:
						item.Description = itemValue
					default:
						fmt.Println("Incorrect type for", itemKey)
					}

				default:
					fmt.Println("Unknown key for Item found in JSON")
				}
			}
			var t []string
			var fin string
			t = append(t, item.Name)
			display_val := int(item.Size/(1000*1000))
			if display_val < 1{
				display_val = int(item.Size/1000)
				fin =strconv.Itoa(display_val)
				fin = fin+" KB"
			}else{
				fin:=strconv.Itoa(display_val)
				fin = fin+" MB"
			}
			t = append(t,fin)
			t = append(t,item.Updated_On)
			t = append(t,item.Description)
		
			retarr = append(retarr,t)
		// Not a JSON object; handle the error
		default:
			fmt.Println("Expecting a JSON object; got something else")
		}
	}
	return retarr

}

func get_storage_list() {

	/* Displays the list of items in user's bashlist*/

	endpoint:=URL+"/getallfiles"
	token:=authentication_handler(false)
	val := "JWT " + token

	header := req.Header{
		"Content-Type":"application/json",
		"Authorization": val,
	}
	r, err := req.Get(endpoint, header)
	if err != nil {
		fmt.Println("Error Connecting. Please check your internet connection.")
		return
	}
	byteVal,_ := r.ToBytes()

	response,err :=jsonparser.GetString(byteVal,"BLCODE")
	if err==nil{
		if response=="INC23"{
			token:=authentication_handler(true)
			val = "JWT " + token
			header = req.Header{
				"Content-Type":"application/json",
				"Authorization": val,
			}

			r, err = req.Get(endpoint, header)
			if err != nil {
				fmt.Println("Error Connecting. Please check your internet connection.")
				return
			}
			byteVal,_ = r.ToBytes()

		}else{
			fmt.Println("Error Connecting. Please check your internet connection.")
			return
		}
	}

	data:=lister(byteVal)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Size", "Last Updated", "Description"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("")
	table.AppendBulk(data) // Add Bulk Data
	table.Render()

}


func show_help() {
	// """ Shows the help page"""
}


func main() {
	return

}