var ws = new WebSocket("ws://localhost:8042/websocket");

function sendUserInput() {
    var input = document.getElementById("inputBox").value;
    // Validate input

    ws.send(input)
}

ws.onmessage = function (event) {
    console.log(event.data);
    var message = JSON.parse(event.data);
    console.log(message);

    switch (message.Kind) {
        case "Finished": break;
        case "Song": playVideo(message.Value); break;
        case "Title": appendToQueue(message.Value); break;
    }
};

function appendToQueue(title) {
    var entry = document.createElement("li");
    entry.appendChild(document.createTextNode(title));

    document.getElementById("songsQueue").appendChild(entry);
}

function informCompletion() {
    ws.send("Finished")
}