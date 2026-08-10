# mapgen

mapgen generates the SVG maps for the
[Regional LoRa Settings](../../content/docs/meshtastic/regional-lora-settings/index.md)
page. It reads county boundaries from `assets/gis/counties.geojson` and writes
two maps:

- `regional-lora-settings.svg` — counties colored by region, from
  `regions.yaml` next to the page (see the sidecar sections below).
- `channel-lora-settings.svg` — counties colored by LoRa channel, from
  `channels.yaml` next to the page.

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

## The generated SVGs are never hand-edited

Both SVGs are generated output. To change anything about the maps — colors,
labels, county assignments — edit `regions.yaml` or `channels.yaml` in
`content/docs/meshtastic/regional-lora-settings/`, re-run
`go -C tools/mapgen run .`, and commit the regenerated SVGs. Never edit
an SVG by hand; the next run will overwrite manual changes.

## Regions sidecar (`regions.yaml`)

`regions.yaml`, next to `index.md`, drives the regional map. The two sidecar
files are the single source of truth for all map colors: mapgen contains no
hardcoded color values and no built-in defaults — every color comes from a
sidecar.

### `regions[]`

Each entry in the `regions` list defines one region:

| Key        | Required | Meaning                                                        |
| ---------- | -------- | -------------------------------------------------------------- |
| `name`     | yes      | Region label text (include number if desired, e.g. `North West #1`). |
| `color`    | no       | Fill color for the region's counties (see resolution below).   |
| `label`    | no       | `{x, y}` position of the region label, plus optional `fontsize` (default 21); omit to skip the label. |
| `counties` | yes      | County names assigned to this region.                          |

Region fill color is read only from `regions[].color`; see the resolution
order below. Unknown keys are ignored, so stale sidecars still parse.

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

## Channels sidecar (`channels.yaml`)

`channels.yaml`, next to `index.md`, drives the channel map. It has the same
shape as the regions sidecar: a `map` block (same keys as above) plus a
`channels` list. Each entry is a group:

| Key        | Required | Meaning                                                        |
| ---------- | -------- | -------------------------------------------------------------- |
| `name`     | yes      | Channel label text. A literal `\n` renders as a line break (e.g. `"LongFast\nMediumFast"` for a county shared by two channels). |
| `color`    | yes      | Fill color for the channel's counties.                         |
| `label`    | no       | `{x, y}` position of the channel label, plus optional `fontsize` (default 21); omit to skip the label. |
| `counties` | no       | County names assigned to this channel.                         |
| `catchall` | no       | At most one channel may set `catchall: true`; see below.       |

### Catchall

A channel with `catchall: true` and no `counties` list colors every county not
claimed by another channel — use it for the default channel so uncovered
counties never render as `unassigned_county`. A catchall group also tells
mapgen the uncovered counties are intentional, suppressing the GeoJSON-side
"matched by no county list" warning (the county-side direction, which catches
spelling drift, always fires). Two catchall groups is a hard error.

## Warnings

mapgen prints warnings to stderr but still generates both SVGs and exits 0.
The checks run on `regions.yaml` and `channels.yaml` alike:

- **Duplicate county**: a county appearing in more than one group's
  `counties` list, or twice in the same list. The warning names the county and
  all groups involved; the **last** group in YAML order wins the fill.
- **Unmatched county**: a listed county matching no GeoJSON `TIGERNAME`
  (spelling drift) always warns. The reverse direction — a GeoJSON feature
  matched by no county list — warns only when no group is a catchall.
