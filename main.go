package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var authFile = "auth.lock"
var authKey = ""

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}

	return b, nil
}

func main() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/cms", func(w http.ResponseWriter, r *http.Request) {
		files, err := ioutil.ReadDir("./cms/stable")
		if err != nil {
			_, _ = w.Write([]byte("Files not found"))
			w.WriteHeader(http.StatusNotFound)
			return
		}

		data := map[string]string{}
		for _, file := range files {
			name := file.Name()
			name = strings.ReplaceAll(name, ".zip", "")
			data[name] = "https://releases.jinya.de/cms/" + file.Name()
		}

		encodedJson, err := json.Marshal(data)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(encodedJson)
	})
	rtr.HandleFunc("/cms/push/{version}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("JinyaAuthKey") == authKey {
			err := bcrypt.CompareHashAndPassword([]byte(authKey), []byte(r.Header.Get("JinyaAuthKey")))
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		vars := mux.Vars(r)
		version := vars["version"]
		if _, err := os.Stat("./cms/stable/"); os.IsNotExist(err) {
			err = os.MkdirAll("./cms/stable/", os.ModePerm)
			if err != nil {
				_, _ = w.Write([]byte(err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		writer, err := os.OpenFile("./cms/stable/"+version+".zip", os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer writer.Close()

		reader := bufio.NewReader(r.Body)
		_, err = reader.WriteTo(writer)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
	rtr.HandleFunc("/cms/{version}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		version := vars["version"]
		file, err := os.OpenFile("./cms/stable/"+version, os.O_RDONLY, os.ModeAppend)
		if err != nil {
			_, _ = w.Write([]byte("File not found"))
			w.WriteHeader(http.StatusNotFound)
			return
		}

		reader := bufio.NewReader(file)
		reader.WriteTo(w)
		w.WriteHeader(http.StatusOK)
	})

	if _, err := os.Stat(authFile); os.IsNotExist(err) {
		key, err := generateRandomBytes(128)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Println(base64.StdEncoding.EncodeToString(key))

		key, err = bcrypt.GenerateFromPassword(key, 13)
		if err != nil {
			log.Fatal(err.Error())
		}

		err = ioutil.WriteFile(authFile, key, os.ModePerm)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	key, err := ioutil.ReadFile(authFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	authKey = string(key)

	log.Println("Serving at localhost:8090...")
	log.Fatal(http.ListenAndServe(":8090", rtr))
}
