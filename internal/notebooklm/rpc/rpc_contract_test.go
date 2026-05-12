package rpc

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRPCConstantsMatchNotebookLMPyV041(t *testing.T) {
	tests := map[string]string{
		"GET_LAST_CONVERSATION_ID": RPCGetLastConversationID,
		"GET_CONVERSATION_TURNS":   RPCGetConversationTurns,
		"GET_SUGGESTED_REPORTS":    RPCGenerateReportSuggestions,
		"POLL_RESEARCH":            RPCPollResearch,
		"GET_USER_TIER":            RPCGetUserTier,
		"GET_INTERACTIVE_HTML":     RPCGetInteractiveHTML,
		"REVISE_SLIDE":             RPCReviseSlide,
		"REMOVE_RECENTLY_VIEWED":   RPCRemoveRecentlyViewed,
		"CHECK_SOURCE_FRESHNESS":   RPCCheckSourceFreshness,
		"ADD_SOURCE_FILE":          RPCAddFileSource,
	}

	want := map[string]string{
		"GET_LAST_CONVERSATION_ID": "hPTbtc",
		"GET_CONVERSATION_TURNS":   "khqZz",
		"GET_SUGGESTED_REPORTS":    "ciyUvf",
		"POLL_RESEARCH":            "e3bVqc",
		"GET_USER_TIER":            "ozz5Z",
		"GET_INTERACTIVE_HTML":     "v9rmvd",
		"REVISE_SLIDE":             "KmcKPe",
		"REMOVE_RECENTLY_VIEWED":   "fejl7e",
		"CHECK_SOURCE_FRESHNESS":   "yR9Yof",
		"ADD_SOURCE_FILE":          "o4cbdc",
	}

	for name, got := range tests {
		if got != want[name] {
			t.Fatalf("%s = %q, want %q", name, got, want[name])
		}
	}
}

func TestGeneratedReportSuggestionsUsesCurrentRPCContract(t *testing.T) {
	root := filepath.Join("..", "..", "..")

	generatedClient := readRepoFile(t, root, "gen", "service", "LabsTailwindOrchestrationService_client.go")
	if !strings.Contains(generatedClient, `ID:         "ciyUvf"`) {
		t.Fatalf("generated service client does not use ciyUvf for GenerateReportSuggestions")
	}
	if strings.Contains(generatedClient, `ID:         "GHsKob"`) {
		t.Fatalf("generated service client still references stale GHsKob report-suggestions RPC")
	}

	generatedEncoder := readRepoFile(t, root, "gen", "method", "LabsTailwindOrchestrationService_GenerateReportSuggestions_encoder.go")
	if !strings.Contains(generatedEncoder, `[[2], %project_id%]`) {
		t.Fatalf("generated report-suggestions encoder does not include the current [[2], notebook_id] argument shape")
	}
}

func readRepoFile(t *testing.T, root string, pathParts ...string) string {
	t.Helper()

	path := filepath.Join(append([]string{root}, pathParts...)...)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}
