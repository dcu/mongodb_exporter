package collector

import(
    //"github.com/prometheus/client_golang/prometheus"
)


//Lock
type ReadWriteLockTimes struct {
    Read                  float64 `bson:"R" type:"counter"`
    Write                 float64 `bson:"W" type:"counter"`
    ReadLower             float64 `bson:"r" type:"counter"`
    WriteLower            float64 `bson:"w" type:"counter"`
}

type LockStats struct {
    TimeLockedMicros    ReadWriteLockTimes  `bson:"timeLockedMicros"`
    TimeAcquiringMicros ReadWriteLockTimes  `bson:"timeAcquiringMicros"`
}


