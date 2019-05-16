package camera

import (
	"math/rand"

	"github.com/go-gl/mathgl/mgl32"

	"rapidengine/configuration"
	"rapidengine/input"
)

type Camera2D struct {
	Speed float32

	SmoothSpeed    float32
	TargetPosition mgl32.Vec3

	ShakeDuration float64
	ShakeStrength float32

	Position  mgl32.Vec3
	UpAxis    mgl32.Vec3
	FrontAxis mgl32.Vec3

	View       mgl32.Mat4
	StaticView mgl32.Mat4

	config *configuration.EngineConfig
}

func NewCamera2D(position mgl32.Vec3, speed float32, config *configuration.EngineConfig) *Camera2D {
	return &Camera2D{
		Position: position,

		UpAxis:    mgl32.Vec3{0, 1, 0},
		FrontAxis: mgl32.Vec3{0, 0, -1},
		StaticView: mgl32.LookAtV(
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0}.Add(mgl32.Vec3{0, 0, -1}),
			mgl32.Vec3{0, 1, 0},
		),

		Speed:       speed,
		SmoothSpeed: 1.0,

		config: config,
	}
}

func (camera2D *Camera2D) Look(delta float64) {
	// Move toward target position
	camera2D.Position = LerpPosition(camera2D.Position, camera2D.TargetPosition, camera2D.SmoothSpeed*float32(delta)*60)

	if camera2D.ShakeDuration > 0 {
		camera2D.Position = camera2D.Position.Add(mgl32.Vec3{
			((rand.Float32() * 2) - 1.0) * camera2D.ShakeStrength,
			((rand.Float32() * 2) - 1.0) * camera2D.ShakeStrength,
			((rand.Float32() * 2) - 1.0) * camera2D.ShakeStrength,
		})
		camera2D.ShakeDuration -= delta
	}

	camera2D.View = mgl32.LookAtV(
		camera2D.Position,
		camera2D.Position.Add(camera2D.FrontAxis),
		camera2D.UpAxis,
	)
}

func (camera2D *Camera2D) GetStaticView() mgl32.Mat4 {
	return camera2D.StaticView
}

func (camera2D *Camera2D) MoveUp() {
	camera2D.TargetPosition = camera2D.TargetPosition.Add(camera2D.UpAxis.Mul(camera2D.Speed))
}

func (camera2D *Camera2D) MoveDown() {
	camera2D.TargetPosition = camera2D.TargetPosition.Sub(camera2D.UpAxis.Mul(camera2D.Speed))
}

func (camera2D *Camera2D) MoveLeft() {
	camera2D.TargetPosition = camera2D.TargetPosition.Sub(camera2D.FrontAxis.Cross(camera2D.UpAxis).Normalize().Mul(camera2D.Speed))
}

func (camera2D *Camera2D) MoveRight() {
	camera2D.TargetPosition = camera2D.TargetPosition.Add(camera2D.FrontAxis.Cross(camera2D.UpAxis).Normalize().Mul(camera2D.Speed))
}

func (camera2D *Camera2D) MoveForward() {
	camera2D.TargetPosition = camera2D.TargetPosition.Add(camera2D.FrontAxis.Mul(camera2D.Speed))
}

func (camera2D *Camera2D) MoveBackward() {
	camera2D.TargetPosition = camera2D.TargetPosition.Sub(camera2D.FrontAxis.Mul(camera2D.Speed))
}

func (camera2D *Camera2D) DefaultControls(inputs *input.Input) {
	if inputs.Keys["w"] {
		camera2D.MoveUp()
	}
	if inputs.Keys["s"] {
		camera2D.MoveDown()
	}
	if inputs.Keys["a"] {
		camera2D.MoveLeft()
	}
	if inputs.Keys["d"] {
		camera2D.MoveRight()
	}
	if inputs.Keys["space"] {
		camera2D.MoveUp()
	}
	if inputs.Keys["shift"] {
		camera2D.MoveDown()
	}
	camera2D.ProcessMouse(inputs.MouseX, inputs.MouseY, inputs.LastMouseX, inputs.LastMouseY)
}

func (camera2D *Camera2D) ChangeYaw(y float32)   {}
func (camera2D *Camera2D) ChangePitch(p float32) {}

func (camera2D *Camera2D) GetFirstViewIndex() *float32 {
	return &camera2D.View[0]
}

func (camera2D *Camera2D) GetPosition() (float32, float32, float32) {
	return ((camera2D.Position.X() / 2) * float32(camera2D.config.ScreenWidth)) + float32(camera2D.config.ScreenWidth/2),
		((camera2D.Position.Y() / 2) * float32(camera2D.config.ScreenHeight)) + float32(camera2D.config.ScreenHeight/2), 0
}

func (camera2D *Camera2D) SetPosition(x, y, z float32) {
	camera2D.TargetPosition = mgl32.Vec3{
		(x - float32(camera2D.config.ScreenWidth/2)) / float32(camera2D.config.ScreenWidth/2),
		(y - float32(camera2D.config.ScreenHeight/2)) / float32(camera2D.config.ScreenHeight/2),
		camera2D.Position.Z(),
	}
}

func (camera2D *Camera2D) SetSpeed(s float32) {
	camera2D.Speed = s
}

func (camera2D *Camera2D) SetSmoothSpeed(s float32) {
	camera2D.SmoothSpeed = s
}

func (camera2D *Camera2D) Shake(duration float64, strength float32) {
	camera2D.ShakeDuration = duration
	camera2D.ShakeStrength = strength
}

func (camera2D *Camera2D) ProcessMouse(mx, my, lmx, lmy float64) {
	return
}

func (camera2D *Camera2D) ChangeRoll(r float32) {}

func LerpPosition(pos1, pos2 mgl32.Vec3, amount float32) mgl32.Vec3 {
	return mgl32.Vec3{
		pos1.X() + amount*(pos2.X()-pos1.X()),
		pos1.Y() + amount*(pos2.Y()-pos1.Y()),
		pos1.Z() + amount*(pos2.Z()-pos1.Z()),
	}
}

func (camera2D *Camera2D) GetFirstModelIndex() *float32 { return nil }
