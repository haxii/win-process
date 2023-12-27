package win_process

import (
	"bytes"
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
	Modules      []string  `json:"Modules"`
	Start        time.Time `json:"-"`
	CreationDate string    `json:"CreationDate"`
}

func GetInfoByName(name string) ([]*Info, error) {
	cmd := fmt.Sprintf("$(Get-WmiObject Win32_Process -Filter \"name= '%s'\")  | Select-Object ProcessId,WorkingSetSize,CommandLine,CreationDate,Modules | ConvertTo-Json -Depth 1", name)
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
	if bytes.HasPrefix(out, []byte("[")) {
		err = json.Unmarshal(out, &infoList)
	} else {
		info := &Info{}
		err = json.Unmarshal(out, &info)
		infoList = append(infoList, info)
	}
	for _, info := range infoList {
		info.Start = formatTime(info.CreationDate)
		for i, module := range info.Modules {
			info.Modules[i] = extractModuleName(module)
		}
	}
	return infoList, err
}

func formatTime(t string) time.Time {
	t, _, _ = strings.Cut(t, "+")
	tt, _ := time.ParseInLocation("20060102150405.999999", t, time.Local)
	return tt
}

var moduleNameReplacer = strings.NewReplacer(`System.Diagnostics.ProcessModule (`, "", ")", "")

func extractModuleName(m string) string {
	return moduleNameReplacer.Replace(m)
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
