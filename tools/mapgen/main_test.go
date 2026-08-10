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

// Characterization: run() regenerates the committed channel SVG
// byte-identically. Mirrors TestRun_RegeneratesCommittedSVG.
func TestRun_RegeneratesCommittedChannelSVG(t *testing.T) {
	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Fatalf("findRepoRoot: %v", err)
	}
	outPath := filepath.Join(repoRoot, channelOutputFile)

	original, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read committed channel SVG: %v", err)
	}
	t.Cleanup(func() {
		if err := os.WriteFile(outPath, original, 0644); err != nil {
			t.Errorf("restore committed channel SVG: %v", err)
		}
	})

	if err := run(); err != nil {
		t.Fatalf("run: %v", err)
	}

	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read regenerated channel SVG: %v", err)
	}
	if !bytes.Equal(got, original) {
		t.Errorf("regenerated channel SVG differs from committed (%d bytes vs %d bytes)", len(got), len(original))
	}
}

// aym-nfs.2: run() writes the channel map alongside the regional map.
// Removes the generated file in cleanup if it did not exist before the test,
// so a failed assertion cannot leave the working tree dirty.
func TestRun_WritesChannelSVG(t *testing.T) {
	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Fatalf("findRepoRoot: %v", err)
	}
	outPath := filepath.Join(repoRoot, channelOutputFile)

	original, readErr := os.ReadFile(outPath)
	t.Cleanup(func() {
		if readErr != nil {
			os.Remove(outPath)
			return
		}
		if err := os.WriteFile(outPath, original, 0644); err != nil {
			t.Errorf("restore channel SVG: %v", err)
		}
	})

	if err := run(); err != nil {
		t.Fatalf("run: %v", err)
	}

	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read channel SVG: %v", err)
	}
	if !bytes.HasPrefix(got, []byte("<svg")) {
		t.Errorf("channel SVG does not start with <svg; got %.40q", got)
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
func squareCountyFixture() (*mapColors, []group, *countyCollection) {
	colors := &mapColors{
		Background:        "none",
		CountyStroke:      "#000000",
		UnassignedCounty:  "#cccccc",
		CountyFillOpacity: 1.0,
		CountyLabel:       "white",
		RegionLabel:       "black",
		RegionLabelHalo:   "white",
	}
	groups := []group{
		{
			Name:     "Test #1",
			Color:    "#ff0000",
			Label:    &regionLabel{X: 100, Y: 200},
			Counties: []string{"Square"},
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
	return colors, groups, gf
}

// groupFixture reuses the square county geography but expresses the
// assignment as a channel-style group, for renderMap tests.
func groupFixture() (*mapColors, []group, *countyCollection) {
	colors, _, gf := squareCountyFixture()
	groups := []group{
		{
			Name:     "MediumFast",
			Color:    "#ff8800",
			Label:    &regionLabel{X: 100, Y: 200},
			Counties: []string{"Square"},
		},
	}
	return colors, groups, gf
}

func TestRenderMap_GroupCountyColored(t *testing.T) {
	colors, groups, gf := groupFixture()
	svg, err := renderMap(colors, groups, gf)
	if err != nil {
		t.Fatalf("renderMap: %v", err)
	}
	if !strings.Contains(svg, `fill="#ff8800"`) {
		t.Errorf("group county missing group fill\n got:\n%s", svg)
	}
}

func TestRenderMap_CatchallColorsUnmatched(t *testing.T) {
	colors, _, gf := twoCountyFixture() // "Square" assigned below; "Tri" unmatched
	groups := []group{
		{Name: "MediumFast", Color: "#ff8800", Counties: []string{"Square"}},
		{Name: "LongFast", Color: "#3b82c4", Catchall: true},
	}
	svg, err := renderMap(colors, groups, gf)
	if err != nil {
		t.Fatalf("renderMap: %v", err)
	}
	if !strings.Contains(svg, `fill="#ff8800"`) {
		t.Errorf("assigned county missing group fill\n got:\n%s", svg)
	}
	if !strings.Contains(svg, `fill="#3b82c4"`) {
		t.Errorf("unmatched county missing catchall fill\n got:\n%s", svg)
	}
	if strings.Contains(svg, `fill="#cccccc"`) {
		t.Errorf("catchall present: no county should use unassigned_county\n got:\n%s", svg)
	}
}

func TestRenderMap_LabelFontsize(t *testing.T) {
	colors, groups, gf := groupFixture()
	groups[0].Label.FontSize = 14
	svg, err := renderMap(colors, groups, gf)
	if err != nil {
		t.Fatalf("renderMap: %v", err)
	}
	if !strings.Contains(svg, `font-size="14"`) {
		t.Errorf("non-default fontsize missing from group label\n got:\n%s", svg)
	}
	if strings.Contains(svg, `font-size="21"`) {
		t.Errorf("default fontsize should not appear when fontsize is set\n got:\n%s", svg)
	}
}

func TestRenderMap_MultiLineGroupName(t *testing.T) {
	colors, groups, gf := groupFixture()
	groups[0].Name = "LongFast\nMediumFast"
	svg, err := renderMap(colors, groups, gf)
	if err != nil {
		t.Fatalf("renderMap: %v", err)
	}
	if !strings.Contains(svg, ">LongFast</tspan>") {
		t.Errorf("missing first line tspan\n got:\n%s", svg)
	}
	if !strings.Contains(svg, ">MediumFast</tspan>") {
		t.Errorf("missing second line tspan\n got:\n%s", svg)
	}
}

func TestRenderMap_SquareCounty(t *testing.T) {
	colors, groups, gf := squareCountyFixture()
	svg, err := renderMap(colors, groups, gf)
	if err != nil {
		t.Fatalf("renderMap: %v", err)
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

	// Region label from group config.
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

func TestRenderMap_BackgroundRect(t *testing.T) {
	colors, groups, gf := squareCountyFixture()
	colors.Background = "#112233"

	svg, err := renderMap(colors, groups, gf)
	if err != nil {
		t.Fatalf("renderMap: %v", err)
	}
	want := `<rect width="1200" height="1000" fill="#112233"/>`
	if !strings.Contains(svg, want) {
		t.Errorf("SVG missing background rect %q\n got:\n%s", want, svg)
	}
}

func TestRenderMap_UnassignedCounty(t *testing.T) {
	colors, _, gf := squareCountyFixture()
	groups := []group{} // no group assignment

	svg, err := renderMap(colors, groups, gf)
	if err != nil {
		t.Fatalf("renderMap: %v", err)
	}
	if !strings.Contains(svg, `fill="#cccccc"`) {
		t.Errorf("expected unassigned fill #cccccc\n got:\n%s", svg)
	}
}

func TestRenderMap_NilGroupLabelSkipped(t *testing.T) {
	colors, groups, gf := squareCountyFixture()
	groups[0].Label = nil

	svg, err := renderMap(colors, groups, gf)
	if err != nil {
		t.Fatalf("renderMap: %v", err)
	}
	if strings.Contains(svg, "Test #1") {
		t.Errorf("expected nil region label to be skipped\n got:\n%s", svg)
	}
}

func TestRenderMap_RejectsInvalidOpacity(t *testing.T) {
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
			colors, groups, gf := squareCountyFixture()
			colors.CountyFillOpacity = tc.opacity
			_, err := renderMap(colors, groups, gf)
			if err == nil {
				t.Fatal("expected error for invalid county_fill_opacity")
			}
			if !strings.Contains(err.Error(), "county_fill_opacity") {
				t.Errorf("error %q should mention county_fill_opacity", err)
			}
		})
	}
}

func TestRenderMap_EmptyGeoJSON(t *testing.T) {
	colors, groups, _ := squareCountyFixture()
	gf := &countyCollection{Type: "FeatureCollection", Features: nil}
	_, err := renderMap(colors, groups, gf)
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
func multiPolygonFixture() (*mapColors, []group, *countyCollection) {
	colors := &mapColors{
		Background:        "none",
		CountyStroke:      "#000000",
		UnassignedCounty:  "#cccccc",
		CountyFillOpacity: 1.0,
		CountyLabel:       "white",
		RegionLabel:       "black",
		RegionLabelHalo:   "white",
	}
	groups := []group{
		{
			Name:     "Test #1",
			Color:    "#ff0000",
			Label:    &regionLabel{X: 100, Y: 200},
			Counties: []string{"Multi"},
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
	return colors, groups, gf
}

func TestRenderMap_MultiPolygon(t *testing.T) {
	colors, groups, gf := multiPolygonFixture()
	svg, err := renderMap(colors, groups, gf)
	if err != nil {
		t.Fatalf("renderMap: %v", err)
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
func twoCountyFixture() (*mapColors, []group, *countyCollection) {
	colors := &mapColors{
		Background:        "none",
		CountyStroke:      "#000000",
		UnassignedCounty:  "#cccccc",
		CountyFillOpacity: 1.0,
		CountyLabel:       "white",
		RegionLabel:       "black",
		RegionLabelHalo:   "white",
	}
	groups := []group{
		{
			Name:     "Test #1",
			Color:    "#ff0000",
			Label:    &regionLabel{X: 100, Y: 200},
			Counties: []string{"Square"},
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
	return colors, groups, gf
}

func TestRenderMap_MixedAssignment(t *testing.T) {
	colors, groups, gf := twoCountyFixture()
	svg, err := renderMap(colors, groups, gf)
	if err != nil {
		t.Fatalf("renderMap: %v", err)
	}

	// Assigned county gets group color.
	if !strings.Contains(svg, `fill="#ff0000"`) {
		t.Errorf("assigned county missing group fill\n got:\n%s", svg)
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

func TestRenderMap_TransparentBackground(t *testing.T) {
	colors, groups, gf := squareCountyFixture()
	colors.Background = "transparent"

	svg, err := renderMap(colors, groups, gf)
	if err != nil {
		t.Fatalf("renderMap: %v", err)
	}
	if strings.Contains(svg, "<rect ") {
		t.Errorf("expected no background rect for background=transparent\n got:\n%s", svg)
	}
}

func TestRenderMap_UnsupportedGeometryType(t *testing.T) {
	colors, groups, _ := squareCountyFixture()
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
	_, err := renderMap(colors, groups, gf)
	if err == nil {
		t.Fatal("expected error for unsupported geometry type")
	}
	if !strings.Contains(err.Error(), "unsupported geometry type") {
		t.Errorf("error %q should mention unsupported geometry type", err)
	}
}

func TestWarnDuplicateCounties_CrossRegion(t *testing.T) {
	groups := []group{
		{Name: "Alpha", Color: "#aa0000", Counties: []string{"Square"}},
		{Name: "Beta", Color: "#00bb00", Counties: []string{"Square"}},
	}
	var buf bytes.Buffer
	warnDuplicateCounties(groups, &buf)
	got := buf.String()
	if !strings.Contains(got, "Square") {
		t.Errorf("warning should name county Square; got %q", got)
	}
	if !strings.Contains(got, "Alpha") || !strings.Contains(got, "Beta") {
		t.Errorf("warning should name both regions; got %q", got)
	}
}

func TestWarnDuplicateCounties_WithinRegion(t *testing.T) {
	groups := []group{
		{Name: "Solo", Color: "#aa0000", Counties: []string{"Square", "Square"}},
	}
	var buf bytes.Buffer
	warnDuplicateCounties(groups, &buf)
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
	groups := []group{
		{Name: "Alpha", Color: "#aa0000", Counties: []string{"Square"}},
		{Name: "Beta", Color: "#00bb00", Counties: []string{"Tri"}},
	}
	var buf bytes.Buffer
	warnDuplicateCounties(groups, &buf)
	if buf.Len() != 0 {
		t.Errorf("expected no warning; got %q", buf.String())
	}
}

func TestWarnUnmatchedCounties_SidecarCountyMissingFromGeoJSON(t *testing.T) {
	_, groups, gf := squareCountyFixture()
	groups[0].Counties = append(groups[0].Counties, "Phantom")

	var buf bytes.Buffer
	warnUnmatchedCounties(groups, gf, &buf)
	got := buf.String()

	if !strings.Contains(got, "Phantom") {
		t.Errorf("warning should name county Phantom; got %q", got)
	}
	if !strings.Contains(got, "Test") {
		t.Errorf("warning should name group Test; got %q", got)
	}
}

func TestWarnUnmatchedCounties_GeoJSONFeatureMatchedByNoCounty(t *testing.T) {
	_, groups, gf := twoCountyFixture() // "Square" assigned; "Tri" is GeoJSON-only

	var buf bytes.Buffer
	warnUnmatchedCounties(groups, gf, &buf)
	got := buf.String()

	if !strings.Contains(got, "Tri") {
		t.Errorf("warning should name TIGERNAME Tri; got %q", got)
	}
	if strings.Contains(got, "Square") {
		t.Errorf("matched county Square should not warn; got %q", got)
	}
}

func TestWarnUnmatchedCounties_CurrentRepoRegionalDataIsClean(t *testing.T) {
	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Fatalf("findRepoRoot: %v", err)
	}
	sc, err := parseSidecar(filepath.Join(repoRoot, regionsFile))
	if err != nil {
		t.Fatalf("parseSidecar regions.yaml: %v", err)
	}
	gf, err := parseGeoJSON(filepath.Join(repoRoot, geoJSONFile))
	if err != nil {
		t.Fatalf("parseGeoJSON: %v", err)
	}

	var buf bytes.Buffer
	warnUnmatchedCounties(sc.Regions, gf, &buf)

	if buf.Len() != 0 {
		t.Errorf("expected no warnings for current repo data; got:\n%s", buf.String())
	}
}

// aym-nfs.2: a group set containing a catchall group suppresses the
// GeoJSON-side "matched by no county list" warning — uncovered counties are
// intentional there. The county-side direction (typo detection) always fires.
func TestWarnUnmatchedCounties_CatchallSuppressesGeoJSONSide(t *testing.T) {
	_, _, gf := twoCountyFixture() // "Tri" is GeoJSON-only
	groups := []group{{
		Name:     "Test",
		Color:    "#aa0000",
		Catchall: true,
		Counties: []string{"Square"},
	}}

	var buf bytes.Buffer
	warnUnmatchedCounties(groups, gf, &buf)

	if strings.Contains(buf.String(), "Tri") {
		t.Errorf("catchall group set should suppress GeoJSON-side warning; got %q", buf.String())
	}
}

// aym-nfs.2: committed channels.yaml produces no duplicate or unmatched
// county warnings; mirrors TestWarnUnmatchedCounties_CurrentRepoDataIsClean.
func TestSidecar_CurrentRepoChannelsClean(t *testing.T) {
	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Fatalf("findRepoRoot: %v", err)
	}
	sc, err := parseSidecar(filepath.Join(repoRoot, channelsFile))
	if err != nil {
		t.Fatalf("parseSidecar: %v", err)
	}
	gf, err := parseGeoJSON(filepath.Join(repoRoot, geoJSONFile))
	if err != nil {
		t.Fatalf("parseGeoJSON: %v", err)
	}

	var buf bytes.Buffer
	warnDuplicateCounties(sc.Channels, &buf)
	warnUnmatchedCounties(sc.Channels, gf, &buf)

	if buf.Len() != 0 {
		t.Errorf("expected no warnings for committed channels.yaml; got:\n%s", buf.String())
	}
}

// aym-nfs.2: every committed channels[].color must be a well-formed color
// string; mirrors the regions.yaml coverage in
// TestParseSidecar_RegionsCurrentRepoDataIsClean.
func TestSidecar_ChannelColorFormat(t *testing.T) {
	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Fatalf("findRepoRoot: %v", err)
	}
	sc, err := parseSidecar(filepath.Join(repoRoot, channelsFile))
	if err != nil {
		t.Fatalf("parseSidecar: %v", err)
	}

	for _, ch := range sc.Channels {
		if !isWellFormedRegionColor(ch.Color) {
			t.Errorf("channel %q color %q is not a well-formed color string", ch.Name, ch.Color)
		}
	}
}

// isWellFormedRegionColor reports whether s is a well-formed regions[].color
// value.
func isWellFormedRegionColor(s string) bool {
	if s == "" {
		return true
	}
	if len(s) != 4 && len(s) != 7 || s[0] != '#' {
		return false
	}
	for _, c := range s[1:] {
		if !('0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F') {
			return false
		}
	}
	return true
}

// aym-jha: regions[].color must be a well-formed color string; the palette
// itself is no longer pinned.
func TestIsWellFormedRegionColor(t *testing.T) {
	cases := []struct {
		name  string
		color string
		want  bool
	}{
		{"rrggbb lowercase", "#ff0000", true},
		{"empty falls back to unassigned_county", "", true},
		{"rgb shorthand", "#f00", true},
		{"rrggbb uppercase", "#FF0000", true},
		{"rrggbb mixed case", "#aB65cD", true},
		{"missing hash", "ff0000", false},
		{"too short", "#12", false},
		{"length five", "#1234", false},
		{"length six", "#12345", false},
		{"too long", "#1234567", false},
		{"non-hex digit", "#gg0000", false},
		{"trailing non-hex", "#12345g", false},
		{"css named color", "red", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isWellFormedRegionColor(tc.color); got != tc.want {
				t.Errorf("isWellFormedRegionColor(%q) = %v, want %v", tc.color, got, tc.want)
			}
		})
	}
}

func TestRenderMap_PhantomCountyStillGenerates(t *testing.T) {
	colors, groups, gf := squareCountyFixture()
	groups[0].Counties = append(groups[0].Counties, "Phantom")

	svg, err := renderMap(colors, groups, gf)
	if err != nil {
		t.Fatalf("renderMap with phantom county: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Errorf("expected SVG output despite phantom county\n got:\n%s", svg)
	}
}

func TestWarnUnmatchedCounties_AllNamesMatch(t *testing.T) {
	_, groups, gf := squareCountyFixture()

	var buf bytes.Buffer
	warnUnmatchedCounties(groups, gf, &buf)

	if buf.Len() != 0 {
		t.Errorf("expected no warning when all names match; got %q", buf.String())
	}
}

func TestRenderMap_DuplicateCountyLastGroupWins(t *testing.T) {
	colors, _, gf := squareCountyFixture()
	groups := []group{
		{Name: "Alpha #1", Color: "#aa0000", Counties: []string{"Square"}},
		{Name: "Beta #2", Color: "#00bb00", Counties: []string{"Square"}},
	}
	svg, err := renderMap(colors, groups, gf)
	if err != nil {
		t.Fatalf("renderMap: %v", err)
	}
	if !strings.Contains(svg, `fill="#00bb00"`) {
		t.Errorf("expected later group color #00bb00; got:\n%s", svg)
	}
	if strings.Contains(svg, `fill="#aa0000"`) {
		t.Errorf("earlier group color should not win; got:\n%s", svg)
	}
}

// The README is the mapgen user documentation (aym-emi); it must cover the
// tool's purpose, how to run it, the sidecar contract, and the rule that
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
		"regions.yaml",
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
		"channels.yaml",
		"channel-lora-settings.svg",
		"catchall",
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
	groups := []group{
		{Name: "Alpha", Color: "#aa0000", Counties: []string{"De Soto"}},
		{Name: "Beta", Color: "#00bb00", Counties: []string{"desoto"}},
	}
	var buf bytes.Buffer
	warnDuplicateCounties(groups, &buf)
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

func TestParseSidecar_Valid(t *testing.T) {
	sc, err := unmarshalSidecar([]byte(`map:
  background: "#112233"
  county_stroke: "#000000"
  unassigned_county: "#cccccc"
  county_fill_opacity: 0.85
  county_label: white
  region_label: black
  region_label_halo: "#f8f8f8"
`))
	if err != nil {
		t.Fatalf("unmarshalSidecar: %v", err)
	}
	m := sc.Map
	if m.Background != "#112233" {
		t.Errorf("background = %q, want #112233", m.Background)
	}
	if m.CountyStroke != "#000000" {
		t.Errorf("county_stroke = %q, want #000000", m.CountyStroke)
	}
	if m.UnassignedCounty != "#cccccc" {
		t.Errorf("unassigned_county = %q, want #cccccc", m.UnassignedCounty)
	}
	if m.CountyFillOpacity != 0.85 {
		t.Errorf("county_fill_opacity = %v, want 0.85", m.CountyFillOpacity)
	}
	if m.CountyLabel != "white" {
		t.Errorf("county_label = %q, want white", m.CountyLabel)
	}
	if m.RegionLabel != "black" {
		t.Errorf("region_label = %q, want black", m.RegionLabel)
	}
	if m.RegionLabelHalo != "#f8f8f8" {
		t.Errorf("region_label_halo = %q, want #f8f8f8", m.RegionLabelHalo)
	}
}

func TestParseSidecar_Channels(t *testing.T) {
	sc, err := unmarshalSidecar([]byte(`map:
  background: none
channels:
  - name: MediumFast
    color: "#ff8800"
    label:
      x: 700
      y: 800
    counties:
      - Miami-Dade
      - Broward
  - name: LongFast
    color: "#0066ff"
    counties:
      - Alachua
`))
	if err != nil {
		t.Fatalf("unmarshalSidecar: %v", err)
	}
	if len(sc.Channels) != 2 {
		t.Fatalf("channels = %d, want 2", len(sc.Channels))
	}
	ch := sc.Channels[0]
	if ch.Name != "MediumFast" || ch.Color != "#ff8800" {
		t.Errorf("channel 0 = %+v", ch)
	}
	if ch.Label == nil || ch.Label.X != 700 || ch.Label.Y != 800 {
		t.Errorf("channel 0 label = %+v, want (700, 800)", ch.Label)
	}
	if len(ch.Counties) != 2 || ch.Counties[0] != "Miami-Dade" || ch.Counties[1] != "Broward" {
		t.Errorf("channel 0 counties = %v, want [Miami-Dade Broward]", ch.Counties)
	}
	if sc.Channels[1].Label != nil {
		t.Errorf("channel without label key should have nil Label, got %+v", sc.Channels[1].Label)
	}
}

func TestParseSidecar_LabelFontsize(t *testing.T) {
	sc, err := unmarshalSidecar([]byte(`channels:
  - name: MediumFast
    label:
      x: 700
      y: 800
      fontsize: 14
  - name: LongFast
    label:
      x: 500
      y: 400
`))
	if err != nil {
		t.Fatalf("unmarshalSidecar: %v", err)
	}
	if got := sc.Channels[0].Label.FontSize; got != 14 {
		t.Errorf("fontsize = %v, want 14", got)
	}
	if got := sc.Channels[1].Label.FontSize; got != 0 {
		t.Errorf("omitted fontsize = %v, want 0", got)
	}
}

func TestParseSidecar_Catchall(t *testing.T) {
	sc, err := unmarshalSidecar([]byte(`channels:
  - name: LongFast
    catchall: true
  - name: MediumFast
`))
	if err != nil {
		t.Fatalf("unmarshalSidecar: %v", err)
	}
	if !sc.Channels[0].Catchall {
		t.Error("catchall: true should parse as true")
	}
	if sc.Channels[1].Catchall {
		t.Error("omitted catchall should default to false")
	}
}

func TestParseSidecar_TwoCatchallsError(t *testing.T) {
	_, err := unmarshalSidecar([]byte(`channels:
  - name: LongFast
    catchall: true
  - name: MediumFast
    catchall: true
`))
	if err == nil {
		t.Fatal("expected error for two catchall groups")
	}
	if !strings.Contains(err.Error(), "catchall") {
		t.Errorf("error %q should mention catchall", err)
	}
}

func TestParseSidecar_MissingFile(t *testing.T) {
	_, err := parseSidecar(filepath.Join(t.TempDir(), "nope.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestParseSidecar_MalformedYAML(t *testing.T) {
	_, err := unmarshalSidecar([]byte("map: [unclosed\n"))
	if err == nil {
		t.Fatal("expected error for malformed YAML")
	}
	if !strings.Contains(err.Error(), "parsing YAML") {
		t.Errorf("error %q should mention parsing YAML", err)
	}
}

// aym-nfs.7: unmarshalSidecar parses a regions.yaml-style document,
// populating Regions on the shared sidecar struct.
func TestParseSidecar_Regions(t *testing.T) {
	sc, err := unmarshalSidecar([]byte(`map:
  background: "#112233"
  county_stroke: "#000000"
  unassigned_county: "#cccccc"
  county_fill_opacity: 0.85
  county_label: white
  region_label: black
  region_label_halo: "#f8f8f8"
regions:
  - name: "North #1"
    color: "#ff0000"
    label:
      x: 100
      y: 200
    counties:
      - Square
      - Tri
  - name: "South #2"
    color: "#00bb00"
    counties:
      - Other
`))
	if err != nil {
		t.Fatalf("unmarshalSidecar: %v", err)
	}
	if len(sc.Regions) != 2 {
		t.Fatalf("regions = %d, want 2", len(sc.Regions))
	}
	r := sc.Regions[0]
	if r.Name != "North #1" || r.Color != "#ff0000" {
		t.Errorf("region 0 = %+v", r)
	}
	if r.Label == nil || r.Label.X != 100 || r.Label.Y != 200 {
		t.Errorf("region 0 label = %+v, want (100, 200)", r.Label)
	}
	if len(r.Counties) != 2 || r.Counties[0] != "Square" || r.Counties[1] != "Tri" {
		t.Errorf("region 0 counties = %v, want [Square Tri]", r.Counties)
	}
	if sc.Regions[1].Label != nil {
		t.Errorf("region without label key should have nil Label, got %+v", sc.Regions[1].Label)
	}
	if sc.Regions[1].Catchall {
		t.Error("region should default Catchall to false")
	}
	if sc.Channels != nil {
		t.Errorf("regions YAML should not populate Channels, got %v", sc.Channels)
	}
}

// aym-nfs.7: the committed regions.yaml parses cleanly via the shared
// sidecar struct and unmarshalSidecar function.
func TestParseSidecar_RegionsCurrentRepoDataIsClean(t *testing.T) {
	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Fatalf("findRepoRoot: %v", err)
	}
	sc, err := parseSidecar(filepath.Join(repoRoot, regionsFile))
	if err != nil {
		t.Fatalf("parseSidecar regions.yaml: %v", err)
	}
	if len(sc.Regions) == 0 {
		t.Fatal("committed regions.yaml should contain at least one region")
	}
	for _, r := range sc.Regions {
		if !isWellFormedRegionColor(r.Color) {
			t.Errorf("region %q color %q is not a well-formed color string", r.Name, r.Color)
		}
		if r.Catchall {
			t.Errorf("region %q should not set catchall", r.Name)
		}
	}
}

func TestParseSidecar_CurrentRepoDataIsClean(t *testing.T) {
	repoRoot, err := findRepoRoot()
	if err != nil {
		t.Fatalf("findRepoRoot: %v", err)
	}
	sc, err := parseSidecar(filepath.Join(repoRoot, channelsFile))
	if err != nil {
		t.Fatalf("parseSidecar channels.yaml: %v", err)
	}
	catchalls := 0
	for _, ch := range sc.Channels {
		if ch.Catchall {
			catchalls++
		}
	}
	if catchalls != 1 {
		t.Errorf("channels.yaml catchall groups = %d, want exactly 1", catchalls)
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
func TestRenderMap_DegenerateGeometryGetsNoLabel(t *testing.T) {
	colors, groups, gf := squareCountyFixture()
	gf.Features = append(gf.Features, &geojson.FeatureOf[countyProps]{
		Type:       "Feature",
		Properties: countyProps{TigerName: "Line"},
		Geometry: orb.Polygon{
			{{-1, 0}, {0, 0}, {1, 0}, {-1, 0}},
		},
	})

	svg, err := renderMap(colors, groups, gf)
	if err != nil {
		t.Fatalf("renderMap: %v", err)
	}
	if strings.Contains(svg, ">Line</text>") {
		t.Errorf("zero-area county should get no label\n got:\n%s", svg)
	}
	if n := strings.Count(svg, "<path "); n != 2 {
		t.Errorf("expected 2 county paths, got %d\n got:\n%s", n, svg)
	}
}

func TestRenderMap_BackgroundVariantsNoRect(t *testing.T) {
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
			colors, groups, gf := squareCountyFixture()
			colors.Background = tc.background

			svg, err := renderMap(colors, groups, gf)
			if err != nil {
				t.Fatalf("renderMap: %v", err)
			}
			if strings.Contains(svg, "<rect ") {
				t.Errorf("background %q should produce no rect\n got:\n%s", tc.background, svg)
			}
		})
	}
}
