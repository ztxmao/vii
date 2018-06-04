package router

import (
	//	"errors"
	//	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type StaticRouter struct {
	Path         string
	DefaultCname string
	DefaultAname string
	routMap      RouteMaps
}

func NewStaticRouter(path, defaultCname, defaultAname string) *StaticRouter {
	t := new(StaticRouter)
	t.Path = path
	t.DefaultCname = defaultCname
	t.DefaultAname = defaultAname
	t.routMap = make(RouteMaps)
	return t
}

func (this *StaticRouter) GetRouter(r *http.Request) (cName, aName string, contollerType reflect.Type, ok bool) {
	cName = this.DefaultCname
	aName = this.DefaultAname
	ok = true
	methods := make([]string, 2)
	if this.Path != "" {
		matchQuery := strings.Replace(r.URL.Path, this.Path, "", 1)
		if matchQuery == r.URL.Path {
			ok = false
		} else {
			methods = strings.Split(strings.Trim(matchQuery, "/"), "/")
			if len(methods) >= 1 && methods[0] != "" {
				cName = methods[0]
			}
			if len(methods) >= 2 && methods[1] != "" {
				aName = methods[1]
			}
		}
	} else {
		methods = strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(methods) >= 1 && methods[0] != "" {
			cName = methods[0]
		}
		if len(methods) >= 2 && methods[1] != "" {
			aName = methods[1]
		}
	}

	aName = strings.ToLower(aName)
	aName = strings.Title(aName)
	aName = aName + METHOD_EXPORT_TAG
	cName = strings.Title(cName) + "Controller"
	if _, ok = this.routMap[cName]; ok {
		contollerType, ok = this.routMap[cName][aName]
	}
	return
}

func (this *StaticRouter) AddController(c interface{}) {
	reflectVal := reflect.ValueOf(c)
	rt := reflectVal.Type()
	ct := reflect.Indirect(reflectVal).Type()
	//	firstParam := strings.TrimSuffix(ct.Name(), "Controller")
	firstParam := ct.Name()
	if _, ok := this.routMap[firstParam]; ok {
		return
	} else {
		this.routMap[firstParam] = make(map[string]reflect.Type)
	}
	var mname string
	for i := 0; i < rt.NumMethod(); i++ {
		mname = rt.Method(i).Name
		if strings.HasSuffix(mname, METHOD_EXPORT_TAG) {
			this.routMap[firstParam][rt.Method(i).Name] = ct
		}
	}
}
