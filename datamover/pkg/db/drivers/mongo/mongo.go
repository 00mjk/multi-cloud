// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongo

import (
	"errors"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/micro/go-log"
	. "github.com/opensds/multi-cloud/dataflow/pkg/model"
)

var adap = &adapter{}

var DataBaseName = "test"
var CollJob = "job"

type MyLock struct {
	LockObj  string    `bson:"lockobj"`
	LockTime time.Time `bson:"locktime"`
}

func Init(host string) *adapter {
	//log.Log("edps:", deps)
	session, err := mgo.Dial(host)
	if err != nil {
		panic(err)
	}
	//defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	adap.s = session

	adap.userID = "unknown"

	return adap
}

func Exit() {
	adap.s.Close()
}

type adapter struct {
	s      *mgo.Session
	userID string
}

func (ad *adapter) UpdateJob(job *Job) error {
	ss := ad.s.Copy()
	defer ss.Close()

	c := ss.DB(DataBaseName).C(CollJob)
	j := Job{}
	err := c.Find(bson.M{"_id": job.Id}).One(&j)
	if err != nil {
		log.Logf("Get job failed before update it, err:%v\n", err)
		return errors.New("Get job failed before update it.")
	}

	if !job.StartTime.IsZero() {
		j.StartTime = job.StartTime
	}
	if !job.EndTime.IsZero() {
		j.EndTime = job.EndTime
	}
	if job.TotalCapacity != 0 {
		j.TotalCapacity = job.TotalCapacity
	}
	if job.TotalCount != 0 {
		j.TotalCount = job.TotalCount
	}
	if job.PassedCount != 0 {
		j.PassedCount = job.PassedCount
	}
	if job.PassedCapacity != 0 {
		j.PassedCapacity = job.PassedCapacity
	}
	if job.Status != "" {
		j.Status = job.Status
	}

	err = c.Update(bson.M{"_id": j.Id}, &j)
	if err != nil {
		log.Fatalf("Update job in database failed, err:%v\n", err)
		return errors.New("Update job in database failed.")
	}

	log.Logf("Update job in database successfully\n", err)
	return nil
}
