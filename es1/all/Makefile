SRCDIR := $(shell pwd)
DSTDIR := dist
GO := go

export index_html
override define index_html
<!DOCTYPE html>
<head>
<link rel="manifest" href="/assets/manifest.json">
</head>
<body>
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
</body>
endef

.PHONY: all clean run dist

all: clean dist

run: clean
	wecty server

clean:
	$(foreach f,wasm_exec.js main.wasm index.html,$(shell if [ -e $f ]; then rm $f; fi))
	@- rm -rf $(DSTDIR)

dist:
	@echo "deploy: $(SRCDIR) -> $(DSTDIR)"
	mkdir -p $(DSTDIR)
	GO111MODULE=on GOOS=js GOARCH=wasm $(GO) build -o $(DSTDIR)/main.wasm .
	echo "$$index_html" > $(DSTDIR)/index.html
	cp $(shell $(GO) env GOROOT)/misc/wasm/wasm_exec.js $(DSTDIR)/wasm_exec.js
	cp -Rf assets $(DSTDIR)/
