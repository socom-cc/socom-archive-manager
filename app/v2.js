function verifyVersion() {
    return socom.version === "v2";
}

async function loadFile(f, root) {
    const path = await root.resolve(f);
    const key = path.join('/');
    const file = await f.getFile();
    const reader = new FileReader();
    reader.onload = function(e) {
        const arrayBuffer = e.target.result;
        const uint8Array = new Uint8Array(arrayBuffer);
        socom.loadArchive(key, uint8Array);
    };
    reader.readAsArrayBuffer(file);
}

async function loadDirectory(d, root) {
    for await (const [k, e] of d.entries()) {
      console.log(k , e);
        switch (e.kind) {
            case "directory":
                await loadDirectory(e, root);
                break;
            case "file":
                if (k.toLowerCase().endsWith(".zdb")) {
                    await loadFile(e, root);
                }
                break;
            default:
                break;
        }
    }
}

async function openDirectory() {
    if (!verifyVersion()) return;

    const dir = await window.showDirectoryPicker();
    console.log(dir);

    await loadDirectory(dir, dir);
}

(async () => {
    socom.download = (idx) => {
        const ret = socom.getBytes(idx);

        const blob = new Blob([ret.data], { type: 'application/octet-stream' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `${ret.name}.bin`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
    };
})();
