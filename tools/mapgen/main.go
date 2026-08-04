package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
	"gopkg.in/yaml.v3"
)

const (
	viewBoxW    = 1200
	viewBoxH    = 1000
	padding     = 20.0
	strokeWidth = 0.4
	outputFile  = "content/docs/meshtastic/regional-lora-settings/regional-lora-settings.svg"
	contentFile = "content/docs/meshtastic/regional-lora-settings/index.md"
	geoJSONFile = "assets/gis/counties.geojson"
)

type frontMatter struct {
	Map     mapColors `yaml:"map"`
	Regions []region  `yaml:"regions"`
}

type mapColors struct {
	Background        string  `yaml:"background"`
	CountyStroke      string  `yaml:"county_stroke"`
	UnassignedCounty  string  `yaml:"unassigned_county"`
	CountyFillOpacity float64 `yaml:"county_fill_opacity"`
	CountyLabel       string  `yaml:"county_label"`
	RegionLabel       string  `yaml:"region_label"`
	RegionLabelHalo   string  `yaml:"region_label_halo"`
}

type regionLabel struct {
	X float64 `yaml:"x"`
	Y float64 `yaml:"y"`
}

type region struct {
	Name     string       `yaml:"name"`
	Color    string       `yaml:"color"`
	Number   string       `yaml:"number"`
	Label    *regionLabel `yaml:"label"`
	Counties []string     `yaml:"counties"`
}

type countyProps struct {
	TigerName string `json:"TIGERNAME"`
}

type countyCollection = geojson.FeatureCollectionOf[countyProps]

type label struct {
	Text string
	X, Y float64
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "mapgen: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	repoRoot, err := findRepoRoot()
	if err != nil {
		return fmt.Errorf("finding repo root: %w", err)
	}

	fm, err := parseFrontmatter(filepath.Join(repoRoot, contentFile))
	if err != nil {
		return fmt.Errorf("parsing frontmatter: %w", err)
	}

	gf, err := parseGeoJSON(filepath.Join(repoRoot, geoJSONFile))
	if err != nil {
		return fmt.Errorf("parsing GeoJSON: %w", err)
	}

	warnDuplicateCounties(fm.Regions, os.Stderr)
	warnUnmatchedCounties(fm, gf, os.Stderr)

	svg, err := buildSVG(fm, gf)
	if err != nil {
		return err
	}

	outPath := filepath.Join(repoRoot, outputFile)
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}
	if err := os.WriteFile(outPath, []byte(svg), 0644); err != nil {
		return fmt.Errorf("writing SVG: %w", err)
	}
	fmt.Printf("Map generated: %s\n", outPath)
	return nil
}

// warnDuplicateCounties writes a warning for each normalized county name that
// appears more than once across region county lists (cross-region or within one
// region). Generation continues; last region in YAML order still wins the fill.
func warnDuplicateCounties(regions []region, w io.Writer) {
	type occurrence struct {
		county string
		region string
	}
	seen := make(map[string][]occurrence)
	order := make([]string, 0)
	for _, r := range regions {
		for _, c := range r.Counties {
			key := normalizeCounty(c)
			if _, ok := seen[key]; !ok {
				order = append(order, key)
			}
			seen[key] = append(seen[key], occurrence{county: c, region: r.Name})
		}
	}
	for _, key := range order {
		occs := seen[key]
		if len(occs) < 2 {
			continue
		}
		regionNames := make([]string, len(occs))
		for i, o := range occs {
			regionNames[i] = o.region
		}
		fmt.Fprintf(w, "mapgen: warning: county %q assigned more than once (regions: %s)\n",
			occs[0].county, strings.Join(regionNames, ", "))
	}
}

// warnUnmatchedCounties writes a warning for each frontmatter county whose
// normalized name matches no GeoJSON TIGERNAME, and for each GeoJSON feature
// matched by no frontmatter county. Generation continues either way; the pair
// of warnings makes spelling drift between the two sources obvious.
func warnUnmatchedCounties(fm *frontMatter, fc *countyCollection, w io.Writer) {
	geoNames := make(map[string]bool)
	for _, f := range fc.Features {
		geoNames[normalizeCounty(f.Properties.TigerName)] = true
	}
	assigned := make(map[string]bool)
	for _, r := range fm.Regions {
		for _, c := range r.Counties {
			key := normalizeCounty(c)
			assigned[key] = true
			if !geoNames[key] {
				fmt.Fprintf(w, "mapgen: warning: county %q (region %q) matches no GeoJSON feature\n", c, r.Name)
			}
		}
	}
	for _, f := range fc.Features {
		if !assigned[normalizeCounty(f.Properties.TigerName)] {
			fmt.Fprintf(w, "mapgen: warning: GeoJSON feature %q matched by no frontmatter county\n", f.Properties.TigerName)
		}
	}
}

// buildSVG generates the regional map SVG from frontmatter and GeoJSON.
// Pure function: no file I/O, no side effects.
func buildSVG(fm *frontMatter, fc *countyCollection) (string, error) {
	if fm.Map.CountyFillOpacity <= 0 || fm.Map.CountyFillOpacity > 1 {
		return "", fmt.Errorf("map.county_fill_opacity must be in (0, 1], got %v", fm.Map.CountyFillOpacity)
	}

	countyColor := make(map[string]string)
	for _, r := range fm.Regions {
		for _, c := range r.Counties {
			countyColor[normalizeCounty(c)] = r.Color
		}
	}

	var bound orb.Bound
	first := true
	for _, feat := range fc.Features {
		if feat.Geometry == nil {
			continue
		}
		if first {
			bound = feat.Geometry.Bound()
			first = false
		} else {
			bound = bound.Union(feat.Geometry.Bound())
		}
	}
	if first {
		return "", fmt.Errorf("no coordinates found in GeoJSON")
	}
	minLon, minLat := bound.Min[0], bound.Min[1]
	maxLon, maxLat := bound.Max[0], bound.Max[1]

	avgLat := (minLat + maxLat) / 2 * math.Pi / 180
	cosLat := math.Cos(avgLat)

	dataW := (maxLon - minLon) * cosLat
	dataH := maxLat - minLat
	scaleX := (float64(viewBoxW) - 2*padding) / dataW
	scaleY := (float64(viewBoxH) - 2*padding) / dataH
	scale := math.Min(scaleX, scaleY)

	offsetX := padding + (float64(viewBoxW)-2*padding-dataW*scale)/2
	offsetY := padding + (float64(viewBoxH)-2*padding-dataH*scale)/2

	transform := func(lon, lat float64) (float64, float64) {
		x := (lon-minLon)*cosLat*scale + offsetX
		y := float64(viewBoxH) - ((lat-minLat)*scale + offsetY)
		return x, y
	}

	// --- Pass 1: build paths AND collect county centroids ---
	var countyLabels []label
	countyCentroid := make(map[string][2]float64)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">`+"\n",
		viewBoxW, viewBoxH, viewBoxW, viewBoxH))
	background := strings.ToLower(strings.TrimSpace(fm.Map.Background))
	if background != "" && background != "none" && background != "transparent" {
		sb.WriteString(fmt.Sprintf(`  <rect width="%d" height="%d" fill="%s"/>`+"\n", viewBoxW, viewBoxH, fm.Map.Background))
	}

	for _, feat := range fc.Features {
		rings, err := geometryRings(feat.Geometry)
		if err != nil {
			return "", fmt.Errorf("feature %q: %w", feat.Properties.TigerName, err)
		}

		fill := countyColor[normalizeCounty(feat.Properties.TigerName)]
		if fill == "" {
			fill = fm.Map.UnassignedCounty
		}

		// Centroid for the county label; degenerate (zero-area) geometries
		// get no label.
		centroid, area := planar.CentroidArea(feat.Geometry)
		if area > 0 {
			cx, cy := transform(centroid[0], centroid[1])
			countyCentroid[normalizeCounty(feat.Properties.TigerName)] = [2]float64{cx, cy}
			countyLabels = append(countyLabels, label{
				Text: feat.Properties.TigerName,
				X:    cx, Y: cy,
			})
		}

		// Build path data.
		var path strings.Builder
		for _, ring := range rings {
			if len(ring) == 0 {
				continue
			}
			x, y := transform(ring[0][0], ring[0][1])
			px, py := int(math.Round(x)), int(math.Round(y))
			path.WriteString(fmt.Sprintf("M%d,%d", px, py))
			for i := 1; i < len(ring); i++ {
				x, y = transform(ring[i][0], ring[i][1])
				nx, ny := int(math.Round(x)), int(math.Round(y))
				if nx != px || ny != py {
					path.WriteString(fmt.Sprintf("L%d,%d", nx, ny))
					px, py = nx, ny
				}
			}
			path.WriteString("Z")
		}

		sb.WriteString(fmt.Sprintf(
			`  <path d="%s" fill="%s" fill-opacity="%.2f" stroke="%s" stroke-width="%.2f" stroke-linejoin="round"/>`+"\n",
			strings.TrimSpace(path.String()), fill, fm.Map.CountyFillOpacity, fm.Map.CountyStroke, strokeWidth))
	}

	// --- County labels (small, white, no stroke) ---
	sb.WriteString(`  <g font-family="Arial, sans-serif" text-anchor="middle" dominant-baseline="middle">` + "\n")
	for _, lbl := range countyLabels {
		sb.WriteString(fmt.Sprintf(
			`    <text x="%.1f" y="%.1f" font-size="9" fill="%s">%s</text>`+"\n",
			lbl.X, lbl.Y, fm.Map.CountyLabel, lbl.Text))
	}
	sb.WriteString("  </g>\n")

	// --- Region labels (positions/numbers from frontmatter; skip if label nil) ---
	sb.WriteString(`  <g font-family="Arial, sans-serif" text-anchor="middle" dominant-baseline="middle" font-weight="bold">` + "\n")
	for _, r := range fm.Regions {
		if r.Label == nil {
			continue
		}
		sb.WriteString(fmt.Sprintf(
			`    <text x="%.1f" y="%.1f" font-size="21" fill="%s" stroke="%s" stroke-width="2.5" paint-order="stroke">%s %s</text>`+"\n",
			r.Label.X, r.Label.Y, fm.Map.RegionLabel, fm.Map.RegionLabelHalo, r.Name, r.Number))
	}
	sb.WriteString("  </g>\n")

	sb.WriteString("</svg>\n")
	return sb.String(), nil
}

func parseFrontmatter(path string) (*frontMatter, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := string(data)
	if !strings.HasPrefix(content, "---") {
		return nil, fmt.Errorf("no YAML frontmatter found")
	}
	end := strings.Index(content[3:], "---")
	if end == -1 {
		return nil, fmt.Errorf("unclosed YAML frontmatter")
	}
	yamlData := content[3 : 3+end]
	var fm frontMatter
	if err := yaml.Unmarshal([]byte(yamlData), &fm); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}
	return &fm, nil
}

func parseGeoJSON(path string) (*countyCollection, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var fc countyCollection
	if err := json.Unmarshal(data, &fc); err != nil {
		return nil, err
	}
	return &fc, nil
}

// geometryRings flattens a Polygon or MultiPolygon into its rings so each can
// be emitted as a closed subpath.
func geometryRings(g orb.Geometry) ([]orb.Ring, error) {
	switch geom := g.(type) {
	case orb.Polygon:
		return geom, nil
	case orb.MultiPolygon:
		var rings []orb.Ring
		for _, poly := range geom {
			rings = append(rings, poly...)
		}
		return rings, nil
	case nil:
		return nil, fmt.Errorf("unsupported geometry type: <nil>")
	default:
		return nil, fmt.Errorf("unsupported geometry type: %s", g.GeoJSONType())
	}
}

func normalizeCounty(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", ""))
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("cannot find repository root")
}
