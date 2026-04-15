package main

// ── Data types ────────────────────────────────────────────────────────────────

type brick struct{ wx, y float64 }
type qblock struct {
	wx, y  float64
	hit    bool
	reward string // "coin" | "mushroom"
}
type coin struct {
	wx, wy float64
	got    bool
	ft     int
}
type pipe struct{ wx, y, w, h, lipH float64 }
type gap struct{ startX, endX float64 }
type mushroom struct {
	wx, wy float64
	vx, vy float64
	active bool
}
type fireball struct {
	wx, wy float64
	vx, vy float64
}
type scorePopup struct {
	x, y float64
	text string
	life int
}
type cloud struct {
	x, y, speed, w, h float64
}

type enemyKind int

const (
	kindGoomba enemyKind = iota
	kindKoopa
	kindRedKoopa
	kindBobomb
)

type enemyState int

const (
	stateWalk enemyState = iota
	stateFlat
	stateShellStill
	stateShell
	stateFuse
)

type enemy struct {
	kind  enemyKind
	state enemyState
	wx, wy float64
	vx    float64
	timer int
}
