// DO NOT EDIT THIS FILE. IT IS AUTOMATICALLY GENERATED.
// package svg feDiffuseLighting is generated from configuration file.
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

// The <feDiffuseLighting> SVG filter primitive lights an image using the alpha
// channel as a bump map
// The resulting image, which is an RGBA opaque image, depends on the light color,
// light position and surface geometry of the input bump map.
type SVGFEDIFFUSELIGHTINGElement struct {
	*Element
}

// Create a new SVGFEDIFFUSELIGHTINGElement element.
// This will create a new element with the tag
// "feDiffuseLighting" during rendering.
func SVG_FEDIFFUSELIGHTING(children ...ElementRenderer) *SVGFEDIFFUSELIGHTINGElement {
	e := NewElement("feDiffuseLighting", children...)
	e.IsSelfClosing = false
	e.Descendants = children

	return &SVGFEDIFFUSELIGHTINGElement{Element: e}
}

func (e *SVGFEDIFFUSELIGHTINGElement) Children(children ...ElementRenderer) *SVGFEDIFFUSELIGHTINGElement {
	e.Descendants = append(e.Descendants, children...)
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) IfChildren(condition bool, children ...ElementRenderer) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.Descendants = append(e.Descendants, children...)
	}
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) TernChildren(condition bool, trueChildren, falseChildren ElementRenderer) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.Descendants = append(e.Descendants, trueChildren)
	} else {
		e.Descendants = append(e.Descendants, falseChildren)
	}
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) Attr(name string, value string) *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set(name, value)
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) Attrs(attrs ...string) *SVGFEDIFFUSELIGHTINGElement {
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

func (e *SVGFEDIFFUSELIGHTINGElement) AttrsMap(attrs map[string]string) *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	for k, v := range attrs {
		e.StringAttributes.Set(k, v)
	}
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) Text(text string) *SVGFEDIFFUSELIGHTINGElement {
	e.Descendants = append(e.Descendants, Text(text))
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) TextF(format string, args ...any) *SVGFEDIFFUSELIGHTINGElement {
	return e.Text(fmt.Sprintf(format, args...))
}

func (e *SVGFEDIFFUSELIGHTINGElement) IfText(condition bool, text string) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.Descendants = append(e.Descendants, Text(text))
	}
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) IfTextF(condition bool, format string, args ...any) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.Descendants = append(e.Descendants, Text(fmt.Sprintf(format, args...)))
	}
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) Escaped(text string) *SVGFEDIFFUSELIGHTINGElement {
	e.Descendants = append(e.Descendants, Escaped(text))
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) IfEscaped(condition bool, text string) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.Descendants = append(e.Descendants, Escaped(text))
	}
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) EscapedF(format string, args ...any) *SVGFEDIFFUSELIGHTINGElement {
	return e.Escaped(fmt.Sprintf(format, args...))
}

func (e *SVGFEDIFFUSELIGHTINGElement) IfEscapedF(condition bool, format string, args ...any) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.Descendants = append(e.Descendants, EscapedF(format, args...))
	}
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) CustomData(key, value string) *SVGFEDIFFUSELIGHTINGElement {
	if e.CustomDataAttributes == nil {
		e.CustomDataAttributes = treemap.New[string, string]()
	}
	e.CustomDataAttributes.Set(key, value)
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) IfCustomData(condition bool, key, value string) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.CustomData(key, value)
	}
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) CustomDataF(key, format string, args ...any) *SVGFEDIFFUSELIGHTINGElement {
	return e.CustomData(key, fmt.Sprintf(format, args...))
}

func (e *SVGFEDIFFUSELIGHTINGElement) IfCustomDataF(condition bool, key, format string, args ...any) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.CustomData(key, fmt.Sprintf(format, args...))
	}
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) CustomDataRemove(key string) *SVGFEDIFFUSELIGHTINGElement {
	if e.CustomDataAttributes == nil {
		return e
	}
	e.CustomDataAttributes.Del(key)
	return e
}

// The input for this filter.
func (e *SVGFEDIFFUSELIGHTINGElement) IN(s string) *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set("in", s)
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) INF(format string, args ...any) *SVGFEDIFFUSELIGHTINGElement {
	return e.IN(fmt.Sprintf(format, args...))
}

func (e *SVGFEDIFFUSELIGHTINGElement) IfIN(condition bool, s string) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.IN(s)
	}
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) IfINF(condition bool, format string, args ...any) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.IN(fmt.Sprintf(format, args...))
	}
	return e
}

// Remove the attribute IN from the element.
func (e *SVGFEDIFFUSELIGHTINGElement) INRemove(s string) *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("in")
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) INRemoveF(format string, args ...any) *SVGFEDIFFUSELIGHTINGElement {
	return e.INRemove(fmt.Sprintf(format, args...))
}

// The 'surfaceScale' attribute indicates the height of the surface when the alpha
// channel is 1.0.
func (e *SVGFEDIFFUSELIGHTINGElement) SURFACE_SCALE(f float64) *SVGFEDIFFUSELIGHTINGElement {
	if e.FloatAttributes == nil {
		e.FloatAttributes = treemap.New[string, float64]()
	}
	e.FloatAttributes.Set("surfaceScale", f)
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) IfSURFACE_SCALE(condition bool, f float64) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.SURFACE_SCALE(f)
	}
	return e
}

// The diffuseConstant attribute represents the proportion of the light that is
// reflected by the surface.
func (e *SVGFEDIFFUSELIGHTINGElement) DIFFUSE_CONSTANT(f float64) *SVGFEDIFFUSELIGHTINGElement {
	if e.FloatAttributes == nil {
		e.FloatAttributes = treemap.New[string, float64]()
	}
	e.FloatAttributes.Set("diffuseConstant", f)
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) IfDIFFUSE_CONSTANT(condition bool, f float64) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.DIFFUSE_CONSTANT(f)
	}
	return e
}

// The kernelUnitLength attribute defines the intended distance in current filter
// units (i.e., units as determined by the value of attribute 'primitiveUnits')
// for dx and dy in the surface normal calculation formulas.
func (e *SVGFEDIFFUSELIGHTINGElement) KERNEL_UNIT_LENGTH(s string) *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set("kernelUnitLength", s)
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) KERNEL_UNIT_LENGTHF(format string, args ...any) *SVGFEDIFFUSELIGHTINGElement {
	return e.KERNEL_UNIT_LENGTH(fmt.Sprintf(format, args...))
}

func (e *SVGFEDIFFUSELIGHTINGElement) IfKERNEL_UNIT_LENGTH(condition bool, s string) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.KERNEL_UNIT_LENGTH(s)
	}
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) IfKERNEL_UNIT_LENGTHF(condition bool, format string, args ...any) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.KERNEL_UNIT_LENGTH(fmt.Sprintf(format, args...))
	}
	return e
}

// Remove the attribute KERNEL_UNIT_LENGTH from the element.
func (e *SVGFEDIFFUSELIGHTINGElement) KERNEL_UNIT_LENGTHRemove(s string) *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("kernelUnitLength")
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) KERNEL_UNIT_LENGTHRemoveF(format string, args ...any) *SVGFEDIFFUSELIGHTINGElement {
	return e.KERNEL_UNIT_LENGTHRemove(fmt.Sprintf(format, args...))
}

// Specifies a unique id for an element
func (e *SVGFEDIFFUSELIGHTINGElement) ID(s string) *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set("id", s)
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) IDF(format string, args ...any) *SVGFEDIFFUSELIGHTINGElement {
	return e.ID(fmt.Sprintf(format, args...))
}

func (e *SVGFEDIFFUSELIGHTINGElement) IfID(condition bool, s string) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.ID(s)
	}
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) IfIDF(condition bool, format string, args ...any) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.ID(fmt.Sprintf(format, args...))
	}
	return e
}

// Remove the attribute ID from the element.
func (e *SVGFEDIFFUSELIGHTINGElement) IDRemove(s string) *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("id")
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) IDRemoveF(format string, args ...any) *SVGFEDIFFUSELIGHTINGElement {
	return e.IDRemove(fmt.Sprintf(format, args...))
}

// Specifies one or more classnames for an element (refers to a class in a style
// sheet)
func (e *SVGFEDIFFUSELIGHTINGElement) CLASS(s ...string) *SVGFEDIFFUSELIGHTINGElement {
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

func (e *SVGFEDIFFUSELIGHTINGElement) IfCLASS(condition bool, s ...string) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.CLASS(s...)
	}
	return e
}

// Remove the attribute CLASS from the element.
func (e *SVGFEDIFFUSELIGHTINGElement) CLASSRemove(s ...string) *SVGFEDIFFUSELIGHTINGElement {
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
func (e *SVGFEDIFFUSELIGHTINGElement) STYLEF(k string, format string, args ...any) *SVGFEDIFFUSELIGHTINGElement {
	return e.STYLE(k, fmt.Sprintf(format, args...))
}

func (e *SVGFEDIFFUSELIGHTINGElement) IfSTYLE(condition bool, k string, v string) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.STYLE(k, v)
	}
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) STYLE(k string, v string) *SVGFEDIFFUSELIGHTINGElement {
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

func (e *SVGFEDIFFUSELIGHTINGElement) IfSTYLEF(condition bool, k string, format string, args ...any) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.STYLE(k, fmt.Sprintf(format, args...))
	}
	return e
}

// Add the attributes in the map to the element.
func (e *SVGFEDIFFUSELIGHTINGElement) STYLEMap(m map[string]string) *SVGFEDIFFUSELIGHTINGElement {
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
func (e *SVGFEDIFFUSELIGHTINGElement) STYLEPairs(pairs ...string) *SVGFEDIFFUSELIGHTINGElement {
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

func (e *SVGFEDIFFUSELIGHTINGElement) IfSTYLEPairs(condition bool, pairs ...string) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.STYLEPairs(pairs...)
	}
	return e
}

// Remove the attribute STYLE from the element.
func (e *SVGFEDIFFUSELIGHTINGElement) STYLERemove(keys ...string) *SVGFEDIFFUSELIGHTINGElement {
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

func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_STORE(v any) *SVGFEDIFFUSELIGHTINGElement {
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

func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_REF(expression string) *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key := "data-ref"

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) IfDATASTAR_REF(condition bool, expression string) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.DATASTAR_REF(expression)
	}
	return e
}

// Remove the attribute DATASTAR_REF from the element.
func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_REFRemove() *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-ref")
	return e
}

// Sets the value of the element

func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_BIND(key string, expression string) *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key = fmt.Sprintf("data-bind-%s", key)

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) IfDATASTAR_BIND(condition bool, key string, expression string) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.DATASTAR_BIND(key, expression)
	}
	return e
}

// Remove the attribute DATASTAR_BIND from the element.
func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_BINDRemove() *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-bind")
	return e
}

// Sets the value of the element

func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_MODEL(expression string) *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key := "data-model"

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) IfDATASTAR_MODEL(condition bool, expression string) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.DATASTAR_MODEL(expression)
	}
	return e
}

// Remove the attribute DATASTAR_MODEL from the element.
func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_MODELRemove() *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-model")
	return e
}

// Sets the textContent of the element

func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_TEXT(expression string) *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key := "data-text"

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) IfDATASTAR_TEXT(condition bool, expression string) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.DATASTAR_TEXT(expression)
	}
	return e
}

// Remove the attribute DATASTAR_TEXT from the element.
func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_TEXTRemove() *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-text")
	return e
}

// Sets the event handler of the element

type SVGFeDiffuseLightingOnMod customDataKeyModifier

// Debounces the event handler
func SVGFeDiffuseLightingOnModDebounce(
	d time.Duration,
) SVGFeDiffuseLightingOnMod {
	return func() string {
		return fmt.Sprintf("debounce_%dms", d.Milliseconds())
	}
}

// Throttles the event handler
func SVGFeDiffuseLightingOnModThrottle(
	d time.Duration,
) SVGFeDiffuseLightingOnMod {
	return func() string {
		return fmt.Sprintf("throttle_%dms", d.Milliseconds())
	}
}

func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_ON(key string, expression string, modifiers ...SVGFeDiffuseLightingOnMod) *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key = fmt.Sprintf("data-on-%s", key)

	customMods := lo.Map(modifiers, func(m SVGFeDiffuseLightingOnMod, i int) customDataKeyModifier {
		return customDataKeyModifier(m)
	})
	key = customDataKey(key, customMods...)
	e.StringAttributes.Set(key, expression)
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) IfDATASTAR_ON(condition bool, key string, expression string, modifiers ...SVGFeDiffuseLightingOnMod) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.DATASTAR_ON(key, expression, modifiers...)
	}
	return e
}

// Remove the attribute DATASTAR_ON from the element.
func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_ONRemove() *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-on")
	return e
}

// Sets the focus of the element

func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_FOCUSSet(b bool) *SVGFEDIFFUSELIGHTINGElement {
	key := "data-focus"
	e.BoolAttributes.Set(key, b)
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_FOCUS() *SVGFEDIFFUSELIGHTINGElement {
	return e.DATASTAR_FOCUSSet(true)
}

// Sets the header of for fetch requests

func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_HEADER(key string, expression string) *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key = fmt.Sprintf("data-header-%s", key)

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) IfDATASTAR_HEADER(condition bool, key string, expression string) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.DATASTAR_HEADER(key, expression)
	}
	return e
}

// Remove the attribute DATASTAR_HEADER from the element.
func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_HEADERRemove() *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-header")
	return e
}

// Sets the indicator selector for fetch requests

func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_FETCH_INDICATOR(expression string) *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key := "data-fetch-indicator"

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) IfDATASTAR_FETCH_INDICATOR(condition bool, expression string) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.DATASTAR_FETCH_INDICATOR(expression)
	}
	return e
}

// Remove the attribute DATASTAR_FETCH_INDICATOR from the element.
func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_FETCH_INDICATORRemove() *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-fetch-indicator")
	return e
}

// Sets the visibility of the element

func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_SHOWSet(b bool) *SVGFEDIFFUSELIGHTINGElement {
	key := "data-show"
	e.BoolAttributes.Set(key, b)
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_SHOW() *SVGFEDIFFUSELIGHTINGElement {
	return e.DATASTAR_SHOWSet(true)
}

// Triggers the callback when the element intersects the viewport

func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_INTERSECTS(expression string) *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key := "data-intersects"

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) IfDATASTAR_INTERSECTS(condition bool, expression string) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.DATASTAR_INTERSECTS(expression)
	}
	return e
}

// Remove the attribute DATASTAR_INTERSECTS from the element.
func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_INTERSECTSRemove() *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-intersects")
	return e
}

// Teleports the element to the given selector

func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_TELEPORTSet(b bool) *SVGFEDIFFUSELIGHTINGElement {
	key := "data-teleport"
	e.BoolAttributes.Set(key, b)
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_TELEPORT() *SVGFEDIFFUSELIGHTINGElement {
	return e.DATASTAR_TELEPORTSet(true)
}

// Scrolls the element into view

func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_SCROLL_INTO_VIEWSet(b bool) *SVGFEDIFFUSELIGHTINGElement {
	key := "data-scroll-into-view"
	e.BoolAttributes.Set(key, b)
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_SCROLL_INTO_VIEW() *SVGFEDIFFUSELIGHTINGElement {
	return e.DATASTAR_SCROLL_INTO_VIEWSet(true)
}

// Setup the ViewTransitionAPI for the element

func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_VIEW_TRANSITION(expression string) *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key := "data-view-transition"

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *SVGFEDIFFUSELIGHTINGElement) IfDATASTAR_VIEW_TRANSITION(condition bool, expression string) *SVGFEDIFFUSELIGHTINGElement {
	if condition {
		e.DATASTAR_VIEW_TRANSITION(expression)
	}
	return e
}

// Remove the attribute DATASTAR_VIEW_TRANSITION from the element.
func (e *SVGFEDIFFUSELIGHTINGElement) DATASTAR_VIEW_TRANSITIONRemove() *SVGFEDIFFUSELIGHTINGElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-view-transition")
	return e
}
