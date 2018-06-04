package router

import (
	"net/http"
	"reflect"
)

type Router interface {
	AddController(interface{})
	GetRouter(r *http.Request) (cName, aName string, contollerType reflect.Type, ok bool)
}

type RouteMaps map[string]map[string]reflect.Type

//controller中以此结尾的方法会参与路由
const METHOD_EXPORT_TAG = "Action"
