package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	viewBoxW    = 1200
	viewBoxH    = 1000
	padding     = 20.0
	strokeColor = "#222222"
	strokeWidth = 0.4
	outputFile = "content/docs/meshtastic/regional-lora-settings/regional-lora-settings.svg"
	contentFile = "content/docs/meshtastic/regional-lora-settings/index.md"
	geoJSONFile = "assets/gis/counties.geojson"
)

type frontMatter struct {
	Regions []region `yaml:"regions"`
}

type region struct {
	Name     string   `yaml:"name"`
	Emoji    string   `yaml:"emoji"`
	Color    string   `yaml:"color"`
	Counties []string `yaml:"counties"`
}

type geoJSON struct {
	Type     string    `json:"type"`
	Features []feature `json:"features"`
}

type feature struct {
	Type       string     `json:"type"`
	Properties properties `json:"properties"`
	Geometry   geometry   `json:"geometry"`
}

type properties struct {
	TigerName string `json:"TIGERNAME"`
}

type geometry struct {
	Type        string          `json:"type"`
	Coordinates json.RawMessage `json:"coordinates"`
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

	// Parse frontmatter → county → colour lookup.
	fm, err := parseFrontmatter(filepath.Join(repoRoot, contentFile))
	if err != nil {
		return fmt.Errorf("parsing frontmatter: %w", err)
	}
	countyColor := make(map[string]string)
	for _, r := range fm.Regions {
		for _, c := range r.Counties {
			countyColor[normalizeCounty(c)] = r.Color
		}
	}

	// Parse GeoJSON.
	gf, err := parseGeoJSON(filepath.Join(repoRoot, geoJSONFile))
	if err != nil {
		return fmt.Errorf("parsing GeoJSON: %w", err)
	}

	// Compute bounding box for all county coordinates.
	var minLon, minLat, maxLon, maxLat float64
	first := true
	for _, feat := range gf.Features {
		rings, err := extractRings(feat.Geometry)
		if err != nil {
			return fmt.Errorf("feature %q: %w", feat.Properties.TigerName, err)
		}
		for _, ring := range rings {
			for _, pt := range ring {
				if first {
					minLon, minLat, maxLon, maxLat = pt[0], pt[1], pt[0], pt[1]
					first = false
				} else {
					if pt[0] < minLon {
						minLon = pt[0]
					}
					if pt[0] > maxLon {
						maxLon = pt[0]
					}
					if pt[1] < minLat {
						minLat = pt[1]
					}
					if pt[1] > maxLat {
						maxLat = pt[1]
					}
				}
			}
		}
	}
	if first {
		return fmt.Errorf("no coordinates found in GeoJSON")
	}

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
		y := float64(viewBoxH) - ((lat-minLat)*scale+offsetY)
		return x, y
	}

	// Build SVG.
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">`+"\n",
		viewBoxW, viewBoxH, viewBoxW, viewBoxH))
	sb.WriteString(fmt.Sprintf(`  <rect width="%d" height="%d" fill="#f8f9fa"/>`+"\n", viewBoxW, viewBoxH))

	for _, feat := range gf.Features {
		rings, err := extractRings(feat.Geometry)
		if err != nil {
			return fmt.Errorf("feature %q: %w", feat.Properties.TigerName, err)
		}

		fill := countyColor[normalizeCounty(feat.Properties.TigerName)]
		if fill == "" {
			fill = "#cccccc" // gray for unmatched counties
		}

		// Build path data with pixel decimation + integer output.
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
			`  <path d="%s" fill="%s" stroke="%s" stroke-width="%.2f" stroke-linejoin="round"/>`+"\n",
			strings.TrimSpace(path.String()), fill, strokeColor, strokeWidth))
	}

	sb.WriteString("</svg>\n")

	outPath := filepath.Join(repoRoot, outputFile)
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}
	if err := os.WriteFile(outPath, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("writing SVG: %w", err)
	}
	fmt.Printf("Map generated: %s\n", outPath)
	return nil
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

func parseGeoJSON(path string) (*geoJSON, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var gf geoJSON
	if err := json.Unmarshal(data, &gf); err != nil {
		return nil, err
	}
	return &gf, nil
}

func extractRings(g geometry) ([][][2]float64, error) {
	switch g.Type {
	case "Polygon":
		var coords [][][]float64
		if err := json.Unmarshal(g.Coordinates, &coords); err != nil {
			return nil, err
		}
		return toPointRings(coords), nil
	case "MultiPolygon":
		var coords [][][][]float64
		if err := json.Unmarshal(g.Coordinates, &coords); err != nil {
			return nil, err
		}
		var rings [][][2]float64
		for _, poly := range coords {
			rings = append(rings, toPointRings(poly)...)
		}
		return rings, nil
	default:
		return nil, fmt.Errorf("unsupported geometry type: %s", g.Type)
	}
}

func toPointRings(coords [][][]float64) [][][2]float64 {
	var rings [][][2]float64
	for _, ring := range coords {
		pts := make([][2]float64, len(ring))
		for i, p := range ring {
			if len(p) >= 2 {
				pts[i] = [2]float64{p[0], p[1]}
			}
		}
		rings = append(rings, pts)
	}
	return rings
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
