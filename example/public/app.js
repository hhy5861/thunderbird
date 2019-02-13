$(function () {
    let conn;
    let msg = $("#msg");
    let log = $("#log");
    let event = $("#event");
    let channel = $("#channel");
    let defaultEventObj = $("#default_event");
    let defaultChannelObj = $("#default_channel");

    function getURL() {
        element = document.head.querySelector("meta[name='thunderbird-url']");
        return element.getAttribute("content")
    }

    function appendLog(msg) {
        let d = log[0]
        let doScroll = d.scrollTop == d.scrollHeight - d.clientHeight
        msg.appendTo(log)
        if (doScroll) {
            d.scrollTop = d.scrollHeight - d.clientHeight
        }
    }

    $("#send").on("click", function () {
        if (!conn) {
            return false;
        }
        if (!msg.val()) {
            return false;
        }

        console.log("msg: ", msg.val());
        conn.perform(defaultChannelObj.val(), defaultEventObj.val(), msg.val());
        msg.val("");

        return false
    });

    $("#unsubscribe").on("click", function () {
        if (!conn) {
            return false;
        }

        conn.unsubscribe(defaultChannelObj.val(), defaultEventObj.val());
        return false
    });

    $("#subscribe").on("click", function () {
        if (!conn) {
            return false;
        }

        if (!channel.val()) {
            return
        }

        if (!event.val()) {
            return
        }

        defaultEventObj.val(event.val());
        defaultChannelObj.val(channel.val());
        conn.subscribe(channel.val(), event.val());
        return false
    });

    conn = Thunderbird.connect(getURL(), function (conn) {
        conn.bind("ACCOUNT_BALANCE", "user_id_1", function (msg) {
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
