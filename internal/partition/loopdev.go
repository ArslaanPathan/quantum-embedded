// Copyright (c) 2026 Arslaan Pathan
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package partition

import (
	"fmt"
	"os/exec"
	"strings"
)

type LoopDevice struct {
	dev string
}

func (l *LoopDevice) Attach(img string) error {
	cmd := exec.Command("sudo", "losetup", "-fP", "--show", img)
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to attach loop device for %s: %w", img, err)
	}
	l.dev = strings.TrimSpace(string(out))
	return nil
}

func (l *LoopDevice) Detach() error {
	if l.dev == "" {
		return nil
	}
	cmd := exec.Command("sudo", "losetup", "-d", l.dev)
	err := cmd.Run()
	l.dev = ""
	return err
}

func (l *LoopDevice) GetDevice() string {
	return l.dev
}
