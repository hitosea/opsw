package utils

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"
)

func Cmd(arg ...string) (string, error) {
	output, err := exec.Command("/bin/bash", arg...).CombinedOutput()
	return string(output), err
}

func CmdSh(arg ...string) (string, error) {
	output, err := exec.Command("/bin/sh", arg...).CombinedOutput()
	return string(output), err
}

func CmdFile(filePath string) {
	_, _ = Cmd("-c", fmt.Sprintf("chmod +x %s", filePath))
	cmdString := exec.Command("/bin/bash", filePath)
	PrintCmdOutput(cmdString)
}

func PrintCmdOutput(cmd *exec.Cmd) {
	cmd.Stdin = os.Stdin

	var wg sync.WaitGroup
	wg.Add(2)
	//捕获标准输出
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("INFO:", err)
		os.Exit(1)
	}
	readout := bufio.NewReader(stdout)
	go func() {
		defer wg.Done()
		GetOutput(readout)
	}()

	//捕获标准错误
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
	readerr := bufio.NewReader(stderr)
	go func() {
		defer wg.Done()
		GetOutput(readerr)
	}()

	//执行命令
	err = cmd.Run()
	if err != nil {
		return
	}
	wg.Wait()
}

func GetOutput(reader *bufio.Reader) {
	var sumOutput string //统计屏幕的全部输出内容
	outputBytes := make([]byte, 200)
	for {
		n, err := reader.Read(outputBytes) //获取屏幕的实时输出(并不是按照回车分割，所以要结合sumOutput)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			sumOutput += err.Error()
		}
		output := string(outputBytes[:n])
		fmt.Print(output) //输出屏幕内容
		sumOutput += output
	}
}

func Exec(cmdStr string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.Command("bash", "-c", cmdStr)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", errors.New("ErrCmdTimeout")
	}
	if err != nil {
		errMsg := ""
		if len(stderr.String()) != 0 {
			errMsg = fmt.Sprintf("stderr: %s", stderr.String())
		}
		if len(stdout.String()) != 0 {
			if len(errMsg) != 0 {
				errMsg = fmt.Sprintf("%s; stdout: %s", errMsg, stdout.String())
			} else {
				errMsg = fmt.Sprintf("stdout: %s", stdout.String())
			}
		}
		return errMsg, err
	}
	return stdout.String(), nil
}

func ExecWithTimeOut(cmdStr string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.Command("bash", "-c", cmdStr)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", errors.New("ErrCmdTimeout")
	}
	if err != nil {
		errMsg := ""
		if len(stderr.String()) != 0 {
			errMsg = fmt.Sprintf("stderr: %s", stderr.String())
		}
		if len(stdout.String()) != 0 {
			if len(errMsg) != 0 {
				errMsg = fmt.Sprintf("%s; stdout: %s", errMsg, stdout.String())
			} else {
				errMsg = fmt.Sprintf("stdout: %s", stdout.String())
			}
		}
		return errMsg, err
	}
	return stdout.String(), nil
}

func ExecCronjobWithTimeOut(cmdStr string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.Command("bash", "-c", cmdStr)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", errors.New("ErrCmdTimeout")
	}

	errMsg := ""
	if len(stderr.String()) != 0 {
		errMsg = fmt.Sprintf("stderr:\n %s", stderr.String())
	}
	if len(stdout.String()) != 0 {
		if len(errMsg) != 0 {
			errMsg = fmt.Sprintf("%s \n\n; stdout:\n %s", errMsg, stdout.String())
		} else {
			errMsg = fmt.Sprintf("stdout\n: %s", stdout.String())
		}
	}
	return errMsg, err
}

func Execf(cmdStr string, a ...interface{}) (string, error) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf(cmdStr, a...))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		errMsg := ""
		if len(stderr.String()) != 0 {
			errMsg = fmt.Sprintf("stderr: %s", stderr.String())
		}
		if len(stdout.String()) != 0 {
			if len(errMsg) != 0 {
				errMsg = fmt.Sprintf("%s; stdout: %s", errMsg, stdout.String())
			} else {
				errMsg = fmt.Sprintf("stdout: %s", stdout.String())
			}
		}
		return errMsg, err
	}
	return stdout.String(), nil
}

func HasNoPasswordSudo() bool {
	cmd2 := exec.Command("sudo", "-n", "ls")
	err2 := cmd2.Run()
	return err2 == nil
}
