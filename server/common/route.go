// Copyright (c) TFG Co. All Rights Reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package common

import (
	"fmt"
	"strings"	
)

// 简单route str处理
type RouteStr struct {
	Server string
	Service string
	Method string

	longStr string
	shortStr string
}

func NewRouteStr(in string) *RouteStr {
	r := &RouteStr{}
	r.Init(in)
	return r
}

func (self *RouteStr) Init(in string) {
	subs := strings.Split(in, ".")
	if len(subs) < 2 {
		return
	}

	serviceIndex := 0
	if len(subs) == 3 {
		self.Server = subs[0]
		serviceIndex = 1
	}
	if len(subs) == 2 {
		self.Server = ""
	}	
	self.Service = subs[serviceIndex]
	self.Method = subs[serviceIndex+1]

	self.longStr = ""
	self.shortStr = ""
}

// 获取长route
func (self *RouteStr) GetLongStr() string {	
	if self.longStr == "" && self.Method != "" {
		self.longStr = fmt.Sprintf("%s.%s.%s",
			self.Server, self.Service, self.Method)
	}
	return self.longStr
}

func (self *RouteStr) GetShortStr() string {
	if self.shortStr == "" && self.Method != "" {
		self.shortStr = fmt.Sprintf("%s.%s",
			self.Service, self.Method)
	}
	return self.shortStr
}