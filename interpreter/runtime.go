package interpreter

import (
	"encoding/binary"
	"fmt"
	"math"
	"toylingo/parser"
	"toylingo/utils"
)

var type_sizes = map[string]int{
	"number": 8,
	"char":   1,
	"bool":   1,
}

type Variable struct {
	pointer Pointer
	vartype string
}

func getValue(variable Variable) interface{} {
	switch variable.vartype {
	case "number":
		return getNumber(variable)
	// case "char":
	// 	return getChar(variable)
	// case "bool":
	// 	return getBool(variable)
	}
	return nil
}

func writeBits(ptr Pointer, value int64, size int) {
	fmt.Println("writeBits", ptr, value, size)
	for i := 0; i < size; i++ {
		HEAP[ptr.address+i] = byte(value & 0xFF)
		value = value >> 8
	}
}

func getNumber(variable Variable) float64 {
	ptr := variable.pointer
	// Take 8 bytes from HEAP starting at ptr.address and convert to float64
	bytes := HEAP[ptr.address : ptr.address+8]
	parsedFloat := math.Float64frombits(binary.LittleEndian.Uint64(bytes))
	fmt.Println("getNumber", parsedFloat)
	return parsedFloat
}

type scopeContext struct {
	scopeType string
	variables map[string]Variable
	functions map[string]parser.TreeNode
}

var contextStack = utils.MakeStack()

func pushScopeContext(label string) scopeContext{
	ctx:=scopeContext{label,make(map[string]Variable),make(map[string]parser.TreeNode)}
	if contextStack.IsEmpty(){
		contextStack.Push(ctx)
		return ctx
	}
	for k,v:=range contextStack.Peek().(scopeContext).variables{
		ctx.variables[k]=v
	}
	for k,v:=range contextStack.Peek().(scopeContext).functions{
		ctx.functions[k]=v
	}
	contextStack.Push(ctx)
	return ctx
}