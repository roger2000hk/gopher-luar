package luar

import (
	"reflect"

	"github.com/yuin/gopher-lua"
)

func funcEvaluate(L *lua.LState, fn reflect.Value) int {
	fnType := fn.Type()
	top := L.GetTop()
	expected := fnType.NumIn()
	variadic := fnType.IsVariadic()
	if (!variadic && top != expected) || (variadic && top < expected-1) {
		panic("invalid number of function arguments")
	}
	args := make([]reflect.Value, top)
	for i := 0; i < L.GetTop(); i++ {
		var hint reflect.Type
		if variadic && i >= expected-1 {
			hint = fnType.In(expected - 1).Elem()
		} else {
			hint = fnType.In(i)
		}
		args[i] = lValueToReflect(L.Get(i+1), hint)
	}
	ret := fn.Call(args)
	for _, val := range ret {
		L.Push(New(L, val.Interface()))
	}
	return len(ret)
}

func funcWrapper(L *lua.LState, fn reflect.Value) *lua.LFunction {
	wrapper := func(L *lua.LState) int {
		return funcEvaluate(L, fn)
	}
	return L.NewFunction(wrapper)
}
