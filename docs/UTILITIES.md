# Utilities Reference (Concise)

This document is a compact overview of Lattice utility classes. It is intended for agents generating HTML and class lists.

## Naming

- Utility tokens are kebab-case.
- Variants use `:` as a separator: `hover:bg-blue-500`, `md:grid-cols-3`.
- Valid class characters: `a-zA-Z0-9-:_/%`.
- If a `classPrefix` is configured, prepend it to every utility.

## Layout & Display

- Display: `block`, `inline`, `inline-block`, `flex`, `grid`, `hidden`, `contents`.
  - Example: `grid gap-6 md:grid-cols-3`
- Position: `relative`, `absolute`, `fixed`, `sticky`.
- Insets: `inset-*`, `top-*`, `right-*`, `bottom-*`, `left-*`.
  - Example: `relative`, `absolute top-2 right-2`

## Sizing

- Width/height: `w-*`, `h-*`, `min-w-*`, `min-h-*`, `max-w-*`, `max-h-*`.
- Container: `container`.
  - Example: `w-64 h-32 max-w-3xl`

## Spacing

- Padding: `p*`, `pt-*`, `pr-*`, `pb-*`, `pl-*`, `px-*`, `py-*`.
- Margin: `m*`, `mt-*`, `mr-*`, `mb-*`, `ml-*`, `mx-*`, `my-*`.
- Gaps: `gap-*`, `gap-x-*`, `gap-y-*`.
  - Example: `px-6 py-4 gap-4`

## Flex & Grid

- Flex: `flex-*`, `items-*`, `justify-*`, `content-*`, `self-*`.
- Grid: `grid-cols-*`, `grid-rows-*`, `col-span-*`, `row-span-*`.
  - Example: `flex items-center justify-between`
  - Example: `grid grid-cols-2 gap-4`

## Typography

- Size/color/align: `text-*`.
- Font family/weight: `font-*`.
- Line height: `leading-*`.
- Letter spacing: `tracking-*`.
- Case/decoration: `uppercase`, `lowercase`, `underline`, `italic`.
- Lists: `list-*`.
  - Example: `text-sm text-ink-700 font-medium`

## Color

- Background: `bg-*`.
- Text: `text-*`.
- Border: `border-*`.
  - Example: `bg-blue-500 text-white border-ink-200`

## Borders & Radius

- Border: `border`, `border-*` (width/style), side variants (`border-t-*`, etc.).
- Radius: `rounded`, `rounded-*`, and corner variants (`rounded-tr`, etc.).
  - Example: `border border-ink-200 rounded-lg`

## Effects

- Shadow: `shadow`, `shadow-*`.
- Opacity: `opacity-*`.
  - Example: `shadow-md opacity-80`

## Overflow & Visibility

- Overflow: `overflow-*`.
- Visibility: `visible`, `invisible`, `sr-only`.

## Object & Aspect

- Object: `object-*`.
- Aspect: `aspect-*`.
  - Example: `aspect-video object-cover`

## Transitions & Transforms

- Transition: `transition`, `transition-*`.
- Duration/Easing/Delay: `duration-*`, `ease-*`, `delay-*`.
- Translate/Rotate/Scale: `translate-x-*`, `translate-y-*`, `rotate-*`, `scale-*`.
  - Example: `transition duration-200 ease-out hover:translate-y-1`

## Interaction

- Cursor: `cursor-*`.
- Pointer events: `pointer-events-*`.
- Selection: `select-*`.
- Isolation: `isolate`.

## Token Scales (source of `*` values)

- Core: `space`, `size`, `radius`, `borderWidth`, `fontSize`, `lineHeight`, `fontWeight`, `letterSpacing`.
- Effects: `shadow`, `opacity`.
- Layout: `z`, `aspect`, `maxWidth`, `maxHeight`, `container`.
- Motion: `duration`, `easing`, `delay`, `translate`, `rotate`, `scale`.

## Minimal Pattern Guidance

- Use semantic groupings: layout → spacing → typography → color → effects.
- Prefer small examples for each utility group.
- Avoid inline styles; use utilities exclusively.
