extends Camera2D

var _velocity := Vector2.ZERO
var camera_speed := 100

func _ready() -> void:
	RenderingServer.set_default_clear_color(Color.BLACK)

func _process(delta: float) -> void:
	self.position += Vector2(Input.get_axis("ui_left", "ui_right"), Input.get_axis("ui_up", "ui_down")).normalized() * camera_speed * delta
