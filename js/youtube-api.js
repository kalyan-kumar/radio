function onYouTubeIframeAPIReady() {
    new YT.Player('player', {
        height: '480',
        width: '854',
        videoId: '{{.ToPlay.Id}}',
        playerVars: {
            start: '{{.ToPlay.Position}}',
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

    }
}
function stopVideo(player) {
    console.log("Player being stopped");
    player.stopVideo();
}