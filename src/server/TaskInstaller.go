package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"
)

// Получение текущего пользователя и его SID
func getUserDetails() (username, sid string, err error) {
	u, err := user.Current()
	if err != nil {
		return "", "", err
	}
	username = u.Username
	sid = u.Uid
	return username, sid, nil
}

// Получение текущего времени в формате ISO 8601
func getCurrentTimeISO8601() string {
	return time.Now().Format(time.RFC3339)
}

func createTask(name, command string) error {
	username, userSID, err := getUserDetails()
	if err != nil {
		return fmt.Errorf("failed to get user details: %v", err)
	}
	currentTime := getCurrentTimeISO8601()

	task := Task{
		Version: "1.2",
		XMLName: xml.Name{Space: "http://schemas.microsoft.com/windows/2004/02/mit/task", Local: "Task"},
		RegistrationInfo: RegistrationInfo{
			Date:        currentTime,
			Author:      username,
			Description: command,
			URI:         "\\" + name,
		},
		Principals: Principals{
			Principal: Principal{
				ID:        "Author",
				UserID:    userSID,
				LogonType: "S4U", // Уточните тип логина в зависимости от требований
				RunLevel:  "HighestAvailable",
			},
		},
		Settings: Settings{
			AllowHardTerminate:         false,
			DisallowStartIfOnBatteries: false,
			StopIfGoingOnBatteries:     false,
			ExecutionTimeLimit:         "PT0S",
			MultipleInstancesPolicy:    "IgnoreNew",
			RestartOnFailure: struct {
				Count    int    `xml:"Count"`
				Interval string `xml:"Interval"`
			}{
				Count:    3,
				Interval: "PT1M",
			},
			StartWhenAvailable: true,
			IdleSettings: struct {
				StopOnIdleEnd bool `xml:"StopOnIdleEnd"`
				RestartOnIdle bool `xml:"RestartOnIdle"`
			}{
				StopOnIdleEnd: false,
				RestartOnIdle: false,
			},
		},
		Triggers: Triggers{
			LogonTrigger: struct{}{}, // Пустая структура для LogonTrigger
		},
		Actions: Actions{
			Context: "Author",
			Exec: Exec{
				Command: command,
			},
		},
	}

	xmlData, err := xml.MarshalIndent(task, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal XML: %v", err)
	}

	xmlFilePath := fmt.Sprintf("./%s.xml", name)
	file, err := os.Create(xmlFilePath)
	if err != nil {
		return fmt.Errorf("failed to create XML file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	if _, err = file.Write(xmlData); err != nil {
		return fmt.Errorf("failed to write XML data to file: %v", err)
	}

	cmd := exec.Command("schtasks", "/create", "/tn", name, "/xml", xmlFilePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create task: %v, output: %s", err, output)
	}

	return nil
}

// Проверяет, существует ли задача с указанным именем
func taskExists(name string) bool {
	out, err := exec.Command("schtasks", "/query", "/tn", name).Output()
	if err != nil {
		// Если задача не найдена, команда вернет ошибку
		return false
	}
	return strings.Contains(string(out), name)
}

func taskInstaller() {
	if !taskExists("crazyFirewallServer") {
		if err := createTask("crazyFirewallServer", `C:\Program Files\crazyfirewall\crazyFirewallServer.exe`); err != nil {
			logger("Error creating server task: %v", err)
		} else {
			logger("Server task created successfully.")
		}
	} else {
		logger("Server task already exists.")
	}

	if !taskExists("crazyFirewallClient") {
		if err := createTask("crazyFirewallClient", `C:\Program Files\crazyfirewall\crazyFirewallClient.exe`); err != nil {
			logger("Error creating client task: %v", err)
		} else {
			logger("Client task created successfully.")
		}
	} else {
		logger("Client task already exists.")
	}
}
