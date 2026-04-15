package main

import "math"

// _updateEnemies processes all enemy AI and Mario-enemy collisions.
// Returns true when Mario takes a hit and _update() should return immediately.
func (g *Game) _updateEnemies(solids [][4]float64, marioL, marioT, marioB, MH float64) bool {
	aliveEnemies := g.enemies[:0]
	for i := range g.enemies {
		e := &g.enemies[i]
		sx := e.wx - g.cam
		if sx < -float64(blk)*5 {
			continue
		}
		if sx > g.W+float64(blk)*3 {
			if e.state == stateWalk {
				e.wx += e.vx
			}
			aliveEnemies = append(aliveEnemies, *e)
			continue
		}

		var stop, keep bool
		switch e.state {
		case stateWalk:
			stop, keep = g._updateWalkingEnemy(e, i, solids, marioL, marioT, marioB, MH, sx, &aliveEnemies)
		case stateShellStill:
			g._updateShellStill(e, solids, marioL, marioT, MH, sx)
			keep = true
		case stateFuse:
			g._updateFuseEnemy(e)
			keep = true
		case stateFlat:
			e.timer--
			keep = e.timer > 0
		case stateShell:
			stop, keep = g._updateRollingShell(e, i, solids, marioL, marioT, marioB, MH, sx, &aliveEnemies)
		}
		if stop {
			return true
		}
		if keep {
			aliveEnemies = append(aliveEnemies, *e)
		}
	}
	g.enemies = aliveEnemies
	return false
}

// _updateWalkingEnemy advances a walking enemy and checks Mario collisions.
// Returns (stop, keep): stop=true means _update() should exit; keep=true means keep the enemy alive.
func (g *Game) _updateWalkingEnemy(e *enemy, i int, solids [][4]float64,
	marioL, marioT, marioB, MH float64, sx float64, alive *[]enemy) (stop, keep bool) {

	e.wx += e.vx
	for _, s := range solids {
		if overlap(e.wx, e.wy, float64(blk), float64(blk), s[0], s[1], s[2], s[3]) {
			if e.vx > 0 {
				e.wx = s[0] - float64(blk)
			} else {
				e.wx = s[0] + s[2]
			}
			e.vx *= -1
			break
		}
	}
	if e.kind == kindRedKoopa {
		g._redKoopaEdgeCheck(e, solids)
	}
	if g.invincible <= 0 && g.stompGrace <= 0 && g.shrinkTimer <= 0 &&
		overlap(marioL+4, marioT, float64(blk)-8, MH, e.wx, e.wy, float64(blk), float64(blk)) {
		if g.mvy > 0 && marioB < e.wy+float64(blk)*0.6 {
			g._stompEnemy(e)
			g.mvy = -10
			g.stompGrace = 25
			g._addScore(200, sx, e.wy-20)
			return false, true
		}
		// Mario takes a hit – preserve all remaining enemies and stop.
		g._takeHit()
		*alive = append(*alive, *e)
		for j := i + 1; j < len(g.enemies); j++ {
			*alive = append(*alive, g.enemies[j])
		}
		g.enemies = *alive
		return true, false
	}
	return false, true
}

func (g *Game) _redKoopaEdgeCheck(e *enemy, solids [][4]float64) {
	aheadX := e.wx + float64(blk)
	if e.vx < 0 {
		aheadX = e.wx - 4
	}
	onSolid := false
	if !g._inGap(aheadX, 4) {
		if e.wy >= g.groundY() {
			onSolid = true
		} else {
			for _, s := range solids {
				if overlap(aheadX, e.wy+float64(blk), 4, 4, s[0], s[1], s[2], s[3]) {
					onSolid = true
					break
				}
			}
		}
	}
	if !onSolid {
		e.vx *= -1
	}
}

func (g *Game) _stompEnemy(e *enemy) {
	switch e.kind {
	case kindGoomba:
		e.state = stateFlat
		e.timer = 18
		e.wy += float64(blk) / 2
	case kindBobomb:
		e.state = stateFuse
		e.timer = 90
		e.vx = 0
	default:
		e.state = stateShellStill
		e.vx = 0
		e.timer = 300
	}
}

// _updateShellStill handles a stopped shell (waiting to reactivate or be kicked).
func (g *Game) _updateShellStill(e *enemy, solids [][4]float64, marioL, marioT, MH float64, sx float64) {
	e.timer--
	if e.timer <= 0 {
		e.state = stateWalk
		e.vx = -1.5
		return
	}
	if g.invincible <= 0 && g.stompGrace <= 0 &&
		overlap(marioL+2, marioT+4, float64(blk)-4, MH-8, e.wx, e.wy, float64(blk), float64(blk)) {
		dir := 10.0
		if g.mwx+float64(blk)/2 > e.wx+float64(blk)/2 {
			dir = -10
		}
		e.state = stateShell
		e.vx = dir
		if g.mwx+float64(blk)/2 < e.wx+float64(blk)/2 {
			g.mwx = e.wx - float64(blk) - 2
		} else {
			g.mwx = e.wx + float64(blk) + 2
		}
		g.stompGrace = 15
		g._addScore(100, sx, e.wy-20)
	}
}

func (g *Game) _updateFuseEnemy(e *enemy) {
	e.timer--
	if e.timer > 0 {
		return
	}
	e.state = stateFlat
	e.timer = 20
	blastR := float64(blk) * 3
	for j := range g.enemies {
		other := &g.enemies[j]
		if other == e || other.state == stateFlat {
			continue
		}
		if math.Abs(other.wx-e.wx) < blastR && math.Abs(other.wy-e.wy) < blastR {
			other.state = stateFlat
			other.timer = 18
			g._addScore(200, other.wx-g.cam, other.wy-20)
		}
	}
	if g.invincible <= 0 && g.shrinkTimer <= 0 {
		if math.Abs(g.mwx-e.wx) < blastR && math.Abs(g.my-e.wy) < blastR {
			g._takeHit()
		}
	}
}

// _updateRollingShell handles a rolling shell.
// Returns (stop, keep): stop=true → exit update; keep=false → drop enemy from alive list.
func (g *Game) _updateRollingShell(e *enemy, i int, solids [][4]float64,
	marioL, marioT, marioB, MH float64, sx float64, alive *[]enemy) (stop, keep bool) {

	e.wx += e.vx
	for j := range g.enemies {
		other := &g.enemies[j]
		if other == e || other.state == stateFlat || other.state == stateShell {
			continue
		}
		if overlap(e.wx, e.wy, float64(blk), float64(blk), other.wx, other.wy, float64(blk), float64(blk)) {
			other.state = stateFlat
			other.timer = 18
			g._addScore(100, other.wx-g.cam, other.wy-20)
		}
	}
	for _, s := range solids {
		if overlap(e.wx, e.wy, float64(blk), float64(blk), s[0], s[1], s[2], s[3]) {
			e.vx *= -1
			break
		}
	}
	if g.invincible <= 0 && g.stompGrace <= 0 &&
		overlap(marioL+4, marioT, float64(blk)-8, MH, e.wx, e.wy, float64(blk), float64(blk)) {
		if g.mvy > 0 && marioB < e.wy+float64(blk)*0.5 {
			e.state = stateShellStill
			e.vx = 0
			e.timer = 300
			g.mvy = -10
			g.stompGrace = 25
			g._addScore(100, sx, e.wy-20)
			return false, true
		}
		// Mario takes a hit – preserve all remaining enemies and stop.
		g._takeHit()
		*alive = append(*alive, *e)
		for j := i + 1; j < len(g.enemies); j++ {
			*alive = append(*alive, g.enemies[j])
		}
		g.enemies = *alive
		return true, false
	}
	// Drop the shell if it has rolled too far off-screen.
	if math.Abs(e.wx-g.mwx) > g.W*2 {
		return false, false
	}
	return false, true
}
