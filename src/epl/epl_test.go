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
  "testing"
)

// func TestThis(t *testing.T) {
//   sources := []string{
//     `1 + 2`,
//   }
//   for _, e := range sources {
//     compileAndValidate(t, e, nil)
//   }
// }

func TestBasicTypes(t *testing.T) {
  var source string
  
  source = `123`
  compileAndValidate(t, source, []token{
    token{span{source, 0, 3}, tokenNumber, float64(123)},
    token{span{source, 6, 0}, tokenEOF, nil},
  })
  
}

func compileAndValidate(test *testing.T, source string, expect []token) {
  fmt.Println(source)
  s := newScanner(source)
  
  for {
    
    t := s.scan()
    fmt.Println("T", t)
    
    if expect != nil {
      
      if len(expect) < 1 {
        test.Errorf("Unexpected end of tokens")
        fmt.Println("!!!")
        return
      }
      
      e := expect[0]
      
      if e.which != t.which {
        test.Errorf("Unexpected token type (%v != %v) in %v", t.which, e.which, source)
        fmt.Println("!!!")
        return
      }
      
      if e.span.excerpt() != t.span.excerpt() {
        test.Errorf("Excerpts do not match (%q != %q) in %v", t.span.excerpt(), e.span.excerpt(), source)
        fmt.Println("!!!")
        return
      }
      
      if e.value != t.value {
        test.Errorf("Values do not match (%v != %v) in %v", t.value, e.value, source)
        fmt.Println("!!!")
        return
      }
      
      expect = expect[1:]
      
    }
    
    if t.which == tokenEOF {
      break
    }else if t.which == tokenError {
      break
    }
    
  }
  
  if expect != nil {
    if len(expect) > 0 {
      test.Errorf("Unexpected end of input (%d tokens remain)", len(expect))
      fmt.Println("!!!")
      return
    }
  }
  
  fmt.Println("---")
}

func TestParse(t *testing.T) {
  
  // basic
  parseAndRun(t, `1`, nil, float64(1))
  parseAndRun(t, `123.456`, nil, float64(123.456))
  parseAndRun(t, `-1`, nil, float64(-1))
  parseAndRun(t, `-123.456`, nil, float64(-123.456))
  parseAndRun(t, `true`, nil, true)
  parseAndRun(t, `false`, nil, false)
  parseAndRun(t, `"abcdef"`, nil, "abcdef")
  
  // string escapes
  parseAndRun(t, `"\t\u2022"`, nil, "\t\u2022")
  parseAndRun(t, `"Joe said \"this is the story...\" and that was that."`, nil, "Joe said \"this is the story...\" and that was that.")
  
  // logic
  parseAndRun(t, `true || true`, nil, true)
  parseAndRun(t, `true || false`, nil, true)
  parseAndRun(t, `false || false`, nil, false)
  parseAndRun(t, `true && true`, nil, true)
  parseAndRun(t, `true && false`, nil, false)
  parseAndRun(t, `false && false`, nil, false)
  parseAndRun(t, `false || true && true`, nil, true)
  parseAndRun(t, `false || true && false`, nil, false)
  parseAndRun(t, `false || false && false`, nil, false)
  
  // arithmetic
  parseAndRun(t, `1 + 2`, nil, float64(3))
  parseAndRun(t, `1.5 + 2.5`, nil, float64(4))
  parseAndRun(t, `10 - 2`, nil, float64(8))
  parseAndRun(t, `10 - 20`, nil, float64(-10))
  parseAndRun(t, `10 / 2`, nil, float64(5))
  parseAndRun(t, `10 * 2`, nil, float64(20))
  parseAndRun(t, `10 * 2 - 1`, nil, float64(19))
  parseAndRun(t, `10 / 2 - 1`, nil, float64(4))
  
  // cases with signs and arithmetic
  parseAndRun(t, `1 + +2`, nil, float64(3))
  parseAndRun(t, `1 + -2`, nil, float64(-1))
  parseAndRun(t, `-1 + -2`, nil, float64(-3))
  parseAndRun(t, `-1 - 2`, nil, float64(-3))
  
  // nesting expressions
  parseAndRun(t, `1 > 2 || 3 > 2`, nil, true)
  parseAndRun(t, `1 > 2 && 3 > 2`, nil, false)
  parseAndRun(t, `1 < 2 || 3 < 2`, nil, true)
  parseAndRun(t, `1 < 2 && 3 > 2`, nil, true)
  
  // parseAndRun(t, `1 - 2 + (2 * 3)`, nil, float64(-7))
  // parseAndRun(t, `true`, nil, true)
  // parseAndRun(t, `false`, nil, false)
  // parseAndRun(t, `nil`, nil, nil)
  // parseAndRun(t, `num`, nil, nil)
  // parseAndRun(t, `num+3`, nil, nil)
  // parseAndRun(t, `num == 3`, nil, nil)
  // parseAndRun(t, `num > 3`, nil, nil)
  // parseAndRun(t, `num < 4 || 1 + 2 < 5`, nil, nil)
  // parseAndRun(t, `"foo" > 3`, nil, nil)
  // parseAndRun(t, `"foo" > "f"`, nil, nil)
  // parseAndRun(t, `"foo" == "foo"`, nil, nil)
  // parseAndRun(t, `foo.bar`, nil, nil)
  // parseAndRun(t, `foo.bar.zar`, nil, nil)
  // parseAndRun(t, `foo.bar.car`, nil, nil)
  // parseAndRun(t, `foo.bar.car.finally`, nil, nil)
}

func parseAndRun(t *testing.T, source string, context interface{}, result interface{}) {
  
  s := newScanner(source)
  p := newParser(s)
  
  if context == nil {
    context = map[string]interface{}{
      "num": 123,
      "foo": map[string]interface{}{
        "bar": map[string]interface{}{
          "zar": "THIS IS IT, BOYS",
          "car": map[string]interface{}{
            "finally": true,
          },
        },
      },
    }
  }
  
  fmt.Printf("--> [%v]\n", source)
  
  x, err := p.parse()
  if err != nil {
    t.Error(fmt.Errorf("[%s] %v", source, err))
    return
  }
  
  y, err := x.Exec(context)
  if err != nil {
    t.Error(fmt.Errorf("[%s] %v", source, err))
    return
  }
  
  fmt.Printf("<-- %v\n", y)
  
  if y != result {
    t.Error(fmt.Errorf("[%s] Expected %v, got %v", source, result, y))
    return
  }
  
}
