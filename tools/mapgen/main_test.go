package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

// Characterization: run() regenerates the committed SVG byte-identically.
// Restores the original file via t.Cleanup so a failed assertion cannot leave
// the working tree dirty.
func TestRun_RegeneratesCommittedSVG(t *testing.T) {
	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Fatalf("findRepoRoot: %v", err)
	}
	outPath := filepath.Join(repoRoot, outputFile)

	original, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read committed SVG: %v", err)
	}
	t.Cleanup(func() {
		if err := os.WriteFile(outPath, original, 0644); err != nil {
			t.Errorf("restore committed SVG: %v", err)
		}
	})

	if err := run(); err != nil {
		t.Fatalf("run: %v", err)
	}

	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read regenerated SVG: %v", err)
	}
	if !bytes.Equal(got, original) {
		t.Errorf("regenerated SVG differs from committed (%d bytes vs %d bytes)", len(got), len(original))
	}
}

// squareCountyFixture is a single unit square centered on the equator so
// cosLat=1 and every transform is hand-computable:
//
//	lon/lat ∈ [-1,1]² → avgLat=0, cosLat=1
//	dataW=dataH=2; scale=min(580,480)=480
//	offsetX=120, offsetY=20
//	transform(lon,lat) = ((lon+1)*480+120, 1000-((lat+1)*480+20))
//	corners → (120,980),(1080,980),(1080,20),(120,20); centroid → (600,500)
func squareCountyFixture() (*frontMatter, *countyCollection) {
	fm := &frontMatter{
		Map: mapColors{
			Background:        "none",
			CountyStroke:      "#000000",
			UnassignedCounty:  "#cccccc",
			CountyFillOpacity: 1.0,
			CountyLabel:       "white",
			RegionLabel:       "black",
			RegionLabelHalo:   "white",
		},
		Regions: []region{
			{
				Name:     "Test #1",
				Color:    "#ff0000",
				Label:    &regionLabel{X: 100, Y: 200},
				Counties: []string{"Square"},
			},
		},
	}

	// Closed ring, CCW: SW → SE → NE → NW → SW
	gf := &countyCollection{
		Type: "FeatureCollection",
		Features: []*geojson.FeatureOf[countyProps]{
			{
				Type:       "Feature",
				Properties: countyProps{TigerName: "Square"},
				Geometry: orb.Polygon{
					{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}},
				},
			},
		},
	}
	return fm, gf
}

func TestBuildSVG_SquareCounty(t *testing.T) {
	fm, gf := squareCountyFixture()
	svg, err := buildSVG(fm, gf)
	if err != nil {
		t.Fatalf("buildSVG: %v", err)
	}

	// Path corners rounded to ints (see squareCountyFixture comment).
	wantPath := `d="M120,980L1080,980L1080,20L120,20L120,980Z"`
	if !strings.Contains(svg, wantPath) {
		t.Errorf("SVG missing expected path %s\n got:\n%s", wantPath, svg)
	}

	// Fill from region color, opacity 1.00, stroke from map colors.
	if !strings.Contains(svg, `fill="#ff0000" fill-opacity="1.00" stroke="#000000"`) {
		t.Errorf("SVG missing expected path styling\n got:\n%s", svg)
	}

	// County label at centroid (600.0, 500.0).
	if !strings.Contains(svg, `<text x="600.0" y="500.0" font-size="9" fill="white">Square</text>`) {
		t.Errorf("SVG missing county label at centroid\n got:\n%s", svg)
	}

	// Region label from frontmatter.
	if !strings.Contains(svg, `<text x="100.0" y="200.0" font-size="21" fill="black" stroke="white"`) {
		t.Errorf("SVG missing region label\n got:\n%s", svg)
	}
	if !strings.Contains(svg, `>Test #1</text>`) {
		t.Errorf("SVG missing region label text\n got:\n%s", svg)
	}

	// Transparent/none background → no rect.
	if strings.Contains(svg, "<rect ") {
		t.Errorf("expected no background rect for background=none\n got:\n%s", svg)
	}
}

func TestBuildSVG_BackgroundRect(t *testing.T) {
	fm, gf := squareCountyFixture()
	fm.Map.Background = "#112233"

	svg, err := buildSVG(fm, gf)
	if err != nil {
		t.Fatalf("buildSVG: %v", err)
	}
	want := `<rect width="1200" height="1000" fill="#112233"/>`
	if !strings.Contains(svg, want) {
		t.Errorf("SVG missing background rect %q\n got:\n%s", want, svg)
	}
}

func TestBuildSVG_UnassignedCounty(t *testing.T) {
	fm, gf := squareCountyFixture()
	fm.Regions = nil // no region assignment

	svg, err := buildSVG(fm, gf)
	if err != nil {
		t.Fatalf("buildSVG: %v", err)
	}
	if !strings.Contains(svg, `fill="#cccccc"`) {
		t.Errorf("expected unassigned fill #cccccc\n got:\n%s", svg)
	}
}

func TestBuildSVG_NilRegionLabelSkipped(t *testing.T) {
	fm, gf := squareCountyFixture()
	fm.Regions[0].Label = nil

	svg, err := buildSVG(fm, gf)
	if err != nil {
		t.Fatalf("buildSVG: %v", err)
	}
	if strings.Contains(svg, "Test #1") {
		t.Errorf("expected nil region label to be skipped\n got:\n%s", svg)
	}
}

func TestBuildSVG_RejectsInvalidOpacity(t *testing.T) {
	cases := []struct {
		name    string
		opacity float64
	}{
		{"zero", 0},
		{"negative", -0.1},
		{"above_one", 1.1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fm, gf := squareCountyFixture()
			fm.Map.CountyFillOpacity = tc.opacity
			_, err := buildSVG(fm, gf)
			if err == nil {
				t.Fatal("expected error for invalid county_fill_opacity")
			}
			if !strings.Contains(err.Error(), "county_fill_opacity") {
				t.Errorf("error %q should mention county_fill_opacity", err)
			}
		})
	}
}

func TestBuildSVG_EmptyGeoJSON(t *testing.T) {
	fm, _ := squareCountyFixture()
	gf := &countyCollection{Type: "FeatureCollection", Features: nil}
	_, err := buildSVG(fm, gf)
	if err == nil {
		t.Fatal("expected error for empty GeoJSON")
	}
	if !strings.Contains(err.Error(), "no coordinates") {
		t.Errorf("error %q should mention no coordinates", err)
	}
}

// multiPolygonFixture is a single feature with two polygons: a large square
// [-1,-1]→[1,1] (area=4, centroid=(0,0)) and a small triangle inside it
// (area=0.045, centroid=(0.65,0.6)). planar.CentroidArea returns the
// area-weighted centroid of both polygons: ((0*4+0.65*0.045)/4.045,
// (0*4+0.6*0.045)/4.045) ≈ (0.00723, 0.00667) → SVG (603.5, 496.8).
// Bbox is still [-1,1]² so the transform is unchanged from squareCountyFixture.
func multiPolygonFixture() (*frontMatter, *countyCollection) {
	fm := &frontMatter{
		Map: mapColors{
			Background:        "none",
			CountyStroke:      "#000000",
			UnassignedCounty:  "#cccccc",
			CountyFillOpacity: 1.0,
			CountyLabel:       "white",
			RegionLabel:       "black",
			RegionLabelHalo:   "white",
		},
		Regions: []region{
			{
				Name:     "Test #1",
				Color:    "#ff0000",
				Label:    &regionLabel{X: 100, Y: 200},
				Counties: []string{"Multi"},
			},
		},
	}

	// MultiPolygon: two polygons, one ring each.
	// Polygon 1 (large): square [-1,-1]→[1,1]
	// Polygon 2 (small): triangle (0.5,0.5)→(0.8,0.5)→(0.65,0.8)
	gf := &countyCollection{
		Type: "FeatureCollection",
		Features: []*geojson.FeatureOf[countyProps]{
			{
				Type:       "Feature",
				Properties: countyProps{TigerName: "Multi"},
				Geometry: orb.MultiPolygon{
					{{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}},
					{{{0.5, 0.5}, {0.8, 0.5}, {0.65, 0.8}, {0.5, 0.5}}},
				},
			},
		},
	}
	return fm, gf
}

func TestBuildSVG_MultiPolygon(t *testing.T) {
	fm, gf := multiPolygonFixture()
	svg, err := buildSVG(fm, gf)
	if err != nil {
		t.Fatalf("buildSVG: %v", err)
	}

	// Both rings concatenated into a single d attribute.
	// Large square: M120,980L1080,980L1080,20L120,20L120,980Z
	// Small triangle: M840,260L984,260L912,116L840,260Z
	if !strings.Contains(svg, `M120,980L1080,980L1080,20L120,20L120,980Z`) {
		t.Errorf("MultiPolygon missing large square path\n got:\n%s", svg)
	}
	if !strings.Contains(svg, `M840,260L984,260L912,116L840,260Z`) {
		t.Errorf("MultiPolygon missing small triangle path\n got:\n%s", svg)
	}

	// Label at the area-weighted centroid (see multiPolygonFixture comment).
	if !strings.Contains(svg, `<text x="603.5" y="496.8" font-size="9" fill="white">Multi</text>`) {
		t.Errorf("MultiPolygon label at wrong position\n got:\n%s", svg)
	}
}

// twoCountyFixture is a single square (assigned to "Test") and a small
// triangle ("Tri", unassigned). Bbox stays [-1,1]², transform unchanged.
func twoCountyFixture() (*frontMatter, *countyCollection) {
	fm := &frontMatter{
		Map: mapColors{
			Background:        "none",
			CountyStroke:      "#000000",
			UnassignedCounty:  "#cccccc",
			CountyFillOpacity: 1.0,
			CountyLabel:       "white",
			RegionLabel:       "black",
			RegionLabelHalo:   "white",
		},
		Regions: []region{
			{
				Name:     "Test #1",
				Color:    "#ff0000",
				Label:    &regionLabel{X: 100, Y: 200},
				Counties: []string{"Square"},
			},
		},
	}

	gf := &countyCollection{
		Type: "FeatureCollection",
		Features: []*geojson.FeatureOf[countyProps]{
			{
				Type:       "Feature",
				Properties: countyProps{TigerName: "Square"},
				Geometry: orb.Polygon{
					{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}},
				},
			},
			{
				Type:       "Feature",
				Properties: countyProps{TigerName: "Tri"},
				Geometry: orb.Polygon{
					{{0.5, 0.5}, {0.8, 0.5}, {0.65, 0.8}, {0.5, 0.5}},
				},
			},
		},
	}
	return fm, gf
}

func TestBuildSVG_MixedAssignment(t *testing.T) {
	fm, gf := twoCountyFixture()
	svg, err := buildSVG(fm, gf)
	if err != nil {
		t.Fatalf("buildSVG: %v", err)
	}

	// Assigned county gets region color.
	if !strings.Contains(svg, `fill="#ff0000"`) {
		t.Errorf("assigned county missing region fill\n got:\n%s", svg)
	}

	// Unassigned county gets the unassigned default.
	if !strings.Contains(svg, `fill="#cccccc"`) {
		t.Errorf("unassigned county missing default fill\n got:\n%s", svg)
	}

	// Both labels present.
	if !strings.Contains(svg, ">Square</text>") {
		t.Errorf("missing Square label\n got:\n%s", svg)
	}
	if !strings.Contains(svg, ">Tri</text>") {
		t.Errorf("missing Tri label\n got:\n%s", svg)
	}
}

func TestBuildSVG_TransparentBackground(t *testing.T) {
	fm, gf := squareCountyFixture()
	fm.Map.Background = "transparent"

	svg, err := buildSVG(fm, gf)
	if err != nil {
		t.Fatalf("buildSVG: %v", err)
	}
	if strings.Contains(svg, "<rect ") {
		t.Errorf("expected no background rect for background=transparent\n got:\n%s", svg)
	}
}

func TestBuildSVG_UnsupportedGeometryType(t *testing.T) {
	fm, _ := squareCountyFixture()
	gf := &countyCollection{
		Type: "FeatureCollection",
		Features: []*geojson.FeatureOf[countyProps]{
			{
				Type:       "Feature",
				Properties: countyProps{TigerName: "Bad"},
				Geometry:   orb.Point{0, 0},
			},
		},
	}
	_, err := buildSVG(fm, gf)
	if err == nil {
		t.Fatal("expected error for unsupported geometry type")
	}
	if !strings.Contains(err.Error(), "unsupported geometry type") {
		t.Errorf("error %q should mention unsupported geometry type", err)
	}
}

func TestWarnDuplicateCounties_CrossRegion(t *testing.T) {
	regions := []region{
		{Name: "Alpha", Color: "#aa0000", Counties: []string{"Square"}},
		{Name: "Beta", Color: "#00bb00", Counties: []string{"Square"}},
	}
	var buf bytes.Buffer
	warnDuplicateCounties(regions, &buf)
	got := buf.String()
	if !strings.Contains(got, "Square") {
		t.Errorf("warning should name county Square; got %q", got)
	}
	if !strings.Contains(got, "Alpha") || !strings.Contains(got, "Beta") {
		t.Errorf("warning should name both regions; got %q", got)
	}
}

func TestWarnDuplicateCounties_WithinRegion(t *testing.T) {
	regions := []region{
		{Name: "Solo", Color: "#aa0000", Counties: []string{"Square", "Square"}},
	}
	var buf bytes.Buffer
	warnDuplicateCounties(regions, &buf)
	got := buf.String()
	if got == "" {
		t.Fatal("expected stderr warning for county repeated within one region")
	}
	if !strings.Contains(got, "Square") {
		t.Errorf("warning should name county Square; got %q", got)
	}
	if !strings.Contains(got, "Solo") {
		t.Errorf("warning should name region Solo; got %q", got)
	}
}

func TestWarnDuplicateCounties_NoDuplicates(t *testing.T) {
	regions := []region{
		{Name: "Alpha", Color: "#aa0000", Counties: []string{"Square"}},
		{Name: "Beta", Color: "#00bb00", Counties: []string{"Tri"}},
	}
	var buf bytes.Buffer
	warnDuplicateCounties(regions, &buf)
	if buf.Len() != 0 {
		t.Errorf("expected no warning; got %q", buf.String())
	}
}

func TestWarnUnmatchedCounties_FrontmatterCountyMissingFromGeoJSON(t *testing.T) {
	fm, gf := squareCountyFixture()
	fm.Regions[0].Counties = append(fm.Regions[0].Counties, "Phantom")

	var buf bytes.Buffer
	warnUnmatchedCounties(fm, gf, &buf)
	got := buf.String()

	if !strings.Contains(got, "Phantom") {
		t.Errorf("warning should name county Phantom; got %q", got)
	}
	if !strings.Contains(got, "Test") {
		t.Errorf("warning should name region Test; got %q", got)
	}
}

func TestWarnUnmatchedCounties_GeoJSONFeatureMatchedByNoCounty(t *testing.T) {
	fm, gf := twoCountyFixture() // "Square" assigned; "Tri" is GeoJSON-only

	var buf bytes.Buffer
	warnUnmatchedCounties(fm, gf, &buf)
	got := buf.String()

	if !strings.Contains(got, "Tri") {
		t.Errorf("warning should name TIGERNAME Tri; got %q", got)
	}
	if strings.Contains(got, "Square") {
		t.Errorf("matched county Square should not warn; got %q", got)
	}
}

func TestWarnUnmatchedCounties_CurrentRepoDataIsClean(t *testing.T) {
	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Fatalf("findRepoRoot: %v", err)
	}
	fm, err := parseFrontmatter(filepath.Join(repoRoot, contentFile))
	if err != nil {
		t.Fatalf("parseFrontmatter: %v", err)
	}
	gf, err := parseGeoJSON(filepath.Join(repoRoot, geoJSONFile))
	if err != nil {
		t.Fatalf("parseGeoJSON: %v", err)
	}

	var buf bytes.Buffer
	warnUnmatchedCounties(fm, gf, &buf)

	if buf.Len() != 0 {
		t.Errorf("expected no warnings for current repo data; got:\n%s", buf.String())
	}
}

func TestBuildSVG_PhantomCountyStillGenerates(t *testing.T) {
	fm, gf := squareCountyFixture()
	fm.Regions[0].Counties = append(fm.Regions[0].Counties, "Phantom")

	svg, err := buildSVG(fm, gf)
	if err != nil {
		t.Fatalf("buildSVG with phantom county: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Errorf("expected SVG output despite phantom county\n got:\n%s", svg)
	}
}

func TestWarnUnmatchedCounties_AllNamesMatch(t *testing.T) {
	fm, gf := squareCountyFixture()

	var buf bytes.Buffer
	warnUnmatchedCounties(fm, gf, &buf)

	if buf.Len() != 0 {
		t.Errorf("expected no warning when all names match; got %q", buf.String())
	}
}

func TestBuildSVG_DuplicateCountyLastRegionWins(t *testing.T) {
	fm, gf := squareCountyFixture()
	fm.Regions = []region{
		{Name: "Alpha #1", Color: "#aa0000", Counties: []string{"Square"}},
		{Name: "Beta #2", Color: "#00bb00", Counties: []string{"Square"}},
	}
	svg, err := buildSVG(fm, gf)
	if err != nil {
		t.Fatalf("buildSVG: %v", err)
	}
	if !strings.Contains(svg, `fill="#00bb00"`) {
		t.Errorf("expected later region color #00bb00; got:\n%s", svg)
	}
	if strings.Contains(svg, `fill="#aa0000"`) {
		t.Errorf("earlier region color should not win; got:\n%s", svg)
	}
}

// The README is the mapgen user documentation (aym-emi); it must cover the
// tool's purpose, how to run it, the frontmatter contract, and the rule that
// the generated SVG is never hand-edited.
func TestREADME_DocumentsMapgen(t *testing.T) {
	data, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatalf("read README.md: %v", err)
	}
	readme := string(data)

	required := []string{
		"GeoJSON",
		"go -C tools/mapgen run .",
		"regions",
		"counties",
		"county_stroke",
		"unassigned_county",
		"county_fill_opacity",
		"county_label",
		"region_label",
		"region_label_halo",
		"background",
		"never hand-edit",
	}
	for _, want := range required {
		if !strings.Contains(readme, want) {
			t.Errorf("README.md missing required content %q", want)
		}
	}
}

func TestNormalizeCounty(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"De Soto", "desoto"},
		{"St. Lucie", "st.lucie"},
		{"Miami-Dade", "miami-dade"},
		{"broward", "broward"},
		{"", ""},
	}
	for _, tc := range cases {
		if got := normalizeCounty(tc.in); got != tc.want {
			t.Errorf("normalizeCounty(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestWarnDuplicateCounties_NormalizedNames(t *testing.T) {
	regions := []region{
		{Name: "Alpha", Color: "#aa0000", Counties: []string{"De Soto"}},
		{Name: "Beta", Color: "#00bb00", Counties: []string{"desoto"}},
	}
	var buf bytes.Buffer
	warnDuplicateCounties(regions, &buf)
	got := buf.String()
	if !strings.Contains(got, "De Soto") {
		t.Errorf("warning should name first-seen county spelling; got %q", got)
	}
	if !strings.Contains(got, "Alpha") || !strings.Contains(got, "Beta") {
		t.Errorf("warning should name both regions; got %q", got)
	}
}

// writeTempFile writes content to a named file in the test's temp dir.
func writeTempFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return path
}

func TestParseFrontmatter_Valid(t *testing.T) {
	path := writeTempFile(t, "index.md", `---
map:
  background: "#112233"
  county_stroke: "#000000"
  unassigned_county: "#cccccc"
  county_fill_opacity: 0.85
  county_label: white
  region_label: black
  region_label_halo: white
regions:
  - name: "Test #1"
    color: "#ff0000"
    label:
      x: 100
      y: 200
    counties:
      - Square
      - Tri
  - name: "NoLabel #2"
    color: "#00bb00"
    counties:
      - Other
---
Body below the frontmatter is ignored.
`)

	fm, err := parseFrontmatter(path)
	if err != nil {
		t.Fatalf("parseFrontmatter: %v", err)
	}
	if fm.Map.Background != "#112233" {
		t.Errorf("background = %q, want #112233", fm.Map.Background)
	}
	if fm.Map.CountyFillOpacity != 0.85 {
		t.Errorf("county_fill_opacity = %v, want 0.85", fm.Map.CountyFillOpacity)
	}
	if len(fm.Regions) != 2 {
		t.Fatalf("regions = %d, want 2", len(fm.Regions))
	}
	r := fm.Regions[0]
	if r.Name != "Test #1" || r.Color != "#ff0000" {
		t.Errorf("region 0 = %+v", r)
	}
	if r.Label == nil || r.Label.X != 100 || r.Label.Y != 200 {
		t.Errorf("region 0 label = %+v, want (100, 200)", r.Label)
	}
	if len(r.Counties) != 2 || r.Counties[0] != "Square" || r.Counties[1] != "Tri" {
		t.Errorf("region 0 counties = %v, want [Square Tri]", r.Counties)
	}
	if fm.Regions[1].Label != nil {
		t.Errorf("region without label key should have nil Label, got %+v", fm.Regions[1].Label)
	}
}

func TestParseFrontmatter_UnknownKeysIgnored(t *testing.T) {
	path := writeTempFile(t, "index.md", `---
regions:
  - name: Test
    emoji: 🌴
    color: "#ff0000"
    counties:
      - Square
---
`)

	fm, err := parseFrontmatter(path)
	if err != nil {
		t.Fatalf("parseFrontmatter with unknown key: %v", err)
	}
	if len(fm.Regions) != 1 || fm.Regions[0].Name != "Test" {
		t.Errorf("regions = %+v, want single region named Test", fm.Regions)
	}
}

func TestParseFrontmatter_Errors(t *testing.T) {
	cases := []struct {
		name    string
		content string
		wantErr string
	}{
		{"no opening delimiter", "no frontmatter here", "no YAML frontmatter found"},
		{"unclosed frontmatter", "---\nmap:\n  background: x\n", "unclosed YAML frontmatter"},
		{"malformed YAML", "---\nmap: [unclosed\n---\n", "parsing YAML"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := writeTempFile(t, "index.md", tc.content)
			_, err := parseFrontmatter(path)
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("error %q should mention %q", err, tc.wantErr)
			}
		})
	}
}

func TestParseFrontmatter_MissingFile(t *testing.T) {
	_, err := parseFrontmatter(filepath.Join(t.TempDir(), "nope.md"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestParseGeoJSON_Valid(t *testing.T) {
	path := writeTempFile(t, "counties.geojson", `{"type":"FeatureCollection","features":[{"type":"Feature","properties":{"TIGERNAME":"Square"},"geometry":{"type":"Polygon","coordinates":[[[-1,-1],[1,-1],[1,1],[-1,1],[-1,-1]]]}}]}`)

	fc, err := parseGeoJSON(path)
	if err != nil {
		t.Fatalf("parseGeoJSON: %v", err)
	}
	if len(fc.Features) != 1 {
		t.Fatalf("features = %d, want 1", len(fc.Features))
	}
	if fc.Features[0].Properties.TigerName != "Square" {
		t.Errorf("TIGERNAME = %q, want Square", fc.Features[0].Properties.TigerName)
	}
	if fc.Features[0].Geometry == nil {
		t.Error("expected geometry to be parsed")
	}
}

func TestParseGeoJSON_Errors(t *testing.T) {
	t.Run("missing file", func(t *testing.T) {
		_, err := parseGeoJSON(filepath.Join(t.TempDir(), "nope.geojson"))
		if err == nil {
			t.Fatal("expected error for missing file")
		}
	})
	t.Run("invalid JSON", func(t *testing.T) {
		path := writeTempFile(t, "bad.geojson", "{not json")
		_, err := parseGeoJSON(path)
		if err == nil {
			t.Fatal("expected error for invalid JSON")
		}
	})
}

func TestGeometryRings(t *testing.T) {
	square := orb.Ring{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}
	tri := orb.Ring{{0, 0}, {0.5, 0}, {0.25, 0.5}, {0, 0}}

	t.Run("polygon returns its rings", func(t *testing.T) {
		rings, err := geometryRings(orb.Polygon{square})
		if err != nil {
			t.Fatalf("geometryRings: %v", err)
		}
		if len(rings) != 1 {
			t.Errorf("rings = %d, want 1", len(rings))
		}
	})
	t.Run("multipolygon flattens rings", func(t *testing.T) {
		rings, err := geometryRings(orb.MultiPolygon{{square}, {tri}})
		if err != nil {
			t.Fatalf("geometryRings: %v", err)
		}
		if len(rings) != 2 {
			t.Errorf("rings = %d, want 2", len(rings))
		}
	})
	t.Run("nil geometry errors", func(t *testing.T) {
		_, err := geometryRings(nil)
		if err == nil {
			t.Fatal("expected error for nil geometry")
		}
		if !strings.Contains(err.Error(), "unsupported geometry type") {
			t.Errorf("error %q should mention unsupported geometry type", err)
		}
	})
	t.Run("non-polygon type errors", func(t *testing.T) {
		_, err := geometryRings(orb.Point{0, 0})
		if err == nil {
			t.Fatal("expected error for Point geometry")
		}
		if !strings.Contains(err.Error(), "Point") {
			t.Errorf("error %q should name the geometry type", err)
		}
	})
}

// chdir switches the working directory for the duration of a test.
func chdir(t *testing.T, dir string) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir %s: %v", dir, err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(orig); err != nil {
			t.Errorf("restore cwd: %v", err)
		}
	})
}

func TestFindRepoRoot_FromNestedSubdir(t *testing.T) {
	// Tests run with cwd = package dir (tools/mapgen), below the repo root.
	root, err := findRepoRoot()
	if err != nil {
		t.Fatalf("findRepoRoot: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".git")); err != nil {
		t.Errorf("returned root %q has no .git: %v", root, err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if cwd == root {
		t.Error("expected findRepoRoot to walk up from the nested package dir")
	}
}

func TestFindRepoRoot_FromRepoRoot(t *testing.T) {
	root, err := findRepoRoot()
	if err != nil {
		t.Fatalf("findRepoRoot: %v", err)
	}
	chdir(t, root)
	got, err := findRepoRoot()
	if err != nil {
		t.Fatalf("findRepoRoot from root: %v", err)
	}
	if got != root {
		t.Errorf("findRepoRoot from root = %q, want %q", got, root)
	}
}

func TestFindRepoRoot_OutsideRepoErrors(t *testing.T) {
	chdir(t, t.TempDir())
	_, err := findRepoRoot()
	if err == nil {
		t.Fatal("expected error outside any repository")
	}
	if !strings.Contains(err.Error(), "cannot find repository root") {
		t.Errorf("error %q should mention cannot find repository root", err)
	}
}

// A collinear ring has zero area; planar.CentroidArea reports area 0, so the
// county still gets a path but no label.
func TestBuildSVG_DegenerateGeometryGetsNoLabel(t *testing.T) {
	fm, gf := squareCountyFixture()
	gf.Features = append(gf.Features, &geojson.FeatureOf[countyProps]{
		Type:       "Feature",
		Properties: countyProps{TigerName: "Line"},
		Geometry: orb.Polygon{
			{{-1, 0}, {0, 0}, {1, 0}, {-1, 0}},
		},
	})

	svg, err := buildSVG(fm, gf)
	if err != nil {
		t.Fatalf("buildSVG: %v", err)
	}
	if strings.Contains(svg, ">Line</text>") {
		t.Errorf("zero-area county should get no label\n got:\n%s", svg)
	}
	if n := strings.Count(svg, "<path "); n != 2 {
		t.Errorf("expected 2 county paths, got %d\n got:\n%s", n, svg)
	}
}

func TestBuildSVG_BackgroundVariantsNoRect(t *testing.T) {
	cases := []struct {
		name       string
		background string
	}{
		{"empty", ""},
		{"whitespace padded none", " None "},
		{"uppercase transparent", "TRANSPARENT"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fm, gf := squareCountyFixture()
			fm.Map.Background = tc.background

			svg, err := buildSVG(fm, gf)
			if err != nil {
				t.Fatalf("buildSVG: %v", err)
			}
			if strings.Contains(svg, "<rect ") {
				t.Errorf("background %q should produce no rect\n got:\n%s", tc.background, svg)
			}
		})
	}
}
