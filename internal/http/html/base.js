{{define "base.js"}}
"use strict";

function pause() {
    return post("/pause", null);
}

function rewind() {
    return post("/rewind", null);
}

function fastforward() {
    return post("/fastforward", null);
}

function play(result) {
    return post("/queue", result);
}

function playAlbum(results) {
    return post("/queue/album", results);
}

function volume(v) {
    return post("/volume", {volume: v});
}

function post(url, body) {
    var xhr = new XMLHttpRequest();
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-type", "application/json");
    var data;
    if (body != null) {
        data = JSON.stringify(body);
    }
    xhr.send(data);
    return false;
}

window.onload = function () {
    var conn;
    conn = new WebSocket("ws://" + document.location.host + "/ws/play-progress");
    conn.onclose = function (evt) {
        
    };

    var play = document.getElementById("play-bar");
    var download = document.getElementById("download-bar");
    var song = document.getElementById("current-song");
    var album = document.getElementById("current-album");
    var artist = document.getElementById("current-artist");
    conn.onmessage = function (evt) {
	var msg = JSON.parse(evt.data);
	if (msg.type == "play progress") {
	    var progress = msg.value;
	    var w = 100 * ( progress.n/progress.total);
		    play.style.width = w + '%';
	} else if (msg.type == "download progress") {
	    var progress = msg.value;
	    var w = 100 * ( progress.n/progress.total);
	    download.style.width = w + '%';
	} else if (msg.type == "play/pause") {
	    var text = msg.value.playing ? 'pause' : 'play';
	    document.getElementById("playpause").innerHTML = text;
	} else if (msg.type == "next song") {
            var result = msg.value;
	    song.innerHTML = result.track.title;
            album.href = "/albums/" + result.album.id;
            album.innerHTML = result.album.title;
            artist.href = "/artists/" + result.artist.id;
            artist.innerHTML = result.artist.name;
	}
    }	
};
{{end}}

