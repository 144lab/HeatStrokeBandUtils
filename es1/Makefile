
GO := go1.12.17
DIST=dist

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

all: clean dist

clean:
	@-rm -rf $(DIST)

dist:
	mkdir -p $(DIST)
	$(GO) generate ./...
	GOOS=js GOARCH=wasm $(GO) build -o $(DIST)/main.wasm .
	cp -f $(shell $(GO) env GOROOT)/misc/wasm/wasm_exec.js $(DIST)/
	cp -Rf assets $(DIST)/
	echo "$$index_html" > $(DIST)/index.html
