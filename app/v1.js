function attachUploadListeners() {
    const handleFile = (file) => {
        const reader = new FileReader();
        reader.onload = function(e) {
            const arrayBuffer = e.target.result;
            const uint8Array = new Uint8Array(arrayBuffer);
            socom.loadArchive(uint8Array);
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
    attachUploadListeners();
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