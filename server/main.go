package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	addr string

	tmpsPath, adminTmpsPath, baseTmpPath, adminBaseTmpPath string
	staticPath, imagesPath, dbPath                         string

	tmps  = NewSyncMap[string, *template.Template]()
	dbMtx *RWMutex[*sql.DB]
)

func init() {
	log.SetFlags(0)

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("error getting this file")
	}
	parentDir := filepath.Dir(filepath.Dir(thisFile))

	tmpsPath = filepath.Join(parentDir, "templates")
	baseTmpPath = filepath.Join(tmpsPath, "base.html")
	adminTmpsPath = filepath.Join(tmpsPath, "admin")
	adminBaseTmpPath = filepath.Join(adminTmpsPath, "base.html")

	staticPath = filepath.Join(parentDir, "static")
	imagesPath = filepath.Join(staticPath, "images")

	dbPath = filepath.Join(parentDir, "db")
}

func main() {
	flag.StringVar(&addr, "addr", "127.0.0.1:8000", "Address to run server on")
	flag.Parse()

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalln("error opening database:", err)
	}
	dbMtx = NewRWMutex(db)

	if err := parseAllTmps(); err != nil {
		log.Fatalln("error parsing templates", err)
	}

	http.HandleFunc("/", handler)
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/forsale", forSaleHandler)
	http.HandleFunc("/forsale/details", forSaleDetailHandler)
	http.HandleFunc("/forsale/buy-success", buySuccessHandler)
	http.HandleFunc("/policies", policiesHandler)
	http.HandleFunc("/admin", adminHandler)
	http.HandleFunc("/admin/parse", parseHandler)
	http.Handle(
		"/static/",
		http.StripPrefix("/static", http.FileServer(http.Dir(staticPath))),
	)
	http.Handle(
		"/static/images/",
		http.StripPrefix("/static/images", http.HandlerFunc(imagesHandler)),
	)
	http.Handle(
		"/api/",
		http.StripPrefix("/api", makeApiHandler()),
	)
	fmt.Println("Starting server on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "" || path == "/" {
		http.Redirect(
			w, r, fmt.Sprintf("http://%s/home", addr), http.StatusFound,
		)
		return
	}
	http.NotFound(w, r)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	ts, _ := tmps.Load("home")
	if err := ts.Execute(w, nil); err != nil {
		log.Println("error executing home template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	ts, _ := tmps.Load("about")
	if err := ts.Execute(w, nil); err != nil {
		log.Println("error executing about template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func forSaleHandler(w http.ResponseWriter, r *http.Request) {
	ts, _ := tmps.Load("forsale")
	if err := ts.Execute(w, nil); err != nil {
		log.Println("error executing forsale template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func forSaleDetailHandler(w http.ResponseWriter, r *http.Request) {
	type PageData struct {
		Frag Frag
	}
	idStr := r.URL.Query().Get("frag_id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id >= uint64(len(frags)) {
		idStr = ""
	}
	if idStr == "" {
		http.NotFound(w, r)
		return
	}
	ts, _ := tmps.Load("forsale-detail")
	if err := ts.Execute(w, PageData{Frag: frags[id]}); err != nil {
		log.Println("error executing forsale-detail template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func buySuccessHandler(w http.ResponseWriter, r *http.Request) {
	ts, _ := tmps.Load("buy-success")
	if err := ts.Execute(w, nil); err != nil {
		log.Println("error executing buy-success template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func policiesHandler(w http.ResponseWriter, r *http.Request) {
	ts, _ := tmps.Load("policies")
	if err := ts.Execute(w, nil); err != nil {
		log.Println("error executing policies template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	ts, _ := tmps.Load("admin/admin")
	if err := ts.Execute(w, nil); err != nil {
		log.Println("error executing admin template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func imagesHandler(w http.ResponseWriter, r *http.Request) {
	//
}

func parseHandler(w http.ResponseWriter, r *http.Request) {
	if err := parseAllTmps(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Parsed at:", time.Now().Format(time.RFC3339))
}

func parseAllTmps() error {
	err := parseHomeTmp()
	if e := parseAboutTmp(); err == nil {
		err = e
	}
	if e := parseForSaleTmp(); err == nil {
		err = e
	}
	if e := parseForSaleDetailTmp(); err == nil {
		err = e
	}
	if e := parsebuySuccessTmp(); err == nil {
		err = e
	}
	if e := parsePoliciesTmp(); err == nil {
		err = e
	}
	if e := parseAdminTmp(); err == nil {
		err = e
	}
	return err
}

func parseHomeTmp() error {
	ts, err := template.New(
		"base.html",
	).Delims("{|", "|}").ParseFiles(
		baseTmpPath, filepath.Join(tmpsPath, "home.html"),
	)
	if err != nil {
		return err
	}
	tmps.Store("home", ts)
	return nil
}

func parseAboutTmp() error {
	ts, err := template.New(
		"base.html",
	).Delims("{|", "|}").ParseFiles(
		baseTmpPath, filepath.Join(tmpsPath, "about.html"),
	)
	if err != nil {
		return err
	}
	tmps.Store("about", ts)
	return nil
}

func parseForSaleTmp() error {
	ts, err := template.New(
		"base.html",
	).Delims("{|", "|}").ParseFiles(
		baseTmpPath, filepath.Join(tmpsPath, "forsale.html"),
	)
	if err != nil {
		return err
	}
	tmps.Store("forsale", ts)
	return nil
}

func parseForSaleDetailTmp() error {
	ts, err := template.New(
		"base.html",
	).Delims("{|", "|}").ParseFiles(
		baseTmpPath, filepath.Join(tmpsPath, "forsale-detail.html"),
	)
	if err != nil {
		return err
	}
	tmps.Store("forsale-detail", ts)
	return nil
}

func parsebuySuccessTmp() error {
	ts, err := template.New(
		"base.html",
	).Delims("{|", "|}").ParseFiles(
		baseTmpPath, filepath.Join(tmpsPath, "buy-success.html"),
	)
	if err != nil {
		return err
	}
	tmps.Store("buy-success", ts)
	return nil
}

func parsePoliciesTmp() error {
	ts, err := template.New(
		"base.html",
	).Delims("{|", "|}").ParseFiles(
		baseTmpPath, filepath.Join(tmpsPath, "policies.html"),
	)
	if err != nil {
		return err
	}
	tmps.Store("policies", ts)
	return nil
}

func parseAdminTmp() error {
	ts, err := template.New(
		"base.html",
	).Delims("{|", "|}").ParseFiles(
		adminBaseTmpPath, filepath.Join(adminTmpsPath, "admin.html"),
	)
	if err != nil {
		return err
	}
	tmps.Store("admin/admin", ts)
	return nil
}
