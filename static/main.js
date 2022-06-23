var socket = new WebSocket("ws://127.0.0.1:8080/websocket?token=someToken");
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

function handleEvent(event) {
    eventActive = true;
    //var new_data = textField.innerHTML = event.data;
    var new_data = event.data.replace(/./g, "<span class='letter'>$&</span>");
    //console.log(new_data);
    textField.innerHTML = new_data;

    /*var animation = */
    anime.timeline({
        loop: false,
    }).add({
        targets: '.textField',
        opacity: 1,
        // easing: "easeInExpo",
        // duration: 100,
    }).add({
        targets: '.textField .letter',
        opacity: [0,1],
        scale: [0.3, 1],
        rotateZ: [180, 0],
        duration: 850,
        easing: "easeOutExpo",
        delay: (el, i) => 20 * (i+1)
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
