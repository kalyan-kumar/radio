var player;

function onYouTubeIframeAPIReady() {
    player = new YT.Player('player', {
        height: '480',
        width: '854',
        videoId: '{{.Id}}',
        playerVars: {
            start: '{{.Position}}',
            controls: 0,
            disablekb: 1,
            rel: 0
        },
        events: {
            onReady: onPlayerReady,
            onStateChange: onPlayerStateChange
        }
    });
}

function onPlayerReady(event) {
    console.log("Player is Ready");
    event.target.playVideo();
}

function onPlayerStateChange(event) {
    console.log("Player state - " + event.data);
    if (event.data === YT.PlayerState.ENDED) {
        informCompletion();
    }
}

function playVideo(Id) {
    console.log("Will play song - " + Id);
    player.loadVideoById({videoId: Id});
}

function stopVideo(player) {
    console.log("Player being stopped");
    player.stopVideo();
}