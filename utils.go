package main 

import "bytes"
import "crypto/aes"
import "crypto/cipher"
import "crypto/rand"
import "crypto/rsa"
import "crypto/sha256"
import "errors"
import "fmt"
import "github.com/pierrre/archivefile/zip"
import "io"
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
	cwd_address := get_cwd()
	cwd := *cwd_address
	path := cwd+"/"+dirname
	exists,err:=object_exists(path)
	if err!=nil{
		fmt.Println("An Unexpected Error Occurred.Please Try Again Later")
		os.Exit(1)
	}
	if !exists{
		fmt.Println(dirname+": No such file or directory")
		return false
	}
	isDir,dirErr := IsDirectory(path)
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

func generate_key_pair()(*rsa.PrivateKey,*rsa.PublicKey,error){
	privkey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err!=nil{
		return nil,nil,err
	}
	pubkey := &privkey.PublicKey
	return &privkey,&pubkey,err
}

func encrypt_with_pubkey(msg []byte,priv_key *rsa.PrivateKey,rec_key *rsa.PublicKey)(*[]byte,*[]byte,error){
    label := []byte("")
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, *rec_key, msg, label)
	if err!=nil{
		return nil,nil,err
	}
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto
	PSSmessage := msg
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, sign_err := rsa.SignPSS(rand.Reader, *priv_key, newhash, hashed, &opts)
	if sign_err!=nil{
		return nil,nil,sign_err
	} 
	return &ciphertext,&signature,nil
}

func decrypt_with_privkey(priv_key *rsa.PrivateKey, sender_pubkey *rsa.PublicKey, ciphertext *[]byte){
	//Decrypt Ciphertext
	hash := sha256.New()
	label := []byte("")
	plainText, err := rsa.DecryptOAEP(hash, rand.Reader, *priv_key, *ciphertext, label)
	if err!=nil{
		return nil,err
	}
	// Verify Signature
	// TODO
	// newhash := crypto.SHA256
	// err = rsa.VerifyPSS(sender_pubkey, newhash, hashed, signature, &opts)

	return &plaintext
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

func dir_to_compressed_bytes(dirname string)(*[]byte,error){
	/* Compresses a directory and converts it to byte array*/

	buf := new(bytes.Buffer)
	progress := func(archivePath string){}
	err := zip.Archive(dirname,buf,progress)
	var arr []byte = buf.Bytes()
	return &arr,err

}

func Encrypt(plaintextptr *[]byte, key *[32]byte) (*[]byte, error) {
	/* Encrypts a byte array with a key*/

	var plaintext []byte = *plaintextptr

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	ciphertext:= gcm.Seal(nonce, nonce, plaintext, nil)
	return &ciphertext,nil
}

func Decrypt(ciphertextptr *[]byte, key *[32]byte) (*[]byte,error) {
	/* Decrypts a byte array with a key*/

	var ciphertext []byte = *ciphertextptr

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	plaintext,plaintexterr:= gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
	return &plaintext,plaintexterr
} 




