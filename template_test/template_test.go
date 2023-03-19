package template_test

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"html/template"
	"testing"
)

func TestHelloWorld(t *testing.T) {

	type user struct {
		Name string
	}
	tmp := template.New("hello-world")
	tpl, err := tmp.Parse(`hello,{{.Name}}`)

	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, &user{Name: "tom"})

	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, `hello,tom`, buffer.String())
}

func TestMapData(t *testing.T) {

	tmp := template.New("hello-world")
	tpl, err := tmp.Parse(`hello,{{.Name}}`)

	require.NoError(t, err)
	buffer := &bytes.Buffer{}

	test := make(map[string]string)
	test["Name"] = "tom"
	err = tpl.Execute(buffer, test)

	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, `hello,tom`, buffer.String())
}

func TestSlice(t *testing.T) {

	tmp := template.New("hello-world")
	tpl, err := tmp.Parse(`hello,{{index . 0}}`)

	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, []string{"tom"})
	require.NoError(t, err)
	assert.Equal(t, `hello,tom`, buffer.String())
}
func TestBasic(t *testing.T) {

	tmp := template.New("hello-world")
	tpl, err := tmp.Parse(`hello,{{.}}`)

	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, "tom")
	require.NoError(t, err)
	assert.Equal(t, `hello,tom`, buffer.String())
}

func TestFuncCall(t *testing.T) {
	tmp := template.New("hello-world")
	tpl, err := tmp.Parse(`
切片长度：{{len .Slice}}
{{printf "%.2f" 1.2345}}
hello,{{.Hello "tom" "Jerry"}}`)

	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, FuncCall{
		Slice: []string{"a", "b"},
	})
	require.NoError(t, err)
	assert.Equal(t, `
切片长度：2
1.23
hello, tom,Jerry`, buffer.String())

}

func TestIfElseBool(t *testing.T) {
	tmp := template.New("hello-world")
	tpl, err := tmp.Parse(`
{{- if and (gt .Age 0) (le .Age 6) }}
儿童 【0，6】
{{ else if and (gt .Age 6) (le .Age 18) }}
少年 【6，18】
{{- else -}}
成人 >18
{{- end -}}
`)

	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, map[string]any{"Age": 19})
	require.NoError(t, err)

	assert.Equal(t, `成  人 >18`, buffer.String())

}
func TestLoop(t *testing.T) {
	tmp := template.New("hello-world")
	tpl, err := tmp.Parse(`
{{- range $ids,$ele := . -}}
{{- $ids -}} -- {{- $ele -}} 
{{- end -}}
`)

	require.NoError(t, err)
	buffer := &bytes.Buffer{}
	err = tpl.Execute(buffer, make([]int, 100))
	require.NoError(t, err)

	assert.Equal(t, `0--01--02--03--04--05--06--07--08--09--010--011--012--013--014--015--016--017--018--019--020--021--022--023--024--025--026--027--028--029--030--031--032--033--034--035--036--037--038--039--040--041--042--043--044--045--046--047--048--049--050--051--052--053--054--055--056--057--058--059--060--061--062--063--064--065--066--067--068--069--070--071--072--073--074--075--076--077--078--079--080--081--082--083--084--085--086--087--088--089--090--091--092--093--094--095--096--097--098--099--0`, buffer.String())

}

type FuncCall struct {
	Slice []string
}

func (f FuncCall) Hello(first string, last string) string {
	return fmt.Sprintf(" %s,%s", first, last)
}
