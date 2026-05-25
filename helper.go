package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func planetPosition(semiMajorAxis, eccentricity, angle float32) rl.Vector3 {
	radius := semiMajorAxis * (1 - eccentricity*eccentricity) / (1 + eccentricity*float32(math.Cos(float64(angle))))
	return rl.Vector3{
		X: radius * float32(math.Cos(float64(angle))),
		Y: 0,
		Z: radius * float32(math.Sin(float64(angle))),
	}
}

func moonPosition(planet rl.Vector3, distance, angle float32) rl.Vector3 {
	return rl.Vector3{
		X: planet.X + distance*float32(math.Cos(float64(angle))),
		Y: distance * 0.12 * float32(math.Sin(float64(angle*1.7))),
		Z: planet.Z + distance*float32(math.Sin(float64(angle))),
	}
}

func camera3D(c OrbitCamera) rl.Camera3D {
	pitch := clamp(c.Pitch, -1.25, 1.25)
	cp := float32(math.Cos(float64(pitch)))

	position := rl.Vector3{
		X: c.Target.X + c.Distance*cp*float32(math.Cos(float64(c.Yaw))),
		Y: c.Target.Y + c.Distance*float32(math.Sin(float64(pitch))),
		Z: c.Target.Z + c.Distance*cp*float32(math.Sin(float64(c.Yaw))),
	}

	return rl.Camera3D{
		Position:   position,
		Target:     c.Target,
		Up:         rl.Vector3{X: 0, Y: 1, Z: 0},
		Fovy:       45,
		Projection: rl.CameraPerspective,
	}
}

func clamp(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func drawEllipticalOrbit(semiMajorAxis, eccentricity float32, alpha uint8) {
	const segments = 160
	semiMinorAxis := semiMajorAxis * float32(math.Sqrt(float64(1-eccentricity*eccentricity)))
	previous := rl.Vector3{
		X: semiMajorAxis * (1 - eccentricity),
		Y: 0,
		Z: 0,
	}
	color := rl.Color{R: 255, G: 255, B: 255, A: alpha}

	for i := int32(1); i <= segments; i++ {
		angle := 2 * math.Pi * float64(i) / segments
		current := rl.Vector3{
			X: semiMajorAxis * (float32(math.Cos(angle)) - eccentricity),
			Y: 0,
			Z: semiMinorAxis * float32(math.Sin(angle)),
		}
		rl.DrawLine3D(previous, current, color)
		previous = current
	}
}

func drawMoonOrbit(center rl.Vector3, radius float32, alpha uint8) {
	const segments = 80
	previous := rl.Vector3{
		X: center.X + radius,
		Y: center.Y,
		Z: center.Z,
	}
	color := rl.Color{R: 255, G: 255, B: 255, A: alpha}

	for i := int32(1); i <= segments; i++ {
		angle := 2 * math.Pi * float64(i) / segments
		current := rl.Vector3{
			X: center.X + radius*float32(math.Cos(angle)),
			Y: center.Y,
			Z: center.Z + radius*float32(math.Sin(angle)),
		}
		rl.DrawLine3D(previous, current, color)
		previous = current
	}
}

func drawSaturnRings(center rl.Vector3, planetRadius float32) {
	for step := int32(0); step < 8; step++ {
		radius := planetRadius*1.45 + float32(step)*planetRadius*0.14
		rl.DrawCircle3D(
			center,
			radius,
			rl.Vector3{X: 1, Y: 0, Z: 0},
			68,
			rl.Color{R: 220, G: 200, B: 145, A: 115},
		)
	}
}

func drawStarBackground(stars []Star, camera OrbitCamera) {
	baseYaw := float32(2.35)
	basePitch := float32(0.55)
	rotationOffset := rl.Vector2{
		X: (camera.Yaw - baseYaw) * 140,
		Y: (camera.Pitch - basePitch) * 180,
	}
	targetOffset := rl.Vector2{
		X: camera.Target.X * 0.08,
		Y: camera.Target.Z * 0.08,
	}

	for _, s := range stars {
		offset := rl.Vector2{
			X: (rotationOffset.X + targetOffset.X) * s.Parallax,
			Y: (rotationOffset.Y + targetOffset.Y) * s.Parallax,
		}
		position := wrapScreenPosition(rl.Vector2{
			X: s.Position.X + offset.X,
			Y: s.Position.Y + offset.Y,
		})

		rl.DrawCircleV(
			position,
			s.Radius,
			rl.Color{R: 255, G: 255, B: 255, A: s.Alpha},
		)
	}
}

func wrapScreenPosition(position rl.Vector2) rl.Vector2 {
	width := float32(screenW)
	height := float32(screenH)

	position.X = float32(math.Mod(float64(position.X), float64(width)))
	position.Y = float32(math.Mod(float64(position.Y), float64(height)))
	if position.X < 0 {
		position.X += width
	}
	if position.Y < 0 {
		position.Y += height
	}

	return position
}

func drawLabel(text string, position rl.Vector3, camera rl.Camera3D, color rl.Color) {
	screen := rl.GetWorldToScreen(position, camera)
	if screen.X < -80 || screen.X > float32(screenW)+80 || screen.Y < -40 || screen.Y > float32(screenH)+40 {
		return
	}
	rl.DrawText(text, int32(screen.X)+8, int32(screen.Y)-6, 12, color)
}

func pickPlanet(camera rl.Camera3D, planetPositions []rl.Vector3, planets []Planet) int {
	ray := rl.GetMouseRay(rl.GetMousePosition(), camera)
	selected := -1
	closest := float32(math.MaxFloat32)

	for i, p := range planets {
		pickRadius := p.Radius * 1.8
		if pickRadius < 8 {
			pickRadius = 8
		}

		collision := rl.GetRayCollisionSphere(ray, planetPositions[i], pickRadius)
		if collision.Hit && collision.Distance < closest {
			selected = i
			closest = collision.Distance
		}
	}

	return selected
}
