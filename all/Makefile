SRCDIR := $(shell pwd)
DSTDIR := dist

export index_html
override define index_html
<!DOCTYPE html>
<script src="wasm_exec.js"></script><script>
(async () => {
	const resp = await fetch('main.wasm');
	if (!resp.ok) {
		const pre = document.createElement('pre');
		pre.innerText = await resp.text();
		document.body.appendChild(pre);
		return;
	}
	const src = await resp.arrayBuffer();
	const go = new Go();
	const result = await WebAssembly.instantiate(src, go.importObject);
	go.run(result.instance);
})();
</script>
endef

run: clean
	wasmserve

clean:
	rm wasm_exec.js main.wasm index.html

update: deploy
	cp $(DSTDIR)/wasm_exec.js ./
	cp $(DSTDIR)/main.wasm ./
	cp $(DSTDIR)/index.html ./

deploy:
	@echo "deploy: $(SRCDIR) -> $(DSTDIR)"
	mkdir -p $(DSTDIR)
	GO111MODULE=on GOOS=js GOARCH=wasm go build -o $(DSTDIR)/main.wasm .
	echo "$$index_html" > $(DSTDIR)/index.html
	cp $(shell go env GOROOT)/misc/wasm/wasm_exec.js $(DSTDIR)/wasm_exec.js
	cp recorder.js $(DSTDIR)/recorder.js
	cp nosleep.min.js $(DSTDIR)/nosleep.min.js
	cp -Rf css $(DSTDIR)/

gopherjs:
	GOPHERJS_GOROOT=$(shell go1.12.17 env GOROOT) GOPATH=$(shell cd ../..;pwd):$(HOME)/go gopherjs build -o main.js .

gopherjs-run:
	GOPHERJS_GOROOT=$(shell go1.12.17 env GOROOT) GOPATH=$(shell cd ../..;pwd):$(HOME)/go \
	gopherjs serve hrm-logger
