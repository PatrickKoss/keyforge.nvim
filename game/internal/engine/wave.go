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

// Per-level wave generators
// Each level has specific enemy pools and scaling

// level1Wave generates waves for Level 1 (Mite, Bug only).
func level1Wave(waveNum int) Wave {
	// 3-4 enemies per wave, Mites and Bugs
	count := 3 + (waveNum-1)/2 // 3, 3, 4, 4, 5
	if count > 5 {
		count = 5
	}
	spawns := make([]Spawn, 0, count)

	for i := range count {
		enemyType := entities.EnemyMite
		if i%2 == 1 && waveNum > 2 {
			enemyType = entities.EnemyBug
		}
		delay := 1.2
		if i == 0 {
			delay = 0
		}
		spawns = append(spawns, Spawn{Type: enemyType, Delay: delay})
	}

	return Wave{Number: waveNum, Spawns: spawns, BonusGold: 15 + waveNum*3}
}

// level2Wave generates waves for Level 2 (Mite, Bug).
func level2Wave(waveNum int) Wave {
	// 3-4 enemies per wave
	count := 3 + waveNum/3
	if count > 5 {
		count = 5
	}
	spawns := make([]Spawn, 0, count)

	miteCount := count / 2
	bugCount := count - miteCount

	for range miteCount {
		delay := 1.0
		if len(spawns) == 0 {
			delay = 0
		}
		spawns = append(spawns, Spawn{Type: entities.EnemyMite, Delay: delay})
	}
	for range bugCount {
		spawns = append(spawns, Spawn{Type: entities.EnemyBug, Delay: 1.0})
	}

	return Wave{Number: waveNum, Spawns: spawns, BonusGold: 18 + waveNum*3}
}

// level3Wave generates waves for Level 3 (Bug, Gremlin).
func level3Wave(waveNum int) Wave {
	// 4-5 enemies per wave
	count := 4 + waveNum/4
	if count > 6 {
		count = 6
	}
	spawns := make([]Spawn, 0, count)

	bugCount := count - waveNum/3
	if bugCount < 2 {
		bugCount = 2
	}
	gremlinCount := count - bugCount

	for range bugCount {
		delay := 0.9
		if len(spawns) == 0 {
			delay = 0
		}
		spawns = append(spawns, Spawn{Type: entities.EnemyBug, Delay: delay})
	}
	for range gremlinCount {
		spawns = append(spawns, Spawn{Type: entities.EnemyGremlin, Delay: 1.2})
	}

	return Wave{Number: waveNum, Spawns: spawns, BonusGold: 22 + waveNum*4}
}

// level4Wave generates waves for Level 4 (Bug, Gremlin, Crawler).
func level4Wave(waveNum int) Wave {
	// 4-5 enemies per wave, introducing Crawler
	count := 4 + waveNum/3
	if count > 6 {
		count = 6
	}
	spawns := make([]Spawn, 0, count)

	// Mix: mostly bugs/gremlins with some crawlers in later waves
	crawlerCount := waveNum / 3
	if crawlerCount > 2 {
		crawlerCount = 2
	}
	gremlinCount := (count - crawlerCount) / 2
	bugCount := count - crawlerCount - gremlinCount

	for range bugCount {
		delay := 0.8
		if len(spawns) == 0 {
			delay = 0
		}
		spawns = append(spawns, Spawn{Type: entities.EnemyBug, Delay: delay})
	}
	for range gremlinCount {
		spawns = append(spawns, Spawn{Type: entities.EnemyGremlin, Delay: 1.0})
	}
	for range crawlerCount {
		spawns = append(spawns, Spawn{Type: entities.EnemyCrawler, Delay: 1.5})
	}

	return Wave{Number: waveNum, Spawns: spawns, BonusGold: 28 + waveNum*5}
}

// level5Wave generates waves for Level 5 (Bug, Gremlin, Specter).
func level5Wave(waveNum int) Wave {
	// 5-6 enemies per wave, introducing Specter
	count := 5 + waveNum/4
	if count > 7 {
		count = 7
	}
	spawns := make([]Spawn, 0, count)

	specterCount := waveNum / 3
	if specterCount > 2 {
		specterCount = 2
	}
	gremlinCount := (count - specterCount) / 2
	bugCount := count - specterCount - gremlinCount

	for range bugCount {
		delay := 0.7
		if len(spawns) == 0 {
			delay = 0
		}
		spawns = append(spawns, Spawn{Type: entities.EnemyBug, Delay: delay})
	}
	for range gremlinCount {
		spawns = append(spawns, Spawn{Type: entities.EnemyGremlin, Delay: 0.9})
	}
	for range specterCount {
		spawns = append(spawns, Spawn{Type: entities.EnemySpecter, Delay: 0.6})
	}

	return Wave{Number: waveNum, Spawns: spawns, BonusGold: 35 + waveNum*5}
}

// level6Wave generates waves for Level 6 (Gremlin, Crawler, Daemon).
func level6Wave(waveNum int) Wave {
	// 5-6 enemies per wave, heavy enemies
	count := 5 + waveNum/4
	if count > 7 {
		count = 7
	}
	spawns := make([]Spawn, 0, count)

	daemonCount := waveNum / 4
	if daemonCount > 2 {
		daemonCount = 2
	}
	crawlerCount := (count - daemonCount) / 3
	gremlinCount := count - daemonCount - crawlerCount

	for range gremlinCount {
		delay := 0.8
		if len(spawns) == 0 {
			delay = 0
		}
		spawns = append(spawns, Spawn{Type: entities.EnemyGremlin, Delay: delay})
	}
	for range crawlerCount {
		spawns = append(spawns, Spawn{Type: entities.EnemyCrawler, Delay: 1.3})
	}
	for range daemonCount {
		spawns = append(spawns, Spawn{Type: entities.EnemyDaemon, Delay: 2.0})
	}

	return Wave{Number: waveNum, Spawns: spawns, BonusGold: 45 + waveNum*6}
}

// level7Wave generates waves for Level 7 (Gremlin, Specter, Daemon).
func level7Wave(waveNum int) Wave {
	// 5-6 enemies per wave, speed and power
	count := 5 + waveNum/4
	if count > 7 {
		count = 7
	}
	spawns := make([]Spawn, 0, count)

	daemonCount := waveNum / 4
	if daemonCount > 2 {
		daemonCount = 2
	}
	specterCount := waveNum / 3
	if specterCount > 2 {
		specterCount = 2
	}
	gremlinCount := count - daemonCount - specterCount
	if gremlinCount < 1 {
		gremlinCount = 1
	}

	for range specterCount {
		delay := 0.5
		if len(spawns) == 0 {
			delay = 0
		}
		spawns = append(spawns, Spawn{Type: entities.EnemySpecter, Delay: delay})
	}
	for range gremlinCount {
		spawns = append(spawns, Spawn{Type: entities.EnemyGremlin, Delay: 0.7})
	}
	for range daemonCount {
		spawns = append(spawns, Spawn{Type: entities.EnemyDaemon, Delay: 1.8})
	}

	return Wave{Number: waveNum, Spawns: spawns, BonusGold: 55 + waveNum*7}
}

// level8Wave generates waves for Level 8 (Crawler, Specter, Daemon).
func level8Wave(waveNum int) Wave {
	// 5-7 enemies per wave, late game mix
	count := 5 + waveNum/3
	if count > 7 {
		count = 7
	}
	spawns := make([]Spawn, 0, count)

	daemonCount := waveNum / 3
	if daemonCount > 3 {
		daemonCount = 3
	}
	specterCount := count / 3
	crawlerCount := count - daemonCount - specterCount
	if crawlerCount < 1 {
		crawlerCount = 1
	}

	for range specterCount {
		delay := 0.4
		if len(spawns) == 0 {
			delay = 0
		}
		spawns = append(spawns, Spawn{Type: entities.EnemySpecter, Delay: delay})
	}
	for range crawlerCount {
		spawns = append(spawns, Spawn{Type: entities.EnemyCrawler, Delay: 1.2})
	}
	for range daemonCount {
		spawns = append(spawns, Spawn{Type: entities.EnemyDaemon, Delay: 1.5})
	}

	return Wave{Number: waveNum, Spawns: spawns, BonusGold: 65 + waveNum*8}
}

// level9Wave generates waves for Level 9 (Specter, Daemon).
func level9Wave(waveNum int) Wave {
	// 5-7 enemies per wave, pre-boss difficulty
	count := 5 + waveNum/3
	if count > 7 {
		count = 7
	}
	spawns := make([]Spawn, 0, count)

	daemonCount := waveNum / 2
	if daemonCount > 4 {
		daemonCount = 4
	}
	specterCount := count - daemonCount

	for range specterCount {
		delay := 0.4
		if len(spawns) == 0 {
			delay = 0
		}
		spawns = append(spawns, Spawn{Type: entities.EnemySpecter, Delay: delay})
	}
	for range daemonCount {
		spawns = append(spawns, Spawn{Type: entities.EnemyDaemon, Delay: 1.3})
	}

	return Wave{Number: waveNum, Spawns: spawns, BonusGold: 75 + waveNum*10}
}

// level10Wave generates waves for Level 10 (All enemies + Boss on wave 10).
func level10Wave(waveNum int) Wave {
	// Final level with all enemy types
	if waveNum == 10 {
		// Boss wave
		return level10BossWave()
	}

	// 6-7 enemies per wave
	count := 6 + waveNum/5
	if count > 7 {
		count = 7
	}
	spawns := make([]Spawn, 0, count)

	// Progressive enemy introduction
	var pool []entities.EnemyType
	switch {
	case waveNum <= 3:
		pool = []entities.EnemyType{entities.EnemyBug, entities.EnemyGremlin, entities.EnemyCrawler}
	case waveNum <= 6:
		pool = []entities.EnemyType{entities.EnemyGremlin, entities.EnemyCrawler, entities.EnemySpecter}
	default:
		pool = []entities.EnemyType{entities.EnemySpecter, entities.EnemyDaemon}
	}

	for i := range count {
		enemyType := pool[i%len(pool)]
		delay := 0.6
		if i == 0 {
			delay = 0
		}
		spawns = append(spawns, Spawn{Type: enemyType, Delay: delay})
	}

	return Wave{Number: waveNum, Spawns: spawns, BonusGold: 80 + waveNum*10}
}

func level10BossWave() Wave {
	// Epic boss wave with supporting enemies
	spawns := []Spawn{
		{Type: entities.EnemySpecter, Delay: 0},
		{Type: entities.EnemySpecter, Delay: 0.3},
		{Type: entities.EnemyGremlin, Delay: 0.5},
		{Type: entities.EnemyCrawler, Delay: 0.8},
		{Type: entities.EnemyDaemon, Delay: 1.0},
		{Type: entities.EnemyBoss, Delay: 2.0},
		{Type: entities.EnemySpecter, Delay: 0.5},
		{Type: entities.EnemyDaemon, Delay: 1.0},
	}

	return Wave{Number: 10, Spawns: spawns, BonusGold: 250}
}
