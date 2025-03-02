package analysis

import "github.com/charmingruby/gdp/internal/storage"

type RetransmissionUnit struct {
	Type             string `json:"type"`
	InitialRWND      int    `json:"initial_rwnd"`
	FinalRWND        int    `json:"final_rwnd"`
	InitialCWND      int    `json:"initial_cwnd"`
	FinalCWND        int    `json:"final_cwnd"`
	InitialSshthresh int    `json:"initial_ssthresh"`
	FinalSshthresh   int    `json:"final_ssthresh"`
}

func SaveRetransmissionData(path string, data []RetransmissionUnit) error {
	return storage.SaveRecord(path, storage.Record{
		Data: data,
	})
}
