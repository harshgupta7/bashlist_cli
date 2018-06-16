package main

import "crypto"
import "crypto/aes"
import "crypto/cipher"
import "crypto/rand"
import "crypto/rsa"
import "crypto/sha256"


func generate_key_pair()(*rsa.PrivateKey,*rsa.PublicKey,error){

	/* Generates RSA Key-Pair*/
	
	privkey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err!=nil{
		return nil,nil,err
	}
	pubkey := &privkey.PublicKey
	return privkey,pubkey,err
}

func encrypt_with_pubkey(msg []byte,priv_key *rsa.PrivateKey,rec_key *rsa.PublicKey)(*[]byte,*[]byte,error){

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

func decrypt_with_privkey(priv_key *rsa.PrivateKey,ciphertext *[]byte)(*[]byte,error){
	/* Decryptes a Ciphertext with PrivateKey. Returns Pointer to decrypted text*/

	//Decrypt Ciphertext
	hash := sha256.New()
	label := []byte("")
	plainText, err := rsa.DecryptOAEP(hash, rand.Reader, priv_key, *ciphertext, label)
	if err!=nil{
		return nil,err
	}
	//TODO:Verify Signature
	return &plainText,nil
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