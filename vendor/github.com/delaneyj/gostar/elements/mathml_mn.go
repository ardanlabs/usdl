// DO NOT EDIT THIS FILE. IT IS AUTOMATICALLY GENERATED.
// package mathml mn is generated from configuration file.
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

// This element is used to display a single number.
type MathMLMNElement struct {
	*Element
}

// Create a new MathMLMNElement element.
// This will create a new element with the tag
// "mn" during rendering.
func MathML_MN(children ...ElementRenderer) *MathMLMNElement {
	e := NewElement("mn", children...)
	e.IsSelfClosing = false
	e.Descendants = children

	return &MathMLMNElement{Element: e}
}

func (e *MathMLMNElement) Children(children ...ElementRenderer) *MathMLMNElement {
	e.Descendants = append(e.Descendants, children...)
	return e
}

func (e *MathMLMNElement) IfChildren(condition bool, children ...ElementRenderer) *MathMLMNElement {
	if condition {
		e.Descendants = append(e.Descendants, children...)
	}
	return e
}

func (e *MathMLMNElement) TernChildren(condition bool, trueChildren, falseChildren ElementRenderer) *MathMLMNElement {
	if condition {
		e.Descendants = append(e.Descendants, trueChildren)
	} else {
		e.Descendants = append(e.Descendants, falseChildren)
	}
	return e
}

func (e *MathMLMNElement) Attr(name string, value string) *MathMLMNElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set(name, value)
	return e
}

func (e *MathMLMNElement) Attrs(attrs ...string) *MathMLMNElement {
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

func (e *MathMLMNElement) AttrsMap(attrs map[string]string) *MathMLMNElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	for k, v := range attrs {
		e.StringAttributes.Set(k, v)
	}
	return e
}

func (e *MathMLMNElement) Text(text string) *MathMLMNElement {
	e.Descendants = append(e.Descendants, Text(text))
	return e
}

func (e *MathMLMNElement) TextF(format string, args ...any) *MathMLMNElement {
	return e.Text(fmt.Sprintf(format, args...))
}

func (e *MathMLMNElement) IfText(condition bool, text string) *MathMLMNElement {
	if condition {
		e.Descendants = append(e.Descendants, Text(text))
	}
	return e
}

func (e *MathMLMNElement) IfTextF(condition bool, format string, args ...any) *MathMLMNElement {
	if condition {
		e.Descendants = append(e.Descendants, Text(fmt.Sprintf(format, args...)))
	}
	return e
}

func (e *MathMLMNElement) Escaped(text string) *MathMLMNElement {
	e.Descendants = append(e.Descendants, Escaped(text))
	return e
}

func (e *MathMLMNElement) IfEscaped(condition bool, text string) *MathMLMNElement {
	if condition {
		e.Descendants = append(e.Descendants, Escaped(text))
	}
	return e
}

func (e *MathMLMNElement) EscapedF(format string, args ...any) *MathMLMNElement {
	return e.Escaped(fmt.Sprintf(format, args...))
}

func (e *MathMLMNElement) IfEscapedF(condition bool, format string, args ...any) *MathMLMNElement {
	if condition {
		e.Descendants = append(e.Descendants, EscapedF(format, args...))
	}
	return e
}

func (e *MathMLMNElement) CustomData(key, value string) *MathMLMNElement {
	if e.CustomDataAttributes == nil {
		e.CustomDataAttributes = treemap.New[string, string]()
	}
	e.CustomDataAttributes.Set(key, value)
	return e
}

func (e *MathMLMNElement) IfCustomData(condition bool, key, value string) *MathMLMNElement {
	if condition {
		e.CustomData(key, value)
	}
	return e
}

func (e *MathMLMNElement) CustomDataF(key, format string, args ...any) *MathMLMNElement {
	return e.CustomData(key, fmt.Sprintf(format, args...))
}

func (e *MathMLMNElement) IfCustomDataF(condition bool, key, format string, args ...any) *MathMLMNElement {
	if condition {
		e.CustomData(key, fmt.Sprintf(format, args...))
	}
	return e
}

func (e *MathMLMNElement) CustomDataRemove(key string) *MathMLMNElement {
	if e.CustomDataAttributes == nil {
		return e
	}
	e.CustomDataAttributes.Del(key)
	return e
}

// Assigns a class name or set of class names to an element
// You may assign the same class name or names to any number of elements
// If you specify multiple class names, they must be separated by whitespace
// characters.
func (e *MathMLMNElement) CLASS(s ...string) *MathMLMNElement {
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

func (e *MathMLMNElement) IfCLASS(condition bool, s ...string) *MathMLMNElement {
	if condition {
		e.CLASS(s...)
	}
	return e
}

// Remove the attribute CLASS from the element.
func (e *MathMLMNElement) CLASSRemove(s ...string) *MathMLMNElement {
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

// This attribute specifies the text directionality of the element, merely
// indicating what direction the text flows when surrounded by text with inherent
// directionality (such as Arabic or Hebrew)
// Possible values are ltr (left-to-right) and rtl (right-to-left).
func (e *MathMLMNElement) DIR(c MathMLMnDirChoice) *MathMLMNElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set("dir", string(c))
	return e
}

type MathMLMnDirChoice string

const (
	// left-to-right
	MathMLMnDir_ltr MathMLMnDirChoice = "ltr"
	// right-to-left
	MathMLMnDir_rtl MathMLMnDirChoice = "rtl"
)

// Remove the attribute DIR from the element.
func (e *MathMLMNElement) DIRRemove(c MathMLMnDirChoice) *MathMLMNElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("dir")
	return e
}

// This attribute specifies whether the element should be rendered using
// displaystyle rules or not
// Possible values are true and false.
func (e *MathMLMNElement) DISPLAYSTYLE(c MathMLMnDisplaystyleChoice) *MathMLMNElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set("displaystyle", string(c))
	return e
}

type MathMLMnDisplaystyleChoice string

const (
	// displaystyle rules
	MathMLMnDisplaystyle_true MathMLMnDisplaystyleChoice = "true"
	// not displaystyle rules
	MathMLMnDisplaystyle_false MathMLMnDisplaystyleChoice = "false"
)

// Remove the attribute DISPLAYSTYLE from the element.
func (e *MathMLMNElement) DISPLAYSTYLERemove(c MathMLMnDisplaystyleChoice) *MathMLMNElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("displaystyle")
	return e
}

// This attribute assigns a name to an element
// This name must be unique in a document.
func (e *MathMLMNElement) ID(s string) *MathMLMNElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set("id", s)
	return e
}

func (e *MathMLMNElement) IDF(format string, args ...any) *MathMLMNElement {
	return e.ID(fmt.Sprintf(format, args...))
}

func (e *MathMLMNElement) IfID(condition bool, s string) *MathMLMNElement {
	if condition {
		e.ID(s)
	}
	return e
}

func (e *MathMLMNElement) IfIDF(condition bool, format string, args ...any) *MathMLMNElement {
	if condition {
		e.ID(fmt.Sprintf(format, args...))
	}
	return e
}

// Remove the attribute ID from the element.
func (e *MathMLMNElement) IDRemove(s string) *MathMLMNElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("id")
	return e
}

func (e *MathMLMNElement) IDRemoveF(format string, args ...any) *MathMLMNElement {
	return e.IDRemove(fmt.Sprintf(format, args...))
}

// This attribute specifies the background color of the element
// Possible values are a color name or a color specification in the format defined
// in the CSS3 Color Module [CSS3COLOR].
func (e *MathMLMNElement) MATHBACKGROUND(s string) *MathMLMNElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set("mathbackground", s)
	return e
}

func (e *MathMLMNElement) MATHBACKGROUNDF(format string, args ...any) *MathMLMNElement {
	return e.MATHBACKGROUND(fmt.Sprintf(format, args...))
}

func (e *MathMLMNElement) IfMATHBACKGROUND(condition bool, s string) *MathMLMNElement {
	if condition {
		e.MATHBACKGROUND(s)
	}
	return e
}

func (e *MathMLMNElement) IfMATHBACKGROUNDF(condition bool, format string, args ...any) *MathMLMNElement {
	if condition {
		e.MATHBACKGROUND(fmt.Sprintf(format, args...))
	}
	return e
}

// Remove the attribute MATHBACKGROUND from the element.
func (e *MathMLMNElement) MATHBACKGROUNDRemove(s string) *MathMLMNElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("mathbackground")
	return e
}

func (e *MathMLMNElement) MATHBACKGROUNDRemoveF(format string, args ...any) *MathMLMNElement {
	return e.MATHBACKGROUNDRemove(fmt.Sprintf(format, args...))
}

// This attribute specifies the color of the element
// Possible values are a color name or a color specification in the format defined
// in the CSS3 Color Module [CSS3COLOR].
func (e *MathMLMNElement) MATHCOLOR(s string) *MathMLMNElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set("mathcolor", s)
	return e
}

func (e *MathMLMNElement) MATHCOLORF(format string, args ...any) *MathMLMNElement {
	return e.MATHCOLOR(fmt.Sprintf(format, args...))
}

func (e *MathMLMNElement) IfMATHCOLOR(condition bool, s string) *MathMLMNElement {
	if condition {
		e.MATHCOLOR(s)
	}
	return e
}

func (e *MathMLMNElement) IfMATHCOLORF(condition bool, format string, args ...any) *MathMLMNElement {
	if condition {
		e.MATHCOLOR(fmt.Sprintf(format, args...))
	}
	return e
}

// Remove the attribute MATHCOLOR from the element.
func (e *MathMLMNElement) MATHCOLORRemove(s string) *MathMLMNElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("mathcolor")
	return e
}

func (e *MathMLMNElement) MATHCOLORRemoveF(format string, args ...any) *MathMLMNElement {
	return e.MATHCOLORRemove(fmt.Sprintf(format, args...))
}

// This attribute specifies the size of the element
// Possible values are a dimension or a dimensionless number.
func (e *MathMLMNElement) MATHSIZE_STR(s string) *MathMLMNElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set("mathsize", s)
	return e
}

func (e *MathMLMNElement) MATHSIZE_STRF(format string, args ...any) *MathMLMNElement {
	return e.MATHSIZE_STR(fmt.Sprintf(format, args...))
}

func (e *MathMLMNElement) IfMATHSIZE_STR(condition bool, s string) *MathMLMNElement {
	if condition {
		e.MATHSIZE_STR(s)
	}
	return e
}

func (e *MathMLMNElement) IfMATHSIZE_STRF(condition bool, format string, args ...any) *MathMLMNElement {
	if condition {
		e.MATHSIZE_STR(fmt.Sprintf(format, args...))
	}
	return e
}

// Remove the attribute MATHSIZE_STR from the element.
func (e *MathMLMNElement) MATHSIZE_STRRemove(s string) *MathMLMNElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("mathsize")
	return e
}

func (e *MathMLMNElement) MATHSIZE_STRRemoveF(format string, args ...any) *MathMLMNElement {
	return e.MATHSIZE_STRRemove(fmt.Sprintf(format, args...))
}

// This attribute declares a cryptographic nonce (number used once) that should be
// used by the server processing the element’s submission, and the resulting
// resource must be delivered with a Content-Security-Policy nonce attribute
// matching the value of the nonce attribute.
func (e *MathMLMNElement) NONCE(s string) *MathMLMNElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set("nonce", s)
	return e
}

func (e *MathMLMNElement) NONCEF(format string, args ...any) *MathMLMNElement {
	return e.NONCE(fmt.Sprintf(format, args...))
}

func (e *MathMLMNElement) IfNONCE(condition bool, s string) *MathMLMNElement {
	if condition {
		e.NONCE(s)
	}
	return e
}

func (e *MathMLMNElement) IfNONCEF(condition bool, format string, args ...any) *MathMLMNElement {
	if condition {
		e.NONCE(fmt.Sprintf(format, args...))
	}
	return e
}

// Remove the attribute NONCE from the element.
func (e *MathMLMNElement) NONCERemove(s string) *MathMLMNElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("nonce")
	return e
}

func (e *MathMLMNElement) NONCERemoveF(format string, args ...any) *MathMLMNElement {
	return e.NONCERemove(fmt.Sprintf(format, args...))
}

// This attribute specifies the script level of the element
// Possible values are an integer between 0 and 7, inclusive.
func (e *MathMLMNElement) SCRIPTLEVEL(i int) *MathMLMNElement {
	if e.IntAttributes == nil {
		e.IntAttributes = treemap.New[string, int]()
	}
	e.IntAttributes.Set("scriptlevel", i)
	return e
}

func (e *MathMLMNElement) IfSCRIPTLEVEL(condition bool, i int) *MathMLMNElement {
	if condition {
		e.SCRIPTLEVEL(i)
	}
	return e
}

// Remove the attribute SCRIPTLEVEL from the element.
func (e *MathMLMNElement) SCRIPTLEVELRemove(i int) *MathMLMNElement {
	if e.IntAttributes == nil {
		return e
	}
	e.IntAttributes.Del("scriptlevel")
	return e
}

// This attribute offers advisory information about the element for which it is
// set.
func (e *MathMLMNElement) STYLEF(k string, format string, args ...any) *MathMLMNElement {
	return e.STYLE(k, fmt.Sprintf(format, args...))
}

func (e *MathMLMNElement) IfSTYLE(condition bool, k string, v string) *MathMLMNElement {
	if condition {
		e.STYLE(k, v)
	}
	return e
}

func (e *MathMLMNElement) STYLE(k string, v string) *MathMLMNElement {
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

func (e *MathMLMNElement) IfSTYLEF(condition bool, k string, format string, args ...any) *MathMLMNElement {
	if condition {
		e.STYLE(k, fmt.Sprintf(format, args...))
	}
	return e
}

// Add the attributes in the map to the element.
func (e *MathMLMNElement) STYLEMap(m map[string]string) *MathMLMNElement {
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
func (e *MathMLMNElement) STYLEPairs(pairs ...string) *MathMLMNElement {
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

func (e *MathMLMNElement) IfSTYLEPairs(condition bool, pairs ...string) *MathMLMNElement {
	if condition {
		e.STYLEPairs(pairs...)
	}
	return e
}

// Remove the attribute STYLE from the element.
func (e *MathMLMNElement) STYLERemove(keys ...string) *MathMLMNElement {
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

// This attribute specifies the position of the current element in the tabbing
// order for the current document
// This value must be a number between 0 and 32767
// User agents should ignore leading zeros.
func (e *MathMLMNElement) TABINDEX(i int) *MathMLMNElement {
	if e.IntAttributes == nil {
		e.IntAttributes = treemap.New[string, int]()
	}
	e.IntAttributes.Set("tabindex", i)
	return e
}

func (e *MathMLMNElement) IfTABINDEX(condition bool, i int) *MathMLMNElement {
	if condition {
		e.TABINDEX(i)
	}
	return e
}

// Remove the attribute TABINDEX from the element.
func (e *MathMLMNElement) TABINDEXRemove(i int) *MathMLMNElement {
	if e.IntAttributes == nil {
		return e
	}
	e.IntAttributes.Del("tabindex")
	return e
}

// Merges the singleton store with the given object

func (e *MathMLMNElement) DATASTAR_STORE(v any) *MathMLMNElement {
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

func (e *MathMLMNElement) DATASTAR_REF(expression string) *MathMLMNElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key := "data-ref"

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *MathMLMNElement) IfDATASTAR_REF(condition bool, expression string) *MathMLMNElement {
	if condition {
		e.DATASTAR_REF(expression)
	}
	return e
}

// Remove the attribute DATASTAR_REF from the element.
func (e *MathMLMNElement) DATASTAR_REFRemove() *MathMLMNElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-ref")
	return e
}

// Sets the value of the element

func (e *MathMLMNElement) DATASTAR_BIND(key string, expression string) *MathMLMNElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key = fmt.Sprintf("data-bind-%s", key)

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *MathMLMNElement) IfDATASTAR_BIND(condition bool, key string, expression string) *MathMLMNElement {
	if condition {
		e.DATASTAR_BIND(key, expression)
	}
	return e
}

// Remove the attribute DATASTAR_BIND from the element.
func (e *MathMLMNElement) DATASTAR_BINDRemove() *MathMLMNElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-bind")
	return e
}

// Sets the value of the element

func (e *MathMLMNElement) DATASTAR_MODEL(expression string) *MathMLMNElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key := "data-model"

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *MathMLMNElement) IfDATASTAR_MODEL(condition bool, expression string) *MathMLMNElement {
	if condition {
		e.DATASTAR_MODEL(expression)
	}
	return e
}

// Remove the attribute DATASTAR_MODEL from the element.
func (e *MathMLMNElement) DATASTAR_MODELRemove() *MathMLMNElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-model")
	return e
}

// Sets the textContent of the element

func (e *MathMLMNElement) DATASTAR_TEXT(expression string) *MathMLMNElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key := "data-text"

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *MathMLMNElement) IfDATASTAR_TEXT(condition bool, expression string) *MathMLMNElement {
	if condition {
		e.DATASTAR_TEXT(expression)
	}
	return e
}

// Remove the attribute DATASTAR_TEXT from the element.
func (e *MathMLMNElement) DATASTAR_TEXTRemove() *MathMLMNElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-text")
	return e
}

// Sets the event handler of the element

type MathMLMnOnMod customDataKeyModifier

// Debounces the event handler
func MathMLMnOnModDebounce(
	d time.Duration,
) MathMLMnOnMod {
	return func() string {
		return fmt.Sprintf("debounce_%dms", d.Milliseconds())
	}
}

// Throttles the event handler
func MathMLMnOnModThrottle(
	d time.Duration,
) MathMLMnOnMod {
	return func() string {
		return fmt.Sprintf("throttle_%dms", d.Milliseconds())
	}
}

func (e *MathMLMNElement) DATASTAR_ON(key string, expression string, modifiers ...MathMLMnOnMod) *MathMLMNElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key = fmt.Sprintf("data-on-%s", key)

	customMods := lo.Map(modifiers, func(m MathMLMnOnMod, i int) customDataKeyModifier {
		return customDataKeyModifier(m)
	})
	key = customDataKey(key, customMods...)
	e.StringAttributes.Set(key, expression)
	return e
}

func (e *MathMLMNElement) IfDATASTAR_ON(condition bool, key string, expression string, modifiers ...MathMLMnOnMod) *MathMLMNElement {
	if condition {
		e.DATASTAR_ON(key, expression, modifiers...)
	}
	return e
}

// Remove the attribute DATASTAR_ON from the element.
func (e *MathMLMNElement) DATASTAR_ONRemove() *MathMLMNElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-on")
	return e
}

// Sets the focus of the element

func (e *MathMLMNElement) DATASTAR_FOCUSSet(b bool) *MathMLMNElement {
	key := "data-focus"
	e.BoolAttributes.Set(key, b)
	return e
}

func (e *MathMLMNElement) DATASTAR_FOCUS() *MathMLMNElement {
	return e.DATASTAR_FOCUSSet(true)
}

// Sets the header of for fetch requests

func (e *MathMLMNElement) DATASTAR_HEADER(key string, expression string) *MathMLMNElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key = fmt.Sprintf("data-header-%s", key)

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *MathMLMNElement) IfDATASTAR_HEADER(condition bool, key string, expression string) *MathMLMNElement {
	if condition {
		e.DATASTAR_HEADER(key, expression)
	}
	return e
}

// Remove the attribute DATASTAR_HEADER from the element.
func (e *MathMLMNElement) DATASTAR_HEADERRemove() *MathMLMNElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-header")
	return e
}

// Sets the indicator selector for fetch requests

func (e *MathMLMNElement) DATASTAR_FETCH_INDICATOR(expression string) *MathMLMNElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key := "data-fetch-indicator"

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *MathMLMNElement) IfDATASTAR_FETCH_INDICATOR(condition bool, expression string) *MathMLMNElement {
	if condition {
		e.DATASTAR_FETCH_INDICATOR(expression)
	}
	return e
}

// Remove the attribute DATASTAR_FETCH_INDICATOR from the element.
func (e *MathMLMNElement) DATASTAR_FETCH_INDICATORRemove() *MathMLMNElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-fetch-indicator")
	return e
}

// Sets the visibility of the element

func (e *MathMLMNElement) DATASTAR_SHOWSet(b bool) *MathMLMNElement {
	key := "data-show"
	e.BoolAttributes.Set(key, b)
	return e
}

func (e *MathMLMNElement) DATASTAR_SHOW() *MathMLMNElement {
	return e.DATASTAR_SHOWSet(true)
}

// Triggers the callback when the element intersects the viewport

func (e *MathMLMNElement) DATASTAR_INTERSECTS(expression string) *MathMLMNElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key := "data-intersects"

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *MathMLMNElement) IfDATASTAR_INTERSECTS(condition bool, expression string) *MathMLMNElement {
	if condition {
		e.DATASTAR_INTERSECTS(expression)
	}
	return e
}

// Remove the attribute DATASTAR_INTERSECTS from the element.
func (e *MathMLMNElement) DATASTAR_INTERSECTSRemove() *MathMLMNElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-intersects")
	return e
}

// Teleports the element to the given selector

func (e *MathMLMNElement) DATASTAR_TELEPORTSet(b bool) *MathMLMNElement {
	key := "data-teleport"
	e.BoolAttributes.Set(key, b)
	return e
}

func (e *MathMLMNElement) DATASTAR_TELEPORT() *MathMLMNElement {
	return e.DATASTAR_TELEPORTSet(true)
}

// Scrolls the element into view

func (e *MathMLMNElement) DATASTAR_SCROLL_INTO_VIEWSet(b bool) *MathMLMNElement {
	key := "data-scroll-into-view"
	e.BoolAttributes.Set(key, b)
	return e
}

func (e *MathMLMNElement) DATASTAR_SCROLL_INTO_VIEW() *MathMLMNElement {
	return e.DATASTAR_SCROLL_INTO_VIEWSet(true)
}

// Setup the ViewTransitionAPI for the element

func (e *MathMLMNElement) DATASTAR_VIEW_TRANSITION(expression string) *MathMLMNElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key := "data-view-transition"

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *MathMLMNElement) IfDATASTAR_VIEW_TRANSITION(condition bool, expression string) *MathMLMNElement {
	if condition {
		e.DATASTAR_VIEW_TRANSITION(expression)
	}
	return e
}

// Remove the attribute DATASTAR_VIEW_TRANSITION from the element.
func (e *MathMLMNElement) DATASTAR_VIEW_TRANSITIONRemove() *MathMLMNElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-view-transition")
	return e
}
