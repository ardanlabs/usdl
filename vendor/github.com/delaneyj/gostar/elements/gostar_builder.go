// DO NOT EDIT THIS FILE. IT IS AUTOMATICALLY GENERATED.
package elements

import (
	"fmt"
	"html"
	"io"
	"strings"

	"github.com/igrmk/treemap/v2"
	"github.com/valyala/bytebufferpool"
	"golang.org/x/exp/constraints"
)

var (
	openBracket       = []byte("<")
	closeBracket      = []byte(">")
	spaceCloseBracket = []byte(" >")
	openSlash         = []byte("</")
	// dataDash          = []byte(" data-")
	equalDblQuote = []byte("=\"")
	dblQuote      = []byte("\"")
	space         = []byte(" ")
)

type ElementRenderer interface {
	Render(w io.Writer) error
}

type ElementRendererFunc func() ElementRenderer

type Element struct {
	Tag                  []byte
	IsSelfClosing        bool
	IntAttributes        *treemap.TreeMap[string, int]
	FloatAttributes      *treemap.TreeMap[string, float64]
	StringAttributes     *treemap.TreeMap[string, string]
	DelimitedStrings     *treemap.TreeMap[string, *DelimitedBuilder[string]]
	KVStrings            *treemap.TreeMap[string, *KVBuilder]
	BoolAttributes       *treemap.TreeMap[string, bool]
	CustomDataAttributes *treemap.TreeMap[string, string]
	Descendants          []ElementRenderer
}

func (e *Element) Attr(name string, value string) *Element {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set(name, value)
	return e
}

func (e *Element) Attrs(attrs ...string) *Element {
	if len(attrs)%2 != 0 {
		panic("attrs must be a multiple of 2")
	}
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	for i := 0; i < len(attrs); i += 2 {
		k := attrs[i]
		v := attrs[i+1]
		e.StringAttributes.Set(k, v)
	}
	return e
}

func (e *Element) AttrsMap(attrs map[string]string) *Element {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	for k, v := range attrs {
		e.StringAttributes.Set(k, v)
	}
	return e
}

func (e *Element) Render(w io.Writer) error {
	w.Write(openBracket)
	w.Write(e.Tag)

	finalKeys := treemap.New[string, string]()

	if e.IntAttributes != nil && e.IntAttributes.Len() > 0 {
		for it := e.IntAttributes.Iterator(); it.Valid(); it.Next() {
			k, v := it.Key(), it.Value()
			finalKeys.Set(k, fmt.Sprint(v))
		}
	}

	if e.FloatAttributes != nil && e.FloatAttributes.Len() > 0 {
		for it := e.FloatAttributes.Iterator(); it.Valid(); it.Next() {
			k, v := it.Key(), it.Value()
			finalKeys.Set(k, fmt.Sprint(v))
		}
	}

	if e.StringAttributes != nil && e.StringAttributes.Len() > 0 {
		for it := e.StringAttributes.Iterator(); it.Valid(); it.Next() {
			k, v := it.Key(), it.Value()
			finalKeys.Set(k, v)
		}
	}

	if e.DelimitedStrings != nil && e.DelimitedStrings.Len() > 0 {
		for it := e.DelimitedStrings.Iterator(); it.Valid(); it.Next() {
			k, v := it.Key(), it.Value()
			buf := bytebufferpool.Get()
			if err := v.Render(buf); err != nil {
				return err
			}
			finalKeys.Set(k, buf.String())
			bytebufferpool.Put(buf)
		}
	}

	if e.KVStrings != nil && e.KVStrings.Len() > 0 {
		for it := e.KVStrings.Iterator(); it.Valid(); it.Next() {
			k, v := it.Key(), it.Value()
			buf := bytebufferpool.Get()
			if err := v.Render(buf); err != nil {
				return err
			}
			finalKeys.Set(k, buf.String())
			bytebufferpool.Put(buf)
		}
	}

	if e.CustomDataAttributes != nil && e.CustomDataAttributes.Len() > 0 {
		for it := e.CustomDataAttributes.Iterator(); it.Valid(); it.Next() {
			k, v := it.Key(), it.Value()
			finalKeys.Set("data-"+k, v)
		}
	}

	if e.BoolAttributes != nil && e.BoolAttributes.Len() > 0 {
		for it := e.BoolAttributes.Iterator(); it.Valid(); it.Next() {
			k, v := it.Key(), it.Value()
			if v {
				finalKeys.Set(k, "")
			}
		}
	}

	if finalKeys.Len() > 0 {
		for it := finalKeys.Iterator(); it.Valid(); it.Next() {
			k, v := it.Key(), it.Value()

			w.Write(space)
			w.Write([]byte(k))

			if v != "" {
				w.Write(equalDblQuote)
				w.Write([]byte(fmt.Sprint(v)))
				w.Write(dblQuote)
			}
		}
	}

	if e.IsSelfClosing {
		w.Write(spaceCloseBracket)
		return nil
	}
	w.Write(closeBracket)

	for _, d := range e.Descendants {
		if d == nil {
			continue
		}
		if err := d.Render(w); err != nil {
			return err
		}
	}

	w.Write(openSlash)
	w.Write(e.Tag)
	w.Write(closeBracket)

	return nil
}

type customDataKeyModifier func() string

func customDataKey(key string, modifiers ...customDataKeyModifier) string {
	sb := strings.Builder{}
	sb.WriteString(key)
	for _, m := range modifiers {
		sb.WriteRune('.')
		sb.WriteString(m())
	}
	return sb.String()
}

type DelimitedBuilder[T constraints.Ordered] struct {
	Delimiter string
	Values    *treemap.TreeMap[T, struct{}]
}

func NewDelimitedBuilder[T constraints.Ordered](delimiter string) *DelimitedBuilder[T] {
	return &DelimitedBuilder[T]{
		Delimiter: delimiter,
		Values:    treemap.New[T, struct{}](),
	}
}

func (d *DelimitedBuilder[T]) Add(values ...T) *DelimitedBuilder[T] {
	for _, v := range values {
		d.Values.Set(v, struct{}{})
	}
	return d
}

func (d *DelimitedBuilder[T]) Remove(values ...T) *DelimitedBuilder[T] {
	for _, v := range values {
		d.Values.Del(v)
	}
	return d
}

func (d *DelimitedBuilder[T]) Render(w io.Writer) error {
	count := 0
	total := d.Values.Len()
	for it := d.Values.Iterator(); it.Valid(); it.Next() {
		b := []byte(fmt.Sprint(it.Key()))
		if _, err := w.Write(b); err != nil {
			return err
		}

		count++
		if count < total {
			w.Write([]byte(d.Delimiter))
		}
	}
	return nil
}

type KVBuilder struct {
	KeyPairDelimiter string
	EntryDelimiter   string
	Values           **treemap.TreeMap[string, string]
}

func NewKVBuilder(keyPairDelimiter, entryDelimiter string) *KVBuilder {
	return &KVBuilder{
		KeyPairDelimiter: keyPairDelimiter,
		EntryDelimiter:   entryDelimiter,
		Values:           new(*treemap.TreeMap[string, string]),
	}
}

func (d *KVBuilder) Add(key, value string) *KVBuilder {
	if *d.Values == nil {
		*d.Values = treemap.New[string, string]()
	}
	(*d.Values).Set(key, value)
	return d
}

func (d *KVBuilder) Remove(key string) *KVBuilder {
	if *d.Values == nil {
		return d
	}
	(*d.Values).Del(key)
	return d
}

func (d *KVBuilder) Render(w io.Writer) error {
	count := 0
	total := (*d.Values).Len()
	for it := (*d.Values).Iterator(); it.Valid(); it.Next() {
		k, v := it.Key(), it.Value()
		w.Write([]byte(k))
		w.Write([]byte(d.KeyPairDelimiter))
		w.Write([]byte(v))
		count++
		if count < total {
			w.Write([]byte(d.EntryDelimiter))
		}
	}
	return nil
}

type TextContent string

func (tc *TextContent) Render(w io.Writer) error {
	_, err := w.Write([]byte(*tc))
	return err
}

func Text(text string) *TextContent {
	return (*TextContent)(&text)
}

func TextF(format string, args ...interface{}) *TextContent {
	return Text(fmt.Sprintf(format, args...))
}

type EscapedContent string

func (ec *EscapedContent) Render(w io.Writer) error {
	_, err := w.Write([]byte(html.EscapeString(string(*ec))))
	return err
}

func Escaped(text string) *EscapedContent {
	return (*EscapedContent)(&text)
}

func EscapedF(format string, args ...interface{}) *EscapedContent {
	return Escaped(fmt.Sprintf(format, args...))
}

type Grouper struct {
	Children []ElementRenderer
}

func (g *Grouper) Render(w io.Writer) error {
	for _, child := range g.Children {
		if err := child.Render(w); err != nil {
			return fmt.Errorf("failed to build element: %w", err)
		}
	}
	return nil
}

func Group(children ...ElementRenderer) *Grouper {
	return &Grouper{
		Children: children,
	}
}

func If(condition bool, children ...ElementRenderer) ElementRenderer {
	if condition {
		return Group(children...)
	}
	return nil
}

func Tern(condition bool, trueChildren, falseChildren ElementRenderer) ElementRenderer {
	if condition {
		return trueChildren
	}
	return falseChildren
}

func Range[T any](values []T, cb func(T) ElementRenderer) ElementRenderer {
	children := make([]ElementRenderer, 0, len(values))
	for _, value := range values {
		children = append(children, cb(value))
	}
	return Group(children...)
}

func RangeI[T any](values []T, cb func(int, T) ElementRenderer) ElementRenderer {
	children := make([]ElementRenderer, 0, len(values))
	for i, value := range values {
		children = append(children, cb(i, value))
	}
	return Group(children...)
}

func DynGroup(childrenFuncs ...ElementRendererFunc) *Grouper {
	children := make([]ElementRenderer, 0, len(childrenFuncs))
	for _, childFunc := range childrenFuncs {
		child := childFunc()
		if child != nil {
			children = append(children, child)
		}
	}
	return &Grouper{
		Children: children,
	}
}

func DynIf(condition bool, childrenFuncs ...ElementRendererFunc) ElementRenderer {
	if condition {
		children := make([]ElementRenderer, 0, len(childrenFuncs))
		for _, childFunc := range childrenFuncs {
			child := childFunc()
			if child != nil {
				children = append(children, child)
			}
		}
		return Group(children...)
	}
	return nil
}

func DynTern(condition bool, trueChildren, falseChildren ElementRendererFunc) ElementRenderer {
	if condition {
		return trueChildren()
	}
	return falseChildren()
}

func NewElement(tag string, children ...ElementRenderer) *Element {
	return &Element{
		Tag:         []byte(tag),
		Descendants: children,
	}
}

func Error(err error) ElementRenderer {
	return Text(err.Error())
}
