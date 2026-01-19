package compile

import (
	"fmt"
	"sort"
	"strings"

	"lcss/internal/config"
	"lcss/internal/extract"
)

type Decl struct {
	Property string
	Value    string
}

type Rule struct {
	Selector string
	Decls    []Decl
	Media    string
}

type variantConfig struct {
	separator   string
	classPrefix string
	responsive  map[string]string
	state       map[string]struct{}
}

func buildUtilities(canonical config.Canonical, classes []string) ([]Rule, []string, []string) {
	variants := buildVariantConfig(canonical.Config)

	rules := make([]Rule, 0, len(classes))
	matched := make([]string, 0, len(classes))
	unknown := make([]string, 0)

	for _, class := range classes {
		rule, ok := matchClass(canonical, variants, class)
		if !ok {
			unknown = append(unknown, class)
			continue
		}
		rules = append(rules, rule)
		matched = append(matched, class)
	}

	return rules, matched, unknown
}

func buildVariantConfig(cfg config.Config) variantConfig {
	responsive := make(map[string]string, len(cfg.Variants.Responsive))
	for _, name := range cfg.Variants.Responsive {
		if value, ok := cfg.Breakpoints[name]; ok {
			responsive[name] = value
		}
	}
	state := make(map[string]struct{}, len(cfg.Variants.State))
	for _, name := range cfg.Variants.State {
		state[name] = struct{}{}
	}

	separator := cfg.Separator
	if separator == "" {
		separator = ":"
	}

	return variantConfig{
		separator:   separator,
		classPrefix: cfg.ClassPrefix,
		responsive:  responsive,
		state:       state,
	}
}

func matchClass(canonical config.Canonical, variants variantConfig, class string) (Rule, bool) {
	parsed, ok := parseClass(variants, class)
	if !ok {
		return Rule{}, false
	}

	decls, ok := matchUtility(parsed.Base, canonical)
	if !ok {
		return Rule{}, false
	}

	selector := "." + escapeClass(class)
	if len(parsed.Pseudos) > 0 {
		selector += strings.Join(parsed.Pseudos, "")
	}

	return Rule{
		Selector: selector,
		Decls:    decls,
		Media:    parsed.Media,
	}, true
}

type parsedClass struct {
	Base    string
	Media   string
	Pseudos []string
}

func parseClass(variants variantConfig, class string) (parsedClass, bool) {
	parts := strings.Split(class, variants.separator)
	if len(parts) == 0 {
		return parsedClass{}, false
	}

	base := parts[len(parts)-1]
	if variants.classPrefix != "" {
		if !strings.HasPrefix(base, variants.classPrefix) {
			return parsedClass{}, false
		}
		base = strings.TrimPrefix(base, variants.classPrefix)
		if base == "" {
			return parsedClass{}, false
		}
	}

	media := ""
	pseudos := make([]string, 0, len(parts)-1)

	for _, variant := range parts[:len(parts)-1] {
		if variant == "" {
			return parsedClass{}, false
		}
		if breakpoint, ok := variants.responsive[variant]; ok {
			if media != "" {
				return parsedClass{}, false
			}
			media = fmt.Sprintf("(min-width: %s)", breakpoint)
			continue
		}
		if _, ok := variants.state[variant]; ok {
			pseudos = append(pseudos, ":"+variant)
			continue
		}
		return parsedClass{}, false
	}

	return parsedClass{
		Base:    base,
		Media:   media,
		Pseudos: pseudos,
	}, true
}

func matchUtility(base string, canonical config.Canonical) ([]Decl, bool) {
	space := canonical.Tokens.Scales["space"]
	size := canonical.Tokens.Scales["size"]
	maxWidth := canonical.Tokens.Scales["maxWidth"]
	maxHeight := canonical.Tokens.Scales["maxHeight"]
	colors := canonical.Tokens.Themes["default"].Colors
	fontSize := canonical.Tokens.Scales["fontSize"]
	lineHeight := canonical.Tokens.Scales["lineHeight"]
	fontWeight := canonical.Tokens.Scales["fontWeight"]
	letterSpacing := canonical.Tokens.Scales["letterSpacing"]
	fonts := canonical.Tokens.Themes["default"].Fonts
	radius := canonical.Tokens.Scales["radius"]
	borderWidth := canonical.Tokens.Scales["borderWidth"]
	shadow := canonical.Tokens.Scales["shadow"]
	zIndex := canonical.Tokens.Scales["z"]
	opacity := canonical.Tokens.Scales["opacity"]
	aspect := canonical.Tokens.Scales["aspect"]
	duration := canonical.Tokens.Scales["duration"]
	easing := canonical.Tokens.Scales["easing"]
	delay := canonical.Tokens.Scales["delay"]
	translate := canonical.Tokens.Scales["translate"]
	rotate := canonical.Tokens.Scales["rotate"]
	scale := canonical.Tokens.Scales["scale"]
	container := canonical.Tokens.Scales["container"]

	if decls, ok := matchSpacing(base, space); ok {
		return decls, true
	}
	if decls, ok := matchSizing(base, space, size, maxWidth, maxHeight, container); ok {
		return decls, true
	}
	if decls, ok := matchDisplay(base); ok {
		return decls, true
	}
	if decls, ok := matchPosition(base, space, size); ok {
		return decls, true
	}
	if decls, ok := matchFlex(base); ok {
		return decls, true
	}
	if decls, ok := matchGrid(base); ok {
		return decls, true
	}
	if decls, ok := matchTypography(base, fonts, fontSize, lineHeight, fontWeight, colors); ok {
		return decls, true
	}
	if decls, ok := matchTypographyExtras(base, letterSpacing); ok {
		return decls, true
	}
	if decls, ok := matchColors(base, colors); ok {
		return decls, true
	}
	if decls, ok := matchBackground(base); ok {
		return decls, true
	}
	if decls, ok := matchBorders(base, colors, borderWidth); ok {
		return decls, true
	}
	if decls, ok := matchRadius(base, radius); ok {
		return decls, true
	}
	if decls, ok := matchShadow(base, shadow); ok {
		return decls, true
	}
	if decls, ok := matchOpacity(base, opacity); ok {
		return decls, true
	}
	if decls, ok := matchZIndex(base, zIndex); ok {
		return decls, true
	}
	if decls, ok := matchOverflow(base); ok {
		return decls, true
	}
	if decls, ok := matchVisibility(base); ok {
		return decls, true
	}
	if decls, ok := matchObject(base); ok {
		return decls, true
	}
	if decls, ok := matchAspect(base, aspect); ok {
		return decls, true
	}
	if decls, ok := matchTransition(base, duration, easing, delay); ok {
		return decls, true
	}
	if decls, ok := matchTransform(base, translate, rotate, scale, space); ok {
		return decls, true
	}
	if decls, ok := matchInteraction(base); ok {
		return decls, true
	}

	return nil, false
}

func matchSpacing(base string, space map[string]string) ([]Decl, bool) {
	key, props := parseSpacing(base)
	if key == "" || len(props) == 0 {
		return nil, false
	}
	if _, ok := space[key]; !ok {
		return nil, false
	}
	value := fmt.Sprintf("var(--space-%s)", key)
	decls := make([]Decl, 0, len(props))
	for _, prop := range props {
		decls = append(decls, Decl{Property: prop, Value: value})
	}
	return decls, true
}

func parseSpacing(base string) (string, []string) {
	switch {
	case strings.HasPrefix(base, "px-"):
		return base[3:], []string{"padding-left", "padding-right"}
	case strings.HasPrefix(base, "py-"):
		return base[3:], []string{"padding-top", "padding-bottom"}
	case strings.HasPrefix(base, "pt-"):
		return base[3:], []string{"padding-top"}
	case strings.HasPrefix(base, "pr-"):
		return base[3:], []string{"padding-right"}
	case strings.HasPrefix(base, "pb-"):
		return base[3:], []string{"padding-bottom"}
	case strings.HasPrefix(base, "pl-"):
		return base[3:], []string{"padding-left"}
	case strings.HasPrefix(base, "p-"):
		return base[2:], []string{"padding"}
	case strings.HasPrefix(base, "mx-"):
		return base[3:], []string{"margin-left", "margin-right"}
	case strings.HasPrefix(base, "my-"):
		return base[3:], []string{"margin-top", "margin-bottom"}
	case strings.HasPrefix(base, "mt-"):
		return base[3:], []string{"margin-top"}
	case strings.HasPrefix(base, "mr-"):
		return base[3:], []string{"margin-right"}
	case strings.HasPrefix(base, "mb-"):
		return base[3:], []string{"margin-bottom"}
	case strings.HasPrefix(base, "ml-"):
		return base[3:], []string{"margin-left"}
	case strings.HasPrefix(base, "m-"):
		return base[2:], []string{"margin"}
	case strings.HasPrefix(base, "gap-x-"):
		return base[6:], []string{"column-gap"}
	case strings.HasPrefix(base, "gap-y-"):
		return base[6:], []string{"row-gap"}
	case strings.HasPrefix(base, "gapx-"):
		return base[5:], []string{"column-gap"}
	case strings.HasPrefix(base, "gapy-"):
		return base[5:], []string{"row-gap"}
	case strings.HasPrefix(base, "gap-"):
		return base[4:], []string{"gap"}
	default:
		return "", nil
	}
}

func matchSizing(base string, space, size, maxWidth, maxHeight, container map[string]string) ([]Decl, bool) {
	if base == "container" {
		key := defaultKey(container, "default", "lg", "xl", "md", "sm")
		if key == "" {
			return nil, false
		}
		return []Decl{
			{Property: "width", Value: "100%"},
			{Property: "margin-left", Value: "auto"},
			{Property: "margin-right", Value: "auto"},
			{Property: "max-width", Value: fmt.Sprintf("var(--container-%s)", key)},
		}, true
	}
	if strings.HasPrefix(base, "w-") {
		key := base[2:]
		if value, ok := sizeValue(key, size, space, true); ok {
			return []Decl{{Property: "width", Value: value}}, true
		}
	}
	if strings.HasPrefix(base, "h-") {
		key := base[2:]
		if value, ok := sizeValue(key, size, space, false); ok {
			return []Decl{{Property: "height", Value: value}}, true
		}
	}
	if strings.HasPrefix(base, "min-w-") {
		key := base[len("min-w-"):]
		if value, ok := sizeValue(key, size, space, true); ok {
			return []Decl{{Property: "min-width", Value: value}}, true
		}
	}
	if strings.HasPrefix(base, "min-h-") {
		key := base[len("min-h-"):]
		if value, ok := sizeValue(key, size, space, false); ok {
			return []Decl{{Property: "min-height", Value: value}}, true
		}
	}
	if strings.HasPrefix(base, "max-w-") {
		key := base[len("max-w-"):]
		if value, ok := sizeValueFromScale(key, maxWidth, size, space, true); ok {
			return []Decl{{Property: "max-width", Value: value}}, true
		}
	}
	if strings.HasPrefix(base, "max-h-") {
		key := base[len("max-h-"):]
		if value, ok := sizeValueFromScale(key, maxHeight, size, space, false); ok {
			return []Decl{{Property: "max-height", Value: value}}, true
		}
	}
	return nil, false
}

func sizeValue(key string, size, space map[string]string, isWidth bool) (string, bool) {
	switch key {
	case "auto":
		return "auto", true
	case "full":
		return "100%", true
	case "screen":
		if isWidth {
			return "100vw", true
		}
		return "100vh", true
	case "min":
		return "min-content", true
	case "max":
		return "max-content", true
	case "fit":
		return "fit-content", true
	default:
		if _, ok := size[key]; ok {
			return fmt.Sprintf("var(--size-%s)", key), true
		}
		if _, ok := space[key]; ok {
			return fmt.Sprintf("var(--space-%s)", key), true
		}
	}
	return "", false
}

func sizeValueFromScale(key string, primary, size, space map[string]string, isWidth bool) (string, bool) {
	if _, ok := primary[key]; ok {
		if isWidth {
			return fmt.Sprintf("var(--max-width-%s)", key), true
		}
		return fmt.Sprintf("var(--max-height-%s)", key), true
	}
	return sizeValue(key, size, space, isWidth)
}

func matchDisplay(base string) ([]Decl, bool) {
	values := map[string]string{
		"block":        "block",
		"inline-block": "inline-block",
		"inline":       "inline",
		"flex":         "flex",
		"inline-flex":  "inline-flex",
		"grid":         "grid",
		"hidden":       "none",
		"contents":     "contents",
	}
	if value, ok := values[base]; ok {
		return []Decl{{Property: "display", Value: value}}, true
	}
	return nil, false
}

func matchFlex(base string) ([]Decl, bool) {
	if base == "flex-row" {
		return []Decl{{Property: "flex-direction", Value: "row"}}, true
	}
	if base == "flex-col" {
		return []Decl{{Property: "flex-direction", Value: "column"}}, true
	}
	switch base {
	case "flex-wrap", "flex-nowrap", "flex-wrap-reverse":
		return []Decl{{Property: "flex-wrap", Value: strings.TrimPrefix(base, "flex-")}}, true
	case "flex-1":
		return []Decl{{Property: "flex", Value: "1 1 0%"}}, true
	case "flex-auto":
		return []Decl{{Property: "flex", Value: "1 1 auto"}}, true
	case "flex-initial":
		return []Decl{{Property: "flex", Value: "0 1 auto"}}, true
	case "flex-none":
		return []Decl{{Property: "flex", Value: "none"}}, true
	case "grow":
		return []Decl{{Property: "flex-grow", Value: "1"}}, true
	case "grow-0":
		return []Decl{{Property: "flex-grow", Value: "0"}}, true
	case "shrink":
		return []Decl{{Property: "flex-shrink", Value: "1"}}, true
	case "shrink-0":
		return []Decl{{Property: "flex-shrink", Value: "0"}}, true
	}
	if strings.HasPrefix(base, "items-") {
		value := strings.TrimPrefix(base, "items-")
		if mapped, ok := mapAlign(value); ok {
			return []Decl{{Property: "align-items", Value: mapped}}, true
		}
	}
	if strings.HasPrefix(base, "justify-") {
		value := strings.TrimPrefix(base, "justify-")
		if mapped, ok := mapJustify(value); ok {
			return []Decl{{Property: "justify-content", Value: mapped}}, true
		}
	}
	if strings.HasPrefix(base, "content-") {
		value := strings.TrimPrefix(base, "content-")
		if mapped, ok := mapJustify(value); ok {
			return []Decl{{Property: "align-content", Value: mapped}}, true
		}
	}
	if strings.HasPrefix(base, "self-") {
		value := strings.TrimPrefix(base, "self-")
		if mapped, ok := mapAlign(value); ok {
			return []Decl{{Property: "align-self", Value: mapped}}, true
		}
	}
	return nil, false
}

func mapAlign(value string) (string, bool) {
	switch value {
	case "start":
		return "flex-start", true
	case "center":
		return "center", true
	case "end":
		return "flex-end", true
	case "stretch":
		return "stretch", true
	case "baseline":
		return "baseline", true
	default:
		return "", false
	}
}

func mapJustify(value string) (string, bool) {
	switch value {
	case "start":
		return "flex-start", true
	case "center":
		return "center", true
	case "end":
		return "flex-end", true
	case "between":
		return "space-between", true
	case "around":
		return "space-around", true
	case "evenly":
		return "space-evenly", true
	default:
		return "", false
	}
}

func matchTypography(base string, fonts, fontSize, lineHeight, fontWeight, colors map[string]string) ([]Decl, bool) {
	if strings.HasPrefix(base, "text-") {
		key := strings.TrimPrefix(base, "text-")
		if _, ok := fontSize[key]; ok {
			return []Decl{{Property: "font-size", Value: fmt.Sprintf("var(--font-size-%s)", key)}}, true
		}
		if _, ok := colors[key]; ok {
			return []Decl{{Property: "color", Value: fmt.Sprintf("var(--color-%s)", key)}}, true
		}
		switch key {
		case "left", "center", "right", "justify":
			return []Decl{{Property: "text-align", Value: key}}, true
		}
	}
	if strings.HasPrefix(base, "leading-") {
		key := strings.TrimPrefix(base, "leading-")
		if _, ok := lineHeight[key]; ok {
			return []Decl{{Property: "line-height", Value: fmt.Sprintf("var(--line-height-%s)", key)}}, true
		}
	}
	if strings.HasPrefix(base, "font-") {
		key := strings.TrimPrefix(base, "font-")
		if _, ok := fonts[key]; ok {
			return []Decl{{Property: "font-family", Value: fmt.Sprintf("var(--font-%s)", key)}}, true
		}
		if _, ok := fontWeight[key]; ok {
			return []Decl{{Property: "font-weight", Value: fmt.Sprintf("var(--font-weight-%s)", key)}}, true
		}
	}
	switch base {
	case "italic":
		return []Decl{{Property: "font-style", Value: "italic"}}, true
	case "not-italic":
		return []Decl{{Property: "font-style", Value: "normal"}}, true
	case "uppercase", "lowercase", "capitalize":
		return []Decl{{Property: "text-transform", Value: base}}, true
	case "normal-case":
		return []Decl{{Property: "text-transform", Value: "none"}}, true
	case "underline":
		return []Decl{{Property: "text-decoration", Value: "underline"}}, true
	case "line-through":
		return []Decl{{Property: "text-decoration", Value: "line-through"}}, true
	case "no-underline":
		return []Decl{{Property: "text-decoration", Value: "none"}}, true
	case "list-none":
		return []Decl{{Property: "list-style", Value: "none"}}, true
	case "list-disc":
		return []Decl{{Property: "list-style", Value: "disc"}}, true
	case "list-decimal":
		return []Decl{{Property: "list-style", Value: "decimal"}}, true
	}
	return nil, false
}

func matchTypographyExtras(base string, letterSpacing map[string]string) ([]Decl, bool) {
	if strings.HasPrefix(base, "tracking-") {
		key := strings.TrimPrefix(base, "tracking-")
		if _, ok := letterSpacing[key]; ok {
			return []Decl{{Property: "letter-spacing", Value: fmt.Sprintf("var(--letter-spacing-%s)", key)}}, true
		}
	}
	return nil, false
}

func matchPosition(base string, space, size map[string]string) ([]Decl, bool) {
	switch base {
	case "static", "relative", "absolute", "fixed", "sticky":
		return []Decl{{Property: "position", Value: base}}, true
	}
	if key, props := parseInset(base); key != "" {
		value, ok := sizeValue(key, size, space, true)
		if !ok {
			return nil, false
		}
		decls := make([]Decl, 0, len(props))
		for _, prop := range props {
			decls = append(decls, Decl{Property: prop, Value: value})
		}
		return decls, true
	}
	return nil, false
}

func parseInset(base string) (string, []string) {
	switch {
	case strings.HasPrefix(base, "inset-x-"):
		return base[len("inset-x-"):], []string{"left", "right"}
	case strings.HasPrefix(base, "inset-y-"):
		return base[len("inset-y-"):], []string{"top", "bottom"}
	case strings.HasPrefix(base, "inset-"):
		return base[len("inset-"):], []string{"top", "right", "bottom", "left"}
	case strings.HasPrefix(base, "top-"):
		return base[len("top-"):], []string{"top"}
	case strings.HasPrefix(base, "right-"):
		return base[len("right-"):], []string{"right"}
	case strings.HasPrefix(base, "bottom-"):
		return base[len("bottom-"):], []string{"bottom"}
	case strings.HasPrefix(base, "left-"):
		return base[len("left-"):], []string{"left"}
	default:
		return "", nil
	}
}

func matchGrid(base string) ([]Decl, bool) {
	if strings.HasPrefix(base, "grid-cols-") {
		key := strings.TrimPrefix(base, "grid-cols-")
		if key == "none" {
			return []Decl{{Property: "grid-template-columns", Value: "none"}}, true
		}
		if value, ok := parsePositiveInt(key); ok {
			return []Decl{{Property: "grid-template-columns", Value: fmt.Sprintf("repeat(%d, minmax(0, 1fr))", value)}}, true
		}
	}
	if strings.HasPrefix(base, "grid-rows-") {
		key := strings.TrimPrefix(base, "grid-rows-")
		if key == "none" {
			return []Decl{{Property: "grid-template-rows", Value: "none"}}, true
		}
		if value, ok := parsePositiveInt(key); ok {
			return []Decl{{Property: "grid-template-rows", Value: fmt.Sprintf("repeat(%d, minmax(0, 1fr))", value)}}, true
		}
	}
	if strings.HasPrefix(base, "col-span-") {
		key := strings.TrimPrefix(base, "col-span-")
		if key == "full" {
			return []Decl{{Property: "grid-column", Value: "1 / -1"}}, true
		}
		if value, ok := parsePositiveInt(key); ok {
			return []Decl{{Property: "grid-column", Value: fmt.Sprintf("span %d / span %d", value, value)}}, true
		}
	}
	if strings.HasPrefix(base, "row-span-") {
		key := strings.TrimPrefix(base, "row-span-")
		if key == "full" {
			return []Decl{{Property: "grid-row", Value: "1 / -1"}}, true
		}
		if value, ok := parsePositiveInt(key); ok {
			return []Decl{{Property: "grid-row", Value: fmt.Sprintf("span %d / span %d", value, value)}}, true
		}
	}
	if strings.HasPrefix(base, "col-start-") {
		key := strings.TrimPrefix(base, "col-start-")
		if value, ok := parsePositiveInt(key); ok {
			return []Decl{{Property: "grid-column-start", Value: fmt.Sprintf("%d", value)}}, true
		}
	}
	if strings.HasPrefix(base, "col-end-") {
		key := strings.TrimPrefix(base, "col-end-")
		if value, ok := parsePositiveInt(key); ok {
			return []Decl{{Property: "grid-column-end", Value: fmt.Sprintf("%d", value)}}, true
		}
	}
	if strings.HasPrefix(base, "row-start-") {
		key := strings.TrimPrefix(base, "row-start-")
		if value, ok := parsePositiveInt(key); ok {
			return []Decl{{Property: "grid-row-start", Value: fmt.Sprintf("%d", value)}}, true
		}
	}
	if strings.HasPrefix(base, "row-end-") {
		key := strings.TrimPrefix(base, "row-end-")
		if value, ok := parsePositiveInt(key); ok {
			return []Decl{{Property: "grid-row-end", Value: fmt.Sprintf("%d", value)}}, true
		}
	}
	if strings.HasPrefix(base, "place-items-") {
		value := strings.TrimPrefix(base, "place-items-")
		if mapped, ok := mapAlign(value); ok {
			return []Decl{{Property: "place-items", Value: mapped}}, true
		}
	}
	if strings.HasPrefix(base, "place-content-") {
		value := strings.TrimPrefix(base, "place-content-")
		if mapped, ok := mapJustify(value); ok {
			return []Decl{{Property: "place-content", Value: mapped}}, true
		}
	}
	return nil, false
}

func matchBackground(base string) ([]Decl, bool) {
	switch base {
	case "bg-cover", "bg-contain":
		return []Decl{{Property: "background-size", Value: strings.TrimPrefix(base, "bg-")}}, true
	case "bg-center", "bg-top", "bg-right", "bg-bottom", "bg-left":
		return []Decl{{Property: "background-position", Value: strings.TrimPrefix(base, "bg-")}}, true
	case "bg-fixed", "bg-local", "bg-scroll":
		return []Decl{{Property: "background-attachment", Value: strings.TrimPrefix(base, "bg-")}}, true
	case "bg-repeat", "bg-no-repeat", "bg-repeat-x", "bg-repeat-y":
		value := strings.TrimPrefix(base, "bg-")
		if value == "no-repeat" {
			return []Decl{{Property: "background-repeat", Value: "no-repeat"}}, true
		}
		if value == "repeat-x" {
			return []Decl{{Property: "background-repeat", Value: "repeat-x"}}, true
		}
		if value == "repeat-y" {
			return []Decl{{Property: "background-repeat", Value: "repeat-y"}}, true
		}
		return []Decl{{Property: "background-repeat", Value: "repeat"}}, true
	}
	return nil, false
}

func matchColors(base string, colors map[string]string) ([]Decl, bool) {
	if strings.HasPrefix(base, "bg-") {
		key := strings.TrimPrefix(base, "bg-")
		if _, ok := colors[key]; ok {
			return []Decl{{Property: "background-color", Value: fmt.Sprintf("var(--color-%s)", key)}}, true
		}
	}
	if strings.HasPrefix(base, "text-") {
		key := strings.TrimPrefix(base, "text-")
		if _, ok := colors[key]; ok {
			return []Decl{{Property: "color", Value: fmt.Sprintf("var(--color-%s)", key)}}, true
		}
	}
	if strings.HasPrefix(base, "border-") {
		key := strings.TrimPrefix(base, "border-")
		if _, ok := colors[key]; ok {
			return []Decl{{Property: "border-color", Value: fmt.Sprintf("var(--color-%s)", key)}}, true
		}
	}
	return nil, false
}

func matchBorders(base string, colors, borderWidth map[string]string) ([]Decl, bool) {
	if base == "border" {
		return []Decl{
			{Property: "border-width", Value: "1px"},
			{Property: "border-style", Value: "solid"},
		}, true
	}
	switch base {
	case "border-solid", "border-dashed", "border-dotted", "border-double", "border-none":
		value := strings.TrimPrefix(base, "border-")
		return []Decl{{Property: "border-style", Value: value}}, true
	}
	if strings.HasPrefix(base, "border-") {
		key := strings.TrimPrefix(base, "border-")
		if _, ok := colors[key]; ok {
			return nil, false
		}
		if _, ok := borderWidth[key]; ok {
			return []Decl{
				{Property: "border-width", Value: fmt.Sprintf("var(--border-width-%s)", key)},
				{Property: "border-style", Value: "solid"},
			}, true
		}
		switch key {
		case "0", "2", "4", "8":
			return []Decl{
				{Property: "border-width", Value: key + "px"},
				{Property: "border-style", Value: "solid"},
			}, true
		}
	}
	if strings.HasPrefix(base, "border-x-") {
		key := strings.TrimPrefix(base, "border-x-")
		if value, ok := borderWidthValue(key, borderWidth); ok {
			return []Decl{
				{Property: "border-left-width", Value: value},
				{Property: "border-right-width", Value: value},
				{Property: "border-style", Value: "solid"},
			}, true
		}
	}
	if strings.HasPrefix(base, "border-y-") {
		key := strings.TrimPrefix(base, "border-y-")
		if value, ok := borderWidthValue(key, borderWidth); ok {
			return []Decl{
				{Property: "border-top-width", Value: value},
				{Property: "border-bottom-width", Value: value},
				{Property: "border-style", Value: "solid"},
			}, true
		}
	}
	if strings.HasPrefix(base, "border-t-") {
		key := strings.TrimPrefix(base, "border-t-")
		if value, ok := borderWidthValue(key, borderWidth); ok {
			return []Decl{{Property: "border-top-width", Value: value}, {Property: "border-style", Value: "solid"}}, true
		}
	}
	if strings.HasPrefix(base, "border-r-") {
		key := strings.TrimPrefix(base, "border-r-")
		if value, ok := borderWidthValue(key, borderWidth); ok {
			return []Decl{{Property: "border-right-width", Value: value}, {Property: "border-style", Value: "solid"}}, true
		}
	}
	if strings.HasPrefix(base, "border-b-") {
		key := strings.TrimPrefix(base, "border-b-")
		if value, ok := borderWidthValue(key, borderWidth); ok {
			return []Decl{{Property: "border-bottom-width", Value: value}, {Property: "border-style", Value: "solid"}}, true
		}
	}
	if strings.HasPrefix(base, "border-l-") {
		key := strings.TrimPrefix(base, "border-l-")
		if value, ok := borderWidthValue(key, borderWidth); ok {
			return []Decl{{Property: "border-left-width", Value: value}, {Property: "border-style", Value: "solid"}}, true
		}
	}
	return nil, false
}

func borderWidthValue(key string, borderWidth map[string]string) (string, bool) {
	if _, ok := borderWidth[key]; ok {
		return fmt.Sprintf("var(--border-width-%s)", key), true
	}
	switch key {
	case "0", "2", "4", "8":
		return key + "px", true
	}
	return "", false
}

func matchRadius(base string, radius map[string]string) ([]Decl, bool) {
	if base == "rounded" {
		key := defaultRadiusKey(radius)
		if key == "" {
			return nil, false
		}
		return []Decl{{Property: "border-radius", Value: fmt.Sprintf("var(--radius-%s)", key)}}, true
	}
	if strings.HasPrefix(base, "rounded-") {
		key := strings.TrimPrefix(base, "rounded-")
		if _, ok := radius[key]; ok {
			return []Decl{{Property: "border-radius", Value: fmt.Sprintf("var(--radius-%s)", key)}}, true
		}
		switch key {
		case "t", "b", "l", "r", "tl", "tr", "bl", "br":
			corner := fmt.Sprintf("var(--radius-%s)", defaultRadiusKey(radius))
			if corner == "var(--radius-)" {
				return nil, false
			}
			return radiusCorners(key, corner), true
		}
	}
	return nil, false
}

func radiusCorners(key, value string) []Decl {
	switch key {
	case "t":
		return []Decl{{Property: "border-top-left-radius", Value: value}, {Property: "border-top-right-radius", Value: value}}
	case "b":
		return []Decl{{Property: "border-bottom-left-radius", Value: value}, {Property: "border-bottom-right-radius", Value: value}}
	case "l":
		return []Decl{{Property: "border-top-left-radius", Value: value}, {Property: "border-bottom-left-radius", Value: value}}
	case "r":
		return []Decl{{Property: "border-top-right-radius", Value: value}, {Property: "border-bottom-right-radius", Value: value}}
	case "tl":
		return []Decl{{Property: "border-top-left-radius", Value: value}}
	case "tr":
		return []Decl{{Property: "border-top-right-radius", Value: value}}
	case "bl":
		return []Decl{{Property: "border-bottom-left-radius", Value: value}}
	case "br":
		return []Decl{{Property: "border-bottom-right-radius", Value: value}}
	default:
		return nil
	}
}

func matchShadow(base string, shadow map[string]string) ([]Decl, bool) {
	if base == "shadow" {
		key := defaultKey(shadow, "default", "md", "sm", "lg", "xl")
		if key == "" {
			return nil, false
		}
		return []Decl{{Property: "box-shadow", Value: fmt.Sprintf("var(--shadow-%s)", key)}}, true
	}
	if strings.HasPrefix(base, "shadow-") {
		key := strings.TrimPrefix(base, "shadow-")
		if _, ok := shadow[key]; ok {
			return []Decl{{Property: "box-shadow", Value: fmt.Sprintf("var(--shadow-%s)", key)}}, true
		}
	}
	if base == "shadow-none" {
		return []Decl{{Property: "box-shadow", Value: "none"}}, true
	}
	return nil, false
}

func matchOpacity(base string, opacity map[string]string) ([]Decl, bool) {
	if strings.HasPrefix(base, "opacity-") {
		key := strings.TrimPrefix(base, "opacity-")
		if _, ok := opacity[key]; ok {
			return []Decl{{Property: "opacity", Value: fmt.Sprintf("var(--opacity-%s)", key)}}, true
		}
	}
	return nil, false
}

func matchZIndex(base string, zIndex map[string]string) ([]Decl, bool) {
	if base == "z-auto" {
		return []Decl{{Property: "z-index", Value: "auto"}}, true
	}
	if strings.HasPrefix(base, "z-") {
		key := strings.TrimPrefix(base, "z-")
		if _, ok := zIndex[key]; ok {
			return []Decl{{Property: "z-index", Value: fmt.Sprintf("var(--z-%s)", key)}}, true
		}
	}
	return nil, false
}

func matchOverflow(base string) ([]Decl, bool) {
	if strings.HasPrefix(base, "overflow-") {
		value := strings.TrimPrefix(base, "overflow-")
		switch value {
		case "auto", "hidden", "visible", "scroll":
			return []Decl{{Property: "overflow", Value: value}}, true
		case "x-auto", "x-hidden", "x-visible", "x-scroll":
			return []Decl{{Property: "overflow-x", Value: strings.TrimPrefix(value, "x-")}}, true
		case "y-auto", "y-hidden", "y-visible", "y-scroll":
			return []Decl{{Property: "overflow-y", Value: strings.TrimPrefix(value, "y-")}}, true
		}
	}
	return nil, false
}

func matchVisibility(base string) ([]Decl, bool) {
	switch base {
	case "visible":
		return []Decl{{Property: "visibility", Value: "visible"}}, true
	case "invisible":
		return []Decl{{Property: "visibility", Value: "hidden"}}, true
	case "sr-only":
		return []Decl{
			{Property: "position", Value: "absolute"},
			{Property: "width", Value: "1px"},
			{Property: "height", Value: "1px"},
			{Property: "padding", Value: "0"},
			{Property: "margin", Value: "-1px"},
			{Property: "overflow", Value: "hidden"},
			{Property: "clip", Value: "rect(0, 0, 0, 0)"},
			{Property: "white-space", Value: "nowrap"},
			{Property: "border", Value: "0"},
		}, true
	}
	return nil, false
}

func matchObject(base string) ([]Decl, bool) {
	if strings.HasPrefix(base, "object-") {
		value := strings.TrimPrefix(base, "object-")
		switch value {
		case "contain", "cover", "fill", "none", "scale-down":
			return []Decl{{Property: "object-fit", Value: value}}, true
		case "center", "top", "right", "bottom", "left":
			return []Decl{{Property: "object-position", Value: value}}, true
		}
	}
	return nil, false
}

func matchAspect(base string, aspect map[string]string) ([]Decl, bool) {
	if strings.HasPrefix(base, "aspect-") {
		key := strings.TrimPrefix(base, "aspect-")
		if _, ok := aspect[key]; ok {
			return []Decl{{Property: "aspect-ratio", Value: fmt.Sprintf("var(--aspect-%s)", key)}}, true
		}
	}
	return nil, false
}

func matchTransition(base string, duration, easing, delay map[string]string) ([]Decl, bool) {
	switch base {
	case "transition":
		return []Decl{{Property: "transition-property", Value: "all"}}, true
	case "transition-colors":
		return []Decl{{Property: "transition-property", Value: "color, background-color, border-color, fill, stroke"}}, true
	case "transition-opacity":
		return []Decl{{Property: "transition-property", Value: "opacity"}}, true
	case "transition-transform":
		return []Decl{{Property: "transition-property", Value: "transform"}}, true
	}
	if strings.HasPrefix(base, "duration-") {
		key := strings.TrimPrefix(base, "duration-")
		if _, ok := duration[key]; ok {
			return []Decl{{Property: "transition-duration", Value: fmt.Sprintf("var(--duration-%s)", key)}}, true
		}
	}
	if strings.HasPrefix(base, "ease-") {
		key := strings.TrimPrefix(base, "ease-")
		if _, ok := easing[key]; ok {
			return []Decl{{Property: "transition-timing-function", Value: fmt.Sprintf("var(--easing-%s)", key)}}, true
		}
	}
	if strings.HasPrefix(base, "delay-") {
		key := strings.TrimPrefix(base, "delay-")
		if _, ok := delay[key]; ok {
			return []Decl{{Property: "transition-delay", Value: fmt.Sprintf("var(--delay-%s)", key)}}, true
		}
	}
	return nil, false
}

func matchTransform(base string, translate, rotate, scale, space map[string]string) ([]Decl, bool) {
	if strings.HasPrefix(base, "translate-x-") {
		key := strings.TrimPrefix(base, "translate-x-")
		if value, ok := transformValue(key, translate, space); ok {
			return []Decl{{Property: "transform", Value: fmt.Sprintf("translateX(%s)", value)}}, true
		}
	}
	if strings.HasPrefix(base, "translate-y-") {
		key := strings.TrimPrefix(base, "translate-y-")
		if value, ok := transformValue(key, translate, space); ok {
			return []Decl{{Property: "transform", Value: fmt.Sprintf("translateY(%s)", value)}}, true
		}
	}
	if strings.HasPrefix(base, "rotate-") {
		key := strings.TrimPrefix(base, "rotate-")
		if _, ok := rotate[key]; ok {
			return []Decl{{Property: "transform", Value: fmt.Sprintf("rotate(%s)", fmt.Sprintf("var(--rotate-%s)", key))}}, true
		}
	}
	if strings.HasPrefix(base, "scale-") {
		key := strings.TrimPrefix(base, "scale-")
		if _, ok := scale[key]; ok {
			return []Decl{{Property: "transform", Value: fmt.Sprintf("scale(%s)", fmt.Sprintf("var(--scale-%s)", key))}}, true
		}
	}
	return nil, false
}

func transformValue(key string, translate, space map[string]string) (string, bool) {
	if _, ok := translate[key]; ok {
		return fmt.Sprintf("var(--translate-%s)", key), true
	}
	if _, ok := space[key]; ok {
		return fmt.Sprintf("var(--space-%s)", key), true
	}
	if key == "full" {
		return "100%", true
	}
	return "", false
}

func matchInteraction(base string) ([]Decl, bool) {
	switch base {
	case "cursor-pointer", "cursor-default", "cursor-text", "cursor-not-allowed":
		return []Decl{{Property: "cursor", Value: strings.TrimPrefix(base, "cursor-")}}, true
	case "pointer-events-none":
		return []Decl{{Property: "pointer-events", Value: "none"}}, true
	case "pointer-events-auto":
		return []Decl{{Property: "pointer-events", Value: "auto"}}, true
	case "select-none", "select-text", "select-all", "select-auto":
		return []Decl{{Property: "user-select", Value: strings.TrimPrefix(base, "select-")}}, true
	case "isolate":
		return []Decl{{Property: "isolation", Value: "isolate"}}, true
	case "isolation-auto":
		return []Decl{{Property: "isolation", Value: "auto"}}, true
	}
	return nil, false
}

func parsePositiveInt(value string) (int, bool) {
	if value == "" {
		return 0, false
	}
	for _, r := range value {
		if r < '0' || r > '9' {
			return 0, false
		}
	}
	parsed := 0
	for _, r := range value {
		parsed = parsed*10 + int(r-'0')
	}
	if parsed <= 0 {
		return 0, false
	}
	return parsed, true
}

func defaultRadiusKey(radius map[string]string) string {
	return defaultKey(radius, "default", "md", "base", "sm")
}

func defaultKey(values map[string]string, preferred ...string) string {
	for _, key := range preferred {
		if _, ok := values[key]; ok {
			return key
		}
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	if len(keys) == 0 {
		return ""
	}
	return keys[0]
}

func renderRules(rules []Rule) string {
	var b strings.Builder
	for _, rule := range rules {
		writeRule(&b, rule)
	}
	return b.String()
}

func writeRule(b *strings.Builder, rule Rule) {
	if rule.Media != "" {
		b.WriteString("@media ")
		b.WriteString(rule.Media)
		b.WriteString(" {\n")
		writeRuleBody(b, rule, "  ")
		b.WriteString("}\n")
		return
	}
	writeRuleBody(b, rule, "")
}

func writeRuleBody(b *strings.Builder, rule Rule, indent string) {
	b.WriteString(indent)
	b.WriteString(rule.Selector)
	b.WriteString(" {\n")
	for _, decl := range rule.Decls {
		b.WriteString(indent)
		b.WriteString("  ")
		b.WriteString(decl.Property)
		b.WriteString(": ")
		b.WriteString(decl.Value)
		b.WriteString(";\n")
	}
	b.WriteString(indent)
	b.WriteString("}\n")
}

func escapeClass(name string) string {
	if name == "" {
		return ""
	}
	var b strings.Builder
	for i, r := range name {
		switch {
		case r == ':':
			b.WriteString("\\:")
		case r == '.' || r == '/' || r == '%':
			b.WriteRune('\\')
			b.WriteRune(r)
		case i == 0 && r >= '0' && r <= '9':
			b.WriteString("\\3")
			b.WriteRune(r)
			b.WriteString(" ")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

type manifest struct {
	Version int      `json:"version"`
	Files   int      `json:"files"`
	Classes []string `json:"classes"`
	Unknown []string `json:"unknown,omitempty"`
}

func buildManifest(result extract.Result, matched, unknown []string) ([]byte, error) {
	data := manifest{
		Version: manifestVersion,
		Files:   result.Files,
		Classes: matched,
		Unknown: unknown,
	}
	return config.MarshalDeterministic(data)
}

func baseCSS(_ config.Canonical) string {
	var b strings.Builder
	b.WriteString("*, *::before, *::after {\n")
	b.WriteString("  box-sizing: border-box;\n")
	b.WriteString("}\n\n")
	b.WriteString("html, body {\n")
	b.WriteString("  height: 100%;\n")
	b.WriteString("}\n\n")
	b.WriteString("body {\n")
	b.WriteString("  margin: 0;\n")
	b.WriteString("  background-color: var(--color-ink-950);\n")
	b.WriteString("  color: var(--color-ink-100);\n")
	b.WriteString("  font-family: var(--font-sans);\n")
	b.WriteString("  line-height: var(--line-height-normal);\n")
	b.WriteString("}\n")
	return b.String()
}
