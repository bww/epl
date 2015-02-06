// 
// Copyright (c) 2015 Brian William Wolter, All rights reserved.
// EPL - A little Embeddable Predicate Language
// 
// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:
// 
//   * Redistributions of source code must retain the above copyright notice, this
//     list of conditions and the following disclaimer.
// 
//   * Redistributions in binary form must reproduce the above copyright notice,
//     this list of conditions and the following disclaimer in the documentation
//     and/or other materials provided with the distribution.
//     
//   * Neither the names of Brian William Wolter, Wolter Group New York, nor the
//     names of its contributors may be used to endorse or promote products derived
//     from this software without specific prior written permission.
//     
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT,
// INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
// LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE
// OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED
// OF THE POSSIBILITY OF SUCH DAMAGE.
// 

package epl

import (
  "io"
  "fmt"
  "reflect"
)

/**
 * Executable context
 */
type runtime struct {
  stdout    io.Writer
}

/**
 * Executable
 */
type executable interface {
  exec(*runtime, interface{})([]interface{}, error)
}

/**
 * Executable
 */
type tree interface {
  add(c executable) *node
}

/**
 * An AST node
 */
type node struct {
  subnodes  []executable
}

/**
 * Execute
 */
func (n *node) exec(runtime *runtime, context interface{}) ([]interface{}, error) {
  return nil, fmt.Errorf("No implementation")
}

/**
 * A program
 */
type emptyNode struct {
  node
}

/**
 * An expression node
 */
type exprNode struct {
  node
}

/**
 * A logical OR node
 */
type logicalOrNode struct {
  node
  left, right executable
}

/**
 * Execute
 */
func (n *logicalOrNode) exec(runtime *runtime, context interface{}) ([]interface{}, error) {
  
  lvi, err := execReturnSingle(n.left, runtime, context)
  if err != nil {
    return nil, err
  }
  lv, err := asBool(lvi)
  if err != nil {
    return nil, err
  }
  
  if lv {
    return []interface{}{true}, nil
  }
  
  rvi, err := execReturnSingle(n.right, runtime, context)
  if err != nil {
    return nil, err
  }
  rv, err := asBool(rvi)
  if err != nil {
    return nil, err
  }
  
  return []interface{}{rv}, nil
}

/**
 * A logical AND node
 */
type logicalAndNode struct {
  node
  left, right executable
}

/**
 * Execute
 */
func (n *logicalAndNode) exec(runtime *runtime, context interface{}) ([]interface{}, error) {
  
  lvi, err := execReturnSingle(n.left, runtime, context)
  if err != nil {
    return nil, err
  }
  lv, err := asBool(lvi)
  if err != nil {
    return nil, err
  }
  
  if !lv {
    return []interface{}{false}, nil
  }
  
  rvi, err := execReturnSingle(n.right, runtime, context)
  if err != nil {
    return nil, err
  }
  rv, err := asBool(rvi)
  if err != nil {
    return nil, err
  }
  
  return []interface{}{rv}, nil
}

/**
 * An arithmetic expression node
 */
type arithmeticNode struct {
  node
  op          token
  left, right executable
}

/**
 * Execute
 */
func (n *arithmeticNode) exec(runtime *runtime, context interface{}) ([]interface{}, error) {
  
  lvi, err := execReturnSingle(n.left, runtime, context)
  if err != nil {
    return nil, err
  }
  lv, err := asNumber(lvi)
  if err != nil {
    return nil, err
  }
  
  rvi, err := execReturnSingle(n.right, runtime, context)
  if err != nil {
    return nil, err
  }
  rv, err := asNumber(rvi)
  if err != nil {
    return nil, err
  }
  
  switch n.op.which {
    case tokenAdd:
      return []interface{}{ lv + rv }, nil
    case tokenSub:
      return []interface{}{ lv - rv }, nil
    case tokenMul:
      return []interface{}{ lv * rv }, nil
    case tokenDiv:
      return []interface{}{ lv / rv }, nil
    default:
      return nil, fmt.Errorf("Invalid operator: %v", n.op)
  }
  
}

/**
 * An identifier expression node
 */
type identNode struct {
  node
  ident string
}

/**
 * Execute
 */
func (n *identNode) exec(runtime *runtime, context interface{}) ([]interface{}, error) {
  switch v := context.(type) {
    case map[string]interface{}:
      return []interface{}{v[n.ident]}, nil
    default:
      return derefProp(context, n.ident)
  }
}

/**
 * Execute
 */
func derefProp(context interface{}, property string) ([]interface{}, error) {
  c, _ := derefValue(reflect.ValueOf(context))
  switch c.Kind() {
    case reflect.Struct:
      return []interface{}{c.FieldByName(property)}, nil
    default:
      return nil, fmt.Errorf("Cannot dereference context: %v (%T)", context, context)
  }
}

/**
 * Dereference a value
 */
func derefValue(value reflect.Value) (reflect.Value, int) {
  v := value
  c := 0
  for ; v.Kind() == reflect.Ptr; {
    v = v.Elem()
    c++
  }
  return v, c
}

/**
 * A literal expression node
 */
type literalNode struct {
  node
  value interface{}
}

/**
 * Execute
 */
func (n *literalNode) exec(runtime *runtime, context interface{}) ([]interface{}, error) {
  return []interface{}{n.value}, nil
}

/**
 * Execute expecting a single return value
 */
func execReturnSingle(e executable, runtime *runtime, context interface{}) (interface{}, error) {
  rva, err := e.exec(runtime, context)
  if err != nil {
    return nil, err
  }
  if len(rva) != 1 {
    return nil, fmt.Errorf("Expected single return value, got %d", len(rva))
  }else{
    return rva[0], nil
  }
}

/**
 * Obtain an interface value as a bool
 */
func asBool(value interface{}) (bool, error) {
  switch v := value.(type) {
    case bool:
      return v, nil
    case uint8:
      return v != 0, nil
    case uint16:
      return v != 0, nil
    case uint32:
      return v != 0, nil
    case uint64:
      return v != 0, nil
    case int8:
      return v != 0, nil
    case int16:
      return v != 0, nil
    case int32:
      return v != 0, nil
    case int64:
      return v != 0, nil
    default:
      return false, fmt.Errorf("Cannot cast %v (%T) to bool", value, value)
  }
}

/**
 * Obtain an interface value as a number
 */
func asNumber(value interface{}) (float64, error) {
  switch v := value.(type) {
    case int:
      return float64(v), nil
    case uint8:
      return float64(v), nil
    case uint16:
      return float64(v), nil
    case uint32:
      return float64(v), nil
    case uint64:
      return float64(v), nil
    case int8:
      return float64(v), nil
    case int16:
      return float64(v), nil
    case int32:
      return float64(v), nil
    case int64:
      return float64(v), nil
    case float32:
      return float64(v), nil
    case float64:
      return v, nil
    default:
      return 0, fmt.Errorf("Cannot cast %v (%T) to numeric", value, value)
  }
}
