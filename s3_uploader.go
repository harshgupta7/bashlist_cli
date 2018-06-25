package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/buger/jsonparser"
)


type AWSField struct {
	Key string
	Value string
}

func FieldsCreator(a *[]byte)([]AWSField){
	j := *a
	var fields []AWSField
	c := make(map[string]interface{})
	e := json.Unmarshal(j, &c)
	if e != nil {
		unexpected_event()
	}
	//i := 0
	for s, _ := range c {
		key := s
		value, _ := jsonparser.GetString(j, key)
		g := AWSField{Key: s, Value: value}
		fields = append(fields, g)
	}
	return fields

}

func Upload(url string, fields []AWSField) error {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for _, f := range fields {
		fw, err := w.CreateFormField(f.Key)
		if err != nil {
			return err
		}
		if _, err := fw.Write([]byte(f.Value)); err != nil {
			return err
		}
	}
	w.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	fmt.Println(req)
	fmt.Println(res)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	return nil
}