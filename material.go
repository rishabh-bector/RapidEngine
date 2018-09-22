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

	color []float32
	shine float32
}

func NewMaterial(program uint32) Material {
	return Material{
		shaderProgram: program,
		shaderType:    SHADER_COLOR,
		color:         []float32{1, 1, 1},
		shine:         0.8,
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
	gl.UseProgram(material.shaderProgram)
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
		gl.Uniform1f(gl.GetUniformLocation(material.shaderProgram, gl.Str("shine\x00")), material.shine)

	case SHADER_CUBEMAP:
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_CUBE_MAP, *material.texture)
		gl.Uniform3fv(gl.GetUniformLocation(material.shaderProgram, gl.Str("materialType\x00")), 1, &SHADER_CUBEMAP_UNI[0])
		gl.Uniform1i(gl.GetUniformLocation(material.shaderProgram, gl.Str("cubeMap\x00")), 0)
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

func (materal *Material) GetColor() []float32 {
	return materal.color
}

func (material *Material) GetTexture() *uint32 {
	return material.texture
}