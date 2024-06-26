package viewer

import (
	"encoding/json"
	"net/http"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

const (
	// VCStack is the name of StackViewer
	VCStack = "stack"
)

// StackViewer collects the stack-stats metrics via `runtime.ReadMemStats()`
type StackViewer struct {
	smgr  *StatsMgr
	graph *charts.Line
}

// NewStackViewer returns the StackViewer instance
// Series: StackSys / StackInuse / MSpanSys / MSpanInuse
func NewStackViewer() Viewer {
	graph := NewBasicView(VCStack)
	graph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Stack"}),
		charts.WithYAxisOpts(opts.YAxis{Name: "Size", AxisLabel: &opts.AxisLabel{Show: true, Formatter: "{value} MB"}}),
	)
	graph.AddSeries("Sys", []opts.LineData{}).
		AddSeries("Inuse", []opts.LineData{}).
		AddSeries("MSpan Sys", []opts.LineData{}).
		AddSeries("MSpan Inuse", []opts.LineData{})

	return &StackViewer{graph: graph}
}

func (vr *StackViewer) SetStatsMgr(smgr *StatsMgr) {
	vr.smgr = smgr
}

func (vr *StackViewer) Name() string {
	return VCStack
}

func (vr *StackViewer) View() *charts.Line {
	return vr.graph
}

func (vr *StackViewer) Serve(w http.ResponseWriter, _ *http.Request) {
	vr.smgr.Tick()

	metrics := Metrics{
		Values: []float64{
			FixedPrecision(float64(memstats.Stats.StackSys)/1024/1024, 2),
			FixedPrecision(float64(memstats.Stats.StackInuse)/1024/1024, 2),
			FixedPrecision(float64(memstats.Stats.MSpanSys)/1024/1024, 2),
			FixedPrecision(float64(memstats.Stats.MSpanInuse)/1024/1024, 2),
		},
		Time: memstats.T,
	}

	bs, _ := json.Marshal(metrics)
	w.Write(bs)
}
