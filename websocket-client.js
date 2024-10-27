const PORT = "5900";
const webSocket = new WebSocket("ws://localhost:5900");

const base64 = {
    decode: s => Uint8Array.from(atob(s), c => c.charCodeAt(0)),
    encode: b => btoa(String.fromCharCode(...new Uint8Array(b))),
    decodeToString: s => new TextDecoder().decode(base64.decode(s)),
    encodeString: s => base64.encode(new TextEncoder().encode(s)),
};

async function sha1(str) {
    const enc = new TextEncoder();
    const hash = await crypto.subtle.digest('SHA-1', enc.encode(str));

    return Array.from(new Uint8Array(hash))
        .map(v => v.toString(16).padStart(2, '0'))
        .join('');
}

const statusElement = document.createElement("p");
document.getElementById("app").prepend(statusElement);

const messageList = document.getElementById("messages");

statusElement.innerText = "Websocket Client: Connecting ⏳⏳⏳";

webSocket.addEventListener("error", event => {
    statusElement.innerText = "Websocket Client: Error ❌"
    console.log(event);
})

webSocket.addEventListener("open", event => {
    statusElement.innerText = "Websocket Client: Connected 😄"
    console.log(event);
})

webSocket.addEventListener("message", event => {
    const li = document.createElement("li");
    li.innerText = event.data;
    messageList.appendChild(document.createElement("li"));
})

const input = document.getElementById("s-ws-key-input");
const submit = document.getElementById("s-ws-key-submit");
const ws_uuid = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11";

submit.addEventListener("mouseup", async () => {
    const key = input.value;

    if (key.length < 1) {
        return;
    }

    const accept = key + ws_uuid;
    const hash = await sha1(accept);

    const base64Accept = base64.encodeString(hash);
    const li = document.createElement("li");
    li.innerText = `Calculated Sec-WebSocket-Key: ${base64Accept}`;
    messageList.appendChild(li);

})

