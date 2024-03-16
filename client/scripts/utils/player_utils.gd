class_name PlayerUtils

static func fromServerPlayer(player):
	player.color = Color().from_hsv(player.color[0] / 360, player.color[1] / 100, player.color[2] / 100)
	return player
