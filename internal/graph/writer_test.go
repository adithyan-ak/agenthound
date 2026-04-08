package graph

import (
	"testing"

	"github.com/adithyan-ak/agenthound/internal/model"
)

func TestEdgeKindEndpointsCoversAllEdgeKinds(t *testing.T) {
	for kind := range model.AllowedEdgeKinds {
		if _, ok := model.EdgeKindEndpoints[kind]; !ok {
			t.Errorf("EdgeKindEndpoints missing entry for edge kind %q", kind)
		}
	}
	for kind := range model.EdgeKindEndpoints {
		if !model.AllowedEdgeKinds[kind] {
			t.Errorf("EdgeKindEndpoints has extra entry for unknown edge kind %q", kind)
		}
	}
}
