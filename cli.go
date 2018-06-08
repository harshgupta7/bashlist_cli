package main 

import "github.com/fatih/color"
import "fmt"
// import "io/ioutil"
// import "github.com/skratchdot/open-golang/open"
import "path/filepath"
import "os"
// import "runtime"
import "bufio"
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
	reader := bufio.NewReader(os.Stdin)
	username, _ := reader.ReadString('\n')
	
	//Ask for Password
	color.Set(color.FgGreen)
	fmt.Print(("Bashlist Password: "))
	color.Unset() 
	pass, _ := gopass.GetPasswd()
	stringPass:=string(pass)
	return username,stringPass

}


func get_account_url(url string)  {
	
}

// func authentication_handler(exp bool) string{

// 	//token has expired, method is asking for a fresh token
// 	if exp{

// 		//get password from secure storage
// 		u,p:=retreive_secret("bashlistcredentials/pass")

// 		//password has been set before:this is not the first time password is being fetched
// 		if p!="X"{

// 			continue

// 		}

// 	//method is asking for the token the first time
// 	else{
// 		return "yy"
// 	}



// }

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
	
	loc:="/bashlistauth"
	endpoint:=URL+loc
	s := []byte(fmt.Sprintf("{%s:%s,%s:%s}", `"email"`,"\""+u+"\"",`"password"`,"\""+p+"\""))
	r, err := req.Post(endpoint, req.BodyJSON(s))
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
	s,p:=up_prompt()
	// setup()
	// s:="harsh"
	// p:="Sairam"
	j,_:=get_token(s,p)
	fmt.Println(j)
	// fmt.Println(USERNAME)
	// USERNAME = "ddr"
	// fmt.Println(USERNAME)
	// save_secret("dsdsdsdsdsds","harsh","Sairam")
	// d,c:=retreive_secret("dsdsdsdsdsds")
	// fmt.Println(d)
	// fmt.Println(c)

	// save_token(j)
	// _, tt,_ := nativeStore.Get("bashlistcredentials")
	// fmt.Println(secret)
	// fmt.Println(tt)
	// fmt.Println(username)


}