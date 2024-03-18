extends Node2D

class_name StarSystem

@onready var name_label: Label = $Name

func set_data(data: Dictionary) -> void:
	print(data)
	name = data.name
	position = Vector2(data.x, data.y)

func _ready() -> void:
	name_label.text = name
