package hw10programoptimization

import (
	"archive/zip"
	"testing"
)

func BenchmarkGetDomainStat(b *testing.B) {
	b.Helper()
	b.StopTimer()

	r, _ := zip.OpenReader("testdata/users.dat.zip")
	defer r.Close()

	data, _ := r.File[0].Open()

	b.StartTimer()
	GetDomainStat(data, "biz")
	b.StopTimer()
}
