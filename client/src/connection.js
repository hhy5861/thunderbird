import Events from "events"
import Channel from "./channel"

export default class Connection extends Events.EventEmitter {
    constructor(url, cb) {
        super();

        this.url = url;
        this.cb = cb;
        this.connect();
    }

    connect() {
        let self = this;

        this.ws = new WebSocket(this.url);
        this.ws.onopen = function (evt) {
            self.cb(self)
        };

        this.ws.onclose = function (evt) {
        };

        this.ws.onerror = function (evt) {
            let data = JSON.parse(evt.data);
            console.log('error', data)
        };

        this.ws.onmessage = function (evt) {
            let data = JSON.parse(evt.data);
            if (data.type === "message") {
                self.emit(data.channel, JSON.parse(data.body))
            }
        }
    }

    send(data) {
        this.ws.send(JSON.stringify(data))
    }

    bind(channel, event, cb) {
        let ch = new Channel(this, channel, event, cb);
        ch.bind();
        this.on(channel, cb);

        return ch
    }

    subscribe(channel, event) {
        let data = {type: "subscribe", channel: channel, event: event};
        this.ws.send(JSON.stringify(data))
    }

    unsubscribe(channel, event) {
        let data = {type: "unsubscribe", channel: channel, event: event};
        this.ws.send(JSON.stringify(data))
    }

    trigger(channel, event, msg) {
        let data = {type: "trigger", channel: channel, event: event, body: JSON.stringify(msg)};
        this.ws.send(JSON.stringify(data))
    }

    perform(channel, event, msg) {
        let data = {type: "message", channel: channel, event: event, body: JSON.stringify(msg)};
        this.ws.send(JSON.stringify(data))
    }
}
