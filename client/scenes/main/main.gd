extends Node

@onready var star_system: PackedScene = preload("res://entities/star_system/star_system.tscn")

func _ready() -> void:
	($HTTPRequest as HTTPRequest).request_completed.connect(_on_data_loaded)
	($HTTPRequest as HTTPRequest).request("http://localhost:3000/world")
	
func _on_data_loaded(result: int, response_code: int, headers: PackedStringArray, body: PackedByteArray) -> void:
	var json: Dictionary = JSON.parse_string(body.get_string_from_utf8())
	for systemData: Dictionary in json.systems:
		var system: StarSystem = star_system.instantiate()
		system.set_data(systemData)
		add_child(system)
	
	for route: PackedInt32Array in json.routes:
		print(route)
