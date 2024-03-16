extends PopupPanel

var StarSystem := preload("res://scenes/star_system.tscn")
var Galaxy := preload("res://scripts/world_server/galaxy.gd")

onready var galaxy: Galaxy = $"../../Galaxy"
onready var resources = $Resources
onready var claimButton = $"VBoxContainer/Control/ClaimSystemButton"
onready var systemResources = $'VBoxContainer/Resources'
onready var systemName = $'VBoxContainer/VBoxContainer/SystemName'

var resourcesIndex = {
	"minerals": 0,
	"energy": 1,
	"food": 2
}

var selectedSystem: StarSystem = null

func _ready():
	galaxy.connect("system_selected", self, "_on_system_selected")
	galaxy.connect("resources_quote", self, "_on_resources_quote")
	
	claimButton.connect("pressed", self, "_claim_system")
	
	self.visible = false
	
func _input(event):
	if (event is InputEventMouseButton) and event.pressed:
		var evLocal = self.make_input_local(event)
		if !Rect2(Vector2(0,0),rect_size).has_point(evLocal.position):
			self.visible = false
		
func _on_system_selected(system: StarSystem):
	self.visible = true
	self.selectedSystem = system
	
	self.popup()
	
# 	Clear old nodes
	for resource in system.resources:
		(systemResources.get_child(resourcesIndex[resource]) as Label).text = resource + ": " + str(system.resources[resource])

	systemName.text = system.name
	update_ui()

func _on_resources_quote(quote):
	update_ui()

func _claim_system():
	galaxy.claim_system(selectedSystem.systemId)
	
func update_ui():
	var hasSufficientFunds = galaxy.has_sufficient_funds("outpost")
	
	var hasNeighbor = false
	if selectedSystem != null:
		for neighborId in selectedSystem.neighbours:
			if galaxy._systems[neighborId].systemOwner == galaxy._playerId:
				hasNeighbor = true
				break
	
	claimButton.disabled = (selectedSystem != null && selectedSystem.systemOwner > 0) || !hasSufficientFunds || !hasNeighbor
	_setClaimTooltip(hasSufficientFunds, hasNeighbor)
	
func _setClaimTooltip(hasSufficientFunds: bool, hasNeighbor: bool):
	var tooltip = "Costs: \n"
	var costs = Costs.getCosts("outpost")
	for resource in costs:
		tooltip += "  - %s: %s\n" % [resource, costs[resource]]
		
	if !hasSufficientFunds:
		tooltip += "\nNot enough resources"
	
	if !hasNeighbor:
		tooltip += "\nNo neighboring systems"
		
	claimButton.hint_tooltip = tooltip
