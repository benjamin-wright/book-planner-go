const go = new Go();
const WASM_URL = "{{ .Path }}";

var {{ .Name }};

if ('instantiateStreaming' in WebAssembly) {
    WebAssembly.instantiateStreaming(fetch(WASM_URL), go.importObject).then(function (obj) {
        {{ .Name }} = obj.instance;
        go.run({{ .Name }});
    })
} else {
    fetch(WASM_URL).then(resp =>
        resp.arrayBuffer()
    ).then(bytes =>
        WebAssembly.instantiate(bytes, go.importObject).then(function (obj) {
            {{ .Name }} = obj.instance;
            go.run({{ .Name }});
        })
    )
}