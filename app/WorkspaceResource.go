package app

import (
	"net/http"

	"github.com/emicklei/go-restful"

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
		proj.Id = name
		proj.Href = entity.Href + "/" + proj.Id
	}

	response.WriteEntity(entity)
}

// POST /projects
// <User><Name>Melissa</Name></User>
//
func (resource *WorkspaceResource) createProject(request *restful.Request, response *restful.Response) {
	entityTemplate := new(model.ProjectTemplate)
	err := request.ReadEntity(entityTemplate)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	_, prjErr := resource.ws.NewProject(entityTemplate.Id)
	if prjErr != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	entity := new(model.Project)
	entity.Id = entityTemplate.Id
	entity.Href = "/projects/" + entity.Id

	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(entity)
}
