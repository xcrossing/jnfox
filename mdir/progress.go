package mdir

import "github.com/cheggaaa/pb/v3"

type progress interface {
	increment()
	finish()
}

type realProgress struct {
	*pb.ProgressBar
}

func (p *realProgress) increment() {
	p.Increment()
}

func (p *realProgress) finish() {
	p.Finish()
}

type fakeProgress struct {
	pb.ProgressBar
}

func (p *fakeProgress) increment() {
}

func (p *fakeProgress) finish() {
}
