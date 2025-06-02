package nodes 


import (
        //"sync"
        //"time"

        //"w5gc.io/wipro5gcore/pkg/smf/pdusmsp/config"
)

type NodeId string

type NodeType                           int8

const (
        CurrentNode                     NodeType = 1
        RemoteNode                      NodeType = 2
)


