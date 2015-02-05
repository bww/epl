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
}

/**
 * Create a parser
 */
func newParser(s *scanner) *parser {
  return &parser{s, make([]token, 0, 2)}
}

/**
 * Obtain a look-ahead token without consuming it
 */
func (p *parser) peek(n int) token {
  var t token
  
  l := len(p.la)
  if n < l {
    return p.la[n]
  }else if n >= cap(p.la) {
    panic("Look-ahead overrun")
  }
  
  p.la = p.la[:l+n]
  for i := l; i < n; i++ {
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
func (p *parser) parse() (*program, error) {
  prog := &program{}
  
  for {
    t := p.next()
    switch t.which {
      
      case tokenEOF:
        return prog, nil
        
      case tokenError:
        return nil, fmt.Errorf("Error: %v", t)
        
      default:
        if n, err := p.parseExpression(prog, t); err != nil {
          return nil, err
        }else{
          prog.add(n)
        }
        
    }
  }
  
}

/**
 * Parse
 */
func (p *parser) parseExpression(parent tree, left token) (executable, error) {
  t := p.next()
  switch t.which {
    
    case tokenEOF:
      return nil, fmt.Errorf("Unexpected end-of-input", t)
      
    case tokenError:
      return nil, fmt.Errorf("Error: %v", t)
      
    case tokenAdd, tokenSub, tokenMul, tokenDiv:
      if n, err := p.parseArithmetic(parent, left, t); err != nil {
        return nil, err
      }else{
        return n, nil
      }
      
    default:
      return nil, fmt.Errorf("Illegal token in expression: %v", t)
      
  }
}

/**
 * Parse an arithmetic expression
 */
func (p *parser) parseArithmetic(parent tree, left, op token) (executable, error) {
  t := p.next()
  switch t.which {
    case tokenEOF:
      return nil, fmt.Errorf("Unexpected end-of-input", t)
    case tokenError:
      return nil, fmt.Errorf("Error: %v", t)
    case tokenNumber:
      return &arithmeticNode{node{}, left, op, t}, nil
    default:
      return nil, fmt.Errorf("Illegal token in arithmetic expression: %v", t)
  }
}
