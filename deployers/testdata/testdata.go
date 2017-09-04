package testdata

import "regexp"

var (
	ListFleetUnitsRegexp = regexp.MustCompile(`-- ggn staging fleetctl -- list-units`)
	ListFleetUnits       = `failed	b4ce7d6	loaded	54604fcc.../127.0.0.1	failed	staging_sleepy_sleepy0.service
failed	91c9aa5	loaded	54604fcc.../127.0.0.1	failed	staging_webhooks_webhooks0.service
failed	dcd3896	loaded	54604fcc.../127.0.0.1	failed	staging_webhooks_webhooks1.service
failed	59c4fc4	loaded	54604fcc.../127.0.0.1	failed	staging_webhooks_webhooks2.service`

	CatUnitRegexp = regexp.MustCompile(`-- ggn staging fleetctl cat staging_webhooks_webhooks\d.service`)
	CatUnit       = `ExecStart=/opt/bin/rkt      --insecure-options=all run      --set-env=TEMPLATER_OVERRIDE='${ATTR_0}'      --set-env=TEMPLATER_OVERRIDE_BASE64='${ATTR_BASE64_0}${ATTR_BASE64_1}'      --set-env=HOSTNAME='webhooks'      --set-env=HOST="%H"      --hostname=webhooks      --dns=10.254.0.3 --dns=10.254.0.4       --dns-search=pp-bourse.par-1.h.blbl.cr       --uuid-file-save=/mnt/sda9/rkt-uuid/pp-bourse/webhooks      --set-env=DOMAINNAME='pp.par-1.h.blbl.cr'      --net='bond0'      --set-env=AIRFLOW_HOME='/opt/webhooks'      aci.blbl.cr/pod-webhooks_aci-go-synapse:1.8.1-1    aci.blbl.cr/pod-webhooks_aci-go-nerve:1.8.1-1    aci.blbl.cr/pod-webhooks_aci-confd:1.8.1-1    aci.blbl.cr/pod-webhooks_aci-embulk:1.8.1-1    aci.blbl.cr/pod-webhooks_aci-zabbix-agent:1.8.1-1    aci.blbl.cr/pod-webhooks_aci-webhooks:1.8.1-1 --exec /usr/local/bin/webhooks -- scheduler ---
	`
)
