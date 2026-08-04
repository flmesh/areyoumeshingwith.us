# mapgen

mapgen generates the regional SVG map for the
[Regional LoRa Settings](../../content/docs/meshtastic/regional-lora-settings/index.md)
page. It reads county boundaries from `assets/gis/counties.geojson` and colors
each county according to the `regions` and `map` keys in that page's YAML
frontmatter, then writes
`content/docs/meshtastic/regional-lora-settings/regional-lora-settings.svg`.

## Running

From the repository root:

```sh
go -C tools/mapgen run .
```

(`go run ./tools/mapgen` does **not** work: `tools/mapgen` is its own Go
module, so the root module does not contain the package. `go -C` changes
directory before running. The binary locates the repo root by walking up to
the `.git` directory, so it can be run from anywhere inside the repo.)

Run tests with:

```sh
go -C tools/mapgen test ./...
```

## The generated SVG is never hand-edited

`regional-lora-settings.svg` is generated output. To change anything about the
map — colors, region labels, county assignments — edit the frontmatter in
`content/docs/meshtastic/regional-lora-settings/index.md`, re-run
`go -C tools/mapgen run .`, and commit the regenerated SVG. Never edit the SVG
by hand; the next run will overwrite manual changes.

## Frontmatter contract

The page frontmatter is the single source of truth for all map colors. mapgen
contains no hardcoded color values and no built-in defaults — every color comes
from the frontmatter.

### `regions[]`

Each entry in the `regions` list defines one region:

| Key        | Required | Meaning                                                        |
| ---------- | -------- | -------------------------------------------------------------- |
| `name`     | yes      | Region label text (include number if desired, e.g. `North West #1`). |
| `color`    | no       | Fill color for the region's counties (see resolution below).   |
| `label`    | no       | `{x, y}` position of the region label; omit to skip the label. |
| `counties` | yes      | County names assigned to this region.                          |

Region fill color is read only from `regions[].color`; see the resolution
order below. Unknown keys are ignored, so stale frontmatter still parses.

County names are matched case-insensitively with spaces removed, so
`De Soto`, `DeSoto`, and `desoto` are equivalent.

### `map.*`

The `map` block controls all non-region colors:

| Key                  | Meaning                                                            |
| -------------------- | ------------------------------------------------------------------ |
| `background`         | Background rect fill. `none`, `transparent`, or empty → no rect.   |
| `county_stroke`      | County border color.                                               |
| `unassigned_county`  | Fill for counties in no region's `counties` list.                  |
| `county_fill_opacity`| County fill opacity; must be in `(0, 1]` or mapgen errors.         |
| `county_label`       | County name label color.                                           |
| `region_label`       | Region label text color.                                           |
| `region_label_halo`  | Region label outline (halo) color.                                 |

Missing keys fail silently — the corresponding SVG attribute is emitted empty —
so the `map:` block is effectively required. `county_fill_opacity` is the one
exception: an out-of-range value (including the zero value from a missing key)
is a hard error.

### County fill resolution order

1. `regions[].color` of the region whose `counties` list contains the county.
2. `map.unassigned_county` if no region claims the county (or the region's
   `color` is empty).

## Warnings

mapgen prints warnings to stderr but still generates the SVG and exits 0:

- **Duplicate county**: a county appearing in more than one region's
  `counties` list, or twice in the same list. The warning names the county and
  all regions involved; the **last** region in YAML order wins the fill.
- **Unmatched county**: a frontmatter county matching no GeoJSON `TIGERNAME`
  (spelling drift), or a GeoJSON feature matched by no frontmatter county.
