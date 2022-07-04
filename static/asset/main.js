var token = new URLSearchParams(window.location.search).get("token")
var socket = new WebSocket("ws://127.0.0.1:8080/websocket?token="+token);
var queue = [];
var eventActive = false;

var textField = document.querySelector(".textField");

function queueProcessor() {
    if (!eventActive && queue.length > 0) {
        // We have an event, trigger it
        var evt = queue.shift();
        console.log("process event", evt, queue.length)
        handleEvent(evt);
    }
}

// process the queue
setInterval(queueProcessor, 500);

function handleEvent(evt) {
    eventActive = true;
    var event = JSON.parse(evt.data);

    /*if (event.event) {

    }*/

    var new_data = textField.innerHTML = event.content;
    //console.log(new_data);
    textField.innerHTML = new_data;

    /*var animation = */
    anime.timeline({
        loop: false,
    }).add({
        targets: '.textField',
        opacity: [0, 1],
        easing: "easeOutExpo",
        translateX: [100, 0],
        //scale: [0.3, 1],
        duration: 1000
    }).add({
        targets: '.textField',
        opacity: 0,
        duration: 800,
        easing: "easeOutExpo",
        delay: 3000,
        complete: function (anim) {
            eventActive = false;
            console.log("animation ended", anim)
        }
    })

    // animation.seek(0);
    // animation.play();
}
socket.onmessage = function (event) {
    console.log("add event", event)
    queue.push(event);
};
