extends Node2D

class_name GalaxyClient

var _client = WebSocketClient.new()

#var url := "wss://world-server-prod.herokuapp.com/"
#var url := "https://world-server-dev.herokuapp.com/"
var url := "ws://localhost:5000"

signal player_joined(data)
signal player_left(data)
signal sync_client(data)
signal system_claimed(data)
signal resources_quote(data)

func _ready():
	_client.connect("connection_established", self, "_connection_established")
	_client.connect("connection_closed", self, "_connection_closed")
	_client.connect("connection_error", self, "_connection_closed")
	_client.connect("data_received", self, "_on_data")

func _process(delta):
	_client.poll()

func send(packet):
	print("send: ", packet)
	_client.get_peer(1).put_packet(JSON.print(packet).to_utf8())

func _connection_established(proto = ""):
	print("Connected ", proto)
	send({
		"command": "set_ready"
	})
	
func _connection_closed(was_clean = false):
	print("Closed clean: ", was_clean)

func _on_data():
	var payload = parse_json(_client.get_peer(1).get_packet().get_string_from_utf8())
	emit_signal(payload.command, payload.data)

func start():	
	var err = _client.connect_to_url(url)
	if err != OK:
		print("Unable to connect")
		set_process(false)

func stop():
	_client.disconnect_from_host()
