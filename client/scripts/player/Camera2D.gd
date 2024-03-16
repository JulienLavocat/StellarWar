extends Camera2D

var velocity := Vector2.ZERO
var camera_speed := 100

func _ready():
	VisualServer.set_default_clear_color(Color.black)

func _process(delta):
	velocity.x = Input.get_axis("ui_left", "ui_right")
	velocity.y = Input.get_axis("ui_up", "ui_down")
	self.position += velocity.normalized() * delta * camera_speed

