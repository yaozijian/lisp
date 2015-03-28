package main

import (
	"errors"
	"fmt"
)

func (rt *lisp_runtime) print(array ...*lisp_runtime_val) (ret interface{}, err error) {

	if len(array) == 0 {
		for k, v := range rt.vars {
			fmt.Printf("%s: %v\n", k, v)
		}
	} else {
		for _, item := range array {
			if len(item.name) == 0 {
				fmt.Printf("%v\n", item.val)
			} else {
				fmt.Printf("%s: %v\n", item.name, item.val)
			}
		}
	}

	return
}

func (rt *lisp_runtime) set(array ...*lisp_runtime_val) (ret interface{}, err error) {

	var varname string

	if len(array) != 2 {
		err = errors.New("set函数只接受两个操作数")
	} else if varname = array[0].name; len(varname) == 0 {
		if varname, _ = array[0].val.(string); len(varname) == 0 {
			err = errors.New("set函数的第一个参数类型不正确")
		}
	}

	if err != nil {
		return
	}

	// 设置第一个参数等于一个已有变量的值
	if len(array[1].name) > 0 {
		rt.vars[varname] = array[1].val
		ret = varname
	} else {
		// 设置第一个参数等于一个字面值
		switch val := array[1].val.(type) {
		case int64:
			rt.vars[varname] = val
			ret = varname
		case string:
			if val == "nil" {
				delete(rt.vars, varname)
			} else {
				fmt.Println(array[1])
				err = fmt.Errorf("set函数的第二个参数类型不正确: %s", val)
			}
		default:
			err = fmt.Errorf("set函数的第二个参数类型不正确")
		}
	}

	return
}

func (rt *lisp_runtime) add(array ...*lisp_runtime_val) (ret interface{}, err error) {
	if len(array) < 2 {
		err = errors.New("add函数的操作数太少")
	} else {
		var retval int64
		for _, val := range array {
			if intval, covok := val.val.(int64); covok {
				retval += intval
			} else {
				err = errors.New("add函数的参数类型不正确")
				break
			}
		}
		if err == nil {
			ret = retval
		}
	}
	return
}

func (rt *lisp_runtime) sub(array ...*lisp_runtime_val) (ret interface{}, err error) {
	if len(array) < 2 {
		err = errors.New("sub函数的操作数太少")
	} else if intval, covok := array[0].val.(int64); covok {
		retval := intval
		for _, val := range array[1:] {
			if intval, covok := val.val.(int64); covok {
				retval -= intval
			} else {
				err = errors.New("sub函数的参数类型不正确")
				break
			}
		}
		if err == nil {
			ret = retval
		}
	} else {
		err = errors.New("sub函数的参数类型不正确")
	}
	return
}

func (rt *lisp_runtime) mul(array ...*lisp_runtime_val) (ret interface{}, err error) {
	if len(array) < 2 {
		err = errors.New("mul函数的操作数太少")
	} else {
		retval := int64(1)
		for _, val := range array {
			if intval, covok := val.val.(int64); covok {
				retval *= intval
			} else {
				err = errors.New("mul函数的参数类型不正确")
				break
			}
		}
		if err == nil {
			ret = retval
		}
	}
	return
}

func (rt *lisp_runtime) div(array ...*lisp_runtime_val) (ret interface{}, err error) {
	if len(array) < 2 {
		err = errors.New("div函数的操作数太少")
	} else if intval, covok := array[0].val.(int64); covok {
		retval := intval
		for _, val := range array[1:] {
			if intval, covok := val.val.(int64); covok {
				retval /= intval
			} else {
				err = errors.New("div函数的参数类型不正确")
				break
			}
		}
		if err == nil {
			ret = retval
		}
	} else {
		err = errors.New("div函数的参数类型不正确")
	}
	return
}
