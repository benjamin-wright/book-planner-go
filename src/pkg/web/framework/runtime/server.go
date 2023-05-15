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

//go:embed templates/base.html
var baseContent string

//go:embed templates/wasm_exec.js
var wasmExec []byte

//go:embed templates/wasm_load.js
var wasmLoad string

type ServerOptions struct {
	Template    string
	Title       string
	HideHeaders bool
	Children    []component.Component
	WASMModules []WASMModule
	PageHandler func(r *http.Request) any
	PostHandler func(w http.ResponseWriter, r *http.Request)
}

type WASMModule struct {
	Name string
	Path string
	Data []byte
}

type staticContext struct {
	Title       string
	Scripts     []string
	CSS         []string
	PathPrefix  string
	HomeURL     string
	LogoutURL   string
	ShowHeaders bool
}

type requestContext struct {
	User string
}

type context struct {
	Static  staticContext
	Request requestContext
	Dynamic any
}

func getPageComponent(options ServerOptions) *template.Template {
	t := template.New("Page")

	pageComponent := component.Component{
		Name:     "page",
		Template: options.Template,
		Children: append([]component.Component{
			header.Get(),
			footer.Get(),
		}, options.Children...),
	}

	pageComponent.Parse(t)
	template.Must(t.Parse(baseContent))

	return t
}

func Run(options ServerOptions) error {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	proxyPrefix := os.Getenv("PROXY_PREFIX")
	hostname := os.Getenv("WEB_HOSTNAME")

	t := getPageComponent(options)

	css := map[string][]byte{}
	js := map[string][]byte{}
	wasm := map[string][]byte{}

	sc := staticContext{
		Title:       options.Title,
		Scripts:     []string{},
		CSS:         []string{},
		PathPrefix:  proxyPrefix,
		ShowHeaders: !options.HideHeaders,
		HomeURL:     "http://" + hostname,
		LogoutURL:   "http://" + hostname + "/logout",
	}

	for _, module := range options.WASMModules {
		// Load the wasm_exec if it's not already there
		if _, ok := js["wasm_exec.js"]; !ok {
			js["wasm_exec.js"] = wasmExec
			sc.Scripts = append(sc.Scripts, "wasm_exec.js")
		}

		// Serve the wasm binary
		wasm[module.Path] = module.Data

		// Parse and serve the specific wasm loader library
		t := template.New("Script")
		template.Must(t.Parse(wasmLoad))

		var script bytes.Buffer
		err := t.Execute(&script, map[string]interface{}{
			"Path": "http://" + hostname + proxyPrefix + "/" + module.Path,
			"Name": module.Name,
		})
		if err != nil {
			return fmt.Errorf("failed to parse wasm load script: %+v", err)
		}

		js[module.Name+"_wasm_load.js"] = script.Bytes()
		sc.Scripts = append(sc.Scripts, module.Name+"_wasm_load.js")
	}

	mux := http.NewServeMux()
	mux.HandleFunc(proxyPrefix, func(w http.ResponseWriter, r *http.Request) {
		zap.S().Infof("HTTP %s %s", r.Method, r.URL.Path)

		if r.Method == "POST" && options.PostHandler != nil {
			options.PostHandler(w, r)
			return
		}

		var data interface{}
		if options.PageHandler != nil {
			data = options.PageHandler(r)
		}

		user := r.Header.Get("X-Auth-User")

		context := context{
			Static: sc,
			Request: requestContext{
				User: user,
			},
			Dynamic: data,
		}

		var response bytes.Buffer

		err := t.Execute(&response, context)
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

	for path, data := range js {
		mux.HandleFunc(proxyPrefix+"/"+path, fileHandler("text/javascript", data))
	}

	for path, data := range wasm {
		mux.HandleFunc(proxyPrefix+"/"+path, fileHandler("application/wasm", data))
	}

	for path, data := range css {
		mux.HandleFunc(proxyPrefix+"/"+path, fileHandler("text/css", data))
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		zap.S().Infof("Not found: %s", r.URL.Path)
		http.NotFound(w, r)
	})

	zap.S().Info("running server...")

	err := http.ListenAndServe("0.0.0.0:80", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		return err
	}

	return nil
}

func fileHandler(contentType string, data []byte) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		zap.S().Infof("Serving file: %s", r.URL.Path)
		w.Header().Set("content-type", contentType)
		w.Write(data)
	}
}

type ServeFile struct {
	Path     string
	Data     []byte
	MimeType string
}

func RunFileServer(files []ServeFile) error {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	proxyPrefix := os.Getenv("PROXY_PREFIX")
	if proxyPrefix == "/" {
		proxyPrefix = ""
	}

	mux := http.NewServeMux()
	for _, file := range files {
		mux.HandleFunc(proxyPrefix+"/"+file.Path, fileHandler(file.MimeType, file.Data))
	}

	zap.S().Info("running server...")

	err := http.ListenAndServe("0.0.0.0:80", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		return err
	}

	return nil
}
