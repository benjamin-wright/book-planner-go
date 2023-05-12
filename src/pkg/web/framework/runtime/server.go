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

//go:embed templates/styles.css
var styles []byte

type ServerOptions struct {
	Template    string
	Title       string
	Children    []component.Component
	WASMModules []WASMModule
	Handler     func(r *http.Request) any
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
	ProxyPrefix string
}

type context struct {
	Static  staticContext
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

	t := getPageComponent(options)

	css := map[string][]byte{
		"styles.css": styles,
	}
	js := map[string][]byte{}
	wasm := map[string][]byte{}

	sc := staticContext{
		Title:       options.Title,
		Scripts:     []string{},
		CSS:         []string{"styles.css"},
		ProxyPrefix: proxyPrefix,
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
			"Path": proxyPrefix + "/" + module.Path,
			"Name": module.Name,
		})
		if err != nil {
			return fmt.Errorf("failed to parse wasm load script: %+v", err)
		}

		js[module.Name+"_wasm_load.js"] = script.Bytes()
		sc.Scripts = append(sc.Scripts, module.Name+"_wasm_load.js")
	}

	mux := http.NewServeMux()
	mux.HandleFunc(sc.ProxyPrefix, func(w http.ResponseWriter, r *http.Request) {
		zap.S().Infof("HTTP %s %s", r.Method, r.URL.Path)

		data := options.Handler(r)

		context := context{
			Static:  sc,
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
		mux.HandleFunc(sc.ProxyPrefix+"/"+path, fileHandler("text/javascript", data))
	}

	for path, data := range wasm {
		mux.HandleFunc(sc.ProxyPrefix+"/"+path, fileHandler("application/wasm", data))
	}

	for path, data := range css {
		mux.HandleFunc(sc.ProxyPrefix+"/"+path, fileHandler("text/css", data))
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

func fileHandler(contentType string, data []byte) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", contentType)
		w.Write(data)
	}
}
