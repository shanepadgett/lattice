package compile

import (
	"math/bits"
	"sort"
	"strconv"
	"strings"

	"lcss/internal/config"
)

const defaultGridSize = 12

func AllClasses(canonical config.Canonical) []string {
	base := baseClasses(canonical)
	if len(base) == 0 {
		return nil
	}

	separator := canonical.Config.Separator
	if separator == "" {
		separator = ":"
	}

	prefix := canonical.Config.ClassPrefix
	if prefix != "" {
		for i, class := range base {
			base[i] = prefix + class
		}
	}

	set := make(map[string]struct{}, len(base))
	for _, class := range base {
		set[class] = struct{}{}
	}

	responsive := canonical.Config.Variants.Responsive
	state := canonical.Config.Variants.State
	stateCombos := orderedSubsets(state)
	if len(responsive) > 0 || len(state) > 0 {
		responsiveOptions := make([]string, 0, len(responsive)+1)
		responsiveOptions = append(responsiveOptions, "")
		responsiveOptions = append(responsiveOptions, responsive...)
		for _, class := range base {
			for _, resp := range responsiveOptions {
				for _, combo := range stateCombos {
					if resp == "" && len(combo) == 0 {
						continue
					}
					parts := make([]string, 0, 1+len(combo))
					if resp != "" {
						parts = append(parts, resp)
					}
					if len(combo) > 0 {
						parts = append(parts, combo...)
					}
					variant := strings.Join(parts, separator) + separator + class
					set[variant] = struct{}{}
				}
			}
		}
	}

	classes := make([]string, 0, len(set))
	for class := range set {
		classes = append(classes, class)
	}
	sort.Strings(classes)
	return classes
}

func orderedSubsets(items []string) [][]string {
	if len(items) == 0 {
		return [][]string{{}}
	}
	count := 1 << len(items)
	combos := make([][]string, 0, count)
	for mask := 0; mask < count; mask++ {
		combo := make([]string, 0, bits.OnesCount(uint(mask)))
		for i, item := range items {
			if mask&(1<<i) != 0 {
				combo = append(combo, item)
			}
		}
		combos = append(combos, combo)
	}
	return combos
}

func baseClasses(canonical config.Canonical) []string {
	set := map[string]struct{}{}
	add := func(name string) {
		if name != "" {
			set[name] = struct{}{}
		}
	}
	addAll := func(prefix string, keys []string) {
		for _, key := range keys {
			add(prefix + key)
		}
	}

	spaceKeys := mapKeys(canonical.Tokens.Scales["space"])
	sizeKeys := mapKeys(canonical.Tokens.Scales["size"])
	maxWidthKeys := mapKeys(canonical.Tokens.Scales["maxWidth"])
	maxHeightKeys := mapKeys(canonical.Tokens.Scales["maxHeight"])
	containerKeys := mapKeys(canonical.Tokens.Scales["container"])
	radiusKeys := mapKeys(canonical.Tokens.Scales["radius"])
	borderWidthKeys := mapKeys(canonical.Tokens.Scales["borderWidth"])
	fontSizeKeys := mapKeys(canonical.Tokens.Scales["fontSize"])
	lineHeightKeys := mapKeys(canonical.Tokens.Scales["lineHeight"])
	fontWeightKeys := mapKeys(canonical.Tokens.Scales["fontWeight"])
	letterSpacingKeys := mapKeys(canonical.Tokens.Scales["letterSpacing"])
	shadowKeys := mapKeys(canonical.Tokens.Scales["shadow"])
	zIndexKeys := mapKeys(canonical.Tokens.Scales["z"])
	opacityKeys := mapKeys(canonical.Tokens.Scales["opacity"])
	aspectKeys := mapKeys(canonical.Tokens.Scales["aspect"])
	durationKeys := mapKeys(canonical.Tokens.Scales["duration"])
	easingKeys := mapKeys(canonical.Tokens.Scales["easing"])
	delayKeys := mapKeys(canonical.Tokens.Scales["delay"])
	translateKeys := mapKeys(canonical.Tokens.Scales["translate"])
	rotateKeys := mapKeys(canonical.Tokens.Scales["rotate"])
	scaleKeys := mapKeys(canonical.Tokens.Scales["scale"])

	colors := canonical.Tokens.Themes["default"].Colors
	colorKeys := mapKeys(colors)
	fonts := canonical.Tokens.Themes["default"].Fonts
	fontKeys := mapKeys(fonts)

	spacingPrefixes := []string{"p-", "px-", "py-", "pt-", "pr-", "pb-", "pl-", "m-", "mx-", "my-", "mt-", "mr-", "mb-", "ml-", "gap-", "gapx-", "gapy-", "gap-x-", "gap-y-"}
	for _, prefix := range spacingPrefixes {
		addAll(prefix, spaceKeys)
	}

	commonSizeKeys := mergeKeys(sizeValueKeys(spaceKeys, sizeKeys))
	addAll("w-", commonSizeKeys)
	addAll("h-", commonSizeKeys)
	addAll("min-w-", commonSizeKeys)
	addAll("min-h-", commonSizeKeys)

	maxWidthAll := mergeKeys(maxWidthKeys, commonSizeKeys)
	addAll("max-w-", maxWidthAll)
	maxHeightAll := mergeKeys(maxHeightKeys, commonSizeKeys)
	addAll("max-h-", maxHeightAll)

	if len(containerKeys) > 0 {
		add("container")
	}

	for _, value := range []string{"block", "inline-block", "inline", "flex", "inline-flex", "grid", "hidden", "contents"} {
		add(value)
	}

	for _, value := range []string{"static", "relative", "absolute", "fixed", "sticky"} {
		add(value)
	}
	insetPrefixes := []string{"inset-", "inset-x-", "inset-y-", "top-", "right-", "bottom-", "left-"}
	for _, prefix := range insetPrefixes {
		addAll(prefix, commonSizeKeys)
	}

	for _, value := range []string{"flex-row", "flex-col", "flex-wrap", "flex-nowrap", "flex-wrap-reverse", "flex-1", "flex-auto", "flex-initial", "flex-none", "grow", "grow-0", "shrink", "shrink-0"} {
		add(value)
	}
	for _, value := range []string{"start", "center", "end", "stretch", "baseline"} {
		add("items-" + value)
		add("self-" + value)
	}
	for _, value := range []string{"start", "center", "end", "between", "around", "evenly"} {
		add("justify-" + value)
		add("content-" + value)
	}

	add("grid-cols-none")
	add("grid-rows-none")
	add("col-span-full")
	add("row-span-full")
	gridColumns := canonical.Config.Build.GridColumns
	if gridColumns == 0 {
		gridColumns = defaultGridSize
	}
	for i := 1; i <= gridColumns; i++ {
		value := fmtInt(i)
		add("grid-cols-" + value)
		add("grid-rows-" + value)
		add("col-span-" + value)
		add("row-span-" + value)
		add("col-start-" + value)
		add("col-end-" + value)
		add("row-start-" + value)
		add("row-end-" + value)
	}
	for _, value := range []string{"start", "center", "end", "stretch", "baseline"} {
		add("place-items-" + value)
	}
	for _, value := range []string{"start", "center", "end", "between", "around", "evenly"} {
		add("place-content-" + value)
	}

	addAll("text-", fontSizeKeys)
	addAll("text-", colorKeys)
	for _, value := range []string{"left", "center", "right", "justify"} {
		add("text-" + value)
	}
	addAll("leading-", lineHeightKeys)
	addAll("font-", fontKeys)
	addAll("font-", fontWeightKeys)
	for _, value := range []string{"italic", "not-italic", "uppercase", "lowercase", "capitalize", "normal-case", "underline", "line-through", "no-underline", "list-none", "list-disc", "list-decimal"} {
		add(value)
	}
	addAll("tracking-", letterSpacingKeys)

	addAll("bg-", colorKeys)
	addAll("text-", colorKeys)
	addAll("border-", colorKeys)

	for _, value := range []string{"bg-cover", "bg-contain", "bg-center", "bg-top", "bg-right", "bg-bottom", "bg-left", "bg-fixed", "bg-local", "bg-scroll", "bg-repeat", "bg-no-repeat", "bg-repeat-x", "bg-repeat-y"} {
		add(value)
	}

	add("border")
	for _, value := range []string{"border-solid", "border-dashed", "border-dotted", "border-double", "border-none"} {
		add(value)
	}
	borderWidthAll := mergeKeys(borderWidthKeys, []string{"0", "2", "4", "8"})
	addAll("border-", borderWidthAll)
	for _, prefix := range []string{"border-x-", "border-y-", "border-t-", "border-r-", "border-b-", "border-l-"} {
		addAll(prefix, borderWidthAll)
	}

	add("rounded")
	addAll("rounded-", radiusKeys)
	for _, value := range []string{"t", "b", "l", "r", "tl", "tr", "bl", "br"} {
		add("rounded-" + value)
	}

	add("shadow")
	add("shadow-none")
	addAll("shadow-", shadowKeys)

	addAll("opacity-", opacityKeys)
	add("z-auto")
	addAll("z-", zIndexKeys)

	for _, value := range []string{"overflow-auto", "overflow-hidden", "overflow-visible", "overflow-scroll", "overflow-x-auto", "overflow-x-hidden", "overflow-x-visible", "overflow-x-scroll", "overflow-y-auto", "overflow-y-hidden", "overflow-y-visible", "overflow-y-scroll"} {
		add(value)
	}

	add("visible")
	add("invisible")
	add("sr-only")

	for _, value := range []string{"object-contain", "object-cover", "object-fill", "object-none", "object-scale-down", "object-center", "object-top", "object-right", "object-bottom", "object-left"} {
		add(value)
	}
	addAll("aspect-", aspectKeys)

	for _, value := range []string{"transition", "transition-colors", "transition-opacity", "transition-transform"} {
		add(value)
	}
	addAll("duration-", durationKeys)
	addAll("ease-", easingKeys)
	addAll("delay-", delayKeys)

	translateAll := mergeKeys(translateKeys, spaceKeys, []string{"full"})
	addAll("translate-x-", translateAll)
	addAll("translate-y-", translateAll)
	addAll("rotate-", rotateKeys)
	addAll("scale-", scaleKeys)

	for _, value := range []string{"cursor-pointer", "cursor-default", "cursor-text", "cursor-not-allowed", "pointer-events-none", "pointer-events-auto", "select-none", "select-text", "select-all", "select-auto", "isolate", "isolation-auto"} {
		add(value)
	}

	classes := make([]string, 0, len(set))
	for class := range set {
		classes = append(classes, class)
	}
	sort.Strings(classes)
	return classes
}

func sizeValueKeys(spaceKeys, sizeKeys []string) []string {
	return mergeKeys([]string{"auto", "full", "screen", "min", "max", "fit"}, spaceKeys, sizeKeys)
}

func mapKeys(values map[string]string) []string {
	if len(values) == 0 {
		return nil
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func mergeKeys(groups ...[]string) []string {
	set := map[string]struct{}{}
	for _, group := range groups {
		for _, value := range group {
			if value == "" {
				continue
			}
			set[value] = struct{}{}
		}
	}
	keys := make([]string, 0, len(set))
	for value := range set {
		keys = append(keys, value)
	}
	sort.Strings(keys)
	return keys
}

func fmtInt(value int) string {
	return strconv.Itoa(value)
}
