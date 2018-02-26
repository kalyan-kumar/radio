var url = window.location.href.replace(/^http:\/\//, '').replace(/radio$/, 'websocket');

var ws = new WebSocket("ws://" + url);

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
        case "Finished":
            break;
        case "Song":
            playVideo(message.Value);
            break;
        case "Title":
            appendToQueue(message.Value);
            break;
    }
};

function appendToQueue(title) {
    var entry = document.createElement("li");
    entry.appendChild(document.createTextNode(title));

    if (YT.PlayerState.PLAYING !== player.getPlayerState()) {
        playVideo(title);
    }

    document.getElementById("songsQueue").appendChild(entry);
}

function informCompletion() {
    ws.send("Finished")
}