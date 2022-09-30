package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func UploadFile(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// menerima file upload dari frontend
		file, fileHandler, fileErr := r.FormFile("uploadImage")

		if fileErr != nil {
			fmt.Println(fileErr.Error())
			json.NewEncoder(w).Encode("File Upload error!")
			return 
		}
		defer file.Close()
		fmt.Printf("Success upload %+v\n", fileHandler.Filename)

		// mengubah nama image
		tempFile, err := ioutil.TempFile("uploads", "image-*"+fileHandler.Filename)

		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("path upload error")
			json.NewEncoder(w).Encode(err.Error())
		}
		defer tempFile.Close()

		// membaca file
		fileBytes, fileByteErr := ioutil.ReadAll(file);

		if fileByteErr != err {
			fmt.Println(fileByteErr.Error())
		}

		// create image temporary file
		tempFile.Write(fileBytes)

		data := tempFile.Name()
		fileName := data[8:]

		ctx := context.WithValue(r.Context(), "dataFile",fileName)
		next.ServeHTTP(w,r.WithContext(ctx))

	}) 
}