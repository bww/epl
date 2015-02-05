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
 * An arithmetic expression node
 */
type arithmeticNode struct {
  node
  op          token
  left, right float64
}

/**
 * Execute
 */
func (n *arithmeticNode) exec(runtime *runtime, context interface{}) ([]interface{}, error) {
  switch n.op.which {
    case tokenAdd:
      return []interface{}{ n.left + n.right }, nil
    case tokenSub:
      return []interface{}{ n.left - n.right }, nil
    case tokenMul:
      return []interface{}{ n.left * n.right }, nil
    case tokenDiv:
      return []interface{}{ n.left / n.right }, nil
    default:
      return nil, fmt.Errorf("Invalid operator: %v", n.op)
  }
}
