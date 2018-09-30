package rapidengine

import "github.com/go-gl/gl/v4.1-core/gl"

const SHADER_COLOR = "SHADER_COLOR"
const SHADER_TEXTURE = "SHADER_TEXTURE"
const SHADER_CUBEMAP = "SHADER_CUBEMAP"

var SHADER_COLOR_UNI = []float32{1, 0, 0}
var SHADER_TEXTURE_UNI = []float32{0, 1, 0}
var SHADER_CUBEMAP_UNI = []float32{0, 0, 1}

type Material struct {
	shaderProgram uint32
	shaderType    string

	texture *uint32

	transparencyEnabled bool
	transparencyTexture *uint32

	color []float32
	shine float32
}

func NewMaterial(program uint32) Material {
	return Material{
		shaderProgram:       program,
		shaderType:          SHADER_COLOR,
		color:               []float32{1, 1, 1},
		shine:               0.8,
		transparencyEnabled: false,
	}
}

func (material *Material) PreRender() {
	switch material.shaderType {
	case SHADER_COLOR:
	case SHADER_TEXTURE:
		gl.BindAttribLocation(material.shaderProgram, 1, gl.Str("tex\x00"))
	}
}

func (material *Material) Render() {
	switch material.shaderType {

	case SHADER_COLOR:
		gl.Uniform3fv(gl.GetUniformLocation(material.shaderProgram, gl.Str("materialType\x00")), 1, &SHADER_COLOR_UNI[0])
		gl.Uniform3fv(gl.GetUniformLocation(material.shaderProgram, gl.Str("color\x00")), 1, &material.color[0])
		gl.Uniform1f(gl.GetUniformLocation(material.shaderProgram, gl.Str("shine\x00")), material.shine)

	case SHADER_TEXTURE:
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, *material.texture)
		gl.Uniform3fv(gl.GetUniformLocation(material.shaderProgram, gl.Str("materialType\x00")), 1, &SHADER_TEXTURE_UNI[0])
		gl.Uniform1i(gl.GetUniformLocation(material.shaderProgram, gl.Str("diffuseMap\x00")), 0)
		gl.Uniform1i(gl.GetUniformLocation(material.shaderProgram, gl.Str("cubeDiffuseMap\x00")), 1)
		gl.Uniform1f(gl.GetUniformLocation(material.shaderProgram, gl.Str("shine\x00")), material.shine)
		if material.transparencyEnabled {
			gl.Uniform1i(gl.GetUniformLocation(material.shaderProgram, gl.Str("transparencyEnabled\x00")), 1)
			gl.ActiveTexture(gl.TEXTURE2)
			gl.BindTexture(gl.TEXTURE_2D, *material.transparencyTexture)
			gl.Uniform1i(gl.GetUniformLocation(material.shaderProgram, gl.Str("transparencyMap\x00")), 2)
		} else {
			gl.Uniform1i(gl.GetUniformLocation(material.shaderProgram, gl.Str("transparencyEnabled\x00")), 0)
		}

	case SHADER_CUBEMAP:
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_CUBE_MAP, *material.texture)
		gl.Uniform3fv(gl.GetUniformLocation(material.shaderProgram, gl.Str("materialType\x00")), 1, &SHADER_CUBEMAP_UNI[0])
		gl.Uniform1i(gl.GetUniformLocation(material.shaderProgram, gl.Str("cubeDiffuseMap\x00")), 1)
		gl.Uniform1f(gl.GetUniformLocation(material.shaderProgram, gl.Str("shine\x00")), material.shine)
	}
}

func (material *Material) BecomeColor(rgb []float32) {
	material.shaderType = SHADER_COLOR
	material.color = rgb
}

func (material *Material) BecomeTexture(t *uint32) {
	material.shaderType = SHADER_TEXTURE
	material.texture = t
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, *t)
	gl.Uniform1i(gl.GetUniformLocation(material.shaderProgram, gl.Str("texture0\x00")), int32(0))
}

func (material *Material) BecomeCubemap(c *uint32) {
	material.shaderType = SHADER_CUBEMAP
	material.texture = c
}

func (material *Material) AttachShader(s uint32) {
	material.shaderProgram = s
}

func (material *Material) AttachTransparency(texture *uint32) {
	material.transparencyEnabled = true
	material.transparencyTexture = texture
}

func (material *Material) RemoveTransparency() {
	material.transparencyEnabled = false
	material.transparencyTexture = nil
}

func (materal *Material) GetColor() []float32 {
	return materal.color
}

func (material *Material) GetTexture() *uint32 {
	return material.texture
}
