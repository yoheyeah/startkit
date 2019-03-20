package apis

import (
	"startkit"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

type DBFunc func(obj interface{}) render.JSON

type ValidatorFunc func(api *API) (bool, error)

type CustomHandler func(api *API) error

type API struct {
	*startkit.Context
	Ctx                  *gin.Context
	StaticPath           string
	RelativePath         string
	Method               string
	Headers              map[string]interface{}
	IDer                 string
	Property             string
	IDerPropertyRelation map[string]string
	Params               []string
	Querys               map[string]string
	Request              interface{}
	CustomHandlers       map[string]CustomHandler // replace the general function Run() error in the binding interface
	ValidatorFuncs       []ValidatorFunc
	DBObject             interface{}
	DBResult             interface{}
	DBCreate             []interface{}
	DBUpdate             interface{}
	Structure            interface{}
	// DBRelated      interface{}
}

func Resp(message, data interface{}) interface{} {
	type HTTPResp struct {
		Message interface{} `json:"message,omitempty"`
		Data    interface{} `json:"data,omitempty"`
	}
	return HTTPResp{
		Message: message,
		Data:    data,
	}
}

type BindingAPI interface {
	Run() error
}

func APIType(api API) BindingAPI {
	switch api.Method {
	case "GET":
		return &GET{API: api}
	case "POST":
		return &POST{API: api}
	case "PUT":
		return &PUT{API: api}
	case "DELETE":
		return &DELETE{API: api}
	default:
		return &OPTIONS{API: api}
	}
}

func New(c *startkit.Context) *API {
	return &API{
		Context: c,
	}
}

func (api *API) ReqPath(method, path string) *API {
	api.Method, api.RelativePath = method, path
	return api
}

func (api *API) Static(path string) *API {
	api.StaticPath = path
	return api
}

func (api *API) ID(id string) *API {
	api.IDer = id
	return api
}

func (api *API) Param(params []string) *API {
	api.Params = append(api.Params, params...)
	return api
}

func (api *API) Query(querys map[string]string) *API {
	api.Querys = querys
	return api
}

func (api *API) Attach(attach string) *API {
	api.Property = attach
	return api
}

func (api *API) Relation(relation map[string]string) *API {
	api.IDerPropertyRelation = relation
	return api
}

func (api *API) RandomID(b bool) *API {
	s := strconv.FormatBool(b)
	if len(api.Querys) == 0 {
		api.Querys = map[string]string{"is_random_id": s}
	} else {
		api.Querys["is_random_id"] = s
	}
	return api
}

func (api *API) Pagination(start, limit string) *API {
	if len(api.Querys) == 0 {
		api.Querys = map[string]string{
			"start": start,
			"limit": limit,
		}
		return api
	}
	api.Querys["start"] = start
	api.Querys["limit"] = limit
	return api
}

func (api *API) Validators(fs ...ValidatorFunc) *API {
	api.ValidatorFuncs = append(api.ValidatorFuncs, fs...)
	return api
}

func (api *API) Model(obj interface{}) *API {
	api.Structure = obj
	return api
}

func (api *API) Table(obj interface{}) *API {
	api.DBObject = obj
	return api
}

func (api *API) Find(dbResult interface{}) *API {
	api.DBResult = dbResult
	return api
}

// func (api *API) Related(dbRelated interface{}) *API {
// 	api.DBRelated = dbRelated
// 	return api
// }

func (api *API) Create(dbCreate []interface{}) *API {
	api.DBCreate = dbCreate
	return api
}

func (api *API) Update(dbUpdate []interface{}) *API {
	api.DBUpdate = dbUpdate
	return api
}

func (api *API) Gin(ctx *gin.Context) *API {
	api.Ctx = ctx
	return api
}

func (api *API) Handle(group *gin.RouterGroup) gin.IRoutes {
	if api.IDer != "" {
		api.RelativePath = api.RelativePath + "/:" + api.IDer
	}
	if api.Property != "" {
		api.RelativePath = api.RelativePath + "/" + api.Property
	}
	if count := len(api.Params); count > 0 {
		for i := 0; i < count; i++ {
			if api.Params[i] != "" {
				api.RelativePath = api.RelativePath + "/:" + api.Params[i]
			}
		}
	}
	return group.Handle(api.Method, api.RelativePath, func(ctx *gin.Context) {
		APIType(*api.Gin(ctx)).Run()
		return
	})
}
