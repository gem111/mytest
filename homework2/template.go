package web

import (
	"bytes"
	"context"
	"html/template"
)

type TemplateEngine interface {

	// Render  渲染页面   按名字来索引 data渲染页面来的数据
	Render(ctx context.Context, tplName string, data any) ([]byte, error)

	//不需要  让具体实现自己管自己的模板  我们只管渲染  只需要一种方法  管理模板本身有点多余  具体的实现具体的人去管就可以了
	//addtemplate

	//渲染页面 数据写入到writer里面   这种也可以  优点就是可以直接输出   缺点就是不好测试
	//Render(ctx context.Context,tplName string,data any,write io.Writer) error

	//用这个context 没有问题  ctx context
}

type GoTemplateEngine struct {
	T *template.Template
}

func (g *GoTemplateEngine) Render(ctx context.Context, tplName string, data any) ([]byte, error) {

	bs := &bytes.Buffer{}
	err := g.T.ExecuteTemplate(bs, tplName, data)
	return bs.Bytes(), err
}

func ParseGlob(pattern string) error {

	return nil
}
