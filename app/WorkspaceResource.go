package app

import (
	"fmt"
	"net/http"
	"strconv"

	"image/color"
	"image/png"

	"github.com/emicklei/go-restful"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/image"
	core "github.com/inkyblackness/shocked-core"
	model "github.com/inkyblackness/shocked-model"
)

type WorkspaceResource struct {
	ws *core.Workspace
}

func NewWorkspaceResource(container *restful.Container, workspace *core.Workspace) *WorkspaceResource {
	resource := &WorkspaceResource{
		ws: workspace}

	service1 := new(restful.WebService)

	service1.
		Path("/ws").
		Doc("Manage workspace").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	service1.Route(service1.GET("").To(resource.getWorkspace).
		// docs
		Doc("get current workspace").
		Operation("getWorkspace").
		Writes(model.Workspace{}))

	container.Add(service1)

	service2 := new(restful.WebService)

	service2.
		Path("/projects").
		Doc("Manage projects").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	service2.Route(service2.GET("").To(resource.getProjects).
		// docs
		Doc("get current projects").
		Operation("getWorkspace").
		Writes(model.Projects{}))

	service2.Route(service2.POST("").To(resource.createProject).
		// docs
		Doc("create a project").
		Operation("createProject").
		Reads(model.ProjectTemplate{}).
		Writes(model.Project{}))

	service2.Route(service2.GET("{project-id}/textures").To(resource.getTextures).
		// docs
		Doc("get textures").
		Operation("getTextures").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Writes(model.Textures{}))

	service2.Route(service2.GET("{project-id}/textures/{texture-id}").To(resource.getTexture).
		// docs
		Doc("get texture").
		Operation("getTexture").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("texture-id", "identifier of the texture").DataType("int")).
		Writes(model.Texture{}))

	service2.Route(service2.PUT("{project-id}/textures/{texture-id}").To(resource.setTexture).
		// docs
		Doc("set texture").
		Operation("setTexture").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("texture-id", "identifier of the texture").DataType("int")).
		Reads(model.TextureProperties{}).
		Writes(model.Texture{}))

	service2.Route(service2.GET("{project-id}/textures/{texture-id}/{texture-size}").To(resource.getTextureImage).
		// docs
		Doc("get texture image").
		Operation("getTextureImage").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("texture-id", "identifier of the texture").DataType("int")).
		Param(service2.PathParameter("texture-size", "Size of the texture").DataType("string")).
		Writes(model.Image{}))

	service2.Route(service2.GET("{project-id}/textures/{texture-id}/{texture-size}/png").To(resource.getTextureImageExport).
		// docs
		Doc("get texture image as PNG").
		Operation("getTextureImageExport").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("texture-id", "identifier of the texture").DataType("int")).
		Param(service2.PathParameter("texture-size", "Size of the texture").DataType("string")).
		Produces("image/png"))

	service2.Route(service2.GET("{project-id}/objects/{class}/{subclass}/{type}").To(resource.getGameObject).
		// docs
		Doc("get game object").
		Operation("getGameObject").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("class", "identifier of the class").DataType("int")).
		Param(service2.PathParameter("subclass", "identifier of the class").DataType("int")).
		Param(service2.PathParameter("type", "identifier of the class").DataType("int")).
		Writes(model.GameObject{}))

	service2.Route(service2.GET("{project-id}/archive/levels").To(resource.getLevels).
		// docs
		Doc("get level list").
		Operation("getLevels").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Writes(model.Levels{}))

	service2.Route(service2.GET("{project-id}/archive/levels/{level-id}").To(resource.getLevel).
		// docs
		Doc("get level information").
		Operation("getLevel").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("level-id", "identifier of the level").DataType("int")).
		Writes(model.Level{}))

	service2.Route(service2.GET("{project-id}/archive/levels/{level-id}/textures").To(resource.getLevelTextures).
		// docs
		Doc("get level textures").
		Operation("getLevelTextures").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("level-id", "identifier of the level").DataType("int")).
		Writes(model.LevelTextures{}))

	service2.Route(service2.PUT("{project-id}/archive/levels/{level-id}/textures").To(resource.setLevelTextures).
		// docs
		Doc("put level textures").
		Operation("setLevelTextures").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("level-id", "identifier of the level").DataType("int")).
		Reads([]string{}).
		Writes(model.LevelTextures{}))

	service2.Route(service2.GET("{project-id}/archive/levels/{level-id}/tiles").To(resource.getLevelTiles).
		// docs
		Doc("get level tiles").
		Operation("getLevelTiles").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("level-id", "identifier of the level").DataType("int")).
		Writes(model.Tiles{}))

	service2.Route(service2.GET("{project-id}/archive/levels/{level-id}/tiles/{y}/{x}").To(resource.getLevelTile).
		// docs
		Doc("get level tile").
		Operation("getLevelTile").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("level-id", "identifier of the level").DataType("int")).
		Param(service2.PathParameter("y", "Y coordinate of the tile").DataType("int")).
		Param(service2.PathParameter("x", "X coordinate of the tile").DataType("int")).
		Writes(model.Tile{}))

	service2.Route(service2.PUT("{project-id}/archive/levels/{level-id}/tiles/{y}/{x}").To(resource.setLevelTile).
		// docs
		Doc("set level tile").
		Operation("setLevelTile").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("level-id", "identifier of the level").DataType("int")).
		Param(service2.PathParameter("y", "Y coordinate of the tile").DataType("int")).
		Param(service2.PathParameter("x", "X coordinate of the tile").DataType("int")).
		Reads(model.TileProperties{}).
		Writes(model.Tile{}))

	service2.Route(service2.GET("{project-id}/archive/levels/{level-id}/objects").To(resource.getLevelObjects).
		// docs
		Doc("get level objects").
		Operation("getLevelObjects").
		Param(service2.PathParameter("project-id", "identifier of the project").DataType("string")).
		Param(service2.PathParameter("level-id", "identifier of the level").DataType("int")).
		Writes(model.LevelObjects{}))

	container.Add(service2)

	return resource
}

// GET /ws
func (resource *WorkspaceResource) getWorkspace(request *restful.Request, response *restful.Response) {
	var entity model.Workspace
	entity.Href = "/"

	entity.Projects.Href = "/projects"

	response.WriteEntity(entity)
}

// GET /projects
func (resource *WorkspaceResource) getProjects(request *restful.Request, response *restful.Response) {
	projectNames := resource.ws.ProjectNames()
	var entity model.Projects
	entity.Href = "/projects"

	entity.Items = make([]model.Identifiable, len(projectNames))
	for index, name := range projectNames {
		proj := &entity.Items[index]
		proj.ID = name
		proj.Href = entity.Href + "/" + proj.ID
	}

	response.WriteEntity(entity)
}

// POST /projects
func (resource *WorkspaceResource) createProject(request *restful.Request, response *restful.Response) {
	entityTemplate := new(model.ProjectTemplate)
	err := request.ReadEntity(entityTemplate)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	_, prjErr := resource.ws.NewProject(entityTemplate.ID)
	if prjErr != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	entity := new(model.Project)
	entity.ID = entityTemplate.ID
	entity.Href = "/projects/" + entity.ID

	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(entity)
}

// GET /projects/{project-id}/textures
func (resource *WorkspaceResource) getTextures(request *restful.Request, response *restful.Response) {
	projectId := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectId)

	if err == nil {
		textures := project.Textures()
		limit := textures.TextureCount()
		var entity model.Textures

		entity.List = make([]model.Texture, limit)
		for id := 0; id < limit; id++ {
			entity.List[id] = resource.textureEntity(project, id)
		}

		response.WriteEntity(entity)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

// GET /projects/{project-id}/textures/{texture-id}
func (resource *WorkspaceResource) getTexture(request *restful.Request, response *restful.Response) {
	projectId := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectId)

	if err == nil {
		textureId, _ := strconv.ParseInt(request.PathParameter("texture-id"), 10, 16)
		entity := resource.textureEntity(project, int(textureId))

		response.WriteEntity(entity)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

// PUT /projects/{project-id}/textures/{texture-id}
func (resource *WorkspaceResource) setTexture(request *restful.Request, response *restful.Response) {
	projectId := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectId)

	if err == nil {
		textureId, _ := strconv.ParseInt(request.PathParameter("texture-id"), 10, 16)
		var properties model.TextureProperties
		err = request.ReadEntity(&properties)
		if err != nil {
			response.AddHeader("Content-Type", "text/plain")
			response.WriteErrorString(http.StatusInternalServerError, err.Error())
			return
		}

		project.Textures().SetProperties(int(textureId), properties)
		entity := resource.textureEntity(project, int(textureId))

		response.WriteEntity(entity)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

func (resource *WorkspaceResource) textureEntity(project *core.Project, textureId int) (entity model.Texture) {
	entity.ID = fmt.Sprintf("%d", textureId)
	entity.Href = "/projects/" + project.Name() + "/textures/" + entity.ID
	entity.Properties = project.Textures().Properties(textureId)
	for _, size := range model.TextureSizes() {
		entity.Images = append(entity.Images, model.Link{Rel: string(size), Href: entity.Href + "/" + string(size)})
	}

	return
}

// GET /projects/{project-id}/textures/{texture-id}/{texture-size}
func (resource *WorkspaceResource) getTextureImage(request *restful.Request, response *restful.Response) {
	projectId := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectId)

	if err == nil {
		textureId, _ := strconv.ParseInt(request.PathParameter("texture-id"), 10, 16)
		textureSize := request.PathParameter("texture-size")
		var entity model.Image

		entity.Href = "/projects/" + projectId + "/textures/" + fmt.Sprintf("%d", textureId) + "/" + textureSize
		bmp := project.Textures().Image(int(textureId), model.TextureSize(textureSize))
		hotspot := bmp.Hotspot()

		entity.Properties.HotspotLeft = hotspot.Min.X
		entity.Properties.HotspotTop = hotspot.Min.Y
		entity.Properties.HotspotRight = hotspot.Max.X
		entity.Properties.HotspotBottom = hotspot.Max.Y

		entity.Formats = []model.Link{model.Link{Rel: "png", Href: entity.Href + "/png"}}

		response.WriteEntity(entity)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

// GET /projects/{project-id}/textures/{texture-id}/{texture-size}/png
func (resource *WorkspaceResource) getTextureImageExport(request *restful.Request, response *restful.Response) {
	projectId := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectId)

	if err == nil {
		textureId, _ := strconv.ParseInt(request.PathParameter("texture-id"), 10, 16)
		textureSize := request.PathParameter("texture-size")
		var palette color.Palette

		bmp := project.Textures().Image(int(textureId), model.TextureSize(textureSize))
		palette, err = project.Palettes().GamePalette()
		image := image.FromBitmap(bmp, palette)

		response.AddHeader("Content-Type", "image/png")
		png.Encode(response.ResponseWriter, image)

	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

// GET /projects/{project-id}/archive/levels
func (resource *WorkspaceResource) getLevels(request *restful.Request, response *restful.Response) {
	projectId := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectId)

	if err == nil {
		var entity model.Levels
		archive := project.Archive()
		levelIDs := archive.LevelIDs()

		entity.Href = "/projects/" + projectId + "/archive/levels"
		for _, id := range levelIDs {
			entry := resource.getLevelEntity(project, archive, id)

			entity.List = append(entity.List, entry)
		}

		response.WriteEntity(entity)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

// GET /projects/{project-id}/archive/levels/{level-id}
func (resource *WorkspaceResource) getLevel(request *restful.Request, response *restful.Response) {
	projectId := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectId)

	if err == nil {
		levelId, _ := strconv.ParseInt(request.PathParameter("level-id"), 10, 16)
		entity := resource.getLevelEntity(project, project.Archive(), int(levelId))

		response.WriteEntity(entity)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

func (resource *WorkspaceResource) getLevelEntity(project *core.Project, archive *core.Archive, levelId int) (entity model.Level) {
	entity.ID = fmt.Sprintf("%d", levelId)
	entity.Href = "/projects/" + project.Name() + "/archive/levels/" + entity.ID
	level := archive.Level(levelId)
	entity.Properties = level.Properties()

	entity.Links = []model.Link{}
	entity.Links = append(entity.Links, model.Link{Rel: "tiles", Href: entity.Href + "/tiles/{y}/{x}"})
	if !entity.Properties.CyberspaceFlag {
		entity.Links = append(entity.Links, model.Link{Rel: "textures", Href: entity.Href + "/textures"})
	}

	return
}

// GET /projects/{project-id}/archive/levels/{level-id}/textures
func (resource *WorkspaceResource) getLevelTextures(request *restful.Request, response *restful.Response) {
	projectId := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectId)

	if err == nil {
		levelId, _ := strconv.ParseInt(request.PathParameter("level-id"), 10, 16)
		level := project.Archive().Level(int(levelId))
		entity := resource.getLevelTexturesEntity(projectId, level)

		response.WriteEntity(entity)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

func (resource *WorkspaceResource) getLevelTexturesEntity(projectId string, level *core.Level) (entity model.LevelTextures) {
	entity.Href = "/projects/" + projectId + "/archive/levels/" + fmt.Sprintf("%d", level.ID()) + "/textures"
	for _, id := range level.Textures() {
		entity.IDs = append(entity.IDs, fmt.Sprintf("%d", id))
	}

	return
}

// PUT /projects/{project-id}/archive/levels/{level-id}/textures
func (resource *WorkspaceResource) setLevelTextures(request *restful.Request, response *restful.Response) {
	projectId := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectId)

	if err == nil {
		levelId, _ := strconv.ParseInt(request.PathParameter("level-id"), 10, 16)

		var idStrings []string
		err = request.ReadEntity(&idStrings)
		if err != nil {
			response.AddHeader("Content-Type", "text/plain")
			response.WriteErrorString(http.StatusInternalServerError, err.Error())
			return
		}

		newIds := make([]int, len(idStrings))
		for index, idString := range idStrings {
			parsedId, _ := strconv.ParseInt(idString, 10, 16)
			newIds[index] = int(parsedId)
		}

		level := project.Archive().Level(int(levelId))
		level.SetTextures(newIds)

		entity := resource.getLevelTexturesEntity(projectId, level)
		response.WriteEntity(entity)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

// GET /projects/{project-id}/archive/levels/{level-id}/tiles
func (resource *WorkspaceResource) getLevelTiles(request *restful.Request, response *restful.Response) {
	projectId := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectId)

	if err == nil {
		levelId, _ := strconv.ParseInt(request.PathParameter("level-id"), 10, 16)
		level := project.Archive().Level(int(levelId))
		var entity model.Tiles

		entity.Table = make([][]model.Tile, 64)
		for y := 0; y < 64; y++ {
			entity.Table[y] = make([]model.Tile, 64)
			for x := 0; x < 64; x++ {
				entity.Table[y][x] = getLevelTileEntity(project, level, x, y)
			}
		}

		response.WriteEntity(entity)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

func getLevelTileEntity(project *core.Project, level *core.Level, x int, y int) (entity model.Tile) {
	entity.Href = "/projects/" + project.Name() + "/archive/levels/" + fmt.Sprintf("%d", level.ID()) +
		"/tiles/" + fmt.Sprintf("%d", y) + "/" + fmt.Sprintf("%d", x)
	entity.Properties = level.TileProperties(int(x), int(y))

	return
}

// GET /projects/{project-id}/archive/levels/{level-id}/tiles/{y}/{x}
func (resource *WorkspaceResource) getLevelTile(request *restful.Request, response *restful.Response) {
	projectId := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectId)

	if err == nil {
		x, _ := strconv.ParseInt(request.PathParameter("x"), 10, 16)
		y, _ := strconv.ParseInt(request.PathParameter("y"), 10, 16)
		levelId, _ := strconv.ParseInt(request.PathParameter("level-id"), 10, 16)
		level := project.Archive().Level(int(levelId))

		response.WriteEntity(getLevelTileEntity(project, level, int(x), int(y)))
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

// PUT /projects/{project-id}/archive/levels/{level-id}/tiles/{y}/{x}
func (resource *WorkspaceResource) setLevelTile(request *restful.Request, response *restful.Response) {
	projectId := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectId)

	if err == nil {
		x, _ := strconv.ParseInt(request.PathParameter("x"), 10, 16)
		y, _ := strconv.ParseInt(request.PathParameter("y"), 10, 16)
		levelId, _ := strconv.ParseInt(request.PathParameter("level-id"), 10, 16)
		level := project.Archive().Level(int(levelId))

		var properties model.TileProperties
		err = request.ReadEntity(&properties)
		if err != nil {
			response.AddHeader("Content-Type", "text/plain")
			response.WriteErrorString(http.StatusInternalServerError, err.Error())
			return
		}

		level.SetTileProperties(int(x), int(y), properties)
		response.WriteEntity(getLevelTileEntity(project, level, int(x), int(y)))
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

// GET /projects/{project-id}/archive/levels/{level-id}/objects
func (resource *WorkspaceResource) getLevelObjects(request *restful.Request, response *restful.Response) {
	projectId := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectId)

	if err == nil {
		levelId, _ := strconv.ParseInt(request.PathParameter("level-id"), 10, 16)
		level := project.Archive().Level(int(levelId))
		hrefBase := "/projects/" + projectId + "/archive/levels/" + fmt.Sprintf("%d", levelId) + "/objects/"
		var entity model.LevelObjects

		entity.Table = level.Objects()
		for i := 0; i < len(entity.Table); i++ {
			entry := &entity.Table[i]
			entry.Href = hrefBase + entry.ID

			entry.Links = append(entry.Links, model.Link{
				Rel:  "static",
				Href: "/projects/" + projectId + "/objects/" + fmt.Sprintf("%d/%d/%d", entry.Class, entry.Subclass, entry.Type)})
		}

		response.WriteEntity(entity)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

// GET /projects/{project-id}/objects/{class}/{subclass}/{type}
func (resource *WorkspaceResource) getGameObject(request *restful.Request, response *restful.Response) {
	projectId := request.PathParameter("project-id")
	project, err := resource.ws.Project(projectId)

	if err == nil {
		classId, _ := strconv.ParseInt(request.PathParameter("class"), 10, 8)
		subclassId, _ := strconv.ParseInt(request.PathParameter("subclass"), 10, 8)
		typeId, _ := strconv.ParseInt(request.PathParameter("type"), 10, 8)
		objId := res.MakeObjectID(res.ObjectClass(classId), res.ObjectSubclass(subclassId), res.ObjectType(typeId))
		entity := resource.objectEntity(project, objId)

		response.WriteEntity(entity)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

func (resource *WorkspaceResource) objectEntity(project *core.Project, objId res.ObjectID) (entity model.GameObject) {
	entity.ID = fmt.Sprintf("%d/%d/%d", objId.Class, objId.Subclass, objId.Type)
	entity.Href = "/projects/" + project.Name() + "/objects/" + entity.ID
	entity.Properties = project.GameObjects().Properties(objId)

	return
}
