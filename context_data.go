package kelly

import (
	"fmt"
)

// contextData 数据存取接口
type contextData interface {
	Set(interface{}, interface{}) contextData        // 存储数据
	Get(interface{}) interface{}                     // 读取数据， 如果不存在， 返回nil
	GetDefault(interface{}, interface{}) interface{} // 读取数据， 如果不存在， 返回默认值
	MustGet(interface{}) interface{}                 // 读取数据， 如果不存在，报错
}

// contextDataMap @contextData基于map的实现
type contextDataMap map[interface{}]interface{}

type contextMapData struct {
	data contextDataMap
}

func (c *contextMapData) Set(key, value interface{}) contextData {
	c.data[key] = value
	return c
}

func (c contextMapData) Get(key interface{}) interface{} {
	if data, ok := c.data[key]; ok {
		return data
	}
	return nil
}

func (c contextMapData) GetDefault(key interface{}, dft interface{}) interface{} {
	if data, ok := c.data[key]; ok {
		return data
	}
	return dft
}

func (c contextMapData) MustGet(key interface{}) interface{} {
	if data, ok := c.data[key]; ok {
		return data
	}
	panic(fmt.Errorf("context key(%s) not exist: %w", key, ErrNoContextData))
}

func newContextMapData() contextData {
	c := &contextMapData{}
	c.data = make(contextDataMap)
	return c
}
