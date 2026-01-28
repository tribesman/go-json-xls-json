package metrics

import (
	"bufio"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Snapshot holds resource counters that we can compute cheaply in-process.
type Snapshot struct {
	start time.Time
	m1    runtime.MemStats
	rssKB int64 // from /proc/self/status (Linux). -1 if unavailable.
}

func Start() Snapshot {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return Snapshot{
		start: time.Now(),
		m1:    m,
		rssKB: readProcSelfRSSKB(),
	}
}

type Result struct {
	Elapsed time.Duration

	// Heap deltas/values from runtime.MemStats.
	AllocDeltaBytes uint64
	TotalAllocDelta uint64
	SysBytes        uint64

	// RSS from /proc/self/status (Linux). -1 if unavailable.
	RSSKB int64
}

func (s Snapshot) Finish() Result {
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	rss2 := readProcSelfRSSKB()

	return Result{
		Elapsed:         time.Since(s.start),
		AllocDeltaBytes: deltaU64(m2.Alloc, s.m1.Alloc),
		TotalAllocDelta: deltaU64(m2.TotalAlloc, s.m1.TotalAlloc),
		SysBytes:        m2.Sys,
		RSSKB:           rss2,
	}
}

func deltaU64(after, before uint64) uint64 {
	if after >= before {
		return after - before
	}
	return 0
}

// readProcSelfRSSKB reads VmRSS from /proc/self/status (Linux).
// Returns -1 if unavailable.
func readProcSelfRSSKB() int64 {
	f, err := os.Open("/proc/self/status")
	if err != nil {
		return -1
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if !strings.HasPrefix(line, "VmRSS:") {
			continue
		}
		// Example: "VmRSS:\t  12345 kB"
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return -1
		}
		v, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return -1
		}
		return v
	}
	return -1
}
