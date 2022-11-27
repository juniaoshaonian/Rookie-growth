package rpc_demo

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"sync/atomic"
)

var MessId uint32

func InitClientProxy(service Service, p Proxy) error {
	serviceVal := reflect.ValueOf(service)
	serviceTyp := reflect.TypeOf(service)
	if serviceTyp.Kind() != reflect.Ptr || serviceTyp.Elem().Kind() != reflect.Struct {
		return errors.New("只支持一级指针")
	}
	serviceVal = serviceVal.Elem()
	serviceTyp = serviceTyp.Elem()
	numField := serviceTyp.NumField()
	for i := 0; i < numField; i++ {
		//篡改每个函数
		fieldTyp := serviceTyp.Field(i)
		fieldVal := serviceVal.Field(i)
		if !fieldVal.CanSet() {
			continue
		}
		fn := reflect.MakeFunc(fieldTyp.Type, func(args []reflect.Value) (results []reflect.Value) {
			// 读取请求
			respType := fieldTyp.Type.Out(0)
			ctx := args[0].Interface().(context.Context)
			res := args[1].Interface()
			data, err := json.Marshal(res)
			if err != nil {
				results = append(results, reflect.Zero(respType))
				results = append(results, reflect.ValueOf(err))
				return
			}
			atomic.AddUint32(&MessId, 1)
			req := &Request{
				BodyLength:  uint32(len(data)),
				MessId:      MessId,
				Compresser:  0,
				Version:     0,
				Serializer:  0,
				ServiceName: service.Name(),
				MethodName:  fieldTyp.Name,
				Arg:         data,
			}
			CalRequestHead(req)

			resp, err := p.Invoke(ctx, req)
			response := reflect.New(respType).Interface()
			json.Unmarshal(resp.Data, response)
			if err != nil {
				results = append(results, reflect.Zero(respType))
				results = append(results, reflect.ValueOf(err))
				return
			}

			results = append(results, reflect.ValueOf(response).Elem())
			results = append(results, reflect.Zero(reflect.TypeOf(new(error)).Elem()))
			return
		})

		fieldVal.Set(fn)
	}
	return nil
}

type Service interface {
	Name() string
}
