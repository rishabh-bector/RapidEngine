package camera

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"rapidengine/configuration"
)

type Camera2D struct {
	Speed float32

	Position  mgl32.Vec3
	UpAxis    mgl32.Vec3
	FrontAxis mgl32.Vec3

	View mgl32.Mat4

	config *configuration.EngineConfig
}

func NewCamera2D(position mgl32.Vec3, speed float32, config *configuration.EngineConfig) *Camera2D {
	return &Camera2D{
		Position:  position,
		UpAxis:    mgl32.Vec3{0, 1, 0},
		FrontAxis: mgl32.Vec3{0, 0, -1},
		Speed:     speed,
		config:    config,
	}
}

func (camera2D *Camera2D) Look() {
	camera2D.View = mgl32.LookAtV(
		camera2D.Position,
		camera2D.Position.Add(camera2D.FrontAxis),
		camera2D.UpAxis,
	)
}

func (camera2D *Camera2D) ProcessInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyW) == glfw.Press {
		camera2D.Position = camera2D.Position.Add(camera2D.UpAxis.Mul(camera2D.Speed))
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		camera2D.Position = camera2D.Position.Sub(camera2D.UpAxis.Mul(camera2D.Speed))
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		camera2D.Position = camera2D.Position.Sub(camera2D.FrontAxis.Cross(camera2D.UpAxis).Normalize().Mul(camera2D.Speed))
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		camera2D.Position = camera2D.Position.Add(camera2D.FrontAxis.Cross(camera2D.UpAxis).Normalize().Mul(camera2D.Speed))
	}
}

func (camera2D *Camera2D) GetFirstViewIndex() *float32 {
	return &camera2D.View[0]
}

func (camera2D *Camera2D) GetPosition() (float32, float32) {
	return ((camera2D.Position.X() / 2) * float32(camera2D.config.ScreenWidth)) + float32(camera2D.config.ScreenWidth/2),
		((camera2D.Position.Y() / 2) * float32(camera2D.config.ScreenHeight)) + float32(camera2D.config.ScreenHeight/2)
}

func (camera2D *Camera2D) SetPosition(x, y int) {
	camera2D.Position = mgl32.Vec3{
		(float32(x) - float32(camera2D.config.ScreenWidth/2)) / float32(camera2D.config.ScreenWidth/2),
		(float32(y) - float32(camera2D.config.ScreenHeight/2)) / float32(camera2D.config.ScreenHeight/2),
		camera2D.Position.Z(),
	}
}

func (camera2D *Camera2D) SetSpeed(s float32) {
	camera2D.Speed = s
}

func (camera2D *Camera2D) ProcessMouse() {
	return
}