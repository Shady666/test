package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	G            = 6.67430e-11 // Гравитационная постоянная
	TimeStep     = 7200        // Шаг времени в секундах (1 час) и ускоренный на 2 для более быстрого движения планет
	ScreenWidth  = 800
	ScreenHeight = 800
	AU           = 1.496e11 // Астрономическая единица (м)
	Scale        = AU / 150 // Масштаб для отображения на экране
)

type Body struct {
	Name      string
	Mass      float64
	Radius    float64
	X, Y      float64 // Физические координаты в метрах
	Vx, Vy    float64 // Скорость в метрах/секунду
	BodyColor [3]uint8
}

var (
	Sun = Body{
		Name:      "Sun",
		Mass:      1.989e30,
		Radius:    20,
		X:         0,
		Y:         0,
		BodyColor: [3]uint8{255, 255, 0},
	}

	Planets = []Body{
		{
			Name:      "Earth",
			Mass:      5.972e24,
			Radius:    10,
			X:         AU,
			Y:         0,
			BodyColor: [3]uint8{0, 0, 255},
		},
		{
			Name:      "Mars",
			Mass:      6.417e23,
			Radius:    8,
			X:         1.524 * AU,
			Y:         0,
			BodyColor: [3]uint8{255, 0, 0},
		},
		{
			Name:      "Venus",
			Mass:      4.867e24,
			Radius:    9,
			X:         0.723 * AU,
			Y:         0,
			BodyColor: [3]uint8{255, 165, 0},
		},
		{
			Name:      "Mercury",
			Mass:      3.301e23,
			Radius:    7,
			X:         0.387 * AU,
			Y:         0,
			BodyColor: [3]uint8{128, 128, 128},
		},
	}
)

func initPlanets() {
	for i := range Planets {
		distance := math.Sqrt(Planets[i].X*Planets[i].X + Planets[i].Y*Planets[i].Y)
		orbitalSpeed := math.Sqrt(G * Sun.Mass / distance)

		Planets[i].Vx = 0
		Planets[i].Vy = orbitalSpeed
	}
}

func computeForces(bodies []Body) []Body {
	newBodies := make([]Body, len(bodies))
	copy(newBodies, bodies)

	for i := range newBodies {
		fx, fy := 0.0, 0.0
		dx := Sun.X - newBodies[i].X
		dy := Sun.Y - newBodies[i].Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist != 0 {
			force := G * newBodies[i].Mass * Sun.Mass / (dist * dist)
			fx += force * dx / dist
			fy += force * dy / dist
		}

		ax := fx / newBodies[i].Mass
		ay := fy / newBodies[i].Mass
		newBodies[i].Vx += ax * TimeStep
		newBodies[i].Vy += ay * TimeStep
	}

	// Обновляем позиции планет
	for i := range newBodies {
		newBodies[i].X += newBodies[i].Vx * TimeStep
		newBodies[i].Y += newBodies[i].Vy * TimeStep
	}

	return newBodies
}

type Game struct{}

func (g *Game) Update() error {
	Planets = computeForces(Planets)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Рисуем Солнце
	ebitenutil.DrawCircle(screen, ScreenWidth/2, ScreenHeight/2, Sun.Radius, Sun.GetColor())

	// Рисуем планеты
	for _, planet := range Planets {
		x := ScreenWidth/2 + planet.X/Scale
		y := ScreenHeight/2 + planet.Y/Scale
		ebitenutil.DrawCircle(screen, x, y, planet.Radius, planet.GetColor())
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (b Body) GetColor() color.Color {
	return color.RGBA{
		R: b.BodyColor[0],
		G: b.BodyColor[1],
		B: b.BodyColor[2],
		A: 255,
	}
}

func main() {
	initPlanets()
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Solar System Simulation")
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
