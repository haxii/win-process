package win_process

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type Info struct {
	ID           int       `json:"ProcessId"`
	WS           int64     `json:"WorkingSetSize"`
	Args         string    `json:"CommandLine"`
	Start        time.Time `json:"-"`
	CreationDate string    `json:"CreationDate"`
}

func GetInfoByName(name string) ([]*Info, error) {
	cmd := fmt.Sprintf("$(Get-WmiObject Win32_Process -Filter \"name= '%s'\")  | Select-Object ProcessId,WorkingSetSize,CommandLine,CreationDate | ConvertTo-Json -Depth 1", name)
	out, err := exec.Command(
		"powershell",
		"-NoProfile",
		"-Command",
		"&{", cmd, "}",
	).Output()
	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, nil
	}
	infoList := make([]*Info, 0)
	err = json.Unmarshal(out, &infoList)
	for _, info := range infoList {
		info.Start = formatTime(info.CreationDate)
	}
	return infoList, err
}

func formatTime(t string) time.Time {
	t, _, _ = strings.Cut(t, "+")
	tt, _ := time.Parse("20060102150405.999999", t)
	return tt
}

func Kill(name string, filter func(info Info) bool) error {
	pList, err := GetInfoByName(name)
	if err != nil {
		return err
	}
	errList := make([]error, 0)
	for _, info := range pList {
		if info == nil {
			continue
		}
		if filter == nil || filter(*info) {
			if info.ID <= 0 {
				continue
			}
			if err = KillProcessByID(info.ID); err != nil {
				errList = append(errList, err)
			}
		}
	}
	if len(errList) > 0 {
		return errors.Join(errList...)
	}
	return nil
}
