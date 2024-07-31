//go:build js && wasm

package main

import (
	"bytes"
	"embed"
	"fmt"
	socomarchive "github.com/socom-cc/socom-archive-manager/pkg"
	"html/template"
	"syscall/js"
)

var done chan struct{}

var archive *socomarchive.SocomArchive
var templates *template.Template

//go:embed templates/*
var templateFolder embed.FS

func loadArchive(this js.Value, p []js.Value) interface{} {
	var err error
	romBytes := p[0]
	length := romBytes.Get("length").Int()
	data := make([]byte, length)
	js.CopyBytesToGo(data, romBytes)
	archive, err = socomarchive.LoadSocomArchive(data)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// TODO reload components

	var h bytes.Buffer
	err = templates.ExecuteTemplate(&h, "components/entry-list", archive.EntryHeaders)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	entryListDiv := js.Global().Get("document").Call("getElementById", "entry-list")
	entryListDiv.Set("innerHTML", h.String())
	return nil
}

func main() {
	templates = template.Must(template.ParseFS(templateFolder, "templates/*.go.tmpl"))
	templates = template.Must(templates.ParseFS(templateFolder, "templates/components/*.go.tmpl"))
	templates = template.Must(templates.ParseFS(templateFolder, "templates/pages/*.go.tmpl"))

	appDiv := js.Global().Get("document").Call("getElementById", "app")

	var doc bytes.Buffer
	err := templates.ExecuteTemplate(&doc, "pages/home", nil)
	if err != nil {
		panic(err)
	}

	appDiv.Set("innerHTML", doc.String())

	socomObj := js.Global().Get("Object").New()
	socomObj.Set("onload", js.FuncOf(loadArchive))
	js.Global().Set("socom", socomObj)

	fmt.Println("Loaded WASM")
	<-done
}
