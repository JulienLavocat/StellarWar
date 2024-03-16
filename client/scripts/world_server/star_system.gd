extends VisibilityNotifier2D

class_name StarSystem

var color := Color.white
var systemOwner: int = -1
var systemId: int
var resources: Dictionary
var neighbours: PoolIntArray
var details: Dictionary

signal on_selected(system)

func init(system) -> void:
	self.systemId = int(system.id)
	self.position = Vector2(system.x, system.y)
	self.name = system.name
	self.systemOwner = system.owner
	self.resources = system.resources
	self.details = system.details
	self.neighbours = []
	for neighborId in system.to:
		self.neighbours.append(int(neighborId))
	
	$Label.text = name
	var area2d = $Area2D
	area2d.connect("mouse_entered", self, "_on_mouse_entered")
	area2d.connect("mouse_exited", self, "_on_mouse_exited")
	area2d.connect("input_event", self, "_on_input_event")
	self.connect("screen_entered", self, "_on_screen_entered")
	self.connect("screen_exited", self, "_on_screen_exited")

func debug():
	return {
		"systemId": self.systemId,
		"position": self.position,
		"name": self.name ,
		"systemOwner": self.systemOwner,
		"resources": self.resources,
		"details": self.details
	}

func set_system_owner(player):
	
	if player == null:
		print(self.systemId, " - ", self.name, " has been freed")
		self.systemOwner = -1
		$Sprite.self_modulate = Color.white
	else:
		print(player.id, " claimed system ", self.systemId, " - ", self.name)
		self.systemOwner = player.id
		$Sprite.self_modulate = player.color

func _on_screen_entered():
#	print(self.name, " screen entered")
	pass
	
func _on_screen_exited():
#	print(self.name, " screen exited")
	pass

func _on_mouse_entered():
#	$Label.visible = true
	pass

func _on_mouse_exited():
#	$Label.visible = false
	pass
	
func _on_input_event(viewport, event, shape_id):
	if (event is InputEventMouseButton && event.pressed):
		self.emit_signal("on_selected", self)

func _ready():
#	$Label.visible = false
	pass # Replace with function body.

