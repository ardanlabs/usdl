// DO NOT EDIT THIS FILE. IT IS AUTOMATICALLY GENERATED.
// package svg feMorphology is generated from configuration file.
// Description:
package elements

import (
	"fmt"
	"html"
	"time"

	"github.com/goccy/go-json"
	"github.com/igrmk/treemap/v2"
	"github.com/samber/lo"
)

// The <feMorphology> SVG filter primitive is used to erode or dilate the input
// image
// It's usefulness lies especially in fattening or thinning effects.
type SVGFEMORPHOLOGYElement struct {
	*Element
}

// Create a new SVGFEMORPHOLOGYElement element.
// This will create a new element with the tag
// "feMorphology" during rendering.
func SVG_FEMORPHOLOGY(children ...ElementRenderer) *SVGFEMORPHOLOGYElement {
	e := NewElement("feMorphology", children...)
	e.IsSelfClosing = false
	e.Descendants = children

	return &SVGFEMORPHOLOGYElement{Element: e}
}

func (e *SVGFEMORPHOLOGYElement) Children(children ...ElementRenderer) *SVGFEMORPHOLOGYElement {
	e.Descendants = append(e.Descendants, children...)
	return e
}

func (e *SVGFEMORPHOLOGYElement) IfChildren(condition bool, children ...ElementRenderer) *SVGFEMORPHOLOGYElement {
	if condition {
		e.Descendants = append(e.Descendants, children...)
	}
	return e
}

func (e *SVGFEMORPHOLOGYElement) TernChildren(condition bool, trueChildren, falseChildren ElementRenderer) *SVGFEMORPHOLOGYElement {
	if condition {
		e.Descendants = append(e.Descendants, trueChildren)
	} else {
		e.Descendants = append(e.Descendants, falseChildren)
	}
	return e
}

func (e *SVGFEMORPHOLOGYElement) Attr(name string, value string) *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set(name, value)
	return e
}

func (e *SVGFEMORPHOLOGYElement) Attrs(attrs ...string) *SVGFEMORPHOLOGYElement {
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

func (e *SVGFEMORPHOLOGYElement) AttrsMap(attrs map[string]string) *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	for k, v := range attrs {
		e.StringAttributes.Set(k, v)
	}
	return e
}

func (e *SVGFEMORPHOLOGYElement) Text(text string) *SVGFEMORPHOLOGYElement {
	e.Descendants = append(e.Descendants, Text(text))
	return e
}

func (e *SVGFEMORPHOLOGYElement) TextF(format string, args ...any) *SVGFEMORPHOLOGYElement {
	return e.Text(fmt.Sprintf(format, args...))
}

func (e *SVGFEMORPHOLOGYElement) IfText(condition bool, text string) *SVGFEMORPHOLOGYElement {
	if condition {
		e.Descendants = append(e.Descendants, Text(text))
	}
	return e
}

func (e *SVGFEMORPHOLOGYElement) IfTextF(condition bool, format string, args ...any) *SVGFEMORPHOLOGYElement {
	if condition {
		e.Descendants = append(e.Descendants, Text(fmt.Sprintf(format, args...)))
	}
	return e
}

func (e *SVGFEMORPHOLOGYElement) Escaped(text string) *SVGFEMORPHOLOGYElement {
	e.Descendants = append(e.Descendants, Escaped(text))
	return e
}

func (e *SVGFEMORPHOLOGYElement) IfEscaped(condition bool, text string) *SVGFEMORPHOLOGYElement {
	if condition {
		e.Descendants = append(e.Descendants, Escaped(text))
	}
	return e
}

func (e *SVGFEMORPHOLOGYElement) EscapedF(format string, args ...any) *SVGFEMORPHOLOGYElement {
	return e.Escaped(fmt.Sprintf(format, args...))
}

func (e *SVGFEMORPHOLOGYElement) IfEscapedF(condition bool, format string, args ...any) *SVGFEMORPHOLOGYElement {
	if condition {
		e.Descendants = append(e.Descendants, EscapedF(format, args...))
	}
	return e
}

func (e *SVGFEMORPHOLOGYElement) CustomData(key, value string) *SVGFEMORPHOLOGYElement {
	if e.CustomDataAttributes == nil {
		e.CustomDataAttributes = treemap.New[string, string]()
	}
	e.CustomDataAttributes.Set(key, value)
	return e
}

func (e *SVGFEMORPHOLOGYElement) IfCustomData(condition bool, key, value string) *SVGFEMORPHOLOGYElement {
	if condition {
		e.CustomData(key, value)
	}
	return e
}

func (e *SVGFEMORPHOLOGYElement) CustomDataF(key, format string, args ...any) *SVGFEMORPHOLOGYElement {
	return e.CustomData(key, fmt.Sprintf(format, args...))
}

func (e *SVGFEMORPHOLOGYElement) IfCustomDataF(condition bool, key, format string, args ...any) *SVGFEMORPHOLOGYElement {
	if condition {
		e.CustomData(key, fmt.Sprintf(format, args...))
	}
	return e
}

func (e *SVGFEMORPHOLOGYElement) CustomDataRemove(key string) *SVGFEMORPHOLOGYElement {
	if e.CustomDataAttributes == nil {
		return e
	}
	e.CustomDataAttributes.Del(key)
	return e
}

// The input for this filter.
func (e *SVGFEMORPHOLOGYElement) IN(s string) *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set("in", s)
	return e
}

func (e *SVGFEMORPHOLOGYElement) INF(format string, args ...any) *SVGFEMORPHOLOGYElement {
	return e.IN(fmt.Sprintf(format, args...))
}

func (e *SVGFEMORPHOLOGYElement) IfIN(condition bool, s string) *SVGFEMORPHOLOGYElement {
	if condition {
		e.IN(s)
	}
	return e
}

func (e *SVGFEMORPHOLOGYElement) IfINF(condition bool, format string, args ...any) *SVGFEMORPHOLOGYElement {
	if condition {
		e.IN(fmt.Sprintf(format, args...))
	}
	return e
}

// Remove the attribute IN from the element.
func (e *SVGFEMORPHOLOGYElement) INRemove(s string) *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("in")
	return e
}

func (e *SVGFEMORPHOLOGYElement) INRemoveF(format string, args ...any) *SVGFEMORPHOLOGYElement {
	return e.INRemove(fmt.Sprintf(format, args...))
}

// The operator attribute defines what type of operation is performed.
func (e *SVGFEMORPHOLOGYElement) OPERATOR(c SVGFeMorphologyOperatorChoice) *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set("operator", string(c))
	return e
}

type SVGFeMorphologyOperatorChoice string

const (
	// The operator attribute defines what type of operation is performed.
	SVGFeMorphologyOperator_erode SVGFeMorphologyOperatorChoice = "erode"
	// The operator attribute defines what type of operation is performed.
	SVGFeMorphologyOperator_dilate SVGFeMorphologyOperatorChoice = "dilate"
)

// Remove the attribute OPERATOR from the element.
func (e *SVGFEMORPHOLOGYElement) OPERATORRemove(c SVGFeMorphologyOperatorChoice) *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("operator")
	return e
}

// The radius attribute indicates the size of the matrix.
func (e *SVGFEMORPHOLOGYElement) RADIUS(f float64) *SVGFEMORPHOLOGYElement {
	if e.FloatAttributes == nil {
		e.FloatAttributes = treemap.New[string, float64]()
	}
	e.FloatAttributes.Set("radius", f)
	return e
}

func (e *SVGFEMORPHOLOGYElement) IfRADIUS(condition bool, f float64) *SVGFEMORPHOLOGYElement {
	if condition {
		e.RADIUS(f)
	}
	return e
}

// Specifies a unique id for an element
func (e *SVGFEMORPHOLOGYElement) ID(s string) *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set("id", s)
	return e
}

func (e *SVGFEMORPHOLOGYElement) IDF(format string, args ...any) *SVGFEMORPHOLOGYElement {
	return e.ID(fmt.Sprintf(format, args...))
}

func (e *SVGFEMORPHOLOGYElement) IfID(condition bool, s string) *SVGFEMORPHOLOGYElement {
	if condition {
		e.ID(s)
	}
	return e
}

func (e *SVGFEMORPHOLOGYElement) IfIDF(condition bool, format string, args ...any) *SVGFEMORPHOLOGYElement {
	if condition {
		e.ID(fmt.Sprintf(format, args...))
	}
	return e
}

// Remove the attribute ID from the element.
func (e *SVGFEMORPHOLOGYElement) IDRemove(s string) *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("id")
	return e
}

func (e *SVGFEMORPHOLOGYElement) IDRemoveF(format string, args ...any) *SVGFEMORPHOLOGYElement {
	return e.IDRemove(fmt.Sprintf(format, args...))
}

// Specifies one or more classnames for an element (refers to a class in a style
// sheet)
func (e *SVGFEMORPHOLOGYElement) CLASS(s ...string) *SVGFEMORPHOLOGYElement {
	if e.DelimitedStrings == nil {
		e.DelimitedStrings = treemap.New[string, *DelimitedBuilder[string]]()
	}
	ds, ok := e.DelimitedStrings.Get("class")
	if !ok {
		ds = NewDelimitedBuilder[string](" ")
		e.DelimitedStrings.Set("class", ds)
	}
	ds.Add(s...)
	return e
}

func (e *SVGFEMORPHOLOGYElement) IfCLASS(condition bool, s ...string) *SVGFEMORPHOLOGYElement {
	if condition {
		e.CLASS(s...)
	}
	return e
}

// Remove the attribute CLASS from the element.
func (e *SVGFEMORPHOLOGYElement) CLASSRemove(s ...string) *SVGFEMORPHOLOGYElement {
	if e.DelimitedStrings == nil {
		return e
	}
	ds, ok := e.DelimitedStrings.Get("class")
	if !ok {
		return e
	}
	ds.Remove(s...)
	return e
}

// Specifies an inline CSS style for an element
func (e *SVGFEMORPHOLOGYElement) STYLEF(k string, format string, args ...any) *SVGFEMORPHOLOGYElement {
	return e.STYLE(k, fmt.Sprintf(format, args...))
}

func (e *SVGFEMORPHOLOGYElement) IfSTYLE(condition bool, k string, v string) *SVGFEMORPHOLOGYElement {
	if condition {
		e.STYLE(k, v)
	}
	return e
}

func (e *SVGFEMORPHOLOGYElement) STYLE(k string, v string) *SVGFEMORPHOLOGYElement {
	if e.KVStrings == nil {
		e.KVStrings = treemap.New[string, *KVBuilder]()
	}
	kv, ok := e.KVStrings.Get("style")
	if !ok {
		kv = NewKVBuilder(":", ";")
		e.KVStrings.Set("style", kv)
	}
	kv.Add(k, v)
	return e
}

func (e *SVGFEMORPHOLOGYElement) IfSTYLEF(condition bool, k string, format string, args ...any) *SVGFEMORPHOLOGYElement {
	if condition {
		e.STYLE(k, fmt.Sprintf(format, args...))
	}
	return e
}

// Add the attributes in the map to the element.
func (e *SVGFEMORPHOLOGYElement) STYLEMap(m map[string]string) *SVGFEMORPHOLOGYElement {
	if e.KVStrings == nil {
		e.KVStrings = treemap.New[string, *KVBuilder]()
	}
	kv, ok := e.KVStrings.Get("style")
	if !ok {
		kv = NewKVBuilder(":", ";")
		e.KVStrings.Set("style", kv)
	}
	for k, v := range m {
		kv.Add(k, v)
	}
	return e
}

// Add pairs of attributes to the element.
func (e *SVGFEMORPHOLOGYElement) STYLEPairs(pairs ...string) *SVGFEMORPHOLOGYElement {
	if len(pairs)%2 != 0 {
		panic("Must have an even number of pairs")
	}
	if e.KVStrings == nil {
		e.KVStrings = treemap.New[string, *KVBuilder]()
	}
	kv, ok := e.KVStrings.Get("style")
	if !ok {
		kv = NewKVBuilder(":", ";")
		e.KVStrings.Set("style", kv)
	}

	for i := 0; i < len(pairs); i += 2 {
		kv.Add(pairs[i], pairs[i+1])
	}

	return e
}

func (e *SVGFEMORPHOLOGYElement) IfSTYLEPairs(condition bool, pairs ...string) *SVGFEMORPHOLOGYElement {
	if condition {
		e.STYLEPairs(pairs...)
	}
	return e
}

// Remove the attribute STYLE from the element.
func (e *SVGFEMORPHOLOGYElement) STYLERemove(keys ...string) *SVGFEMORPHOLOGYElement {
	if e.KVStrings == nil {
		return e
	}
	kv, ok := e.KVStrings.Get("style")
	if !ok {
		return e
	}
	for _, k := range keys {
		kv.Remove(k)
	}
	return e
}

// Merges the singleton store with the given object

func (e *SVGFEMORPHOLOGYElement) DATASTAR_STORE(v any) *SVGFEMORPHOLOGYElement {
	if e.CustomDataAttributes == nil {
		e.CustomDataAttributes = treemap.New[string, string]()
	}
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	e.CustomDataAttributes.Set("store", html.EscapeString(string(b)))
	return e
}

// Sets the reference of the element

func (e *SVGFEMORPHOLOGYElement) DATASTAR_REF(expression string) *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key := "data-ref"

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *SVGFEMORPHOLOGYElement) IfDATASTAR_REF(condition bool, expression string) *SVGFEMORPHOLOGYElement {
	if condition {
		e.DATASTAR_REF(expression)
	}
	return e
}

// Remove the attribute DATASTAR_REF from the element.
func (e *SVGFEMORPHOLOGYElement) DATASTAR_REFRemove() *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-ref")
	return e
}

// Sets the value of the element

func (e *SVGFEMORPHOLOGYElement) DATASTAR_BIND(key string, expression string) *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key = fmt.Sprintf("data-bind-%s", key)

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *SVGFEMORPHOLOGYElement) IfDATASTAR_BIND(condition bool, key string, expression string) *SVGFEMORPHOLOGYElement {
	if condition {
		e.DATASTAR_BIND(key, expression)
	}
	return e
}

// Remove the attribute DATASTAR_BIND from the element.
func (e *SVGFEMORPHOLOGYElement) DATASTAR_BINDRemove() *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-bind")
	return e
}

// Sets the value of the element

func (e *SVGFEMORPHOLOGYElement) DATASTAR_MODEL(expression string) *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key := "data-model"

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *SVGFEMORPHOLOGYElement) IfDATASTAR_MODEL(condition bool, expression string) *SVGFEMORPHOLOGYElement {
	if condition {
		e.DATASTAR_MODEL(expression)
	}
	return e
}

// Remove the attribute DATASTAR_MODEL from the element.
func (e *SVGFEMORPHOLOGYElement) DATASTAR_MODELRemove() *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-model")
	return e
}

// Sets the textContent of the element

func (e *SVGFEMORPHOLOGYElement) DATASTAR_TEXT(expression string) *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key := "data-text"

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *SVGFEMORPHOLOGYElement) IfDATASTAR_TEXT(condition bool, expression string) *SVGFEMORPHOLOGYElement {
	if condition {
		e.DATASTAR_TEXT(expression)
	}
	return e
}

// Remove the attribute DATASTAR_TEXT from the element.
func (e *SVGFEMORPHOLOGYElement) DATASTAR_TEXTRemove() *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-text")
	return e
}

// Sets the event handler of the element

type SVGFeMorphologyOnMod customDataKeyModifier

// Debounces the event handler
func SVGFeMorphologyOnModDebounce(
	d time.Duration,
) SVGFeMorphologyOnMod {
	return func() string {
		return fmt.Sprintf("debounce_%dms", d.Milliseconds())
	}
}

// Throttles the event handler
func SVGFeMorphologyOnModThrottle(
	d time.Duration,
) SVGFeMorphologyOnMod {
	return func() string {
		return fmt.Sprintf("throttle_%dms", d.Milliseconds())
	}
}

func (e *SVGFEMORPHOLOGYElement) DATASTAR_ON(key string, expression string, modifiers ...SVGFeMorphologyOnMod) *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key = fmt.Sprintf("data-on-%s", key)

	customMods := lo.Map(modifiers, func(m SVGFeMorphologyOnMod, i int) customDataKeyModifier {
		return customDataKeyModifier(m)
	})
	key = customDataKey(key, customMods...)
	e.StringAttributes.Set(key, expression)
	return e
}

func (e *SVGFEMORPHOLOGYElement) IfDATASTAR_ON(condition bool, key string, expression string, modifiers ...SVGFeMorphologyOnMod) *SVGFEMORPHOLOGYElement {
	if condition {
		e.DATASTAR_ON(key, expression, modifiers...)
	}
	return e
}

// Remove the attribute DATASTAR_ON from the element.
func (e *SVGFEMORPHOLOGYElement) DATASTAR_ONRemove() *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-on")
	return e
}

// Sets the focus of the element

func (e *SVGFEMORPHOLOGYElement) DATASTAR_FOCUSSet(b bool) *SVGFEMORPHOLOGYElement {
	key := "data-focus"
	e.BoolAttributes.Set(key, b)
	return e
}

func (e *SVGFEMORPHOLOGYElement) DATASTAR_FOCUS() *SVGFEMORPHOLOGYElement {
	return e.DATASTAR_FOCUSSet(true)
}

// Sets the header of for fetch requests

func (e *SVGFEMORPHOLOGYElement) DATASTAR_HEADER(key string, expression string) *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key = fmt.Sprintf("data-header-%s", key)

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *SVGFEMORPHOLOGYElement) IfDATASTAR_HEADER(condition bool, key string, expression string) *SVGFEMORPHOLOGYElement {
	if condition {
		e.DATASTAR_HEADER(key, expression)
	}
	return e
}

// Remove the attribute DATASTAR_HEADER from the element.
func (e *SVGFEMORPHOLOGYElement) DATASTAR_HEADERRemove() *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-header")
	return e
}

// Sets the indicator selector for fetch requests

func (e *SVGFEMORPHOLOGYElement) DATASTAR_FETCH_INDICATOR(expression string) *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key := "data-fetch-indicator"

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *SVGFEMORPHOLOGYElement) IfDATASTAR_FETCH_INDICATOR(condition bool, expression string) *SVGFEMORPHOLOGYElement {
	if condition {
		e.DATASTAR_FETCH_INDICATOR(expression)
	}
	return e
}

// Remove the attribute DATASTAR_FETCH_INDICATOR from the element.
func (e *SVGFEMORPHOLOGYElement) DATASTAR_FETCH_INDICATORRemove() *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-fetch-indicator")
	return e
}

// Sets the visibility of the element

func (e *SVGFEMORPHOLOGYElement) DATASTAR_SHOWSet(b bool) *SVGFEMORPHOLOGYElement {
	key := "data-show"
	e.BoolAttributes.Set(key, b)
	return e
}

func (e *SVGFEMORPHOLOGYElement) DATASTAR_SHOW() *SVGFEMORPHOLOGYElement {
	return e.DATASTAR_SHOWSet(true)
}

// Triggers the callback when the element intersects the viewport

func (e *SVGFEMORPHOLOGYElement) DATASTAR_INTERSECTS(expression string) *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key := "data-intersects"

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *SVGFEMORPHOLOGYElement) IfDATASTAR_INTERSECTS(condition bool, expression string) *SVGFEMORPHOLOGYElement {
	if condition {
		e.DATASTAR_INTERSECTS(expression)
	}
	return e
}

// Remove the attribute DATASTAR_INTERSECTS from the element.
func (e *SVGFEMORPHOLOGYElement) DATASTAR_INTERSECTSRemove() *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-intersects")
	return e
}

// Teleports the element to the given selector

func (e *SVGFEMORPHOLOGYElement) DATASTAR_TELEPORTSet(b bool) *SVGFEMORPHOLOGYElement {
	key := "data-teleport"
	e.BoolAttributes.Set(key, b)
	return e
}

func (e *SVGFEMORPHOLOGYElement) DATASTAR_TELEPORT() *SVGFEMORPHOLOGYElement {
	return e.DATASTAR_TELEPORTSet(true)
}

// Scrolls the element into view

func (e *SVGFEMORPHOLOGYElement) DATASTAR_SCROLL_INTO_VIEWSet(b bool) *SVGFEMORPHOLOGYElement {
	key := "data-scroll-into-view"
	e.BoolAttributes.Set(key, b)
	return e
}

func (e *SVGFEMORPHOLOGYElement) DATASTAR_SCROLL_INTO_VIEW() *SVGFEMORPHOLOGYElement {
	return e.DATASTAR_SCROLL_INTO_VIEWSet(true)
}

// Setup the ViewTransitionAPI for the element

func (e *SVGFEMORPHOLOGYElement) DATASTAR_VIEW_TRANSITION(expression string) *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key := "data-view-transition"

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *SVGFEMORPHOLOGYElement) IfDATASTAR_VIEW_TRANSITION(condition bool, expression string) *SVGFEMORPHOLOGYElement {
	if condition {
		e.DATASTAR_VIEW_TRANSITION(expression)
	}
	return e
}

// Remove the attribute DATASTAR_VIEW_TRANSITION from the element.
func (e *SVGFEMORPHOLOGYElement) DATASTAR_VIEW_TRANSITIONRemove() *SVGFEMORPHOLOGYElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-view-transition")
	return e
}
