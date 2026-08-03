package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
func squareCountyFixture() (*frontMatter, *geoJSON) {
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
				Name:     "Test",
				Color:    "#ff0000",
				Number:   "#1",
				Label:    &regionLabel{X: 100, Y: 200},
				Counties: []string{"Square"},
			},
		},
	}

	// Closed ring, CCW: SW → SE → NE → NW → SW
	coords, _ := json.Marshal([][][]float64{
		{
			{-1, -1},
			{1, -1},
			{1, 1},
			{-1, 1},
			{-1, -1},
		},
	})
	gf := &geoJSON{
		Type: "FeatureCollection",
		Features: []feature{
			{
				Type:       "Feature",
				Properties: properties{TigerName: "Square"},
				Geometry: geometry{
					Type:        "Polygon",
					Coordinates: coords,
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
	gf := &geoJSON{Type: "FeatureCollection", Features: nil}
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
// (area=0.045, centroid=(0.65,0.6)). The centroid-picking logic selects the
// larger ring, so the label stays at the square's centroid (600,500).
// Bbox is still [-1,1]² so the transform is unchanged from squareCountyFixture.
func multiPolygonFixture() (*frontMatter, *geoJSON) {
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
				Name:     "Test",
				Color:    "#ff0000",
				Number:   "#1",
				Label:    &regionLabel{X: 100, Y: 200},
				Counties: []string{"Multi"},
			},
		},
	}

	// MultiPolygon: two polygons, one ring each.
	// Polygon 1 (large): square [-1,-1]→[1,1]
	// Polygon 2 (small): triangle (0.5,0.5)→(0.8,0.5)→(0.65,0.8)
	ring1 := [][]float64{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}}
	ring2 := [][]float64{{0.5, 0.5}, {0.8, 0.5}, {0.65, 0.8}, {0.5, 0.5}}
	coords, _ := json.Marshal([][][][]float64{{ring1}, {ring2}})
	gf := &geoJSON{
		Type: "FeatureCollection",
		Features: []feature{
			{
				Type:       "Feature",
				Properties: properties{TigerName: "Multi"},
				Geometry: geometry{
					Type:        "MultiPolygon",
					Coordinates: coords,
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

	// Label at largest ring's centroid (square centre (0,0) → (600,500)).
	if !strings.Contains(svg, `<text x="600.0" y="500.0" font-size="9" fill="white">Multi</text>`) {
		t.Errorf("MultiPolygon label at wrong position\n got:\n%s", svg)
	}
}

// twoCountyFixture is a single square (assigned to "Test") and a small
// triangle ("Tri", unassigned). Bbox stays [-1,1]², transform unchanged.
func twoCountyFixture() (*frontMatter, *geoJSON) {
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
				Name:     "Test",
				Color:    "#ff0000",
				Number:   "#1",
				Label:    &regionLabel{X: 100, Y: 200},
				Counties: []string{"Square"},
			},
		},
	}

	squareCoords, _ := json.Marshal([][][]float64{
		{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}, {-1, -1}},
	})
	triCoords, _ := json.Marshal([][][]float64{
		{{0.5, 0.5}, {0.8, 0.5}, {0.65, 0.8}, {0.5, 0.5}},
	})
	gf := &geoJSON{
		Type: "FeatureCollection",
		Features: []feature{
			{
				Type:       "Feature",
				Properties: properties{TigerName: "Square"},
				Geometry:   geometry{Type: "Polygon", Coordinates: squareCoords},
			},
			{
				Type:       "Feature",
				Properties: properties{TigerName: "Tri"},
				Geometry:   geometry{Type: "Polygon", Coordinates: triCoords},
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
	coords, _ := json.Marshal([]float64{0, 0})
	gf := &geoJSON{
		Type: "FeatureCollection",
		Features: []feature{
			{
				Type:       "Feature",
				Properties: properties{TigerName: "Bad"},
				Geometry:   geometry{Type: "Point", Coordinates: coords},
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
		{Name: "Alpha", Color: "#aa0000", Number: "#1", Counties: []string{"Square"}},
		{Name: "Beta", Color: "#00bb00", Number: "#2", Counties: []string{"Square"}},
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
