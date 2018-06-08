package main 

import "github.com/fatih/color"
import "fmt"
// import "io/ioutil"
// import "github.com/skratchdot/open-golang/open"
import "path/filepath"
import "os"
import "encoding/json"
// import "runtime"

// import "bufio"
import "github.com/howeyc/gopass"
import "github.com/imroc/req"
// import "encoding/json"

import "github.com/docker/docker-credential-helpers/credentials"
import "github.com/docker/docker-credential-helpers/osxkeychain"


//FOR MAC

var URL string = "http://127.0.0.1:5000"

var nativeStore = osxkeychain.Osxkeychain{}

func up_prompt() (string,string) {
	//Ask for Username
	color.Set(color.FgGreen)
	fmt.Print(("Bashlist Email: "))
	color.Unset() 
	// reader := bufio.NewReader(os.Stdin)
	// username, _ := reader.ReadString('\n')

	// fmt.Print("Enter text: ")
    var input string
    fmt.Scanln(&input)
    // fmt.Print(input)
	
	//Ask for Password
	color.Set(color.FgGreen)
	fmt.Print(("Bashlist Password: "))
	color.Unset() 
	pass, _ := gopass.GetPasswd()
	string_pass := string(pass)
	return input,string_pass

}


func get_account_url(url string)  {
	
}

type User struct{
	Email string
	Password string
}


func incorrect_auth_loop() string {

	for {
		u1,p1 := up_prompt()
		t1,s1 := get_token(u1,p1)
		if s1==true{
			save_secret("bashlistcredentials/pass",u1,p1)
			save_secret("bashlistcredentials/token",u1,t1)
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
		u,p:=retreive_secret("bashlistcredentials/pass")

		//password has been set before:this is not the first time password is being fetched
		if p!="X"{
			token,success := get_token(u,p)
			//Auth succeeded
			if success{
				save_secret("bashlistcredentials/token",u,token)
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

		u,t:=retreive_secret("bashlistcredentials/token")
		//not the first time token will be called
		if t!="X"{
			return t

		// token is not set
		} else {

			u1,p1 := retreive_secret("bashlistcredentials/pass")
			
			//token is not set but password is set
			if p1!="X" {
				token,success := get_token(u1,p1)
				//Auth succeeded
				if success{
					save_secret("bashlistcredentials/token",u,token)
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
	// """Gets account URL and opens it in browser"""
	return
}

func upload_file() {
	// """Uploads a file"""
	return
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

func upload_directory() {
	// """Uploads a directory"""
	return
}

// func download_file(s string) {
// 	// """Downloads a file"""
// 	msg := "Fetching "+ s
// 	r, err := req.Get(url)
// 	if err != nil{
// 		fmt.Println("Error downloading file. Please check your connection")
// 	}
// 	r.ToFile(s)
// 	fmt.Println("Done.")
// 	return
// }

func download_directory() {
	// """Downloads a directory"""
	return
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
	authentication_handler(false)

	// save_token(j)
	// _, tt,_ := nativeStore.Get("bashlistcredentials")
	// fmt.Println(secret)
	// fmt.Println(tt)
	// fmt.Println(username)





}