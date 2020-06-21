window.onUpdate = function() {};

var socket = new WebSocket(`ws://${location.host}/ws`);
setInterval(function() {
    socket.send('update');
}, 1000);
socket.onmessage = function(data) {
    let res = JSON.parse(data.data);
    if (res){
        window.onUpdate(res);
    }
};

castles=[
    "Castle", "Rampart", "Tower", "Inferno", "Necropolis", "Dungeon", "Stronghold", "Fortress", "Conflux", "Cove", "RandomTown"
];
colors = ["red", "blue", "biege", "green", "orange", "purple", "teal", "pink"];