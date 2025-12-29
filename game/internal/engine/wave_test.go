package engine

import (
	"testing"

	"github.com/keyforge/keyforge/internal/entities"
)

func TestGetWaveEasy(t *testing.T) {
	for waveNum := 1; waveNum <= 3; waveNum++ {
		wave := GetWave(waveNum)
		if wave.Number != waveNum {
			t.Errorf("Wave %d has wrong number: %d", waveNum, wave.Number)
		}
		if len(wave.Spawns) == 0 {
			t.Errorf("Wave %d has no spawns", waveNum)
		}
		// Easy waves should only have bugs
		for _, spawn := range wave.Spawns {
			if spawn.Type != entities.EnemyBug {
				t.Errorf("Easy wave %d spawns non-bug: %v", waveNum, spawn.Type)
			}
		}
	}
}

func TestGetWaveMedium(t *testing.T) {
	for waveNum := 4; waveNum <= 6; waveNum++ {
		wave := GetWave(waveNum)
		if wave.Number != waveNum {
			t.Errorf("Wave %d has wrong number: %d", waveNum, wave.Number)
		}
		// Medium waves have bugs and gremlins
		hasBug := false
		hasGremlin := false
		for _, spawn := range wave.Spawns {
			if spawn.Type == entities.EnemyBug {
				hasBug = true
			}
			if spawn.Type == entities.EnemyGremlin {
				hasGremlin = true
			}
		}
		if !hasBug {
			t.Errorf("Medium wave %d should have bugs", waveNum)
		}
		if !hasGremlin {
			t.Errorf("Medium wave %d should have gremlins", waveNum)
		}
	}
}

func TestGetWaveHard(t *testing.T) {
	for waveNum := 7; waveNum <= 9; waveNum++ {
		wave := GetWave(waveNum)
		if wave.Number != waveNum {
			t.Errorf("Wave %d has wrong number: %d", waveNum, wave.Number)
		}
		// Hard waves should have daemons
		hasDaemon := false
		for _, spawn := range wave.Spawns {
			if spawn.Type == entities.EnemyDaemon {
				hasDaemon = true
			}
		}
		if !hasDaemon {
			t.Errorf("Hard wave %d should have daemons", waveNum)
		}
	}
}

func TestGetWaveBoss(t *testing.T) {
	wave := GetWave(10)
	hasBoss := false
	for _, spawn := range wave.Spawns {
		if spawn.Type == entities.EnemyBoss {
			hasBoss = true
		}
	}
	if !hasBoss {
		t.Error("Boss wave should have a boss")
	}
}

func TestWaveSpawnDelays(t *testing.T) {
	// First spawn should always have delay 0
	for waveNum := 1; waveNum <= 10; waveNum++ {
		wave := GetWave(waveNum)
		if len(wave.Spawns) == 0 {
			t.Errorf("Wave %d has no spawns", waveNum)
			continue
		}
		if wave.Spawns[0].Delay != 0 {
			t.Errorf("Wave %d first spawn should have delay 0, got %v", waveNum, wave.Spawns[0].Delay)
		}
	}
}

func TestWaveBonusGold(t *testing.T) {
	// All waves should have positive bonus gold
	for waveNum := 1; waveNum <= 10; waveNum++ {
		wave := GetWave(waveNum)
		if wave.BonusGold <= 0 {
			t.Errorf("Wave %d should have positive bonus gold, got %d", waveNum, wave.BonusGold)
		}
	}

	// Later waves should generally give more gold
	wave1 := GetWave(1)
	wave10 := GetWave(10)
	if wave10.BonusGold <= wave1.BonusGold {
		t.Errorf("Wave 10 gold (%d) should be more than wave 1 (%d)",
			wave10.BonusGold, wave1.BonusGold)
	}
}

func TestLevel1Wave(t *testing.T) {
	for waveNum := 1; waveNum <= 5; waveNum++ {
		wave := level1Wave(waveNum)
		t.Run("wave number matches", func(t *testing.T) {
			if wave.Number != waveNum {
				t.Errorf("Expected wave number %d, got %d", waveNum, wave.Number)
			}
		})

		t.Run("spawn count 3-5", func(t *testing.T) {
			if len(wave.Spawns) < 3 || len(wave.Spawns) > 5 {
				t.Errorf("Wave %d should have 3-5 spawns, got %d", waveNum, len(wave.Spawns))
			}
		})

		t.Run("only mites and bugs", func(t *testing.T) {
			for _, spawn := range wave.Spawns {
				if spawn.Type != entities.EnemyMite && spawn.Type != entities.EnemyBug {
					t.Errorf("Wave %d spawns invalid enemy: %v", waveNum, spawn.Type)
				}
			}
		})
	}
}

func TestLevel2Wave(t *testing.T) {
	for waveNum := 1; waveNum <= 6; waveNum++ {
		wave := level2Wave(waveNum)

		t.Run("spawn count 3-5", func(t *testing.T) {
			if len(wave.Spawns) < 3 || len(wave.Spawns) > 5 {
				t.Errorf("Wave %d should have 3-5 spawns, got %d", waveNum, len(wave.Spawns))
			}
		})

		t.Run("only mites and bugs", func(t *testing.T) {
			for _, spawn := range wave.Spawns {
				if spawn.Type != entities.EnemyMite && spawn.Type != entities.EnemyBug {
					t.Errorf("Wave %d spawns invalid enemy: %v", waveNum, spawn.Type)
				}
			}
		})
	}
}

func TestLevel3Wave(t *testing.T) {
	for waveNum := 1; waveNum <= 7; waveNum++ {
		wave := level3Wave(waveNum)

		t.Run("spawn count 4-6", func(t *testing.T) {
			if len(wave.Spawns) < 4 || len(wave.Spawns) > 6 {
				t.Errorf("Wave %d should have 4-6 spawns, got %d", waveNum, len(wave.Spawns))
			}
		})

		t.Run("only bugs and gremlins", func(t *testing.T) {
			for _, spawn := range wave.Spawns {
				if spawn.Type != entities.EnemyBug && spawn.Type != entities.EnemyGremlin {
					t.Errorf("Wave %d spawns invalid enemy: %v", waveNum, spawn.Type)
				}
			}
		})
	}
}

func TestLevel4Wave(t *testing.T) {
	for waveNum := 1; waveNum <= 7; waveNum++ {
		wave := level4Wave(waveNum)

		t.Run("spawn count 4-6", func(t *testing.T) {
			if len(wave.Spawns) < 4 || len(wave.Spawns) > 6 {
				t.Errorf("Wave %d should have 4-6 spawns, got %d", waveNum, len(wave.Spawns))
			}
		})

		t.Run("only bugs, gremlins, crawlers", func(t *testing.T) {
			for _, spawn := range wave.Spawns {
				if spawn.Type != entities.EnemyBug &&
					spawn.Type != entities.EnemyGremlin &&
					spawn.Type != entities.EnemyCrawler {
					t.Errorf("Wave %d spawns invalid enemy: %v", waveNum, spawn.Type)
				}
			}
		})
	}
}

func TestLevel5Wave(t *testing.T) {
	for waveNum := 1; waveNum <= 8; waveNum++ {
		wave := level5Wave(waveNum)

		t.Run("spawn count 5-7", func(t *testing.T) {
			if len(wave.Spawns) < 5 || len(wave.Spawns) > 7 {
				t.Errorf("Wave %d should have 5-7 spawns, got %d", waveNum, len(wave.Spawns))
			}
		})

		t.Run("only bugs, gremlins, specters", func(t *testing.T) {
			for _, spawn := range wave.Spawns {
				if spawn.Type != entities.EnemyBug &&
					spawn.Type != entities.EnemyGremlin &&
					spawn.Type != entities.EnemySpecter {
					t.Errorf("Wave %d spawns invalid enemy: %v", waveNum, spawn.Type)
				}
			}
		})
	}
}

func TestLevel6Wave(t *testing.T) {
	for waveNum := 1; waveNum <= 8; waveNum++ {
		wave := level6Wave(waveNum)

		t.Run("spawn count 5-7", func(t *testing.T) {
			if len(wave.Spawns) < 5 || len(wave.Spawns) > 7 {
				t.Errorf("Wave %d should have 5-7 spawns, got %d", waveNum, len(wave.Spawns))
			}
		})

		t.Run("only gremlins, crawlers, daemons", func(t *testing.T) {
			for _, spawn := range wave.Spawns {
				if spawn.Type != entities.EnemyGremlin &&
					spawn.Type != entities.EnemyCrawler &&
					spawn.Type != entities.EnemyDaemon {
					t.Errorf("Wave %d spawns invalid enemy: %v", waveNum, spawn.Type)
				}
			}
		})
	}
}

func TestLevel7Wave(t *testing.T) {
	for waveNum := 1; waveNum <= 9; waveNum++ {
		wave := level7Wave(waveNum)

		t.Run("spawn count 5-7", func(t *testing.T) {
			if len(wave.Spawns) < 5 || len(wave.Spawns) > 7 {
				t.Errorf("Wave %d should have 5-7 spawns, got %d", waveNum, len(wave.Spawns))
			}
		})

		t.Run("only gremlins, specters, daemons", func(t *testing.T) {
			for _, spawn := range wave.Spawns {
				if spawn.Type != entities.EnemyGremlin &&
					spawn.Type != entities.EnemySpecter &&
					spawn.Type != entities.EnemyDaemon {
					t.Errorf("Wave %d spawns invalid enemy: %v", waveNum, spawn.Type)
				}
			}
		})
	}
}

func TestLevel8Wave(t *testing.T) {
	for waveNum := 1; waveNum <= 9; waveNum++ {
		wave := level8Wave(waveNum)

		t.Run("spawn count 5-7", func(t *testing.T) {
			if len(wave.Spawns) < 5 || len(wave.Spawns) > 7 {
				t.Errorf("Wave %d should have 5-7 spawns, got %d", waveNum, len(wave.Spawns))
			}
		})

		t.Run("only crawlers, specters, daemons", func(t *testing.T) {
			for _, spawn := range wave.Spawns {
				if spawn.Type != entities.EnemyCrawler &&
					spawn.Type != entities.EnemySpecter &&
					spawn.Type != entities.EnemyDaemon {
					t.Errorf("Wave %d spawns invalid enemy: %v", waveNum, spawn.Type)
				}
			}
		})
	}
}

func TestLevel9Wave(t *testing.T) {
	for waveNum := 1; waveNum <= 10; waveNum++ {
		wave := level9Wave(waveNum)

		t.Run("spawn count 5-7", func(t *testing.T) {
			if len(wave.Spawns) < 5 || len(wave.Spawns) > 7 {
				t.Errorf("Wave %d should have 5-7 spawns, got %d", waveNum, len(wave.Spawns))
			}
		})

		t.Run("only specters and daemons", func(t *testing.T) {
			for _, spawn := range wave.Spawns {
				if spawn.Type != entities.EnemySpecter && spawn.Type != entities.EnemyDaemon {
					t.Errorf("Wave %d spawns invalid enemy: %v", waveNum, spawn.Type)
				}
			}
		})
	}
}

func TestLevel10Wave(t *testing.T) {
	t.Run("regular waves 1-9", func(t *testing.T) {
		for waveNum := 1; waveNum <= 9; waveNum++ {
			wave := level10Wave(waveNum)
			if len(wave.Spawns) < 6 || len(wave.Spawns) > 7 {
				t.Errorf("Wave %d should have 6-7 spawns, got %d", waveNum, len(wave.Spawns))
			}
		}
	})

	t.Run("boss wave 10", func(t *testing.T) {
		wave := level10Wave(10)
		hasBoss := false
		for _, spawn := range wave.Spawns {
			if spawn.Type == entities.EnemyBoss {
				hasBoss = true
			}
		}
		if !hasBoss {
			t.Error("Wave 10 should have a boss")
		}
		if wave.BonusGold < 200 {
			t.Errorf("Boss wave should have bonus gold >= 200, got %d", wave.BonusGold)
		}
	})
}

func TestLevel10BossWave(t *testing.T) {
	wave := level10BossWave()

	t.Run("has boss", func(t *testing.T) {
		hasBoss := false
		for _, spawn := range wave.Spawns {
			if spawn.Type == entities.EnemyBoss {
				hasBoss = true
			}
		}
		if !hasBoss {
			t.Error("Boss wave must have a boss")
		}
	})

	t.Run("has support enemies", func(t *testing.T) {
		if len(wave.Spawns) < 5 {
			t.Errorf("Boss wave should have support enemies, got %d spawns", len(wave.Spawns))
		}
	})

	t.Run("high bonus gold", func(t *testing.T) {
		if wave.BonusGold < 200 {
			t.Errorf("Boss wave should have high bonus gold, got %d", wave.BonusGold)
		}
	})
}

func TestWaveSpawnCount(t *testing.T) {
	// Test that all per-level waves stay within 3-7 enemy range (plus boss wave)
	registry := NewLevelRegistry()

	for _, level := range registry.GetAll() {
		t.Run(level.Name, func(t *testing.T) {
			for waveNum := 1; waveNum <= level.TotalWaves; waveNum++ {
				wave := level.WaveFunc(waveNum)

				// Allow boss waves to have more enemies
				isBossWave := false
				for _, spawn := range wave.Spawns {
					if spawn.Type == entities.EnemyBoss {
						isBossWave = true
						break
					}
				}

				if !isBossWave {
					if len(wave.Spawns) < 3 {
						t.Errorf("Wave %d has too few spawns: %d (min 3)", waveNum, len(wave.Spawns))
					}
					if len(wave.Spawns) > 7 {
						t.Errorf("Wave %d has too many spawns: %d (max 7)", waveNum, len(wave.Spawns))
					}
				}
			}
		})
	}
}
