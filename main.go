package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-version"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"io"
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

type Build struct {
	Version string
	Link    string
	Created string
}

func main() {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("./templates/homepage.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error rendering homepage"))
			return
		}

		io.Copy(w, file)
	})
	rtr.HandleFunc("/cms/unstable", func(w http.ResponseWriter, r *http.Request) {
		sendFileOverview(w, r, "Unstable")
	})
	rtr.HandleFunc("/cms/unstable/push/{version}", func(w http.ResponseWriter, r *http.Request) {
		pushNewVersion(w, r, "unstable")
	})
	rtr.HandleFunc("/cms/unstable/{version}", func(w http.ResponseWriter, r *http.Request) {
		downloadFile(w, r, "unstable")
	})

	rtr.HandleFunc("/cms", func(w http.ResponseWriter, r *http.Request) {
		sendFileOverview(w, r, "Stable")
	})
	rtr.HandleFunc("/cms/push/{version}", func(w http.ResponseWriter, r *http.Request) {
		pushNewVersion(w, r, "stable")
	})
	rtr.HandleFunc("/cms/{version}", func(w http.ResponseWriter, r *http.Request) {
		downloadFile(w, r, "stable")
	})

	rtr.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

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

		err = os.WriteFile(authFile, key, os.ModePerm)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	key, err := os.ReadFile(authFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	authKey = string(key)

	log.Println("Serving at localhost:8090...")
	log.Fatal(http.ListenAndServe(":8090", rtr))
}

func downloadFile(w http.ResponseWriter, r *http.Request, stability string) {
	vars := mux.Vars(r)
	ver := vars["version"]
	file, err := os.OpenFile("./cms/"+stability+"/"+ver, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		_, _ = w.Write([]byte("File not found"))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	defer file.Close()

	io.Copy(w, file)
	w.WriteHeader(http.StatusOK)
}

func pushNewVersion(w http.ResponseWriter, r *http.Request, stability string) {
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
	ver := vars["version"]
	if _, err := os.Stat("./cms/" + stability); os.IsNotExist(err) {
		err = os.MkdirAll("./cms/"+stability, os.ModePerm)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	writer, err := os.OpenFile("./cms/"+stability+"/"+ver+".zip", os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer writer.Close()

	_, err = io.Copy(writer, r.Body)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func sendFileOverview(w http.ResponseWriter, r *http.Request, buildType string) {
	stability := strings.ToLower(buildType)
	basePath := "./cms/" + stability + "/"
	if stability == "stable" {
		stability = ""
	} else {
		stability += "/"
	}
	files, err := os.ReadDir(basePath)
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

	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		tmpl, err := template.New("page").ParseFiles("./templates/builds.gohtml")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		sort.Sort(sort.Reverse(version.Collection(versions)))
		builds := make([]Build, len(versions))
		for i, ver := range versions {
			stat, err := os.Stat(basePath + ver.Original() + ".zip")
			var created string
			if err != nil {
				created = ""
			} else {
				created = stat.ModTime().Format("2006-01-02")
			}
			builds[i] = Build{
				Link:    "https://releases.jinya.de/cms/" + stability + ver.Original() + ".zip",
				Version: ver.Original(),
				Created: created,
			}
		}
		tmpl.ExecuteTemplate(w, "page", struct {
			Builds    []Build
			Stability string
		}{
			Builds:    builds,
			Stability: buildType,
		})

		w.Header().Add("Content-Type", "text/html")
	} else {
		data := make([]string, len(versions))
		for i, ver := range versions {
			data[i] = fmt.Sprintf("\"%s\": \"%s\"", ver.Original(), "https://releases.jinya.de/cms/"+stability+ver.Original()+".zip")
		}

		json := strings.Join(data, ",")
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf("{%s}", json)))
	}
}
