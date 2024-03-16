extends Node

var Costs := preload("res://scripts/costs.gd")

onready var _galaxy: Galaxy = $"../Galaxy"

func print_players():
	Console.write_line(JSON.print(_galaxy.get_players(), "\t"))

func print_systems(systemIds = ""):
	var ids = systemIds.split(",")
	for id in ids:
		Console.write_line(JSON.print(_galaxy.get_system(int(id)).debug(), "\t"))

func print_costs():
	Console.write_line(JSON.print(Costs._costs, "\t"))

func _ready():
	Console.add_command("players", self, "print_players")\
		.register()
	
	Console.add_command("costs", self, "print_costs")\
		.register()
		
	Console.add_command("systems", self, "print_systems")\
		.add_argument('ids', TYPE_STRING)\
		.register()
