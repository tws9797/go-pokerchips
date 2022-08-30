var app = new Vue({
  el: '#app',
  data: {
    ws: null,
    serverUrl: "ws://localhost:8080/ws",
    roomInput: null,
    room: null,
    user: {
      name: ""
    },
    users: []
  },
  methods: {
    createRoom() {
      axios.post("http://localhost:8080/api/room", {
        name: this.user.name
      }).then(res => {
        this.room = {
          id: res.data.data.id,
          name: res.data.data.uri,
          messages: []
        }
        this.connectToWebsocket()
      })
    },
    connectToWebsocket() {
      this.ws = new WebSocket(this.serverUrl + "?name=" + this.user.name + "&room=" + this.room.id);
      this.ws.addEventListener('open', (event) => { this.onWebsocketOpen(event) });
      this.ws.addEventListener('message', (event) => { this.handleNewMessage(event) });
    },
    onWebsocketOpen() {
      console.log("connected to WS!");
    },
    handleNewMessage(event) {

      let data = event.data;
      data = data.split(/\r?\n/);

      console.log(data)

      for (let i = 0; i < data.length; i++) {
        let msg = JSON.parse(data[i]);

        // display the message in the correct room.
        const room = this.room;
        if (typeof room !== "undefined") {
          room.messages.push(msg);
        }
      }
    },
    sendMessage(room) {
      console.log(room)
      // send message to correct room.
      if (room.newMessage !== "") {
        this.ws.send(JSON.stringify({
          action: 'send-message',
          message: room.newMessage,
          target:  room.name
        }));
        room.newMessage = "";
      }
    },
    joinRoom() {
      this.ws.send(JSON.stringify({ action: 'join-room', message: this.roomInput }));
      this.messages = [];
      this.room = ({ "name": this.roomInput, "messages": [] });
      this.roomInput = "";
    },
    leaveRoom(room) {
      this.ws.send(JSON.stringify({ action: 'leave-room', message: room.name }));
    }
  }
})