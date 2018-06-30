package main

import "crypto"
import "crypto/aes"
import "crypto/cipher"
import "crypto/md5"
import "crypto/rand"
import "crypto/rsa"
import "crypto/sha256"
import "crypto/x509"
import "encoding/base64"
import "encoding/pem"
import "errors"
import "fmt"
import "golang.org/x/crypto/pbkdf2"
import "io"

const saltlen = 8
const keylen = 32
const iterations = 10000


func AuthPassFromPassword(password string) string{

	sha_256 := sha256.New()
	sha_256.Write([]byte(password))
	hashval:= fmt.Sprintf("%x",sha_256.Sum(nil))
	return hashval

}

func EncryptWithPubKey(msg []byte,priv_key *rsa.PrivateKey,rec_key *rsa.PublicKey)(*[]byte,*[]byte,error){

	/* Encrypts and signs message with receivers private key*/


    label := []byte("")
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, rec_key, msg, label)
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
	signature, sign_err := rsa.SignPSS(rand.Reader, priv_key, newhash, hashed, &opts)
	if sign_err!=nil{
		return nil,nil,sign_err
	} 
	return &ciphertext,&signature,nil
}

func DecryptWithPrivKey(priv_key *rsa.PrivateKey,ciphertext *[]byte)(*[]byte,error){
	/* Decryptes a Ciphertext with PrivateKey. Returns Pointer to decrypted text*/

	//Decrypt Ciphertext
	hash := sha256.New()
	label := []byte("")
	plainText, err := rsa.DecryptOAEP(hash, rand.Reader, priv_key, *ciphertext, label)
	if err!=nil{
		return nil,err
	}
	//TODO:Verify Signature -> Pushed to V2
	return &plainText,nil
}

func EncryptObject(plaintextptr *[]byte, keyptr *[]byte,encrypted chan *[]byte ) () {
	/* Encrypts a byte array with a key*/

	var plaintext = *plaintextptr
	key := *keyptr
	block, err := aes.NewCipher(key[:])
	if err != nil {
		encrypted<-nil
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		encrypted<-nil
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		encrypted<-nil
	}

	ciphertext:= gcm.Seal(nonce, nonce, plaintext, nil)
	encrypted<-&ciphertext
	close(encrypted)
}

func DecryptObject(ciphertextptr *[]byte, keyptr *[]byte) (*[]byte,error) {
	/* Decrypts a byte array with a key*/

	var ciphertext = *ciphertextptr
	key := *keyptr
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


func ExportRsaPrivateKeyAsPemStr(privkey *rsa.PrivateKey) string {
    privkey_bytes := x509.MarshalPKCS1PrivateKey(privkey)
    privkey_pem := pem.EncodeToMemory(
            &pem.Block{
                    Type:  "RSA PRIVATE KEY",
                    Bytes: privkey_bytes,
            },
    )
    return string(privkey_pem)
}

func ParseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
    block, _ := pem.Decode([]byte(privPEM))
    if block == nil {
            return nil, errors.New("failed to parse PEM block containing the key")
    }

    priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
            return nil, err
    }

    return priv, nil
}

func ExportRsaPublicKeyAsPemStr(pubkey *rsa.PublicKey) (string, error) {
    pubkey_bytes, err := x509.MarshalPKIXPublicKey(pubkey)
    if err != nil {
            return "", err
    }
    pubkey_pem := pem.EncodeToMemory(
            &pem.Block{
                    Type:  "RSA PUBLIC KEY",
                    Bytes: pubkey_bytes,
            },
    )

    return string(pubkey_pem), nil
}

func ParseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
    block, _ := pem.Decode([]byte(pubPEM))
    if block == nil {
            return nil, errors.New("failed to parse PEM block containing the key")
    }

    pub, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
            return nil, err
    }

    switch pub := pub.(type) {
    case *rsa.PublicKey:
            return pub, nil
    default:
            break // fall through
    }
    return nil, errors.New("Key type is not RSA")
}

func encrypt_secret(plaintextptr *string, password string) (*string,error) {

	plaintext:=*plaintextptr

    header := make([]byte, saltlen + aes.BlockSize)

    salt := header[:saltlen]
    if _, err := io.ReadFull(rand.Reader, salt); err != nil {
        return nil,err
    }

    iv := header[saltlen:aes.BlockSize+saltlen]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil,err
    }

    key := pbkdf2.Key([]byte(password), salt, iterations, keylen, md5.New)

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil,err
    }

    ciphertext := make([]byte, len(header) + len(plaintext))
    copy(ciphertext, header)

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize+saltlen:], []byte(plaintext))
    val := base64Encode(ciphertext)
    return &val,nil
}

func decrypt_secret(encryptedptr *string, password string) (*string,error) {

	encrypted := *encryptedptr

    a, err := base64Decode([]byte(encrypted))
    if err != nil {
      return nil,err
    }
    ciphertext := a
    salt := ciphertext[:saltlen]
    iv := ciphertext[saltlen:aes.BlockSize+saltlen]
    key := pbkdf2.Key([]byte(password), salt, iterations, keylen, md5.New)

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil,err
    }

    if len(ciphertext) < aes.BlockSize {
        return nil,nil
    }

    decrypted := ciphertext[saltlen+aes.BlockSize:]
    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(decrypted, decrypted)

    val := string(decrypted)
    return &val,nil
}
func base64Encode(src []byte) string {
    return base64.StdEncoding.EncodeToString(src)
}

func base64Decode(src []byte) ([]byte, error) {
    return base64.StdEncoding.DecodeString(string(src))
}

