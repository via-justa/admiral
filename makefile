VERSION=

DOT:= .
DASH:= _
# replace . with _
ver = $(subst $(DOT),$(DASH),$(VERSION))

.PHONY: build 

build:
	@GOOS=windows GOORCH=386 go build -v -ldflags="-s -w -X 'main.appVersion=$(VERSION)'" -o admiral_windows_386_$(ver).exe \
	&& upx admiral_windows_386_$(ver).exe \
	&& chmod +x admiral_windows_386_$(ver).exe
	@GOOS=linux GOORCH=amd64 go build -v -ldflags="-s -w -X 'main.appVersion=$(VERSION)'" -o admiral_linux_amd64_$(ver) \
	&& upx admiral_linux_amd64_$(ver) \
	&& chmod +x admiral_linux_amd64_$(ver)
	@GOOS=darwin GOORCH=386 go build -v -ldflags="-s -w -X 'main.appVersion=$(VERSION)'" -o admiral_darwin_386_$(ver) \
	&& upx admiral_darwin_386_$(ver) \
	&& chmod +x admiral_linux_amd64_$(ver)
%:
	@:
