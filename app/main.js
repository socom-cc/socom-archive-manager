async function loadWasm() {
    if (WebAssembly && !WebAssembly.instantiateStreaming) { // polyfill
        WebAssembly.instantiateStreaming = async (resp, importObject) => {
            const source = await (await resp).arrayBuffer();
            return await WebAssembly.instantiate(source, importObject);
        };
    }

    const go = new Go();

    const result = await WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject);
    go.run(result.instance);
}

function attachUploadListeners() {
    const handleFile = (file) => {
        const reader = new FileReader();
        reader.onload = function(e) {
            const arrayBuffer = e.target.result;
            const uint8Array = new Uint8Array(arrayBuffer);
            socom.onload(uint8Array);
        };
        reader.readAsArrayBuffer(file);
    };

    const dragOverListener = (event) => {
        event.preventDefault();
    };

    const dragLeaveListener = (event) => {
        event.preventDefault();
    };

    const dropListener = (event) => {
        event.preventDefault();
        if (event.dataTransfer.files.length > 0) {
            const file = event.dataTransfer.files[0];
            handleFile(file);
        }
    };

    document.addEventListener('dragover', dragOverListener);
    document.addEventListener('dragleave', dragLeaveListener);
    document.addEventListener('drop', dropListener);

    return { dragOverListener, dragLeaveListener, dropListener };
}

(async () => {
    await loadWasm();
    attachUploadListeners();
})()