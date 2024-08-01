//go:build js && wasm

package main

import (
	"bytes"
	"embed"
	"fmt"
	socomarchive "github.com/socom-cc/socom-archive-manager/pkg"
	"html/template"
	"slices"
	"strings"
	"syscall/js"
)

var done chan struct{}

// Variables for V1
var archive *socomarchive.SocomArchive

// Variables for V2
var activeKey string
var archives map[string]*socomarchive.SocomArchive

var templates *template.Template

//go:embed templates/*
var templateFolder embed.FS

func loadArchiveV1(this js.Value, p []js.Value) interface{} {
	var err error
	dataBytes := p[0]
	length := dataBytes.Get("length").Int()
	data := make([]byte, length)
	js.CopyBytesToGo(data, dataBytes)
	archive, err = socomarchive.LoadSocomArchive(data)
	if err != nil {
		fmt.Println(err)
		return nil
	}

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

func loadArchiveV2(this js.Value, p []js.Value) interface{} {
	var err error
	key := p[0].String()
	key = strings.ToLower(key)
	dataBytes := p[1]
	length := dataBytes.Get("length").Int()
	data := make([]byte, length)
	js.CopyBytesToGo(data, dataBytes)
	archive, err = socomarchive.LoadSocomArchive(data)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	archives[key] = archive

	var keys []string
	for k := range archives {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	var doc bytes.Buffer
	err = templates.ExecuteTemplate(&doc, "components/v2/sidebar", keys)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	sidebar := js.Global().Get("document").Call("getElementById", "sidebar")
	sidebar.Set("innerHTML", doc.String())

	return nil
}

func openArchive(this js.Value, p []js.Value) interface{} {
	if p[0].Type() != js.TypeString {
		fmt.Println("arg 0 not string")
		return nil
	}

	key := p[0].String()
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, "\\", "/")

	a, ok := archives[key]
	if !ok {
		fmt.Println("invalid key", key)
		return nil
	}

	activeKey = key

	var h bytes.Buffer
	err := templates.ExecuteTemplate(&h, "components/entry-list-v2", a.EntryHeaders)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	entryListDiv := js.Global().Get("document").Call("getElementById", "entry-list")
	entryListDiv.Set("innerHTML", h.String())
	return nil
}

func getBytesV1(this js.Value, p []js.Value) interface{} {
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

func getBytesV2(this js.Value, p []js.Value) interface{} {
	if p[0].Type() != js.TypeNumber {
		return nil
	}

	idx := p[0].Int()

	a, ok := archives[activeKey]
	if !ok {
		fmt.Println("invalid key")
		return nil
	}

	data, err := a.GetEntry(idx)
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

func runV1(this js.Value, p []js.Value) interface{} {
	archive = nil
	appDiv := js.Global().Get("document").Call("getElementById", "app")

	var doc bytes.Buffer
	err := templates.ExecuteTemplate(&doc, "pages/v1", nil)
	if err != nil {
		panic(err)
	}

	socomObj := js.Global().Get("Object").New()
	socomObj.Set("version", "v2")
	socomObj.Set("loadArchive", js.FuncOf(loadArchiveV1))
	socomObj.Set("getBytes", js.FuncOf(getBytesV1))
	js.Global().Set("socom", socomObj)

	appDiv.Set("innerHTML", doc.String())
	script := js.Global().Get("document").Call("createElement", "script")
	script.Set("src", "v1.js")
	appDiv.Call("appendChild", script)

	fmt.Println("Loaded WASM V1")

	return nil
}

func runV2(this js.Value, p []js.Value) interface{} {
	archives = make(map[string]*socomarchive.SocomArchive)
	appDiv := js.Global().Get("document").Call("getElementById", "app")

	var doc bytes.Buffer
	err := templates.ExecuteTemplate(&doc, "pages/v2", nil)
	if err != nil {
		panic(err)
	}

	socomObj := js.Global().Get("Object").New()
	socomObj.Set("version", "v2")
	socomObj.Set("getBytes", js.FuncOf(getBytesV2))
	socomObj.Set("loadArchive", js.FuncOf(loadArchiveV2))
	socomObj.Set("openArchive", js.FuncOf(openArchive))
	js.Global().Set("socom", socomObj)

	appDiv.Set("innerHTML", doc.String())
	script := js.Global().Get("document").Call("createElement", "script")
	script.Set("src", "v2.js")
	appDiv.Call("appendChild", script)

	fmt.Println("Loaded WASM V2")

	return nil
}

func main() {
	templates = template.Must(template.ParseFS(templateFolder, "templates/*.go.tmpl"))
	templates = template.Must(templates.ParseFS(templateFolder, "templates/components/*.go.tmpl"))
	templates = template.Must(templates.ParseFS(templateFolder, "templates/components/v2/*.go.tmpl"))
	templates = template.Must(templates.ParseFS(templateFolder, "templates/pages/*.go.tmpl"))

	runV2(js.Value{}, nil)

	js.Global().Set("runV2", js.FuncOf(runV2))
	js.Global().Set("runV1", js.FuncOf(runV1))

	<-done
}
