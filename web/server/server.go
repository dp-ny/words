package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"

	"../../games"

	"github.com/julienschmidt/httprouter"
)

var port = flag.Int("port", 9000, "the port on which to serve")

var partials = "web/views/partials/*.html"
var wordsTemplate = "words.html"

var manager *games.Manager
var templates map[string]*template.Template

func init() {
	manager = games.NewManager()
	templates = make(map[string]*template.Template)
	loadTemplates("web/views")
	loadTemplates("web/views/errors")
	loadTemplates("web/views/partials")
}

func loadTemplates(dir string) {
	views, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, t := range views {
		if t.IsDir() {
			continue
		}
		partials, err := template.New(t.Name()).ParseGlob(partials)
		templates[t.Name()] = partials
		if err != nil {
			panic(err)
		}
		name := path.Join(dir, t.Name())
		templates[t.Name()].ParseFiles(name)
	}
}

func main() {
	flag.Parse()
	router := httprouter.New()
	router.GET("/", Homepage)
	router.GET("/words", Words)
	router.GET("/words/:id", WordsView)
	router.GET("/healthy", Healthy)
	router.GET("/d/:path", Default)
	router.ServeFiles("/public/*filepath", http.Dir("web/public"))

	fmt.Printf("Starting server on port: %d\n", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), router)
	if err != nil {
		fmt.Printf("Unable to start server: %v\n", err)
	}
}

// Homepage is the default landing page for the app
func Homepage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	executeTemplate(w, "index.html", d("Title", "home"))
}

// Default is the landing page for non-configured routes in the app
func Default(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	templateName := lowercaseFirstChar(p.ByName("path"))
	templateFile := templateName
	if !strings.HasSuffix(templateFile, ".html") {
		templateFile = templateFile + ".html"
	}
	executeTemplate(w, templateFile, d("Title", templateName))
}

// Words handles the page for the words app supported on this page
func Words(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var board [][]string
	board = append(board, []string{"W", "O", "R", "D"})
	board = append(board, []string{"S", "#", "#", "#"})
	board = append(board, []string{"#", "B", "Y", "#"})
	board = append(board, []string{"D", "P", "N", "Y"})

	wordsBoard(w, r, board, "")
}

// WordsView retrieves a new game, if requested, or an existing game
func WordsView(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if id == "new" {
		WordsNew(w, r, p)
		return
	}
	game, ok := manager.Get(id)
	if !ok {
		notFound(w, fmt.Sprintf("Unable to find %s", id))
		return
	}
	board := game.Board.ToStringArray()
	wordsBoard(w, r, board, game.JsonTime())
}

func wordsBoard(w http.ResponseWriter, r *http.Request, board [][]string, time string) {
	executeTemplate(w, wordsTemplate, d("Title", "Words", "Words", board, "Time", time))
}

// WordsNew retrieves a new words game to be displayed
func WordsNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	game, err := manager.NewGame()
	if err != nil {
		serverError(w, err)
		return
	}
	board := game.Board.ToStringArray()
	html, err := templateString("_wordsTable.html", d("Words", board))
	if err != nil {
		serverError(w, err)
		return
	}
	jsonResponse(w, d("id", game.ID, "time", game.JsonTime(), "html", html))
}

func executeTemplate(w io.Writer, t string, d map[string]interface{}) {
	// d["Bootstrap"] = "/public/css/bootstrap.css"
	d["Bootstrap"] = "//maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css"
	err := executeTemplateInternal(w, t, d)
	if err != nil {
		serverError(w, err)
	}
}

func executeTemplateInternal(w io.Writer, t string, d map[string]interface{}) error {
	d["Bootstrap"] = "//maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css"
	tmpl, ok := templates[t]
	if !ok {
		return fmt.Errorf("Unable to execute template: %s with %v", t, d)
	}
	return tmpl.ExecuteTemplate(w, t, d)
}

// Healthy returns success for any health checkers of this app
func Healthy(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("good"))
}

func serverError(w io.Writer, err error) {
	internalErr := executeTemplateInternal(w, "500.html", d("Error", err.Error()))
	if internalErr != nil {
		fmt.Fprintf(os.Stderr, "Something went horribly wrong. Error: %s, trying to show error: %s\n", internalErr.Error(), err.Error())
	}
}

func notFound(w io.Writer, msg string) {
	internalErr := executeTemplateInternal(w, "404.html", d("Error", msg))
	if internalErr != nil {
		fmt.Fprintf(os.Stderr, "Something went horribly wrong. Error: %s, trying to show error (%d): %s\n", internalErr.Error(), 404, msg)
	}
}

func jsonResponse(w http.ResponseWriter, data map[string]interface{}) {
	encoder := json.NewEncoder(w)
	w.Header().Add("Content-Type", "application/json")
	encoder.Encode(data)
}

func templateString(t string, d map[string]interface{}) (string, error) {
	w := new(bytes.Buffer)
	err := executeTemplateInternal(w, t, d)
	if err != nil {
		return "#", err
	}
	return w.String(), nil
}

// d makes data based on key-value pairs where the key is always a string
func d(datas ...interface{}) map[string]interface{} {
	if len(datas)%2 != 0 {
		panic("d must only be called with key value pairs")
	}
	m := make(map[string]interface{})
	for i := 0; i < len(datas); i += 2 {
		k := datas[i]
		s, ok := k.(string)
		if !ok {
			panic("d must only be called with string as the first of a key-value pair")
		}
		m[s] = datas[i+1]
	}
	return m
}

func lowercaseFirstChar(str string) string {
	return strings.ToLower(string(str[0])) + str[1:]
}
