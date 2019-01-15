package main

import (
	"context"
	"runtime"

	dt "github.com/dylhunn/dragontoothmg"
)

var JobChannel = make(chan Job, 1000000)

type Worker struct {
	ID int64
}

type Job struct {
	Ctx        context.Context
	Args       Args
	Return     chan Return
	AllowedIDs int64
}

type Return struct {
	Score    int
	BestMove dt.Move
	Pv       []dt.Move
}

type Args struct {
	Board              *dt.Board
	Depth, Alpha, Beta int
	MoveList           []dt.Move
}

func (w *Worker) Work() {
	runtime.LockOSThread()
	for {
		w.DoOnce()
	}
}

func (w *Worker) DoOnce() {
	job := <-JobChannel
	if job.AllowedIDs&w.ID != w.ID {
		JobChannel <- job
		return
	}

	var ret Return
	ret.Score, ret.BestMove, ret.Pv = w.negaMax(ctx, job.Args.Board, job.Args.Depth, job.Args.Alpha, job.Args.Beta, job.Args.MoveList, job.AllowedIDs)
	job.Return <- ret
}
