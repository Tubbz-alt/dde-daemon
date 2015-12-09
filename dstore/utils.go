package dstore

import (
	"dbus/com/deepin/lastore"
	"fmt"
	"pkg.deepin.io/lib/dbus"
	"time"
)

const (
	DStoreDBusDest = "com.deepin.lastore"
	DStoreDBusPath = "/com/deepin/lastore"

	JobStatusSucceed = "succeed"
	JobStatusFailed  = "failed"
	JobStatusEnd     = "end"
)

const (
	jobTypeInstall = "install"
	jobTypeRemove  = "remove"
)

func newDStoreManager() (*lastore.Manager, error) {
	return lastore.NewManager(DStoreDBusDest, DStoreDBusPath)
}

func destroyDStoreManager(manager *lastore.Manager) {
	if manager == nil {
		return
	}
	lastore.DestroyManager(manager)
}

func newDStoreJob(jobPath dbus.ObjectPath) (*lastore.Job, error) {
	return lastore.NewJob(DStoreDBusDest, jobPath)
}

func destroyDStoreJob(job *lastore.Job) {
	if job == nil {
		return
	}
	lastore.DestroyJob(job)
}

func waitJobDone(jobPath dbus.ObjectPath, jobType string, timeout <-chan time.Time, result *(chan error)) {
	job, err := newDStoreJob(jobPath)
	if err != nil {
		*result <- err
		return
	}
	defer destroyDStoreJob(job)

	quit := make(chan struct{})

	job.Status.ConnectChanged(func() {
		if jobPath != job.Path || job.Type.Get() != jobType {
			return
		}

		status := job.Status.Get()
		switch status {
		case JobStatusSucceed, JobStatusEnd:
			*result <- nil
			close(quit)
			return
		case JobStatusFailed:
			*result <- fmt.Errorf(job.Description.Get())
			close(quit)
			return
		}
	})

	select {
	case <-quit:
		return
	case <-timeout:
		*result <- fmt.Errorf("Do job '%v - %v' timeout",
			jobType, job.Packages.Get())
		return
	}
}