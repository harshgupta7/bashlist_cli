package main 

import "bytes"
import "crypto/rand"
import "fmt"
import "github.com/pierrre/archivefile/zip"
import "os"
import "path/filepath"
import (
	"errors"
	"github.com/fatih/color"
	"log"
)


func get_code_path()(string){
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
            log.Fatal(err)
    }
    return dir
}

func unexpected_event(){
	color.Red("Bashlist encountered an unexpected error. Please try again later.")
	os.Exit(1)
}

func generate_random_string(length int)(string,error){

	var chars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	if length == 0 {
		return "",nil
	}
	clen := len(chars)
	if clen < 2 || clen > 256 {
		return "",errors.New("Insufficient length size")
		
	}
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) 
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			return "",err
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				continue
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b),nil
			}
		}
	}
}

func object_exists(path string) (bool, error) {
	/* Checks whether a path is a valid path*/
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}

func IsDirectory(path string) (bool, error) {

	/* Checks whether a object is directory or a file*/
    fileInfo, err := os.Stat(path)
    return fileInfo.IsDir(), err
}


func directory_exists(dirname string, context string)(bool){
	/*Checks whether a directory exists in the cwd or not*/
	cwd_address := get_cwd()
	cwd := *cwd_address
	path := cwd+"/"+dirname
	exists,err:=object_exists(path)
	if err!=nil{
		unexpected_event()
	}
	if !exists && context=="pull"{
		return false
	}
	if !exists && context=="push"{
		color.Cyan(dirname+": No such file or directory")
		return false
	}
	isDir,dirErr := IsDirectory(path)
	if dirErr!=nil{
		unexpected_event()
	}
	if isDir && context=="pull"{
		return true
	}else if !isDir && context=="pull"{
		return false
	}
	if !isDir && context=="push"{
		fmt.Println(dirname+": Not a directory. Only directories can be pushed to bashlist.")
		fmt.Println("Place "+dirname+" inside a directory to push.")
		return false
	}
	return false
}


func get_size(path string) (int64, error) {
	/*Gets the size of object*/
	//WARNING: MUST CHECK FOR EXISTENCE BEFORE USING. FAILURE WILL RESULT IN SEG FAULT.
    var size int64
    err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
        if !info.IsDir() {
            size += info.Size()
        }
        return err
    })
    return size, err
}

//func get_object_count(directory string)int{
//	/* Counts number of objects in directory*/
//	files,err := ioutil.ReadDir(directory)
//	if err!=nil{
//		fmt.Println("An Unexpected Error Occurred. Please try again later")
//		os.Exit(1)
//	}
//	return len(files)
//}

func get_cwd()*string{
	/* Gets the current working directory*/
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err!=nil{
    	fmt.Println("An Unexpected Error Occurred. Please try again later")
		os.Exit(1)
    }
    return &dir
}


func dir_to_compressed_bytes(dirname string,done chan *[]byte)() {
	/* Compresses a directory and converts it to byte array*/
	donesig := color.New(color.FgGreen).SprintFunc()
	buf := new(bytes.Buffer)
	progress := func(archivePath string) {
		fmt.Printf("Processing: %s....%s\n", archivePath, donesig("OK"))
	}
	err := zip.Archive(dirname, buf, progress)
	if err != nil {
		color.Red("Bashlist encountered an unexpected error while processing %s", dirname)
	}

	var arr= buf.Bytes()
	done <- &arr

	close(done)

}





