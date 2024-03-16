extends VBoxContainer

var resourcesIndex = {
	"minerals": 0,
	"energy": 1,
	"food": 2
}

onready var galaxy: Galaxy = $"../../Galaxy"

func _ready():
	galaxy.connect("resources_quote", self, "_on_resource_quote")
	
func _on_resource_quote(quote: Dictionary):
	for resource in quote:
		(get_child(resourcesIndex[resource]) as Label).text = resource + ": " + str(quote[resource])
