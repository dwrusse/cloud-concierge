package pyscriptexec

import (
	"bytes"
	"fmt"
	"os/exec"
)

// ExecutePythonScript is a generic function for executing a python script
// from within the python_scripts directory.
func (pse *pyScriptExec) ExecutePythonScript(name string, otherArgs []string) error {
	path := fmt.Sprintf("/python_scripts/%v/main.py", name)
	args := []string{path}
	args = append(args, otherArgs...)

	cmd := exec.Command("python3", args...)

	// Setting up logging objects
	var out bytes.Buffer
	cmd.Stdout = &out

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf(
			"%v\n\nError:\n%vStdOut:\n%v",
			err, stderr.String(), out.String(),
		)
	}
	fmt.Printf("\n%v Output:\n\n%v\n", path, out.String())
	return nil
}

// RunNLPEngine is a function that wraps ExecutePythonScript to execute
// python_scripts/nlpengine/main.py
func (pse *pyScriptExec) RunNLPEngine() error {
	err := pse.ExecutePythonScript("nlpengine", []string{})
	if err != nil {
		return err
	}
	return nil
}

// RunStateOfCloudReport is a function that wraps ExecutePythonScript to execute
// python_scripts/state_of_cloud_report/main.py
func (pse *pyScriptExec) RunStateOfCloudReport(uniqueID string, jobName string) error {
	jobArgs := []string{
		"--job_name", jobName,
		"--job_unique_id", uniqueID,
	}
	err := pse.ExecutePythonScript("state_of_cloud_report", jobArgs)
	if err != nil {
		return err
	}

	return nil
}
