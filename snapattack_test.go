package snapattack

import (
	"context"
	"encoding/json"
	"os"
	"testing"
)

func TestSignatureExport(t *testing.T) {
	c := NewClient(os.Getenv("SNAPATTACK_API_KEY"))

	filter := Filter{}
	err := json.Unmarshal([]byte(`{
		"op": "and",
		"items": [
		  {
			"op": "or",
			"items": [
			  {
				"op": "equals",
				"field": "gap",
				"value": "true"
			  },
			  {
				"op": "equals",
				"field": "validated",
				"value": "true"
			  }
			]
		  },
		  {
			"op": "in",
			"field": "visibility",
			"value": [
			  "Published"
			]
		  },
		  {
			"op": "for_each",
			"items": [
			  {
				"op": "in",
				"field": "ranks.rank",
				"value": [
				  "Highest",
				  "High",
				  "Medium"
				]
			  },
			  {
				"op": "equals",
				"field": "ranks.organization_id",
				"value": 168,
				"case_sensitive": true
			  }
			]
		  },
		  {
			"op": "for_each",
			"items": [
			  {
				"op": "in",
				"field": "severities.severity_name",
				"value": [
				  "Highest",
				  "High",
				  "Medium"
				]
			  },
			  {
				"op": "equals",
				"field": "severities.organization_id",
				"value": 168,
				"case_sensitive": true
			  }
			]
		  },
		  {
			"op": "any",
			"field": "analytic_compilation_targets",
			"value": [
			  83
			]
		  }
		]
	  }`), &filter)
	if err != nil {
		t.Errorf("json.Unmarshal() filter: %v", err)
		return
	}

	signatures, err := c.ExportSignatures(context.Background(), filter, Targets.LimaCharlie)
	if err != nil {
		t.Errorf("ExportSignatures(): %v", err)
		return
	}
	if len(signatures) < 5 {
		t.Errorf("unexpected signatures: %#v", signatures)
	}
}
