package main

import "encoding/xml"

type RegistrationInfo struct {
	Date        string `xml:"Date"`
	Author      string `xml:"Author"`
	Description string `xml:"Description"`
	URI         string `xml:"URI"`
}

type Settings struct {
	AllowHardTerminate         bool   `xml:"AllowHardTerminate"`
	DisallowStartIfOnBatteries bool   `xml:"DisallowStartIfOnBatteries"`
	StopIfGoingOnBatteries     bool   `xml:"StopIfGoingOnBatteries"`
	ExecutionTimeLimit         string `xml:"ExecutionTimeLimit"`
	MultipleInstancesPolicy    string `xml:"MultipleInstancesPolicy"`
	RestartOnFailure           struct {
		Count    int    `xml:"Count"`
		Interval string `xml:"Interval"`
	} `xml:"RestartOnFailure"`
	StartWhenAvailable bool `xml:"StartWhenAvailable"`
	IdleSettings       struct {
		StopOnIdleEnd bool `xml:"StopOnIdleEnd"`
		RestartOnIdle bool `xml:"RestartOnIdle"`
	} `xml:"IdleSettings"`
}

type Principal struct {
	ID        string `xml:"id,attr"`
	UserID    string `xml:"UserId"`
	LogonType string `xml:"LogonType"`
	RunLevel  string `xml:"RunLevel"`
}

type Principals struct {
	Principal Principal `xml:"Principal"`
}

type Triggers struct {
	LogonTrigger struct{} `xml:"LogonTrigger"`
}

type Task struct {
	XMLName          xml.Name         `xml:"http://schemas.microsoft.com/windows/2004/02/mit/task Task"`
	Version          string           `xml:"version,attr"`
	RegistrationInfo RegistrationInfo `xml:"RegistrationInfo"`
	Principals       Principals       `xml:"Principals"`
	Settings         Settings         `xml:"Settings"`
	Triggers         Triggers         `xml:"Triggers"`
	Actions          Actions          `xml:"Actions"`
}

type Exec struct {
	Command string `xml:"Command"`
}

type Actions struct {
	Context string `xml:"Context,attr"`
	Exec    Exec   `xml:"Exec"`
}
