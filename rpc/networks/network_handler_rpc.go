package networks

import (
	"mmschainnewaccount/config"
	"mmschainnewaccount/rpc/model"
	"net/http"
	"strconv"

	"github.com/prestonTao/libp2parea/nodeStore"

	coreconfig "github.com/prestonTao/libp2parea/config"
)

func NetworkInfo(rj *model.RpcJson, w http.ResponseWriter, r *http.Request) (res []byte, err error) {
	netAddr := nodeStore.NodeSelf.IdInfo.Id.B58String()
	isSuper := nodeStore.NodeSelf.IsSuper
	nodes := nodeStore.GetLogicNodes()
	ids := make([]string, 0)
	for _, one := range nodes {
		ids = append(ids, one.B58String())
	}

	m := make(map[string]interface{})
	m["netaddr"] = netAddr
	m["issuper"] = isSuper
	m["superNodes"] = ids

	m["webaddr"] = config.WebAddr + ":" + strconv.Itoa(int(config.WebPort))
	m["tcpaddr"] = coreconfig.Init_LocalIP + ":" + strconv.Itoa(int(coreconfig.Init_LocalPort))

	res, err = model.Tojson(m)
	return
}
