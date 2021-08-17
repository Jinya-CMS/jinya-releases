package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-version"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
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
	rtr.HandleFunc("/cms/unstable", func(w http.ResponseWriter, r *http.Request) {
		files, err := ioutil.ReadDir("./cms/unstable")
		if err != nil {
			_, _ = w.Write([]byte("Files not found"))
			w.WriteHeader(http.StatusNotFound)
			return
		}

		versions := make([]*version.Version, len(files))
		for i, file := range files {
			name := file.Name()
			name = strings.ReplaceAll(name, ".zip", "")
			ver, _ := version.NewVersion(name)
			versions[i] = ver
			fmt.Println(name)
		}

		json := ""
		if len(versions) > 0 {
			sort.Sort(version.Collection(versions))
			data := make([]string, len(versions))
			for i, ver := range versions {
				data[i] = fmt.Sprintf("\"%s\": \"%s\"", ver.Original(), "https://releases.jinya.de/cms/unstable/"+ver.Original()+".zip")
			}

			json = strings.Join(data, ",")
		}

		w.Write([]byte(fmt.Sprintf("{%s}", json)))
	})
	rtr.HandleFunc("/cms/unstable/push/{version}", func(w http.ResponseWriter, r *http.Request) {
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
		if _, err := os.Stat("./cms/unstable/"); os.IsNotExist(err) {
			err = os.MkdirAll("./cms/unstable/", os.ModePerm)
			if err != nil {
				_, _ = w.Write([]byte(err.Error()))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		writer, err := os.OpenFile("./cms/unstable/"+version+".zip", os.O_CREATE|os.O_WRONLY, 0777)
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
	rtr.HandleFunc("/cms/unstable/{version}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		version := vars["version"]
		file, err := os.OpenFile("./cms/unstable/"+version, os.O_RDONLY, os.ModeAppend)
		if err != nil {
			_, _ = w.Write([]byte("File not found"))
			w.WriteHeader(http.StatusNotFound)
			return
		}

		reader := bufio.NewReader(file)
		reader.WriteTo(w)
		w.WriteHeader(http.StatusOK)
	})

	rtr.HandleFunc("/cms", func(w http.ResponseWriter, r *http.Request) {
		files, err := ioutil.ReadDir("./cms/stable")
		if err != nil {
			_, _ = w.Write([]byte("Files not found"))
			w.WriteHeader(http.StatusNotFound)
			return
		}

		versions := make([]*version.Version, len(files))
		for i, file := range files {
			name := file.Name()
			name = strings.ReplaceAll(name, ".zip", "")
			ver, _ := version.NewVersion(name)
			versions[i] = ver
		}

		sort.Sort(version.Collection(versions))

		data := make([]string, len(versions))
		for i, ver := range versions {
			data[i] = fmt.Sprintf("\"%s\": \"%s\"", ver.Original(), "https://releases.jinya.de/cms/"+ver.Original()+".zip")
		}

		json := strings.Join(data, ",")

		w.Write([]byte(fmt.Sprintf("{%s}", json)))
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
