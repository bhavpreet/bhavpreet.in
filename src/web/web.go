package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type Config struct {
	Host         string
	Port         string
	ReadTimeout  int
	WriteTimeout int
	RootPath     string
}

var conf Config

// Default configuration
// var conf Config = Config{
// 	Host:         "",
// 	Port:         "8080",
// 	ReadTimeout:  15,
// 	WriteTimeout: 15,
// 	RootPath:   "/tmp",
// }

type Project struct {
	Type        string   `json:"type"`
	Name        string   `json:"name"`
	Link        string   `json:"link"`          // Internal link, later constucted a /p.Type/p.Link
	MainImage   string   `json:"main_image"`    // relative URL
	Description string   `json:"description"`   // Project Description
	Images      []string `json:"images"`        // list of images for display in carousel
	ExtLink     string   `json:"external_link"` // case of external link
}

type Category struct {
	Name     string    `json:"name"`
	Projects []Project `json:"projects"`
}

type MainData struct {
	MainImage string     `json:"main_image"`
	Cat       []Category `json:"categories"`
}

var Info = MainData{}
var g_tmpl = template.New("")
var url_rx *regexp.Regexp

func MainHandler(w http.ResponseWriter, r *http.Request) {
	g_tmpl.ExecuteTemplate(w, "page_main", struct{ Info MainData }{Info: Info})
	//t.Execute(w, "")
}

func ProjHandler(w http.ResponseWriter, r *http.Request) {
	var (
		_cat       string
		_proj_link string
		proj       *Project
	)

	z := url_rx.FindAllStringSubmatch(r.URL.Path, -1)
	_cat = z[0][1]
	_proj_link = z[0][2]

	// find project and assign and send
	for _, c := range Info.Cat {
		if strings.ToLower(c.Name) == strings.ToLower(_cat) {
			for _, p := range c.Projects {
				if strings.ToLower(p.Link) == strings.ToLower(_proj_link) {
					proj = &p
					// fmt.Println("XXXXXXXXX Found project", proj.Name)
					break
				}
			}
		}
		if proj != nil {
			break
		}
	}

	g_tmpl.ExecuteTemplate(w, "page_proj", struct {
		Info MainData
		Proj *Project
	}{Info: Info, Proj: proj})
	//t.Execute(w, "")
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	g_tmpl.ExecuteTemplate(w, "page_about", struct{ Info MainData }{Info: Info})
	//t.Execute(w, "")
}

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	g_tmpl.ExecuteTemplate(w, "page_contact", struct{ Info MainData }{Info: Info})
	//t.Execute(w, "")
}

var Urls = map[string]string{}

func UrlFor(name string) string {
	u, ok := Urls[name]
	if !ok {
		u = ""
	}
	// fmt.Println("URL FOR: ", name, u)
	return u
}

var fmap = template.FuncMap{
	"url_for": UrlFor,
}

func AddHandlerFunc(name string, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	Urls[name] = pattern
	// fmt.Println("URLS: ", Urls)
	http.HandleFunc(pattern, handler)
}

func run_web() {
	// e := syscall.Chroot("~/noise/gogo/")
	// if e != nil {
	// 	log.Fatal(e)
	// }

	AddHandlerFunc("main", "/", MainHandler)
	for _, c := range Info.Cat {
		for _, p := range c.Projects {
			AddHandlerFunc(p.Link, fmt.Sprintf("/%s/%s",
				strings.ToLower(c.Name),
				strings.ToLower(p.Link)),
				ProjHandler)
		}
	}

	AddHandlerFunc("about", "/about", AboutHandler)
	AddHandlerFunc("contact", "/contact", ContactHandler)

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir(conf.RootPath+"/static/"))))
	srv := &http.Server{
		Handler: nil, // http.DefaultServeMux
		Addr:    fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: time.Duration(conf.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func ReadConfig() {
}

func init() {
	flag.StringVar(&conf.Host, "host", "", "Host name to bind on")
	flag.StringVar(&conf.Port, "port", "8080", "Port to bind on")
	flag.IntVar(&conf.ReadTimeout, "rtimeout", 15, "Read Timeout")
	flag.IntVar(&conf.WriteTimeout, "wtimeout", 15, "Write Timeout")
	flag.StringVar(&conf.RootPath, "root", ".", "Path to web ROOT")

	flag.Parse()

	port := os.Getenv("PORT") // For Heroku, port is supplied as environ ??
	if port != "" {
		conf.Port = port
	}
	// compile templates
	layout_tmpl := conf.RootPath + "/templates/layout.html"
	index_tmpl := conf.RootPath + "/templates/index.html"
	proj_tmpl := conf.RootPath + "/templates/proj.html"
	about_tmpl := conf.RootPath + "/templates/about.html"
	contact_tmpl := conf.RootPath + "/templates/contact.html"
	_, err := g_tmpl.Funcs(fmap).ParseFiles(index_tmpl, proj_tmpl, layout_tmpl, about_tmpl, contact_tmpl)
	if err != nil {
		fmt.Println(err)
	}

	r, err := ioutil.ReadFile("data.json")
	err = json.Unmarshal(r, &Info)
	if err != nil {
		fmt.Println("Error reading data.json", err)
		os.Exit(-1)
	}
	// fmt.Printf("%v", Info.Cat)
	// for _, c := range Info.Cat {
	// 	for _, p := range c.Projects {
	// 		fmt.Println(p.Name)
	// 		if p.Type == "internal" {
	// 			fmt.Println(p.Description)
	// 		} else {
	// 			fmt.Println(p.ExtLink)
	// 		}
	// 	}
	// }

	url_rx = regexp.MustCompile("^/(.*)/(.*)$")

}

func main() {
	run_web()
}
