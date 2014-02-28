
var x = 100,
    y = 100,
    client_color = '#'+Math.floor(Math.random()*16777215).toString(16),
    speed = 1,
    width = 500,
    height = 500;

var left, right, up, down;
left = right = up = down = false;
var edging = false;

var websocket;

function key_down(e) {
    switch (e.keyCode) {
        case 40:
            down = true;
            break;

        case 38:
            up = true;
            break;

        case 37:
            left = true;
            break;

        case 39:
            right = true;
            break;
    }
}

function key_up(e) {
    switch (e.keyCode) {
        case 40:
            down = false;
            break;

        case 38:
            up = false;
            break;

        case 37:
            left = false;
            break;

        case 39:
            right = false;
            break;
    }
}

function update_position() {
    // Respond to user input
    if (down) {
        y += speed;
    }
    if (up) {
        y -= speed;
    }
    if (left) {
        x -= speed;
    }
    if (right) {
        x += speed;
    }

    //Don't let them out of the box
    if (x > width){
        x = width;
    }
    if (x < 0){
        x = 0;
    }
    if (y > height){
        y = height;
    }
    if (y < 0){
        y = 0;
    }

    if (x > 0 && y > 0 && x < width && y < height) {
        edging = false;
    }
    else {
        edging = true;
    }
}

function redraw() {

    var c = document.getElementById("dot_canvas");
    var context = c.getContext("2d");
    context.fillStyle = "#1F8FCE";
    if (edging) {
        context.fillStyle = '#'+Math.floor(Math.random()*16777215).toString(16);
    }
    context.fillRect(0, 0, width, height);

    var radius = 15;
    context.beginPath();
    context.arc(x, y, radius, 0, 2 * Math.PI, false);
    context.fillStyle = 'green';
    context.fill();
    context.lineWidth = 5;
    context.strokeStyle = '#003300';
    context.stroke();
}

function say_stuff() {
    console.log("Sending Client Info");

    client = {
        x: x,
        y: y,
        color: client_color,
    };

    console.log(client);

    websocket.send(JSON.stringify(client));
}

$(function(){
    window.onkeydown = key_down;
    window.onkeyup = key_up;
    setInterval(update_position, 10);
    setInterval(redraw, 10);

    websocket = new WebSocket("ws://localhost:8000/ws");
    websocket.onclose = function(e) {
        console.log('closing');
    };

    websocket.onmessage = function(e) {
        console.log(e.data);
    };

    setInterval(say_stuff, 1000);

});



