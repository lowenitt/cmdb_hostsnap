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

package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/wusendong/cmdb_hostsnap/publiser"
	"github.com/wusendong/cmdb_hostsnap/storage/redis"
)

type Hostsnap struct {
	configfile string
	pub        publiser.Publiser
	ctx        context.Context
	cancel     func()

	confLock sync.RWMutex
}

type HostsnapConfig struct {
	Redis   redis.Config `json:"redis"`
	Channel string       `json:"channel"`
}

func NewHostsnap(configfile string) (*Hostsnap, error) {
	conf, err := readConfig(configfile)
	if err != nil {
		return nil, fmt.Errorf("read config file %s error %v", configfile, err)
	}

	logrus.Infof("NewHostsnap with config: %#v ", conf)
	pub, err := publiser.NewRedisPubliser(conf.Channel, conf.Redis)
	if err != nil {
		return nil, fmt.Errorf("NewRedisPubliser error %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &Hostsnap{ctx: ctx, cancel: cancel, pub: pub, configfile: configfile}, nil
}

func readConfig(configfile string) (*HostsnapConfig, error) {
	file, err := os.OpenFile(configfile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	conf := new(HostsnapConfig)
	if err := json.NewDecoder(file).Decode(conf); err != nil {
		return nil, err
	}
	return conf, err
}

func (snap *Hostsnap) Run() error {
	logrus.Info("hostsnap started")
	defer logrus.Info("hostsnap stoped")

	ticker := time.NewTicker(time.Second * 5)

	event := map[string]interface{}{}
	err := json.Unmarshal([]byte(hostsnapExample), &event)
	if err != nil {
		return fmt.Errorf("prepare event failed: %v", err)
	}

	for {
		select {
		case <-snap.ctx.Done():
			return nil
		case <-ticker.C:
		}

		err := snap.publishEvent(event)
		if err != nil {
			logrus.Errorf("publish error %v", err)
		}

		logrus.Info("sended event")
	}
}

func (snap *Hostsnap) Stop() error {
	logrus.Info("hostsnap stop")
	snap.pub.Close()
	snap.cancel()
	return nil
}

func (snap *Hostsnap) Reload() error {
	conf, err := readConfig(snap.configfile)
	if err != nil {
		return fmt.Errorf("read config file %s error %v", snap.configfile, err)
	}
	pub, err := publiser.NewRedisPubliser(conf.Channel, conf.Redis)
	if err != nil {
		return fmt.Errorf("NewRedisPubliser error %v", err)
	}
	snap.pub.Close()

	logrus.Infof("reload with config: %#v", conf)
	snap.confLock.Lock()
	snap.pub = pub
	snap.confLock.Unlock()

	return nil
}

func (snap *Hostsnap) publishEvent(event map[string]interface{}) error {
	snap.confLock.RLock()
	err := snap.pub.PublishEvent(event)
	snap.confLock.RUnlock()
	return err
}

const hostsnapExample = `{
	"ip": "192.168.1.7",
	"bizid": 0,
	"cloudid": 0,
	"data": {
		"timezone": 8,
		"datetime": "2017-09-19 16:57:07",
		"utctime": "2017-09-19 08:57:07",
		"country": "Asia",
		"city": "Shanghai",
		"cpu": {
			"cpuinfo": [
				{
					"cpu": 0,
					"vendorID": "GenuineIntel",
					"family": "6",
					"model": "63",
					"stepping": 2,
					"physicalID": "0",
					"coreID": "0",
					"cores": 1,
					"modelName": "Intel(R) Xeon(R) CPU E5-26xx v3",
					"mhz": 2294.01,
					"cacheSize": 4096
				}
			],
			"per_usage": [
				3.0232169701043103
			],
			"total_usage": 3.0232169701043103,
			"per_stat": [
				{
					"cpu": "cpu0",
					"user": 5206.09,
					"system": 6107.04,
					"idle": 337100.84,
					"nice": 6.68,
					"iowait": 528.24,
					"irq": 0.02,
					"softirq": 13.48,
					"steal": 0,
					"guest": 0,
					"guestNice": 0,
					"stolen": 0
				}
			],
			"total_stat": {
				"cpu": "cpu-total",
				"user": 5206.09,
				"system": 6107.04,
				"idle": 337100.84,
				"nice": 6.68,
				"iowait": 528.24,
				"irq": 0.02,
				"softirq": 13.48,
				"steal": 0,
				"guest": 0,
				"guestNice": 0,
				"stolen": 0
			}
		},
		"env": {
			"crontab": [
				{
					"user": "root",
					"content": "#secu-tcs-agent monitor, install at Fri Sep 15 16:12:02 CST 2017\n* * * * * /usr/local/sa/agent/secu-tcs-agent-mon-safe.sh /usr/local/sa/agent \u003e /dev/null 2\u003e\u00261\n*/1 * * * * /usr/local/qcloud/stargate/admin/start.sh \u003e /dev/null 2\u003e\u00261 \u0026\n*/20 * * * * /usr/sbin/ntpdate ntpupdate.tencentyun.com \u003e/dev/null \u0026\n*/1 * * * * cd /usr/local/gse/gseagent; ./cron_agent.sh 1\u003e/dev/null 2\u003e\u00261\n"
				}
			],
			"host": "127.0.0.1  localhost  localhost.localdomain  VM_0_31_centos\n::1         localhost localhost.localdomain localhost6 localhost6.localdomain6\n",
			"route": "Kernel IP routing table\nDestination     Gateway         Genmask         Flags Metric Ref    Use Iface\n10.0.0.0        0.0.0.0         255.255.255.0   U     0      0        0 eth0\n169.254.0.0     0.0.0.0         255.255.0.0     U     1002   0        0 eth0\n0.0.0.0         10.0.0.1        0.0.0.0         UG    0      0        0 eth0\n"
		},
		"disk": {
			"diskstat": {
				"vda1": {
					"major": 252,
					"minor": 1,
					"readCount": 24347,
					"mergedReadCount": 570,
					"writeCount": 696357,
					"mergedWriteCount": 4684783,
					"readBytes": 783955968,
					"writeBytes": 22041231360,
					"readSectors": 1531164,
					"writeSectors": 43049280,
					"readTime": 80626,
					"writeTime": 12704736,
					"iopsInProgress": 0,
					"ioTime": 822057,
					"weightedIoTime": 12785026,
					"name": "vda1",
					"serialNumber": "",
					"speedIORead": 0,
					"speedByteRead": 0,
					"speedIOWrite": 2.9,
					"speedByteWrite": 171144.53333333333,
					"util": 0.0025666666666666667,
					"avgrq_sz": 115.26436781609195,
					"avgqu_sz": 0.06568333333333334,
					"await": 22.649425287356323,
					"svctm": 0.8850574712643678
				}
			},
			"partition": [
				{
					"device": "/dev/vda1",
					"mountpoint": "/",
					"fstype": "ext3",
					"opts": "rw,noatime,acl,user_xattr"
				}
			],
			"usage": [
				{
					"path": "/",
					"fstype": "ext2/ext3",
					"total": 52843638784,
					"free": 47807447040,
					"used": 2351915008,
					"usedPercent": 4.4507060113962345,
					"inodesTotal": 3276800,
					"inodesUsed": 29554,
					"inodesFree": 3247246,
					"inodesUsedPercent": 0.9019165039062501
				}
			]
		},
		"load": {
			"load_avg": {
				"load1": 0,
				"load5": 0,
				"load15": 0
			}
		},
		"mem": {
			"meminfo": {
				"total": 1044832256,
				"available": 805912576,
				"used": 238919680,
				"usedPercent": 22.866797864249705,
				"free": 92041216,
				"active": 521183232,
				"inactive": 352964608,
				"wired": 0,
				"buffers": 110895104,
				"cached": 602976256,
				"writeback": 0,
				"dirty": 151552,
				"writebacktmp": 0
			},
			"vmstat": {
				"total": 0,
				"used": 0,
				"free": 0,
				"usedPercent": 0,
				"sin": 0,
				"sout": 0
			}
		},
		"net": {
			"interface": [
				{
					"mtu": 65536,
					"name": "lo",
					"hardwareaddr": "28:31:52:1d:c6:0a",
					"flags": [
						"up",
						"loopback"
					],
					"addrs": [
						{
							"addr": "192.168.1.7/8"
						}
					]
				},
				{
					"mtu": 1500,
					"name": "eth0",
					"hardwareaddr": "52:54:00:19:2e:e8",
					"flags": [
						"up",
						"broadcast",
						"multicast"
					],
					"addrs": [
						{
							"addr": "192.168.1.2/24"
						}
					]
				}
			],
			"dev": [
				{
					"name": "lo",
					"speedSent": 0,
					"speedRecv": 0,
					"speedPacketsSent": 0,
					"speedPacketsRecv": 0,
					"bytesSent": 604,
					"bytesRecv": 604,
					"packetsSent": 2,
					"packetsRecv": 2,
					"errin": 0,
					"errout": 0,
					"dropin": 0,
					"dropout": 0,
					"fifoin": 0,
					"fifoout": 0
				},
				{
					"name": "eth0",
					"speedSent": 574,
					"speedRecv": 214,
					"speedPacketsSent": 3,
					"speedPacketsRecv": 2,
					"bytesSent": 161709123,
					"bytesRecv": 285910298,
					"packetsSent": 1116625,
					"packetsRecv": 1167796,
					"errin": 0,
					"errout": 0,
					"dropin": 0,
					"dropout": 0,
					"fifoin": 0,
					"fifoout": 0
				}
			],
			"netstat": {
				"established": 2,
				"syncSent": 1,
				"synRecv": 0,
				"finWait1": 0,
				"finWait2": 0,
				"timeWait": 0,
				"close": 0,
				"closeWait": 0,
				"lastAck": 0,
				"listen": 2,
				"closing": 0
			},
			"protocolstat": [
				{
					"protocol": "udp",
					"stats": {
						"inDatagrams": 176253,
						"inErrors": 0,
						"noPorts": 1,
						"outDatagrams": 199569,
						"rcvbufErrors": 0,
						"sndbufErrors": 0
					}
				}
			]
		},
		"system": {
			"info": {
				"hostname": "VM_0_31_centos",
				"uptime": 348315,
				"bootTime": 1505463112,
				"procs": 142,
				"os": "linux",
				"platform": "centos",
				"platformFamily": "rhel",
				"platformVersion": "6.2",
				"kernelVersion": "2.6.32-504.30.3.el6.x86_64",
				"virtualizationSystem": "",
				"virtualizationRole": "",
				"hostid": "96D0F4CA-2157-40E6-BF22-6A7CD9B6EB8C",
				"systemtype": "64-bit"
			}
		}
	}
}`
