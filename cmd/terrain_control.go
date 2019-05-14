package cmd

import (
	"fmt"
	"rapidengine/configuration"
	"rapidengine/geometry"
	"rapidengine/material"
	"rapidengine/terrain"

	"github.com/go-gl/mathgl/mgl32"
)

type TerrainControl struct {
	engine *Engine

	terrainEnabled bool
	root           *terrain.Terrain

	foliages []*terrain.Foliage

	waters []*terrain.Water
}

func NewTerrainControl() TerrainControl {
	return TerrainControl{}
}

func (tc *TerrainControl) Initialize(engine *Engine) {
	tc.engine = engine
}

func (tc *TerrainControl) Update() {
	if tc.terrainEnabled {
		tc.engine.Renderer.RenderChild(tc.root.TChild)

		for _, f := range tc.foliages {
			tc.engine.Renderer.RenderChild(f.FChild)
		}

		for _, w := range tc.waters {
			tc.engine.Renderer.RenderChild(w.WChild)
		}
	}
}

func (tc *TerrainControl) InstanceFoliage(f *terrain.Foliage) {
	tc.foliages = append(tc.foliages, f)
}

func (tc *TerrainControl) InstanceWater(w *terrain.Water) {
	tc.waters = append(tc.waters, w)
}

func (tc *TerrainControl) NewTerrain(width int, height int, vertices int) *terrain.Terrain {
	t := terrain.NewTerrain(width, height)

	t.TChild = tc.engine.ChildControl.NewChild3D()

	t.TChild.AttachModel(
		geometry.Model{
			Meshes:    []geometry.Mesh{geometry.NewPlane(width, height, vertices, nil, 1)},
			Materials: map[int]material.Material{0: tc.engine.MaterialControl.NewTerrainMaterial()},
		},
	)

	t.TChild.Model.Meshes[0].TesselationEnabled = true

	t.TChild.SetPosition(0, 1, 0)
	t.TChild.PreRender(tc.engine.Renderer.MainCamera)

	tc.terrainEnabled = true
	tc.root = &t

	return &t
}

func (tc *TerrainControl) NewPlanetaryTerrain(width int, height int, vertices int) *terrain.Terrain {
	t := terrain.NewTerrain(width, height)

	t.TChild = tc.engine.ChildControl.NewChild3D()

	t.TChild.AttachMaterial(tc.engine.MaterialControl.NewTerrainMaterial())
	//t.TChild.AttachMesh(geometry.NewPlane(width, height, vertices, nil, 1))
	t.TChild.AttachMesh(geometry.LoadObj("../rapidengine/assets/obj/sphere.obj", 10000))
	t.TChild.SetInstanceRenderDistance(1000000000)

	t.TChild.PreRender(tc.engine.Renderer.MainCamera)

	tc.terrainEnabled = true
	tc.root = &t

	return &t
}

func (tc *TerrainControl) NewFoliage(width int, height int, instances int) *terrain.Foliage {
	f := terrain.NewFoliage(width, height)

	f.FChild = tc.engine.ChildControl.NewChild3D()

	f.FChild.AttachModel(
		tc.engine.GeometryControl.LoadModel("./billboard.obj", tc.engine.MaterialControl.NewFoliageMaterial()),
	)

	f.FChild.Model.EnableInstancing(instances)
	f.FChild.Model.Meshes[0].InstancingEnabled = true
	f.FChild.Model.Meshes[0].NumInstances = instances
	f.FChild.SetInstanceRenderDistance(100000)

	f.FChild.Model.Meshes[0].ComputeTangents()

	f.FChild.Model.Meshes[0].TexCoordsEnabled = true
	f.FChild.Model.Meshes[0].NormalsEnabled = true
	f.FChild.Model.Meshes[0].TangentsEnabled = true
	f.FChild.Model.Meshes[0].BitangentsEnabled = true

	f.FChild.PreRender(tc.engine.Renderer.MainCamera)

	return &f
}

func (tc *TerrainControl) NewWater(width int, height int, vertices int) *terrain.Water {
	w := terrain.NewWater(width, height)

	w.WChild = tc.engine.ChildControl.NewChild3D()

	w.WChild.AttachMaterial(tc.engine.MaterialControl.NewWaterMaterial())
	w.WChild.AttachMesh(geometry.NewPlane(width, height, vertices, nil, 1))

	w.WChild.PreRender(tc.engine.Renderer.MainCamera)

	return &w
}

func (terrainControl *TerrainControl) NewSkyBox(
	path string,
	name string,
	ext string,

	shaderControl *ShaderControl,
	textureControl *TextureControl,

	config *configuration.EngineConfig,
) *terrain.SkyBox {

	shaderControl.GetShader("skybox").Bind()

	textureControl.NewCubeMap(
		fmt.Sprintf("%s/%s/%s_LF.%s", path, name, name, ext),
		fmt.Sprintf("%s/%s/%s_RT.%s", path, name, name, ext),
		fmt.Sprintf("%s/%s/%s_UP.%s", path, name, name, ext),
		fmt.Sprintf("%s/%s/%s_DN.%s", path, name, name, ext),
		fmt.Sprintf("%s/%s/%s_FR.%s", path, name, name, ext),
		fmt.Sprintf("%s/%s/%s_BK.%s", path, name, name, ext),

		"skybox",
	)

	cmaterial := terrainControl.engine.MaterialControl.NewCubemapMaterial()
	cmaterial.CubeDiffuseMap = textureControl.GetTexture("skybox").Addr

	indices := []uint32{}
	for i := 0; i < len(terrain.SkyBoxVertices); i++ {
		indices = append(indices, uint32(i))
	}

	vao := geometry.NewVertexArray(terrain.SkyBoxVertices, indices)
	vao.AddVertexAttribute(geometry.CubeTextures, 1, 2)

	return terrain.NewSkyBox(
		cmaterial,
		vao,
		mgl32.Perspective(
			mgl32.DegToRad(45),
			float32(config.ScreenWidth)/float32(config.ScreenHeight),
			0.1, 100,
		),
		mgl32.Ident4(),
		[]*material.ShaderProgram{
			terrainControl.engine.ShaderControl.GetShader("standard"),
		},
	)
}
