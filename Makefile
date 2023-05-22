OUTDIR ?= _output

.PHONY: build-cli
build-cli: $(OUTDIR)/kept-linux-$(GOARCH).sha256
build-cli: $(OUTDIR)/kept-darwin-$(GOARCH).sha256

.PHONY: install-cli
install-cli:
	go install -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) ./main.go

.PHONY: kept-linux-$(GOARCH)
kept-linux-$(GOARCH): $(OUTDIR)/kept-linux-$(GOARCH)
$(OUTDIR)/kept-linux-$(GOARCH): $(SOURCES)
	CGO_ENABLED=0 GOARCH=$(GOARCH) GOOS=linux go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ main.go
	upx $@

$(OUTDIR)/kept-linux-$(GOARCH).sha256: $(SOURCES) $(OUTDIR)/kept-linux-$(GOARCH)
	shasum -a 256 $(OUTDIR)/kept-linux-$(GOARCH) > $(OUTDIR)/kept-linux-$(GOARCH).sha256

.PHONY: kept-darwin-$(GOARCH)
kept-darwin-$(GOARCH): $(OUTDIR)/kept-darwin-$(GOARCH)
$(OUTDIR)/kept-darwin-$(GOARCH): $(SOURCES)
	CGO_ENABLED=0 GOARCH=$(GOARCH) GOOS=darwin go build -ldflags=$(LDFLAGS) -gcflags=$(GCFLAGS) -o $@ main.go

$(OUTDIR)/kept-darwin-$(GOARCH).sha256: $(SOURCES) $(OUTDIR)/kept-darwin-$(GOARCH)
	shasum -a 256 $(OUTDIR)/kept-darwin-$(GOARCH) > $(OUTDIR)/kept-darwin-$(GOARCH).sha256
