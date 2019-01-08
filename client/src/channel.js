export default class Channel {
    constructor(conn, name, event, cb) {
        this.conn = conn;
        this.name = name;
        this.event = event;
        this.cb = cb
    }

    bind() {
        let data = {type: "subscribe", channel: this.name, event: this.event};
        this.conn.send(data)
    }
}
