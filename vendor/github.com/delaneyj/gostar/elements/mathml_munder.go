// DO NOT EDIT THIS FILE. IT IS AUTOMATICALLY GENERATED.
// package mathml munder is generated from configuration file.
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

// This element is used to display an expression with an underbar.
type MathMLMUNDERElement struct {
	*Element
}

// Create a new MathMLMUNDERElement element.
// This will create a new element with the tag
// "munder" during rendering.
func MathML_MUNDER(children ...ElementRenderer) *MathMLMUNDERElement {
	e := NewElement("munder", children...)
	e.IsSelfClosing = false
	e.Descendants = children

	return &MathMLMUNDERElement{Element: e}
}

func (e *MathMLMUNDERElement) Children(children ...ElementRenderer) *MathMLMUNDERElement {
	e.Descendants = append(e.Descendants, children...)
	return e
}

func (e *MathMLMUNDERElement) IfChildren(condition bool, children ...ElementRenderer) *MathMLMUNDERElement {
	if condition {
		e.Descendants = append(e.Descendants, children...)
	}
	return e
}

func (e *MathMLMUNDERElement) TernChildren(condition bool, trueChildren, falseChildren ElementRenderer) *MathMLMUNDERElement {
	if condition {
		e.Descendants = append(e.Descendants, trueChildren)
	} else {
		e.Descendants = append(e.Descendants, falseChildren)
	}
	return e
}

func (e *MathMLMUNDERElement) Attr(name string, value string) *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set(name, value)
	return e
}

func (e *MathMLMUNDERElement) Attrs(attrs ...string) *MathMLMUNDERElement {
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

func (e *MathMLMUNDERElement) AttrsMap(attrs map[string]string) *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	for k, v := range attrs {
		e.StringAttributes.Set(k, v)
	}
	return e
}

func (e *MathMLMUNDERElement) Text(text string) *MathMLMUNDERElement {
	e.Descendants = append(e.Descendants, Text(text))
	return e
}

func (e *MathMLMUNDERElement) TextF(format string, args ...any) *MathMLMUNDERElement {
	return e.Text(fmt.Sprintf(format, args...))
}

func (e *MathMLMUNDERElement) IfText(condition bool, text string) *MathMLMUNDERElement {
	if condition {
		e.Descendants = append(e.Descendants, Text(text))
	}
	return e
}

func (e *MathMLMUNDERElement) IfTextF(condition bool, format string, args ...any) *MathMLMUNDERElement {
	if condition {
		e.Descendants = append(e.Descendants, Text(fmt.Sprintf(format, args...)))
	}
	return e
}

func (e *MathMLMUNDERElement) Escaped(text string) *MathMLMUNDERElement {
	e.Descendants = append(e.Descendants, Escaped(text))
	return e
}

func (e *MathMLMUNDERElement) IfEscaped(condition bool, text string) *MathMLMUNDERElement {
	if condition {
		e.Descendants = append(e.Descendants, Escaped(text))
	}
	return e
}

func (e *MathMLMUNDERElement) EscapedF(format string, args ...any) *MathMLMUNDERElement {
	return e.Escaped(fmt.Sprintf(format, args...))
}

func (e *MathMLMUNDERElement) IfEscapedF(condition bool, format string, args ...any) *MathMLMUNDERElement {
	if condition {
		e.Descendants = append(e.Descendants, EscapedF(format, args...))
	}
	return e
}

func (e *MathMLMUNDERElement) CustomData(key, value string) *MathMLMUNDERElement {
	if e.CustomDataAttributes == nil {
		e.CustomDataAttributes = treemap.New[string, string]()
	}
	e.CustomDataAttributes.Set(key, value)
	return e
}

func (e *MathMLMUNDERElement) IfCustomData(condition bool, key, value string) *MathMLMUNDERElement {
	if condition {
		e.CustomData(key, value)
	}
	return e
}

func (e *MathMLMUNDERElement) CustomDataF(key, format string, args ...any) *MathMLMUNDERElement {
	return e.CustomData(key, fmt.Sprintf(format, args...))
}

func (e *MathMLMUNDERElement) IfCustomDataF(condition bool, key, format string, args ...any) *MathMLMUNDERElement {
	if condition {
		e.CustomData(key, fmt.Sprintf(format, args...))
	}
	return e
}

func (e *MathMLMUNDERElement) CustomDataRemove(key string) *MathMLMUNDERElement {
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
func (e *MathMLMUNDERElement) CLASS(s ...string) *MathMLMUNDERElement {
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

func (e *MathMLMUNDERElement) IfCLASS(condition bool, s ...string) *MathMLMUNDERElement {
	if condition {
		e.CLASS(s...)
	}
	return e
}

// Remove the attribute CLASS from the element.
func (e *MathMLMUNDERElement) CLASSRemove(s ...string) *MathMLMUNDERElement {
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
func (e *MathMLMUNDERElement) DIR(c MathMLMunderDirChoice) *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set("dir", string(c))
	return e
}

type MathMLMunderDirChoice string

const (
	// left-to-right
	MathMLMunderDir_ltr MathMLMunderDirChoice = "ltr"
	// right-to-left
	MathMLMunderDir_rtl MathMLMunderDirChoice = "rtl"
)

// Remove the attribute DIR from the element.
func (e *MathMLMUNDERElement) DIRRemove(c MathMLMunderDirChoice) *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("dir")
	return e
}

// This attribute specifies whether the element should be rendered using
// displaystyle rules or not
// Possible values are true and false.
func (e *MathMLMUNDERElement) DISPLAYSTYLE(c MathMLMunderDisplaystyleChoice) *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set("displaystyle", string(c))
	return e
}

type MathMLMunderDisplaystyleChoice string

const (
	// displaystyle rules
	MathMLMunderDisplaystyle_true MathMLMunderDisplaystyleChoice = "true"
	// not displaystyle rules
	MathMLMunderDisplaystyle_false MathMLMunderDisplaystyleChoice = "false"
)

// Remove the attribute DISPLAYSTYLE from the element.
func (e *MathMLMUNDERElement) DISPLAYSTYLERemove(c MathMLMunderDisplaystyleChoice) *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("displaystyle")
	return e
}

// This attribute assigns a name to an element
// This name must be unique in a document.
func (e *MathMLMUNDERElement) ID(s string) *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set("id", s)
	return e
}

func (e *MathMLMUNDERElement) IDF(format string, args ...any) *MathMLMUNDERElement {
	return e.ID(fmt.Sprintf(format, args...))
}

func (e *MathMLMUNDERElement) IfID(condition bool, s string) *MathMLMUNDERElement {
	if condition {
		e.ID(s)
	}
	return e
}

func (e *MathMLMUNDERElement) IfIDF(condition bool, format string, args ...any) *MathMLMUNDERElement {
	if condition {
		e.ID(fmt.Sprintf(format, args...))
	}
	return e
}

// Remove the attribute ID from the element.
func (e *MathMLMUNDERElement) IDRemove(s string) *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("id")
	return e
}

func (e *MathMLMUNDERElement) IDRemoveF(format string, args ...any) *MathMLMUNDERElement {
	return e.IDRemove(fmt.Sprintf(format, args...))
}

// This attribute specifies the background color of the element
// Possible values are a color name or a color specification in the format defined
// in the CSS3 Color Module [CSS3COLOR].
func (e *MathMLMUNDERElement) MATHBACKGROUND(s string) *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set("mathbackground", s)
	return e
}

func (e *MathMLMUNDERElement) MATHBACKGROUNDF(format string, args ...any) *MathMLMUNDERElement {
	return e.MATHBACKGROUND(fmt.Sprintf(format, args...))
}

func (e *MathMLMUNDERElement) IfMATHBACKGROUND(condition bool, s string) *MathMLMUNDERElement {
	if condition {
		e.MATHBACKGROUND(s)
	}
	return e
}

func (e *MathMLMUNDERElement) IfMATHBACKGROUNDF(condition bool, format string, args ...any) *MathMLMUNDERElement {
	if condition {
		e.MATHBACKGROUND(fmt.Sprintf(format, args...))
	}
	return e
}

// Remove the attribute MATHBACKGROUND from the element.
func (e *MathMLMUNDERElement) MATHBACKGROUNDRemove(s string) *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("mathbackground")
	return e
}

func (e *MathMLMUNDERElement) MATHBACKGROUNDRemoveF(format string, args ...any) *MathMLMUNDERElement {
	return e.MATHBACKGROUNDRemove(fmt.Sprintf(format, args...))
}

// This attribute specifies the color of the element
// Possible values are a color name or a color specification in the format defined
// in the CSS3 Color Module [CSS3COLOR].
func (e *MathMLMUNDERElement) MATHCOLOR(s string) *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set("mathcolor", s)
	return e
}

func (e *MathMLMUNDERElement) MATHCOLORF(format string, args ...any) *MathMLMUNDERElement {
	return e.MATHCOLOR(fmt.Sprintf(format, args...))
}

func (e *MathMLMUNDERElement) IfMATHCOLOR(condition bool, s string) *MathMLMUNDERElement {
	if condition {
		e.MATHCOLOR(s)
	}
	return e
}

func (e *MathMLMUNDERElement) IfMATHCOLORF(condition bool, format string, args ...any) *MathMLMUNDERElement {
	if condition {
		e.MATHCOLOR(fmt.Sprintf(format, args...))
	}
	return e
}

// Remove the attribute MATHCOLOR from the element.
func (e *MathMLMUNDERElement) MATHCOLORRemove(s string) *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("mathcolor")
	return e
}

func (e *MathMLMUNDERElement) MATHCOLORRemoveF(format string, args ...any) *MathMLMUNDERElement {
	return e.MATHCOLORRemove(fmt.Sprintf(format, args...))
}

// This attribute specifies the size of the element
// Possible values are a dimension or a dimensionless number.
func (e *MathMLMUNDERElement) MATHSIZE_STR(s string) *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set("mathsize", s)
	return e
}

func (e *MathMLMUNDERElement) MATHSIZE_STRF(format string, args ...any) *MathMLMUNDERElement {
	return e.MATHSIZE_STR(fmt.Sprintf(format, args...))
}

func (e *MathMLMUNDERElement) IfMATHSIZE_STR(condition bool, s string) *MathMLMUNDERElement {
	if condition {
		e.MATHSIZE_STR(s)
	}
	return e
}

func (e *MathMLMUNDERElement) IfMATHSIZE_STRF(condition bool, format string, args ...any) *MathMLMUNDERElement {
	if condition {
		e.MATHSIZE_STR(fmt.Sprintf(format, args...))
	}
	return e
}

// Remove the attribute MATHSIZE_STR from the element.
func (e *MathMLMUNDERElement) MATHSIZE_STRRemove(s string) *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("mathsize")
	return e
}

func (e *MathMLMUNDERElement) MATHSIZE_STRRemoveF(format string, args ...any) *MathMLMUNDERElement {
	return e.MATHSIZE_STRRemove(fmt.Sprintf(format, args...))
}

// This attribute declares a cryptographic nonce (number used once) that should be
// used by the server processing the element’s submission, and the resulting
// resource must be delivered with a Content-Security-Policy nonce attribute
// matching the value of the nonce attribute.
func (e *MathMLMUNDERElement) NONCE(s string) *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}
	e.StringAttributes.Set("nonce", s)
	return e
}

func (e *MathMLMUNDERElement) NONCEF(format string, args ...any) *MathMLMUNDERElement {
	return e.NONCE(fmt.Sprintf(format, args...))
}

func (e *MathMLMUNDERElement) IfNONCE(condition bool, s string) *MathMLMUNDERElement {
	if condition {
		e.NONCE(s)
	}
	return e
}

func (e *MathMLMUNDERElement) IfNONCEF(condition bool, format string, args ...any) *MathMLMUNDERElement {
	if condition {
		e.NONCE(fmt.Sprintf(format, args...))
	}
	return e
}

// Remove the attribute NONCE from the element.
func (e *MathMLMUNDERElement) NONCERemove(s string) *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("nonce")
	return e
}

func (e *MathMLMUNDERElement) NONCERemoveF(format string, args ...any) *MathMLMUNDERElement {
	return e.NONCERemove(fmt.Sprintf(format, args...))
}

// This attribute specifies the script level of the element
// Possible values are an integer between 0 and 7, inclusive.
func (e *MathMLMUNDERElement) SCRIPTLEVEL(i int) *MathMLMUNDERElement {
	if e.IntAttributes == nil {
		e.IntAttributes = treemap.New[string, int]()
	}
	e.IntAttributes.Set("scriptlevel", i)
	return e
}

func (e *MathMLMUNDERElement) IfSCRIPTLEVEL(condition bool, i int) *MathMLMUNDERElement {
	if condition {
		e.SCRIPTLEVEL(i)
	}
	return e
}

// Remove the attribute SCRIPTLEVEL from the element.
func (e *MathMLMUNDERElement) SCRIPTLEVELRemove(i int) *MathMLMUNDERElement {
	if e.IntAttributes == nil {
		return e
	}
	e.IntAttributes.Del("scriptlevel")
	return e
}

// This attribute offers advisory information about the element for which it is
// set.
func (e *MathMLMUNDERElement) STYLEF(k string, format string, args ...any) *MathMLMUNDERElement {
	return e.STYLE(k, fmt.Sprintf(format, args...))
}

func (e *MathMLMUNDERElement) IfSTYLE(condition bool, k string, v string) *MathMLMUNDERElement {
	if condition {
		e.STYLE(k, v)
	}
	return e
}

func (e *MathMLMUNDERElement) STYLE(k string, v string) *MathMLMUNDERElement {
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

func (e *MathMLMUNDERElement) IfSTYLEF(condition bool, k string, format string, args ...any) *MathMLMUNDERElement {
	if condition {
		e.STYLE(k, fmt.Sprintf(format, args...))
	}
	return e
}

// Add the attributes in the map to the element.
func (e *MathMLMUNDERElement) STYLEMap(m map[string]string) *MathMLMUNDERElement {
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
func (e *MathMLMUNDERElement) STYLEPairs(pairs ...string) *MathMLMUNDERElement {
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

func (e *MathMLMUNDERElement) IfSTYLEPairs(condition bool, pairs ...string) *MathMLMUNDERElement {
	if condition {
		e.STYLEPairs(pairs...)
	}
	return e
}

// Remove the attribute STYLE from the element.
func (e *MathMLMUNDERElement) STYLERemove(keys ...string) *MathMLMUNDERElement {
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
func (e *MathMLMUNDERElement) TABINDEX(i int) *MathMLMUNDERElement {
	if e.IntAttributes == nil {
		e.IntAttributes = treemap.New[string, int]()
	}
	e.IntAttributes.Set("tabindex", i)
	return e
}

func (e *MathMLMUNDERElement) IfTABINDEX(condition bool, i int) *MathMLMUNDERElement {
	if condition {
		e.TABINDEX(i)
	}
	return e
}

// Remove the attribute TABINDEX from the element.
func (e *MathMLMUNDERElement) TABINDEXRemove(i int) *MathMLMUNDERElement {
	if e.IntAttributes == nil {
		return e
	}
	e.IntAttributes.Del("tabindex")
	return e
}

// Merges the singleton store with the given object

func (e *MathMLMUNDERElement) DATASTAR_STORE(v any) *MathMLMUNDERElement {
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

func (e *MathMLMUNDERElement) DATASTAR_REF(expression string) *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key := "data-ref"

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *MathMLMUNDERElement) IfDATASTAR_REF(condition bool, expression string) *MathMLMUNDERElement {
	if condition {
		e.DATASTAR_REF(expression)
	}
	return e
}

// Remove the attribute DATASTAR_REF from the element.
func (e *MathMLMUNDERElement) DATASTAR_REFRemove() *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-ref")
	return e
}

// Sets the value of the element

func (e *MathMLMUNDERElement) DATASTAR_BIND(key string, expression string) *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key = fmt.Sprintf("data-bind-%s", key)

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *MathMLMUNDERElement) IfDATASTAR_BIND(condition bool, key string, expression string) *MathMLMUNDERElement {
	if condition {
		e.DATASTAR_BIND(key, expression)
	}
	return e
}

// Remove the attribute DATASTAR_BIND from the element.
func (e *MathMLMUNDERElement) DATASTAR_BINDRemove() *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-bind")
	return e
}

// Sets the value of the element

func (e *MathMLMUNDERElement) DATASTAR_MODEL(expression string) *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key := "data-model"

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *MathMLMUNDERElement) IfDATASTAR_MODEL(condition bool, expression string) *MathMLMUNDERElement {
	if condition {
		e.DATASTAR_MODEL(expression)
	}
	return e
}

// Remove the attribute DATASTAR_MODEL from the element.
func (e *MathMLMUNDERElement) DATASTAR_MODELRemove() *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-model")
	return e
}

// Sets the textContent of the element

func (e *MathMLMUNDERElement) DATASTAR_TEXT(expression string) *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key := "data-text"

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *MathMLMUNDERElement) IfDATASTAR_TEXT(condition bool, expression string) *MathMLMUNDERElement {
	if condition {
		e.DATASTAR_TEXT(expression)
	}
	return e
}

// Remove the attribute DATASTAR_TEXT from the element.
func (e *MathMLMUNDERElement) DATASTAR_TEXTRemove() *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-text")
	return e
}

// Sets the event handler of the element

type MathMLMunderOnMod customDataKeyModifier

// Debounces the event handler
func MathMLMunderOnModDebounce(
	d time.Duration,
) MathMLMunderOnMod {
	return func() string {
		return fmt.Sprintf("debounce_%dms", d.Milliseconds())
	}
}

// Throttles the event handler
func MathMLMunderOnModThrottle(
	d time.Duration,
) MathMLMunderOnMod {
	return func() string {
		return fmt.Sprintf("throttle_%dms", d.Milliseconds())
	}
}

func (e *MathMLMUNDERElement) DATASTAR_ON(key string, expression string, modifiers ...MathMLMunderOnMod) *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key = fmt.Sprintf("data-on-%s", key)

	customMods := lo.Map(modifiers, func(m MathMLMunderOnMod, i int) customDataKeyModifier {
		return customDataKeyModifier(m)
	})
	key = customDataKey(key, customMods...)
	e.StringAttributes.Set(key, expression)
	return e
}

func (e *MathMLMUNDERElement) IfDATASTAR_ON(condition bool, key string, expression string, modifiers ...MathMLMunderOnMod) *MathMLMUNDERElement {
	if condition {
		e.DATASTAR_ON(key, expression, modifiers...)
	}
	return e
}

// Remove the attribute DATASTAR_ON from the element.
func (e *MathMLMUNDERElement) DATASTAR_ONRemove() *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-on")
	return e
}

// Sets the focus of the element

func (e *MathMLMUNDERElement) DATASTAR_FOCUSSet(b bool) *MathMLMUNDERElement {
	key := "data-focus"
	e.BoolAttributes.Set(key, b)
	return e
}

func (e *MathMLMUNDERElement) DATASTAR_FOCUS() *MathMLMUNDERElement {
	return e.DATASTAR_FOCUSSet(true)
}

// Sets the header of for fetch requests

func (e *MathMLMUNDERElement) DATASTAR_HEADER(key string, expression string) *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key = fmt.Sprintf("data-header-%s", key)

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *MathMLMUNDERElement) IfDATASTAR_HEADER(condition bool, key string, expression string) *MathMLMUNDERElement {
	if condition {
		e.DATASTAR_HEADER(key, expression)
	}
	return e
}

// Remove the attribute DATASTAR_HEADER from the element.
func (e *MathMLMUNDERElement) DATASTAR_HEADERRemove() *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-header")
	return e
}

// Sets the indicator selector for fetch requests

func (e *MathMLMUNDERElement) DATASTAR_FETCH_INDICATOR(expression string) *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key := "data-fetch-indicator"

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *MathMLMUNDERElement) IfDATASTAR_FETCH_INDICATOR(condition bool, expression string) *MathMLMUNDERElement {
	if condition {
		e.DATASTAR_FETCH_INDICATOR(expression)
	}
	return e
}

// Remove the attribute DATASTAR_FETCH_INDICATOR from the element.
func (e *MathMLMUNDERElement) DATASTAR_FETCH_INDICATORRemove() *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-fetch-indicator")
	return e
}

// Sets the visibility of the element

func (e *MathMLMUNDERElement) DATASTAR_SHOWSet(b bool) *MathMLMUNDERElement {
	key := "data-show"
	e.BoolAttributes.Set(key, b)
	return e
}

func (e *MathMLMUNDERElement) DATASTAR_SHOW() *MathMLMUNDERElement {
	return e.DATASTAR_SHOWSet(true)
}

// Triggers the callback when the element intersects the viewport

func (e *MathMLMUNDERElement) DATASTAR_INTERSECTS(expression string) *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key := "data-intersects"

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *MathMLMUNDERElement) IfDATASTAR_INTERSECTS(condition bool, expression string) *MathMLMUNDERElement {
	if condition {
		e.DATASTAR_INTERSECTS(expression)
	}
	return e
}

// Remove the attribute DATASTAR_INTERSECTS from the element.
func (e *MathMLMUNDERElement) DATASTAR_INTERSECTSRemove() *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-intersects")
	return e
}

// Teleports the element to the given selector

func (e *MathMLMUNDERElement) DATASTAR_TELEPORTSet(b bool) *MathMLMUNDERElement {
	key := "data-teleport"
	e.BoolAttributes.Set(key, b)
	return e
}

func (e *MathMLMUNDERElement) DATASTAR_TELEPORT() *MathMLMUNDERElement {
	return e.DATASTAR_TELEPORTSet(true)
}

// Scrolls the element into view

func (e *MathMLMUNDERElement) DATASTAR_SCROLL_INTO_VIEWSet(b bool) *MathMLMUNDERElement {
	key := "data-scroll-into-view"
	e.BoolAttributes.Set(key, b)
	return e
}

func (e *MathMLMUNDERElement) DATASTAR_SCROLL_INTO_VIEW() *MathMLMUNDERElement {
	return e.DATASTAR_SCROLL_INTO_VIEWSet(true)
}

// Setup the ViewTransitionAPI for the element

func (e *MathMLMUNDERElement) DATASTAR_VIEW_TRANSITION(expression string) *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		e.StringAttributes = treemap.New[string, string]()
	}

	key := "data-view-transition"

	e.StringAttributes.Set(key, expression)
	return e
}

func (e *MathMLMUNDERElement) IfDATASTAR_VIEW_TRANSITION(condition bool, expression string) *MathMLMUNDERElement {
	if condition {
		e.DATASTAR_VIEW_TRANSITION(expression)
	}
	return e
}

// Remove the attribute DATASTAR_VIEW_TRANSITION from the element.
func (e *MathMLMUNDERElement) DATASTAR_VIEW_TRANSITIONRemove() *MathMLMUNDERElement {
	if e.StringAttributes == nil {
		return e
	}
	e.StringAttributes.Del("data-view-transition")
	return e
}
