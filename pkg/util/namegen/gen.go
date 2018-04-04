package namegen

import (
	"math/rand"
	"time"
)

var (
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func ColoredPhysics() string {
	i := rnd.Intn(len(Colors))
	j := rnd.Intn(len(Physicists))
	return Colors[i] + "-" + Physicists[j]
}

func Aster() string {
	return Asteroids[rnd.Intn(len(Asteroids))]
}