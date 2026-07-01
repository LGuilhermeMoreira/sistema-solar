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

func vector3Add(a, b rl.Vector3) rl.Vector3 {
	return rl.Vector3{X: a.X + b.X, Y: a.Y + b.Y, Z: a.Z + b.Z}
}

func vector3Subtract(a, b rl.Vector3) rl.Vector3 {
	return rl.Vector3{X: a.X - b.X, Y: a.Y - b.Y, Z: a.Z - b.Z}
}

func vector3Scale(v rl.Vector3, scale float32) rl.Vector3 {
	return rl.Vector3{X: v.X * scale, Y: v.Y * scale, Z: v.Z * scale}
}

func vector3Length(v rl.Vector3) float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

func vector3Normalize(v rl.Vector3) rl.Vector3 {
	length := vector3Length(v)
	if length == 0 {
		return rl.Vector3{}
	}
	return vector3Scale(v, 1/length)
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

func genRingMesh(innerRadius, outerRadius float32, segments int32) rl.Mesh {
	front := make([]float32, segments*18)
	frontUV := make([]float32, segments*12)

	for i := int32(0); i < segments; i++ {
		angle := 2 * math.Pi * float64(i) / float64(segments)
		nextAngle := 2 * math.Pi * float64(i+1) / float64(segments)
		cos := float32(math.Cos(angle))
		sin := float32(math.Sin(angle))
		nextCos := float32(math.Cos(nextAngle))
		nextSin := float32(math.Sin(nextAngle))

		front[i*18] = innerRadius * cos
		front[i*18+1] = 0
		front[i*18+2] = innerRadius * sin
		front[i*18+3] = outerRadius * cos
		front[i*18+4] = 0
		front[i*18+5] = outerRadius * sin
		front[i*18+6] = innerRadius * nextCos
		front[i*18+7] = 0
		front[i*18+8] = innerRadius * nextSin
		front[i*18+9] = innerRadius * nextCos
		front[i*18+10] = 0
		front[i*18+11] = innerRadius * nextSin
		front[i*18+12] = outerRadius * cos
		front[i*18+13] = 0
		front[i*18+14] = outerRadius * sin
		front[i*18+15] = outerRadius * nextCos
		front[i*18+16] = 0
		front[i*18+17] = outerRadius * nextSin

		frontUV[i*12] = 0
		frontUV[i*12+1] = float32(i) / float32(segments)
		frontUV[i*12+2] = 1
		frontUV[i*12+3] = float32(i) / float32(segments)
		frontUV[i*12+4] = 0
		frontUV[i*12+5] = float32(i+1) / float32(segments)
		frontUV[i*12+6] = 0
		frontUV[i*12+7] = float32(i+1) / float32(segments)
		frontUV[i*12+8] = 1
		frontUV[i*12+9] = float32(i) / float32(segments)
		frontUV[i*12+10] = 1
		frontUV[i*12+11] = float32(i+1) / float32(segments)
	}

	// Duplica com winding invertido -> visível de cima e de baixo
	vertices := make([]float32, len(front)*2)
	texCoords := make([]float32, len(frontUV)*2)
	normals := make([]float32, len(front)*2)

	copy(vertices, front)
	copy(texCoords, frontUV)
	for i := 0; i < len(front)/3; i++ {
		normals[i*3+1] = 1
	}

	bv0 := len(front)
	bt0 := len(frontUV)
	for i := int32(0); i < segments; i++ {
		fv, ft := i*18, i*12
		bv, bt := bv0+int(i*18), bt0+int(i*12)

		copy(vertices[bv:bv+3], front[fv:fv+3])
		copy(vertices[bv+3:bv+6], front[fv+6:fv+9])
		copy(vertices[bv+6:bv+9], front[fv+3:fv+6])
		copy(vertices[bv+9:bv+12], front[fv+9:fv+12])
		copy(vertices[bv+12:bv+15], front[fv+15:fv+18])
		copy(vertices[bv+15:bv+18], front[fv+12:fv+15])

		copy(texCoords[bt:bt+2], frontUV[ft:ft+2])
		copy(texCoords[bt+2:bt+4], frontUV[ft+4:ft+6])
		copy(texCoords[bt+4:bt+6], frontUV[ft+2:ft+4])
		copy(texCoords[bt+6:bt+8], frontUV[ft+6:ft+8])
		copy(texCoords[bt+8:bt+10], frontUV[ft+10:ft+12])
		copy(texCoords[bt+10:bt+12], frontUV[ft+8:ft+10])
	}
	for i := len(front) / 3; i < len(vertices)/3; i++ {
		normals[i*3+1] = -1
	}

	mesh := rl.Mesh{}
	mesh.VertexCount = int32(len(vertices) / 3)
	mesh.TriangleCount = int32(len(vertices) / 9)
	mesh.Vertices = &vertices[0]
	mesh.Texcoords = &texCoords[0]
	mesh.Normals = &normals[0]

	rl.UploadMesh(&mesh, false)
	return mesh
}

func drawSaturnRings(center rl.Vector3, planetRadius float32, ringTexture rl.Texture2D, shader rl.Shader) {
	const segments = 128
	innerRadius := planetRadius * 1.3
	outerRadius := planetRadius * 2.4

	mesh := genRingMesh(innerRadius, outerRadius, segments)
	material := rl.LoadMaterialDefault() // sem material.Shader = shader
	if ringTexture.ID > 0 {
		rl.SetMaterialTexture(&material, rl.MapDiffuse, ringTexture)
	}

	transform := rl.MatrixTranslate(center.X, center.Y, center.Z)
	rl.DrawMesh(mesh, material, transform)
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

		rl.DrawPixelV(
			position,
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

func loadPlanetAssets(p *Planet, shader rl.Shader) {
	p.Mesh = rl.GenMeshSphere(1.0, 32, 32)
	p.Material = rl.LoadMaterialDefault()
	p.Material.Shader = shader

	if p.TexturePath != "" {
		tex := rl.LoadTexture(p.TexturePath)
		if tex.ID > 0 {
			p.Texture = tex
			p.Material.GetMap(rl.MapDiffuse).Texture = tex
		}
	}
}

func drawPlanet(p Planet, position rl.Vector3, angle float32) {
	if p.Texture.ID > 0 {
		mat := new(rl.Material)
		*mat = p.Material
		rl.SetMaterialTexture(mat, rl.MapDiffuse, p.Texture)
		transform := rl.MatrixScale(p.Radius, p.Radius, p.Radius)
		rotation := rl.MatrixRotateY(angle)
		transform = rl.MatrixMultiply(transform, rotation)
		translation := rl.MatrixTranslate(position.X, position.Y, position.Z)
		transform = rl.MatrixMultiply(transform, translation)
		rl.DrawMesh(p.Mesh, *mat, transform)
	} else {
		rl.DrawSphere(position, p.Radius, p.Color)
		rl.DrawSphereWires(position, p.Radius*1.02, 10, 18, rl.Color{R: 255, G: 255, B: 255, A: 35})
	}
}

func drawSun(mesh rl.Mesh, material rl.Material, hasTexture bool, angle float32) {
	const sunRadius = float32(22)
	if hasTexture {
		transform := rl.MatrixScale(sunRadius, sunRadius, sunRadius)
		rotation := rl.MatrixRotateY(angle)
		transform = rl.MatrixMultiply(transform, rotation)
		rl.DrawMesh(mesh, material, transform)
	} else {
		rl.DrawSphere(rl.Vector3{X: 0, Y: 0, Z: 0}, sunRadius, rl.Color{R: 255, G: 220, B: 80, A: 255})
		rl.DrawSphereWires(rl.Vector3{X: 0, Y: 0, Z: 0}, 26, 16, 24, rl.Color{R: 255, G: 190, B: 40, A: 100})
	}
}

func loadCometAssets(comet *Comet, texturePath string, shader rl.Shader) {
	comet.Mesh = rl.GenMeshSphere(1.0, 32, 32)
	comet.Material = rl.LoadMaterialDefault()
	comet.Material.Shader = shader
	if texturePath != "" {
		tex := rl.LoadTexture(texturePath)
		if tex.ID > 0 {
			comet.Texture = tex
			rl.SetMaterialTexture(&comet.Material, rl.MapDiffuse, tex)
		}
	}
}

func launchComet(earth Planet, simulationSpeed float32, existing Comet) Comet {
	const timeToImpact = float32(3.2)

	futureAngle := earth.Angle + earth.Speed*simulationSpeed*timeToImpact
	impactPoint := planetPosition(earth.SemiMajorAxis, earth.Eccentricity, futureAngle)
	start := vector3Add(impactPoint, rl.Vector3{X: -260, Y: 95, Z: -210})
	velocity := vector3Scale(vector3Subtract(impactPoint, start), 1/timeToImpact)

	return Comet{
		Active:   true,
		Position: start,
		Velocity: velocity,
		Radius:   0.9,
		MaxAge:   timeToImpact + 2,
		Texture:  existing.Texture,
		Mesh:     existing.Mesh,
		Material: existing.Material,
	}
}

func updateComet(comet *Comet, target rl.Vector3, targetRadius float32, dt float32) bool {
	if !comet.Active {
		return false
	}

	comet.Position = vector3Add(comet.Position, vector3Scale(comet.Velocity, dt))
	comet.Age += dt
	if vector3Length(vector3Subtract(target, comet.Position)) <= targetRadius+comet.Radius {
		comet.Active = false
		return true
	}
	if comet.Age >= comet.MaxAge {
		comet.Active = false
	}
	return false
}

func drawComet(comet *Comet) {
	if !comet.Active {
		return
	}

	direction := vector3Normalize(comet.Velocity)
	tailEnd := vector3Subtract(comet.Position, vector3Scale(direction, 48))
	tailWideA := vector3Add(tailEnd, rl.Vector3{X: 0, Y: 9, Z: 0})
	tailWideB := vector3Add(tailEnd, rl.Vector3{X: 0, Y: -6, Z: 0})

	if comet.Texture.ID > 0 {
		mat := new(rl.Material)
		*mat = comet.Material
		rl.SetMaterialTexture(mat, rl.MapDiffuse, comet.Texture)
		transform := rl.MatrixScale(comet.Radius, comet.Radius, comet.Radius)
		translation := rl.MatrixTranslate(comet.Position.X, comet.Position.Y, comet.Position.Z)
		transform = rl.MatrixMultiply(transform, translation)
		rl.DrawMesh(comet.Mesh, *mat, transform)
	} else {
		rl.DrawSphere(comet.Position, comet.Radius, rl.Color{R: 238, G: 238, B: 225, A: 255})
	}

	rl.DrawLine3D(comet.Position, tailEnd, rl.Color{R: 255, G: 175, B: 70, A: 210})
	rl.DrawLine3D(comet.Position, tailWideA, rl.Color{R: 255, G: 220, B: 135, A: 135})
	rl.DrawLine3D(comet.Position, tailWideB, rl.Color{R: 110, G: 185, B: 255, A: 120})
}

func loadMoonAssets(m *Moon, p Planet, shader rl.Shader) {
	m.Mesh = rl.GenMeshSphere(1.0, 32, 32)
	m.Material = rl.LoadMaterialDefault()
	m.Material.Shader = shader

	neutralColor := rl.Color{R: 232, G: 228, B: 214, A: 255}
	m.Color = neutralColor
	m.Material.GetMap(rl.MapDiffuse).Color = neutralColor

	image := rl.GenImageColor(2, 2, neutralColor)
	if image != nil {
		m.Texture = rl.LoadTextureFromImage(image)
		rl.UnloadImage(image)
	}
	if m.Texture.ID > 0 {
		rl.SetMaterialTexture(&m.Material, rl.MapDiffuse, m.Texture)
	}
}

func drawMoon(m Moon, position rl.Vector3, angle float32) {
	mat := new(rl.Material)
	*mat = m.Material
	if m.Texture.ID > 0 {
		rl.SetMaterialTexture(mat, rl.MapDiffuse, m.Texture)
	}

	transform := rl.MatrixScale(m.Radius, m.Radius, m.Radius)
	rotation := rl.MatrixRotateY(0)
	transform = rl.MatrixMultiply(transform, rotation)
	translation := rl.MatrixTranslate(position.X, position.Y, position.Z)
	transform = rl.MatrixMultiply(transform, translation)

	rl.DrawMesh(m.Mesh, *mat, transform)
}
