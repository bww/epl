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
  "fmt"
)

/**
 * A parser
 */
type parser struct {
  scanner   *scanner
  la        []token
  nodes     []executable
}

/**
 * Create a parser
 */
func newParser(s *scanner) *parser {
  return &parser{s, make([]token, 0, 2), make([]executable, 0, 3)}
}

/**
 * Obtain a look-ahead token without consuming it
 */
func (p *parser) peek(n int) token {
  var t token
  
  l := len(p.la)
  if n < l {
    return p.la[n]
  }else if n + 1 > cap(p.la) {
    panic(fmt.Errorf("Look-ahead overrun: %d >= %d", n + 1, cap(p.la)))
  }
  
  p.la = p.la[:n+1]
  for i := l; i < n + 1; i++ {
    t = p.scanner.scan()
    p.la[i] = t
  }
  
  return t
}

/**
 * Consume the next token
 */
func (p *parser) next() token {
  l := len(p.la)
  if l < 1 {
    return p.scanner.scan()
  }else{
    t := p.la[0]
    for i := 1; i < l; i++ { p.la[i-1] = p.la[i] }
    p.la = p.la[:l-1]
    return t
  }
}

/**
 * Parse
 */
func (p *parser) parse() (executable, error) {
  return p.parseExpression()
}

/**
 * Parse
 */
func (p *parser) parseExpression() (executable, error) {
  return p.parseArithmetic()
  /*
  left := p.peek(0)
  switch left.which {
    case tokenEOF:
      return p.parsePrimary()
    case tokenError:
      return nil, fmt.Errorf("Error: %v", left)
  }
  
  op := p.peek(1)
  switch op.which {
    case tokenEOF:
      return nil, fmt.Errorf("Unexpected end-of-input")
    case tokenError:
      return nil, fmt.Errorf("Error: %v", op)
    case tokenAdd, tokenSub, tokenMul, tokenDiv:
      return p.parseArithmetic()
    default:
      return nil, fmt.Errorf("Illegal token in expression: %v", op)
  }
  */
}

/**
 * Parse an arithmetic expression
 */
func (p *parser) parseArithmetic() (executable, error) {
  
  left, err := p.parsePrimary()
  if err != nil {
    return nil, err
  }
  
  op := p.peek(0)
  switch op.which {
    case tokenError:
      return nil, fmt.Errorf("Error: %v", op)
    case tokenAdd, tokenSub, tokenMul, tokenDiv:
      break // valid tokens
    default:
      return left, nil
  }
  
  p.next()
  right, err := p.parseArithmetic()
  if err != nil {
    return nil, err
  }
  
  return &arithmeticNode{node{}, op, left, right}, nil
}

/**
 * Parse a primary expression
 */
func (p *parser) parsePrimary() (executable, error) {
  t := p.next()
  switch t.which {
    case tokenEOF:
      return nil, fmt.Errorf("Unexpected end-of-input")
    case tokenError:
      return nil, fmt.Errorf("Error: %v", t)
    case tokenLParen:
      return p.parseParen()
    case tokenIdentifier:
      return &identNode{node{}, t.value.(string)}, nil
    case tokenNumber, tokenString:
      return &literalNode{node{}, t.value}, nil
    case tokenTrue:
      return &literalNode{node{}, true}, nil
    case tokenFalse:
      return &literalNode{node{}, false}, nil
    case tokenNil:
      return &literalNode{node{}, nil}, nil
    default:
      return nil, fmt.Errorf("Illegal token in primary expression: %v", t)
  }
}

/**
 * Parse a (sub-expression)
 */
func (p *parser) parseParen() (executable, error) {
  
  e, err := p.parseExpression()
  if err != nil {
    return nil, err
  }
  
  t := p.next()
  if t.which != tokenRParen {
    return nil, fmt.Errorf("Expected ')' but found %v", t)
  }
  
  return e, nil
}
