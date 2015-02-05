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
  "os"
  "fmt"
  "testing"
)

func TestThis(t *testing.T) {
  
  sources := []string{
    `1 + 2`,
  }
  
  for _, e := range sources {
    compileAndValidate(t, e, nil)
  }
  
}

/*
func TestBasicEscaping(t *testing.T) {
  var source string
  
  source = `\foo`
  compileAndValidate(t, source, []token{
    token{span{source, 0, 4}, tokenVerbatim, source},
    token{span{source, 4, 0}, tokenEOF, nil},
  })
  
  source = `\@`
  compileAndValidate(t, source, []token{
    token{span{source, 1, 1}, tokenVerbatim, "@"},
    token{span{source, 2, 0}, tokenEOF, nil},
  })
  
  source = `x\@`
  compileAndValidate(t, source, []token{
    token{span{source, 0, 1}, tokenVerbatim, "x"},
    token{span{source, 2, 1}, tokenVerbatim, "@"},
    token{span{source, 3, 0}, tokenEOF, nil},
  })
  
  source = `\\\@`
  compileAndValidate(t, source, []token{
    token{span{source, 0, 1}, tokenVerbatim, "\\"},
    token{span{source, 3, 1}, tokenVerbatim, "@"},
    token{span{source, 4, 0}, tokenEOF, nil},
  })
  
  source = `\@\\`
  compileAndValidate(t, source, []token{
    token{span{source, 1, 1}, tokenVerbatim, "@"},
    token{span{source, 2, 1}, tokenVerbatim, "\\"},
    token{span{source, 4, 0}, tokenEOF, nil},
  })
  
  source = `\\`
  compileAndValidate(t, source, []token{
    token{span{source, 1, 1}, tokenVerbatim, "\\"},
    token{span{source, 2, 0}, tokenEOF, nil},
  })
  
  source = `\`
  compileAndValidate(t, source, []token{
    token{span{source, 0, 1}, tokenVerbatim, "\\"},
    token{span{source, 1, 0}, tokenEOF, nil},
  })
  
  source = `foo\`
  compileAndValidate(t, source, []token{
    token{span{source, 0, 4}, tokenVerbatim, "foo\\"},
    token{span{source, 4, 0}, tokenEOF, nil},
  })
  
}
*/

func TestBasicTypes(t *testing.T) {
  var source string
  
  source = `123`
  compileAndValidate(t, source, []token{
    token{span{source, 1, 3}, tokenNumber, float64(123)},
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
  parseAndRun(t, `1+2`)
  parseAndRun(t, `1`)
  parseAndRun(t, `"Yes" + 2`)
  parseAndRun(t, `true`)
  parseAndRun(t, `false`)
  parseAndRun(t, `nil`)
}

func parseAndRun(t *testing.T, source string) {
  
  s := newScanner(source)
  p := newParser(s)
  r := &runtime{os.Stdout}
  
  x, err := p.parse()
  if err != nil {
    t.Error(err)
  }
  
  fmt.Println("-----> %v", x)
  
  y, err := x.exec(r, nil)
  if err != nil {
    t.Error(err)
  }else{
    fmt.Println("-----> %v", y)
  }
  
}
