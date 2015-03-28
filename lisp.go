package main

import (
	"bufio"
	"container/list"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	left_brace     = "("
	right_brace    = ")"
	operator_add   = "+"
	operator_sub   = "-"
	operator_mul   = "*"
	operator_div   = "/"
	operator_set   = "set"
	operator_print = "print"
)

type lisp_runtime struct {
	vars map[string]interface{}
}

type lisp_runtime_val struct {
	name string
	val  interface{}
}

type lisp_runtime_func func(rt *lisp_runtime, array ...*lisp_runtime_val) (ret interface{}, err error)

func (rt *lisp_runtime) init() {
	rt.vars = make(map[string]interface{})
	rt.vars[operator_print] = (*lisp_runtime).print
	rt.vars[operator_set] = (*lisp_runtime).set
	rt.vars[operator_add] = (*lisp_runtime).add
	rt.vars[operator_sub] = (*lisp_runtime).sub
	rt.vars[operator_mul] = (*lisp_runtime).mul
	rt.vars[operator_div] = (*lisp_runtime).div
}

func main() {

	var lisp lisp_runtime
	lisp.init()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(">")
		line, _, _ := reader.ReadLine()
		codeline := strings.TrimSpace(string(line))

		switch {
		case len(codeline) == 0:
		case codeline == "exit":
			return
		default:
			// 解析行,输出结果
			val, err := lisp.parseLine(codeline)
			if err == nil {
				if val != nil {
					switch retval := val.(type) {
					case string:
						//  结果是字符串时，可能表示某个变量名
						if rtobj, _ := lisp.vars[retval]; rtobj != nil {
							fmt.Printf("%s: %v\n", retval, rtobj)
						} else {
							fmt.Println(retval)
						}
					default:
						fmt.Println(val)
					}
				}
			} else {
				fmt.Println(err)
			}
		}
	}
}

func (rt *lisp_runtime) parseLine(line string) (val interface{}, err error) {

	line = strings.TrimSpace(line)

	var token string
	var state int

	stack := list.New()

	for _, char := range line {
		str := string(char)
		switch {
		// 左括号
		case str == left_brace:
			state = -1
			stack.PushBack(left_brace)
		// 右括号
		case str == right_brace:
			if len(token) > 0 {
				stack.PushBack(token)
				token = ""
			}
			// 遇到右括号时计算列表的值
			state = 0
			val, err = rt.evalLastList(stack)
		default:
			state = 1
			// 空格是参数分隔符
			if int32(char) == 32 {
				if len(token) > 0 {
					stack.PushBack(token)
					token = ""
				}
			} else {
				token += str
			}
		}
	}

	if state != 0 {
		err = errors.New("代码行必须以右括号结束")
		val = nil
	}

	return
}

// 对列表进行估值
func (rt *lisp_runtime) evalLastList(stack *list.List) (val interface{}, err error) {

	err = errors.New("语法错误: 找不到匹配的左括号")

	if stack == nil || stack.Len() == 0 {
		return
	}

	var listbegin *list.Element

	ptr := stack.Back()
	for ptr != nil {
		if item, _ := ptr.Value.(string); item == left_brace {

			// 记录列表开始位置
			listbegin = ptr

			// 列表的第一个元素是运算符(函数)
			operator := ptr.Next()

			if operator == nil {
				// 空列表是合法的
				err = nil
				break
			} else if fn, exist := rt.vars[operator.Value.(string)]; !exist {
				err = errors.New(fmt.Sprintf("符号 %s 未定义", operator.Value.(string)))
				break
			} else if rtfunc, _ := fn.(func(*lisp_runtime, ...*lisp_runtime_val) (interface{}, error)); rtfunc == nil {
				err = errors.New(fmt.Sprintf("符号 %s 不是函数", operator.Value.(string)))
				break
			} else {
				val, err = rt.callRuntimeFunc(rtfunc, operator.Next())
				break
			}
		} else {
			ptr = ptr.Prev()
		}
	}

	// 如果没有错误,则弹出最后的列表,压入估值结果
	if err == nil {

		ptr = stack.Back()
		for ptr != listbegin {
			stack.Remove(ptr)
			ptr = stack.Back()
		}
		stack.Remove(ptr)

		switch retval := val.(type) {
		case nil:
		case int64:
			stack.PushBack(fmt.Sprintf("%d", retval))
		default:
			stack.PushBack(retval)
		}
	}

	return
}

// 调用运行时函数
func (rt *lisp_runtime) callRuntimeFunc(rtfunc lisp_runtime_func, params *list.Element) (val interface{}, err error) {

	array := []*lisp_runtime_val{}

	for ptr := params; ptr != nil; ptr = ptr.Next() {
		strarg, _ := ptr.Value.(string)
		runtime_val := &lisp_runtime_val{}
		// 参数是已经存在的运行时变量
		if rtobj, exist := rt.vars[strarg]; exist {
			runtime_val.name = strarg
			runtime_val.val = rtobj
			// 参数是整数
		} else if intarg, converr := strconv.ParseInt(strarg, 10, 64); converr == nil {
			runtime_val.val = intarg
			// 参数是字符串
		} else {
			runtime_val.val = strarg
		}
		array = append(array, runtime_val)
	}

	val, err = rtfunc(rt, array...)

	return
}
