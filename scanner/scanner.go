// @program: json
// @author: edte
// @create: 2020-07-22 18:10
package scanner

import (
	"errors"
)

type scan interface {
	CheckValid() (error, bool)
	error(string)
	begin() bool
	checkAll(i int) bool
	checkObject(i int) bool
	checkArray(i int) bool
	checkNumber(i int) bool
	checkString(i int) bool
	checkBool(i int) bool
	checkTrue(i int) bool
	checkFalse(i int) bool
	checkNull(i int) bool
	checkColon(i int) bool
	checkComma(i int) bool
}

type scanner struct {
	data []byte
	len  int
	err  error
}

func newScanner(data []byte) *scanner {
	return &scanner{
		data: data,
		len:  len(data),
		err:  nil,
	}
}

func New(data []byte) *scanner {
	return newScanner(data)
}

func (s scanner) CheckValid() (bool, error) {
	return s.begin(), s.err
}

func (s scanner) error(msg string) {
	s.err = errors.New(msg)
}

func (s scanner) begin() bool {
	for i := 0; i < s.len; i++ {
		switch s.data[i] {
		case ' ', '\t', '\n', '\r':
			continue
		case '{', '[', '"', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 't', 'f', 'n':
			if ok := s.checkAll(i); !ok {
				return false
			}
			for ; i < s.len; i++ {
				switch s.data[i] {
				case ' ', '\t', '\n', '\r':
					continue
				default:
					return false
				}
			}
			return true
		default:
			return false
		}
	}
	return false
}

func (s scanner) checkAll(i int) bool {
	for ; i < s.len; i++ {
		switch s.data[i] {
		case ' ', '\t', '\n', '\r':
			continue
		case '{':
			return s.checkObject(i + 1)
		case '[':
			return s.checkArray(i + 1)
		case '"':
			return s.checkString(i + 1)
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return s.checkNumber(i + 1)
		case 't', 'f':
			return s.checkBool(i + 1)
		case 'n':
			return s.checkNull(i + 1)
		default:
			return false
		}
	}
	return false
}

func (s scanner) checkObject(i int) bool {
	for ; i < s.len; i++ {
		switch s.data[i] {
		case ' ', '\t', '\n', '\r':
			continue
		case '}':
			return false
		case '"':
		key:
			if ok := s.checkString(i + 1); !ok {
				return false
			}
			if ok := s.checkColon(i); !ok {
				return false
			}
			if ok := s.checkAll(i); !ok {
				return false
			}
			if ok := s.checkComma(i, '}'); !ok {
				return false
			}
			if s.data[i] == '}' {
				return true
			}
			i++
			for ; i < s.len; i++ {
				switch s.data[i] {
				default:
					return false
				case ' ', '\t', '\n', '\r':
					continue
				case '"':
					goto key
				}
			}
			return false
		default:
			return false
		}
	}
	return false
}

func (s scanner) checkColon(i int) bool {
	for ; i < s.len; i++ {
		switch s.data[i] {
		default:
			return false
		case ' ', '\t', '\n', '\r':
			continue
		case ':':
			return true
		}
	}
	return false
}

func (s scanner) checkArray(i int) bool {
	for ; i < s.len; i++ {
		switch s.data[i] {
		case ' ', '\t', '\n', '\r':
			continue
		case ']':
			return true
		default:
			for ; i < s.len; i++ {
				if s.data[i] == ']' {
					return true
				}
				if ok := s.checkAll(i); !ok {
					return false
				}
				if ok := s.checkComma(i, ']'); !ok {
					return false
				}
			}
		}
	}
	return false
}

func (s scanner) checkComma(i int, end byte) bool {
	for ; i < s.len; i++ {
		switch s.data[i] {
		case ' ', '\t', '\n', '\r':
			continue
		case ',':
			return true
		case end:
			return true
		default:
			return false
		}
	}
	return false
}

// checkNumber
// 正数
// 负数
// 整数
// 小数
// e 表达式
func (s scanner) checkNumber(i int) bool {
	if i == s.len {
		s.error("expected number")
		return true
	}

	// 如果数字有符号，则字符前进一位，忽略掉符号
	if s.data[i] == '-' || s.data[i] == '+' {
		i++
	}

	if i == s.len {
		s.error("expected number")
		return false
	}

	// 判断是否是纯数字
	for ; i < s.len; i++ {
		if !(s.data[i] >= '0' && s.data[i] <= '9') {
			s.error("expected number")
			return false
		}
	}

	// 小数
	if s.data[i] == '.' {
		i++
		if i == s.len {
			s.error("expected number")
			return false
		}
		for ; i < s.len; i++ {
			if !(s.data[i] >= '0' && s.data[i] <= '9') {
				s.error("expected number")
				return false
			}
		}
	}
	return true
}

// 普通
// todo: 转义字符
func (s scanner) checkString(i int) bool {
	for ; i < s.len; i++ {
		if s.data[i] == '"' {
			return true
		}
	}
	s.error("expected string")
	return false
}

func (s scanner) checkBool(i int) bool {
	if s.data[i] == 'a' {
		return s.checkFalse(i)
	}
	return s.checkTrue(i)
}

func (s scanner) checkTrue(i int) bool {
	if i+3 <= s.len && s.data[i] == 'r' && s.data[i+1] == 'u' && s.data[i+2] == 'e' {
		return true
	}
	s.error("expected true")
	return false
}

func (s scanner) checkFalse(i int) bool {
	if i+4 <= s.len && s.data[i] == 'a' && s.data[i+1] == 'l' && s.data[i+2] == 's' && s.data[i+3] == 'e' {
		return true
	}
	s.error("expected false")
	return false
}

func (s scanner) checkNull(i int) bool {
	if i+3 <= s.len && s.data[i] == 'u' && s.data[i+1] == 'l' && s.data[i+2] == 'l' {
		return true
	}
	s.error("expected null")
	return false
}
