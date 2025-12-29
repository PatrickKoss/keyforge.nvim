package engine

import "github.com/keyforge/keyforge/internal/entities"

// Spawn defines a single enemy spawn in a wave.
type Spawn struct {
	Type  entities.EnemyType
	Delay float64 // seconds after previous spawn
}

// Wave defines a wave of enemies.
type Wave struct {
	Number    int
	Spawns    []Spawn
	BonusGold int
}

// GetWave returns the wave configuration for a given wave number.
func GetWave(waveNum int) Wave {
	// Generate increasingly difficult waves
	switch {
	case waveNum <= 3:
		return easyWave(waveNum)
	case waveNum <= 6:
		return mediumWave(waveNum)
	case waveNum <= 9:
		return hardWave(waveNum)
	default:
		return bossWave(waveNum)
	}
}

func easyWave(num int) Wave {
	count := 3 + num*2 // 5, 7, 9 enemies
	spawns := make([]Spawn, count)
	for i := range count {
		spawns[i] = Spawn{
			Type:  entities.EnemyBug,
			Delay: 1.0,
		}
	}
	spawns[0].Delay = 0 // first spawn is immediate
	return Wave{
		Number:    num,
		Spawns:    spawns,
		BonusGold: 20 + num*5,
	}
}

func mediumWave(num int) Wave {
	bugCount := 5
	gremlinCount := num - 2 // 2, 3, 4 gremlins
	spawns := make([]Spawn, 0, bugCount+gremlinCount)

	// Mix bugs and gremlins
	for range bugCount {
		delay := 1.0
		if len(spawns) == 0 {
			delay = 0
		}
		spawns = append(spawns, Spawn{Type: entities.EnemyBug, Delay: delay})
	}
	for range gremlinCount {
		spawns = append(spawns, Spawn{Type: entities.EnemyGremlin, Delay: 1.5})
	}

	return Wave{
		Number:    num,
		Spawns:    spawns,
		BonusGold: 30 + num*5,
	}
}

func hardWave(num int) Wave {
	bugCount := 3
	gremlinCount := 4
	daemonCount := num - 5 // 2, 3, 4 daemons
	spawns := make([]Spawn, 0, bugCount+gremlinCount+daemonCount)

	// Bugs first
	for range bugCount {
		delay := 0.8
		if len(spawns) == 0 {
			delay = 0
		}
		spawns = append(spawns, Spawn{Type: entities.EnemyBug, Delay: delay})
	}
	// Then gremlins
	for range gremlinCount {
		spawns = append(spawns, Spawn{Type: entities.EnemyGremlin, Delay: 1.0})
	}
	// Then daemons
	for range daemonCount {
		spawns = append(spawns, Spawn{Type: entities.EnemyDaemon, Delay: 2.0})
	}

	return Wave{
		Number:    num,
		Spawns:    spawns,
		BonusGold: 50 + num*10,
	}
}

func bossWave(num int) Wave {
	// Boss wave with supporting enemies
	spawns := []Spawn{
		{Type: entities.EnemyBug, Delay: 0},
		{Type: entities.EnemyBug, Delay: 0.5},
		{Type: entities.EnemyGremlin, Delay: 0.5},
		{Type: entities.EnemyGremlin, Delay: 0.5},
		{Type: entities.EnemyDaemon, Delay: 1.0},
		{Type: entities.EnemyBoss, Delay: 2.0},
		{Type: entities.EnemyGremlin, Delay: 1.0},
		{Type: entities.EnemyGremlin, Delay: 0.5},
		{Type: entities.EnemyDaemon, Delay: 1.5},
	}

	return Wave{
		Number:    num,
		Spawns:    spawns,
		BonusGold: 200,
	}
}
