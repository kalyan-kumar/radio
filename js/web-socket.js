var ws = new WebSocket("ws://localhost:8042/websocket");

function sendUserInput() {
    var input = document.getElementById("inputBox").value;
    // Validate input

    ws.send(input)
}

ws.onmessage = function (event) {
    var entry = document.createElement("li");
    entry.appendChild(document.createTextNode(event.data));

    document.getElementById("songsQueue").appendChild(entry);
};