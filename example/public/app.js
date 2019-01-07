$(function () {
    var conn
    var msg = $("#msg")
    var log = $("#log")

    function getURL() {
        element = document.head.querySelector("meta[name='thunderbird-url']");
        return element.getAttribute("content")
    }

    function appendLog(msg) {
        var d = log[0]
        var doScroll = d.scrollTop == d.scrollHeight - d.clientHeight
        msg.appendTo(log)
        if (doScroll) {
            d.scrollTop = d.scrollHeight - d.clientHeight
        }
    }

    $("#form,#send").submit(function (evt) {
        if (!conn) {
            return false;
        }
        if (!msg.val()) {
            return false;
        }

        conn.perform("room", "abc", msg.val());
        msg.val("");

        return false
    });

    $("#form,#unsubscribe").submit(function (evt) {
        if (!conn) {
            return false;
        }

        conn.unsubscribe("room", "abc");
        return false
    });

    conn = Thunderbird.connect(getURL(), function (conn) {
        conn.subscribe("room", "abc", function (msg) {
            appendLog($("<div/>").text(msg))
        })
    })

    //if (window["WebSocket"]) {
    //new Thunderbird("");
    //conn = new WebSocket();
    //conn.onclose = function(evt) {
    //appendLog($("<div><b>Connection closed.</b></div>"))
    //}
    //conn.onmessage = function(evt) {
    //appendLog($("<div/>").text(evt.data))
    //}
    //} else {
    //appendLog($("<div><b>Your browser does not support WebSockets.</b></div>"))
    //}
});
