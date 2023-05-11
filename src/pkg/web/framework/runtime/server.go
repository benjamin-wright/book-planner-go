package runtime

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"go.uber.org/zap"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/components/footer"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/components/header"
	"ponglehub.co.uk/book-planner-go/src/pkg/web/framework/component"
)

//go:generate sh -c "cp ${DOLLAR}(tinygo env TINYGOROOT)/targets/wasm_exec.js templates/wasm_exec.js"

type Server struct {
	server *http.ServeMux
	t      *template.Template
	wasm   map[string][]byte
	js     map[string][]byte
	css    map[string][]byte
	static staticContext
}

//go:embed templates/base.html
var baseContent string

func NewServer(content string, title string, children ...component.Component) *Server {
	t := template.New("Page")
	prefix := os.Getenv("PATH_PREFIX")

	pageComponent := component.Component{
		Name:     "page",
		Template: content,
		Children: append([]component.Component{
			header.Get(),
			footer.Get(),
		}, children...),
	}

	pageComponent.Parse(t)
	template.Must(t.Parse(baseContent))

	return &Server{
		server: http.NewServeMux(),
		t:      t,
		wasm:   map[string][]byte{},
		js:     map[string][]byte{},
		css: map[string][]byte{
			"styles.css": styles,
		},
		static: staticContext{
			Title:   title,
			Scripts: []string{},
			CSS:     []string{prefix + "/styles.css"},
		},
	}
}

//go:embed templates/wasm_exec.js
var wasmExec []byte

//go:embed templates/wasm_load.js
var wasmLoad string

//go:embed templates/styles.css
var styles []byte

func (s *Server) AddWASMModule(name string, path string, data []byte) error {
	prefix := os.Getenv("PATH_PREFIX")

	// Load the wasm_exec if it's not already there
	if _, ok := s.js["wasm_exec.js"]; !ok {
		s.js["wasm_exec.js"] = wasmExec
		s.static.Scripts = append(s.static.Scripts, prefix+"/wasm_exec.js")
	}

	// Serve the wasm binary
	s.wasm[path] = data

	// Parse and serve the specific wasm loader library
	t := template.New("Script")
	template.Must(t.Parse(wasmLoad))

	var script bytes.Buffer
	err := t.Execute(&script, map[string]interface{}{
		"Path": prefix + "/" + path,
		"Name": name,
	})
	if err != nil {
		return fmt.Errorf("failed to parse wasm load script: %+v", err)
	}

	s.js[name+"_wasm_load.js"] = script.Bytes()
	s.static.Scripts = append(s.static.Scripts, prefix+"/"+name+"_wasm_load.js")

	return nil
}

type staticContext struct {
	Title   string
	Scripts []string
	CSS     []string
}

type Context struct {
	Static  staticContext
	Dynamic any
}

func (s *Server) Run(handler func(r *http.Request) any) {
	prefix := os.Getenv("PATH_PREFIX")

	mux := http.NewServeMux()
	mux.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		zap.S().Infof("HTTP %s %s", r.Method, r.URL.Path)

		data := handler(r)

		context := Context{
			Static:  s.static,
			Dynamic: data,
		}

		var response bytes.Buffer

		err := s.t.Execute(&response, context)
		if err != nil {
			zap.S().Errorf("Failed to render response: %+v", err)
			w.WriteHeader(500)
			return
		}

		_, err = w.Write(response.Bytes())
		if err != nil {
			zap.S().Errorf("Failed to write response: %+v", err)
		}
	})

	for path, data := range s.js {
		mux.HandleFunc(prefix+"/"+path, fileHandler("text/javascript", data))
	}

	for path, data := range s.wasm {
		mux.HandleFunc(prefix+"/"+path, fileHandler("application/wasm", data))
	}

	for path, data := range s.css {
		mux.HandleFunc(prefix+"/"+path, fileHandler("text/css", data))
	}

	zap.S().Info("running server...")

	err := http.ListenAndServe("0.0.0.0:80", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func fileHandler(contentType string, data []byte) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", contentType)
		w.Write(data)
	}
}
