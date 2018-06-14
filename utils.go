package main 

import "fmt"
import "io/ioutil"
import "os"
import "path/filepath"



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


func bucket_exists(dirname string)(bool){
	/*Checks whether a directory exists in the cwd or not*/
	exists,err:=object_exists(dirname)
	if err!=nil{
		fmt.Println("An Unexpected Error Occurred.Please Try Again Later")
		os.Exit(1)
	}
	if !exists{
		fmt.Println(dirname+": No such file or directory")
		return false
	}
	isDir,dirErr := IsDirectory(dirname)
	if dirErr!=nil{
		fmt.Println("An Unexpected Error Occurred.Please Try Again Later")
		os.Exit(1)
	}
	if !isDir{
		fmt.Println(dirname+": Not a directory. Only directories can be pushed to bashlist.")
		os.Exit(1)
	}
	return true
}

func get_size(path string) (int64, error) {
	/*Gets the size of directory*/
    var size int64
    err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
        if !info.IsDir() {
            size += info.Size()
        }
        return err
    })
    return size, err
}

func get_object_count(directory string)int{
	/* Counts number of objects in directory*/
	files,err := ioutil.ReadDir(directory)
	if err!=nil{
		fmt.Println("An Unexpected Error Occurred. Please try again later")
		os.Exit(1)
	}
	return len(files)
}

func get_cwd()*string{
	/* Gets the current working directory*/
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err!=nil{
    	fmt.Println("An Unexpected Error Occurred. Please try again later")
		os.Exit(1)
    }
    return &dir
}



