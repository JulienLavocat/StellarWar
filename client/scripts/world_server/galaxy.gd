extends Node2D

const PlayerUtils = preload("res://scripts/utils/player_utils.gd")
const Costs = preload("res://scripts/costs.gd")

var _star_system = preload("res://scenes/star_system.tscn")

class_name Galaxy

var lines := []
var _systems := {}
var _players := {}
var _playerId: int = -1
var currentResources := {}

signal system_selected(system)
signal resources_quote(quote)

func _ready():
	WorldServer.connect("player_joined", self, "_player_joined")
	WorldServer.connect("player_left", self, "_player_left")
	WorldServer.connect("sync_client", self, "_sync_client")
	WorldServer.connect("system_claimed", self, "_system_claimed")
	WorldServer.connect("resources_quote", self, "_resources_quote")
	
	WorldServer.start()

func _draw():
	for line in lines:
		draw_line(line[0], line[1], Color.aqua)

func _exit_tree():
	WorldServer.stop()

func _sync_client(data):
	
	Costs.setCosts(data.costs)
	
	_init_players(data.players, data.self)
	_draw_map(data.galaxy)
	
	var homeSystem = _systems.get(int(data.self.homeSystem)) as StarSystem
	homeSystem.set_system_owner(_players.get(_playerId))
	$"../Camera".position = homeSystem.position

func _player_joined(player):
	var playerid = int(player.id)
	print(playerid, " joined")
	_players[playerid] = PlayerUtils.fromServerPlayer(player)
	_system_claimed({
		"playerId": playerid,
		"systemId": player.homeSystem
	})

func _player_left(data: Dictionary):
	var playerid = int(data.id)
	print(playerid, " left")
	_players.erase(playerid)
	for systemId in data.ownedSystems:
		_systems.get(systemId).set_system_owner(null)

func _system_claimed(claimData: Dictionary):
	var player = _players.get(int(claimData.playerId))
	_systems.get(int(claimData.systemId)).set_system_owner(player)

func _system_selected(system):
	emit_signal("system_selected", system)

func _init_players(players: Dictionary, current: Dictionary):
	_playerId = int(current.id)
	for player in players.values():
		_players[int(player.id)] = PlayerUtils.fromServerPlayer(player)

	_players[_playerId] = PlayerUtils.fromServerPlayer(current)
	print("Current player: ", _players[_playerId])

func _draw_map(data):
	for system in data.systems:
		var system_instance: StarSystem = _star_system.instance()
		system_instance.init(system)
		if system.owner != -1:
			system_instance.set_system_owner(_players.get(int(system.owner)))
		system_instance.connect("on_selected", self, "_system_selected")
		add_child(system_instance)
		var systemPosition = Vector2(system.x, system.y)
		_systems[int(system.id)] = system_instance

	for route in data.routes:
		var from = _systems.get(int(route[0])).position
		var to = _systems.get(int(route[1])).position
		lines.append([from, to])
	update()

func _resources_quote(quote: Dictionary): 
	self.currentResources = quote
	self.emit_signal("resources_quote", quote)

func get_players():
	return _players
	
func get_system(id):
	return _systems.get(id)

func claim_system(systemId):
	self.emit_signal("resources_quote", currentResources)
	remove_resources(Costs.getCosts("outpost"))
	WorldServer.send({
			"command": "claim_system",
			"data": systemId
		})

func remove_resources(resources: Dictionary):
	for resource in resources:
		currentResources[resource] -= resources[resource]
	self.emit_signal("resources_quote", currentResources)

func has_sufficient_funds(structure: String):
	var costs = Costs.getCosts(structure)
	for resource in costs:
		if currentResources[resource] < costs[resource]:
			return false
	return true
