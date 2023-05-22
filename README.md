# Book Planner

A golang-based microservices app for planning novels and books. A test bed for trying:
- traefik forward auth
- custom operator-driven database and migration management
- golang-based microservices monorepo
- golang-based MVC using html/template lib to minimise javascript payload
- maybe: golang-based wasm front-end modules using tinygo and generators
- k3d and tilt-driven development environment

Prime directive is minimal resource footprint in terms of development overhead and runtime overhead

## Gotchas

While tinygo are sorting out their GC issues:

```go
//go:generate sh -c "cp ${DOLLAR}(tinygo env TINYGOROOT)/targets/wasm_exec.js templates/wasm_exec.js"
```

then swap in:
```js
// func finalizeRef(v ref)
"syscall/js.finalizeRef": (v_addr) => {
    // Note: TinyGo does not support finalizers so this is only called
    // for one specific case, by js.go:jsString.
    const id = mem().getUint32(v_addr, true);
    this._goRefCounts[id]--;
    if (this._goRefCounts[id] === 0) {
        const v = this._values[id];
        this._values[id] = null;
        this._ids.delete(v);
        this._idPool.push(id);
    }
},
```

## TODO

- swap out pgx.conn for a connection pool (support restarts, etc)
- make services idempotent
- event-driven-architecture with nats
- request-response paradigm but trap the response somehow...
- nextjs micro-frontend (ref: https://blog.logrocket.com/build-monorepo-next-js/#benefits-monorepo-next-js)
