class_name Costs

const _costs = {}

static func setCosts(costs: Dictionary):
	_costs.clear()
	for resource in costs:
		_costs[resource] = costs[resource]

static func getCosts(name: String):
	return _costs[name]
