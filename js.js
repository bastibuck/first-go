const socket = new WebSocket("ws://localhost:8081/ws");

socket.onopen = function (event) {
  // Handle connection open
};

socket.onmessage = function (event) {
  // Handle received message
  const message = JSON.parse(event.data);

  if (message.type === "event_signup") {
    const messageElement = document.createElement("div");
    messageElement.innerText = `New pax ${message.payload.new_pax} for event id: ${message.payload.event_id}`;
    document.body.appendChild(messageElement);
  }
};

socket.onclose = function (event) {
  // Handle connection close
};

function sendMessage(message) {
  socket.send(message);
}
