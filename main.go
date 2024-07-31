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

func getBytes(this js.Value, p []js.Value) interface{} {
	if p[0].Type() != js.TypeNumber {
		return nil
	}

	idx := p[0].Int()

	data, err := archive.GetEntry(idx)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	dataBytes := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(dataBytes, data)

	retObj := js.Global().Get("Object").New()
	retObj.Set("data", dataBytes)
	retObj.Set("name", archive.EntryHeaders[idx].Name)

	return retObj
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
	script := js.Global().Get("document").Call("createElement", "script")
	script.Set("src", "v1.js")
	appDiv.Call("appendChild", script)

	socomObj := js.Global().Get("Object").New()
	socomObj.Set("onload", js.FuncOf(loadArchive))
	socomObj.Set("getBytes", js.FuncOf(getBytes))
	js.Global().Set("socom", socomObj)

	fmt.Println("Loaded WASM")
	<-done
}
