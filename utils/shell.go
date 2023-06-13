package utils

import (
	"context"
	"os/exec"
	"time"
)

func ExecCmd(cmdName string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, cmdName, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}

func ExecCmdInDir(cmdName string, dir string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, cmdName, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}

func ExecCmdWithTimeout(second int32, cmdName string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(second)*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, cmdName, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}

func ExecCmdInDirWithTimeout(second int8, cmdName string, dir string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(second)*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, cmdName, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}
