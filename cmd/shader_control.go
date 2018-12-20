package cmd

import "rapidengine/material"

type ShaderControl struct {
	programs map[string]*material.ShaderProgram
}

func NewShaderControl() ShaderControl {
	return ShaderControl{make(map[string]*material.ShaderProgram)}
}

func (shaderControl *ShaderControl) BindShader(name string) {
	shaderControl.programs[name].Bind()
}

func (shaderControl *ShaderControl) Initialize() {
	shaderControl.programs = map[string]*material.ShaderProgram{
		"basic":    &material.BasicProgram,
		"standard": &material.StandardProgram,
		"skybox":   &material.SkyBoxProgram,
		"terrain":  &material.TerrainProgram,
		"planet":   &material.PlanetProgram,
		"foliage":  &material.FoliageProgram,
		"water":    &material.WaterProgram,
		"sun":      &material.SunProgram,

		"post_final":          &material.PostFinalProgram,
		"post_hdr":            &material.PostHDRProgram,
		"post_horizontal":     &material.PostHorizontalProgram,
		"post_vertical":       &material.PostVerticalProgram,
		"post_prescattering":  &material.PostPreScatteringProgram,
		"post_postscattering": &material.PostPostScatteringProgram,
		"post_prebloom":       &material.PostPreBloomProgram,
		"post_postbloom":      &material.PostPostBloomProgram,
	}
	for _, prog := range shaderControl.programs {
		prog.Compile()
	}
}

func (shaderControl *ShaderControl) GetShader(name string) *material.ShaderProgram {
	return shaderControl.programs[name]
}
