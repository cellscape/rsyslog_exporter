package main

import (
	"fmt"
	"testing"
)

func testHelper(t *testing.T, line []byte, testCase []*testUnit) {
	exporter := newRsyslogExporter()
	exporter.handleStatLine(line)

	for _, k := range exporter.keys() {
		t.Logf("have key: '%s'", k)
	}

	for _, item := range testCase {
		p, err := exporter.get(item.key())
		if err != nil {
			t.Error(err)
		}

		if want, got := item.Val, p.promValue(); want != got {
			t.Errorf("%s: want '%f', got '%f'", item.Name, want, got)
		}
	}

	exporter.handleStatLine(line)

	for _, item := range testCase {
		p, err := exporter.get(item.key())
		if err != nil {
			t.Error(err)
		}

		var wanted float64
		switch p.Type {
		case counter:
			wanted = item.Val
		case gauge:
			wanted = item.Val
		default:
			t.Errorf("%d is not a valid metric type", p.Type)
			continue
		}

		if want, got := wanted, p.promValue(); want != got {
			t.Errorf("%s: want '%f', got '%f'", item.Name, want, got)
		}
	}
}

type testUnit struct {
	Name string
	Val  float64
	LabelValue string
}

func (t *testUnit) key() string {
	return fmt.Sprintf("%s.%s", t.Name, t.LabelValue)
}

func TestHandleLineWithAction(t *testing.T) {
	tests := []*testUnit{
		&testUnit{
			Name: "action_processed",
			Val:  100000,
			LabelValue: "test_action",
		},
		&testUnit{
			Name: "action_failed",
			Val:  2,
			LabelValue: "test_action",
		},
		&testUnit{
			Name: "action_suspended",
			Val:  1,
			LabelValue: "test_action",
		},
		&testUnit{
			Name: "action_suspended_duration",
			Val:  1000,
			LabelValue: "test_action",
		},
		&testUnit{
			Name: "action_resumed",
			Val:  1,
			LabelValue: "test_action",
		},
	}

	actionLog := []byte(`{"name":"test_action","processed":100000,"failed":2,"suspended":1,"suspended.duration":1000,"resumed":1}`)
	testHelper(t, actionLog, tests)
}

func TestHandleLineWithResource(t *testing.T) {
	tests := []*testUnit{
		&testUnit{
			Name: "resource_utime",
			Val:  10,
			LabelValue: "resource-usage",
		},
		&testUnit{
			Name: "resource_stime",
			Val:  20,
			LabelValue: "resource-usage",
		},
		&testUnit{
			Name: "resource_maxrss",
			Val:  30,
			LabelValue: "resource-usage",
		},
		&testUnit{
			Name: "resource_minflt",
			Val:  40,
			LabelValue: "resource-usage",
		},
		&testUnit{
			Name: "resource_majflt",
			Val:  50,
			LabelValue: "resource-usage",
		},
		&testUnit{
			Name: "resource_inblock",
			Val:  60,
			LabelValue: "resource-usage",
		},
		&testUnit{
			Name: "resource_oublock",
			Val:  70,
			LabelValue: "resource-usage",
		},
		&testUnit{
			Name: "resource_nvcsw",
			Val:  80,
			LabelValue: "resource-usage",
		},
		&testUnit{
			Name: "resource_nivcsw",
			Val:  90,
			LabelValue: "resource-usage",
		},
	}

	resourceLog := []byte(`{"name":"resource-usage","utime":10,"stime":20,"maxrss":30,"minflt":40,"majflt":50,"inblock":60,"oublock":70,"nvcsw":80,"nivcsw":90}`)
	testHelper(t, resourceLog, tests)
}

func TestHandleLineWithInput(t *testing.T) {
	tests := []*testUnit{
		&testUnit{
			Name: "input_submitted",
			Val:  1000,
			LabelValue: "test_input",
		},
	}

	inputLog := []byte(`{"name":"test_input", "origin":"imuxsock", "submitted":1000}`)
	testHelper(t, inputLog, tests)
}

func TestHandleLineWithQueue(t *testing.T) {
	tests := []*testUnit{
		&testUnit{
			Name: "queue_size",
			Val:  10,
			LabelValue: "main Q",
		},
		&testUnit{
			Name: "queue_enqueued",
			Val:  20,
			LabelValue: "main Q",
		},
		&testUnit{
			Name: "queue_full",
			Val:  30,
			LabelValue: "main Q",
		},
		&testUnit{
			Name: "queue_discarded_full",
			Val:  40,
			LabelValue: "main Q",
		},
		&testUnit{
			Name: "queue_discarded_not_full",
			Val:  50,
			LabelValue: "main Q",
		},
		&testUnit{
			Name: "queue_max_size",
			Val:  60,
			LabelValue: "main Q",
		},
	}

	queueLog = []byte(`{"name":"main Q","size":10,"enqueued":20,"full":30,"discarded.full":40,"discarded.nf":50,"maxqsize":60}`)
	testHelper(t, queueLog, tests)
}

func TestHandleUnknown(t *testing.T) {
	unknownLog := []byte(`{"a":"b"}`)

	exporter := newRsyslogExporter()
	exporter.handleStatLine(unknownLog)

	if want, got := 0, len(exporter.keys()); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}
}
