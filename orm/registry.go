package orm

import (
	"mytest/orm/internal/errs"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

// Registry 元数据注册中心的抽象
type Registry interface {
	// Get 查找元数据
	Get(val any) (*Model, error)
	// Register 注册一个模型
	Register(val any, opts ...ModelOpt) (*Model, error)
}

// registry 基于标签和接口的实现
// 目前来看，我们只有一个实现，所以暂时可以维持私有
type registry struct {
	models sync.Map
}

func NewRegistry() Registry {
	return &registry{}
}

// Get 查找元数据模型
func (r *registry) Get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*Model), nil
	}
	return r.Register(val)
}

func (r *registry) Register(val any, opts ...ModelOpt) (*Model, error) {
	m, err := r.parseModel(val)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		err = opt(m)
		if err != nil {
			return nil, err
		}
	}

	typ := reflect.TypeOf(val)
	r.models.Store(typ, m)
	return m, nil
}

// parseModel 支持从标签中提取自定义设置
// 标签形式 orm:"key1=value1,key2=value2"
func (r *registry) parseModel(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	if typ.Kind() != reflect.Ptr ||
		typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	typ = typ.Elem()

	// 获得字段的数量
	numField := typ.NumField()
	fds := make(map[string]*field, numField)
	cols := make(map[string]*field, numField)
	for i := 0; i < numField; i++ {
		fdType := typ.Field(i)
		tags, err := r.parseTag(fdType.Tag)
		if err != nil {
			return nil, err
		}
		colName := tags[tagKeyColumn]
		if colName == "" {
			colName = underscoreName(fdType.Name)
		}
		fd := &field{
			colName: colName,
			goName:  fdType.Name,
			typ:     fdType.Type,
		}
		fds[fdType.Name] = fd
		cols[colName] = fd
	}
	var tableName string
	if tn, ok := val.(TableName); ok {
		tableName = tn.TableName()
	}

	if tableName == "" {
		tableName = underscoreName(typ.Name())
	}

	return &Model{
		tableName: tableName,
		fieldMap:  fds,
		colMap:    cols,
	}, nil
}

func (r *registry) parseTag(tag reflect.StructTag) (map[string]string, error) {
	ormTag := tag.Get("orm")
	if ormTag == "" {
		// 返回一个空的 map，这样调用者就不需要判断 nil 了
		return map[string]string{}, nil
	}
	// 这个初始化容量就是我们支持的 key 的数量，
	// 现在只有一个，所以我们初始化为 1
	res := make(map[string]string, 1)

	// 接下来就是字符串处理了
	pairs := strings.Split(ormTag, ",")
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			return nil, errs.NewErrInvalidTagContent(pair)
		}
		res[kv[0]] = kv[1]
	}
	return res, nil
}

// underscoreName 驼峰转字符串命名
func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}

	}
	return string(buf)
}

var models = map[reflect.Type]*Model{}
