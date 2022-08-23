package sqlx

import (
	"database/sql"
	"strings"
)

type Matcher func(any) bool
type RoutingStrategy func([]*sql.DB) *sql.DB

type DBCluster struct {
	masterSlaves []*MasterSlave
}

type MasterSlave struct {
	master  *sql.DB
	slaves  []*sql.DB
	match   Matcher
	routing RoutingStrategy
}

func (cluster *DBCluster) Sharding(factor any) *MasterSlave {
	for _, ms := range cluster.masterSlaves {
		if ms.match(factor) {
			return ms
		}
	}

	return cluster.masterSlaves[0]
}

func (receiver *MasterSlave) DB(sql string) *sql.DB {
	if "select" == strings.ToLower(sql[0:6]) {
		return receiver.Slave()
	} else {
		return receiver.Master()
	}

}

func (receiver *MasterSlave) Master() *sql.DB {
	return receiver.master
}

func (receiver *MasterSlave) Slave() *sql.DB {
	return receiver.routing(receiver.slaves)
}
