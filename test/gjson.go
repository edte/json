// @program: json
// @author: edte
// @create: 2020-07-23 11:57
package test

// Result represents a json value that is returned from Get().
type Result struct {
	// Type is the json type
	Type Type
	// Raw is the raw json
	Raw string
	// Str is the json string
	Str string
	// Num is the json number
	Num float64
	// Index of raw value in original json, zero means index unknown
	Index int
}

type Type int

const (
	// Null is a null json value
	Null Type = iota
	// False is a json false boolean
	False
	// Number is json number
	Number
	// String is a json string
	String
	// True is a json true boolean
	True
	// JSON is a raw block of JSON
	JSON
)

// String returns a string representation of the type.
func (t Type) String() string {
	switch t {
	default:
		return ""
	case Null:
		return "Null"
	case False:
		return "False"
	case Number:
		return "Number"
	case String:
		return "String"
	case True:
		return "True"
	case JSON:
		return "JSON"
	}
}


/*func Get(json, path string) Result {
	if len(path) > 1 {
		if path[0] == '[' || path[0] == '{' {
			// using a subselector path
			kind := path[0]
			var ok bool
			var subs []subSelector
			subs, path, ok = parseSubSelectors(path)
			if ok {
				if len(path) == 0 || (path[0] == '|' || path[0] == '.') {
					var b []byte
					b = append(b, kind)
					var i int
					for _, sub := range subs {
						res := Get(json, sub.path)
						if res.Exists() {
							if i > 0 {
								b = append(b, ',')
							}
							if kind == '{' {
								if len(sub.name) > 0 {
									if sub.name[0] == '"' && Valid(sub.name) {
										b = append(b, sub.name...)
									} else {
										b = appendJSONString(b, sub.name)
									}
								} else {
									last := nameOfLast(sub.path)
									if isSimpleName(last) {
										b = appendJSONString(b, last)
									} else {
										b = appendJSONString(b, "_")
									}
								}
								b = append(b, ':')
							}
							var raw string
							if len(res.Raw) == 0 {
								raw = res.String()
								if len(raw) == 0 {
									raw = "null"
								}
							} else {
								raw = res.Raw
							}
							b = append(b, raw...)
							i++
						}
					}
					b = append(b, kind+2)
					var res Result
					res.Raw = string(b)
					res.Type = JSON
					if len(path) > 0 {
						res = res.Get(path[1:])
					}
					res.Index = 0
					return res
				}
			}
		}
	}

	var i int
	var c = &parseContext{json: json}
	if len(path) >= 2 && path[0] == '.' && path[1] == '.' {
		c.lines = true
		parseArray(c, 0, path[2:])
	} else {
		for ; i < len(c.json); i++ {
			if c.json[i] == '{' {
				i++
				parseObject(c, i, path)
				break
			}
			if c.json[i] == '[' {
				i++
				parseArray(c, i, path)
				break
			}
		}
	}
	if c.piped {
		res := c.value.Get(c.pipe)
		res.Index = 0
		return res
	}
	fillIndex(json, c)
	return c.value
}
*/