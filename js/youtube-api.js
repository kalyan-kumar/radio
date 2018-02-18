var player;

function onYouTubeIframeAPIReady() {
    player = new YT.Player('player', {
        height: '480',
        width: '854',
        videoId: 'GDQob4AOCsQ',
        events: {
            onReady: onPlayerReady,
            onStateChange: onPlayerStateChange
        }
    });
}

function onPlayerReady(event) {
    console.log("onPlayerReady event -");
    console.log(event.target);
    event.target.playVideo();
}

var done = false;
function onPlayerStateChange(event) {
    console.log("onPlayerStateChange event - ");
    console.log(event.target);
    if (event.data === YT.PlayerState.PLAYING && !done) {
        setTimeout(stopVideo, 6000);
        done = true;
    }
}
function stopVideo() {
    player.stopVideo();
}