package main 


import "github.com/fatih/color"
import "os"
import "fmt"
import "io/ioutil"
import "github.com/skratchdot/open-golang/open"
import "path/filepath"
import "os"
import "fmt"

var PLATFORM int
var USERNAME string

func setup() {
	"""Saves username and password"""
	return
}

func is_setup() {
	"""checks whether username and password exist"""
}

func open_account_url() {
	"""Gets account URL and opens it in browser"""
	return
}

func upload_file() {
	"""Uploads a file"""
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
	"""Uploads a directory"""
	return
}

func download_file() {
	"""Downloads a file"""
	return
}

func download_directory() {
	"""Downloads a directory"""
	return
}


func get_token() {
	"""Gets a fresh auth token from server"""
	return
}


func get_storage_list() {
	""" Gets a list of stored objects for the user"""
	return
}

func display_files(){
	""" Displays list of objects for the user in a pretty format"""
	return
}

func show_help() {
	""" Shows the help page"""
}


func main() {

}