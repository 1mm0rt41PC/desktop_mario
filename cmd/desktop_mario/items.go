package main

// _updateCoins handles coin collection and floating animation.
func (g *Game) _updateCoins(solids [][4]float64, MT, MH float64) {
	aliveCoins := g.coins[:0]
	for i := range g.coins {
		c := &g.coins[i]
		sx := c.wx - g.cam
		if sx < -float64(blk)*2 {
			continue
		}
		if !c.got {
			if overlap(g.mwx, MT, float64(blk), MH, c.wx, c.wy, 8*pixelScale, float64(blk)) {
				c.got = true
				c.ft = 15
				g._addScore(100, sx, c.wy)
			}
		} else {
			c.ft--
			c.wy -= 5
			if c.ft <= 0 {
				continue
			}
		}
		aliveCoins = append(aliveCoins, *c)
	}
	g.coins = aliveCoins
}

// _updateMushrooms handles mushroom movement, physics, and collection.
func (g *Game) _updateMushrooms(solids [][4]float64, MT, MH float64) {
	aliveMushrooms := g.mushrooms[:0]
	for i := range g.mushrooms {
		m := &g.mushrooms[i]
		sx := m.wx - g.cam
		if sx < -float64(blk)*2 {
			continue
		}
		if m.active {
			m.vy += 1.0
			if m.vy > 10 {
				m.vy = 10
			}
			m.wx += m.vx
			m.wy += m.vy
			if m.wy >= g.groundY() {
				m.wy = g.groundY()
				m.vy = 0
			}
			for _, s := range solids {
				if overlap(m.wx, m.wy, float64(blk), float64(blk), s[0], s[1], s[2], s[3]) {
					if m.vx > 0 {
						m.wx = s[0] - float64(blk)
					} else {
						m.wx = s[0] + s[2]
					}
					m.vx *= -1
					break
				}
			}
			if overlap(g.mwx, MT, float64(blk), MH, m.wx, m.wy, float64(blk), float64(blk)) {
				m.active = false
				if !g.isBig {
					g.isBig = true
					g._addScore(1000, sx, m.wy-20)
				}
				continue
			}
		}
		aliveMushrooms = append(aliveMushrooms, *m)
	}
	g.mushrooms = aliveMushrooms
}

// _updateFireballs handles fireball physics, collision, and enemy hits.
func (g *Game) _updateFireballs(solids [][4]float64) {
	fbSz := float64(8 * pixelScale)
	aliveFireballs := g.fireballs[:0]
	for i := range g.fireballs {
		fb := &g.fireballs[i]
		fb.wy += fb.vy
		fb.wx += fb.vx
		fb.vy += 1.0
		if fb.wy >= g.groundY() {
			fb.wy = g.groundY()
			fb.vy = -8.0
		}
		hitSolid := false
		for _, s := range solids {
			if overlap(fb.wx, fb.wy, fbSz, fbSz, s[0], s[1], s[2], s[3]) {
				if fb.vy <= 0 || fb.wy+fbSz > s[1]+fbSz/2 {
					hitSolid = true
				} else {
					fb.wy = s[1] - fbSz
					fb.vy = -8.0
				}
				break
			}
		}
		if hitSolid {
			continue
		}
		fbsx := fb.wx - g.cam
		if fbsx < -float64(blk)*2 || fbsx > g.W+float64(blk)*2 || fb.wy > g.H+50 {
			continue
		}
		killed := false
		for j := range g.enemies {
			e := &g.enemies[j]
			if e.state == stateFlat {
				continue
			}
			if overlap(fb.wx, fb.wy, fbSz, fbSz, e.wx, e.wy, float64(blk), float64(blk)) {
				e.state = stateFlat
				e.timer = 18
				g._addScore(200, e.wx-g.cam, e.wy-20)
				killed = true
				break
			}
		}
		if killed {
			continue
		}
		aliveFireballs = append(aliveFireballs, *fb)
	}
	g.fireballs = aliveFireballs
}

// _updatePopups advances score popup animations.
func (g *Game) _updatePopups() {
	alivePopups := g.popups[:0]
	for i := range g.popups {
		p := &g.popups[i]
		p.y -= 3
		p.life--
		if p.life > 0 {
			alivePopups = append(alivePopups, *p)
		}
	}
	g.popups = alivePopups
}
