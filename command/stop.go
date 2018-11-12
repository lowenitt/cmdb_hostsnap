/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package command

import (
	"fmt"
	"os"
	"syscall"

	"gopkg.in/urfave/cli.v1"

	"github.com/wusendong/cmdb_hostsnap/pidfile"
)

func StopCmd() cli.Command {
	return cli.Command{
		Name:  "stop",
		Usage: "stop the hostsnap process",
		Action: func(c *cli.Context) {
			pid, err := pidfile.ReadPid()
			if err != nil {
				fmt.Fprintf(os.Stderr, "read pid file error %v", err)
			}
			proc, err := os.FindProcess(pid)
			if err != nil {
				fmt.Fprintf(os.Stderr, "pid %d not exist", pid)
			}

			signal := syscall.SIGQUIT
			err = proc.Signal(signal)
			if err != nil {
				fmt.Fprintf(os.Stderr, "signal %v to pid %d faile: %v", signal, pid, err)
			}
			fmt.Print("stop success")
		},
	}
}
