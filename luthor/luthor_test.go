package luthor

import "testing"

func Test_Load(t *testing.T) {
	cfg, err := LoadData([]string{"../../terra/modules/aws_eventbridge"}, nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log(cfg.Modules)
}
