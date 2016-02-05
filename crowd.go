package crowd

import (
	"errors"
	"github.com/eaciit/toolkit"
	"reflect"
)

type FnCrowd func(x interface{}) interface{}

var Self FnCrowd = func(x interface{}) interface{} {
	return x
}

func _fn(f FnCrowd) FnCrowd {
	if f == nil {
		return Self
	} else {
		return f
	}
}

const (
	CrowdMap   string = "map"
	CrowdSlice string = "slice"
)

type Crowd struct {
	data     interface{}
	dataType string

	Error error
}

func From(data interface{}) *Crowd {
	c := new(Crowd)

	isPtr := false
	isSlice := false

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		isPtr = true
		if reflect.Indirect(v).Kind() == reflect.Slice {
			isSlice = true
		}
	}

	if isPtr && isSlice {
		c.data = data
	} else {
		c.Error = errors.New("From: Data is not a pointer of slice")
	}
	return c
}

func (c *Crowd) Len() int {
	v := indirect(c.data)
	return v.Len()
}

func (c *Crowd) Item(idx int) interface{} {
	return toolkit.SliceItem(c.data, idx)
}

func (c *Crowd) Min(fn FnCrowd) interface{} {
	var min interface{}
	l := c.Len()

	min, _ = toolkit.GetEmptySliceElement(c.data)
	fn = _fn(fn)
	for i := 0; i < l; i++ {
		item := fn(c.Item(i))
		if i == 0 {
			min = item
		} else if toolkit.Compare(min, item, "$gt") {
			min = item
		}
	}
	return min
}

func (c *Crowd) Max(fn FnCrowd) interface{} {
	var max interface{}
	l := c.Len()

	max, _ = toolkit.GetEmptySliceElement(c.data)
	fn = _fn(fn)
	for i := 0; i < l; i++ {
		item := fn(c.Item(i))
		if i == 0 {
			max = item
		} else if toolkit.Compare(max, item, "$lt") {
			max = item
		}
	}
	return max
}

func (c *Crowd) Sum(fn FnCrowd) float64 {
	l := c.Len()

	ret, _ := toolkit.GetEmptySliceElement(c.data)
	//toolkit.Println("Value: ", ret, reflect.TypeOf(ret).String())
	if !toolkit.IsNumber(ret) {
		return 0
	}

	fn = _fn(fn)
	sum := float64(0)
	for i := 0; i < l; i++ {
		item := toolkit.ToFloat64(fn(c.Item(i)), 4, toolkit.RoundingAuto)
		sum += item
	}
	//e := toolkit.Serde(sum, &ret, "json")

	return sum
}

func (c *Crowd) Avg(fn FnCrowd) float64 {
	l := c.Len()
	if l == 0 {
		return 0
	}
	ret, _ := toolkit.GetEmptySliceElement(c.data)
	//toolkit.Println("Value: ", ret, reflect.TypeOf(ret).String())
	if !toolkit.IsNumber(ret) {
		return 0
	}

	fn = _fn(fn)
	sum := float64(0)
	for i := 0; i < l; i++ {
		item := toolkit.ToFloat64(fn(c.Item(i)), 4, toolkit.RoundingAuto)
		sum += item
	}
	//e := toolkit.Serde(sum, &ret, "json")

	return sum / float64(l)
}
