Steps Taken:

Followed this link to setup prometheus server locally : https://serverspace.io/support/help/install-prometheus-ubuntu-20-04/
( Currently using the contained prometheus.yml file for setup )

Followed this link to setup grafana w/ prometheus : https://medium.com/devops-dudes/install-prometheus-on-ubuntu-18-04-a51602c6256b

Go here to view prom server & metrics : http://localhost:9090/graph?g0.expr=&g0.tab=0&g0.stacked=0&g0.show_exemplars=0&g0.range_input=1h
Go here to view grafana server : http://localhost:3000/d/b647efb8-4a43-4e07-81de-7cd2ef1a636a/local-geth-monitoring?orgId=1&from=now-30m&to=now&var-geth_l1_miner_job=geth-l1-miner&var-Chains=All&refresh=30s#

Future : 
https://stackoverflow.com/questions/46916328/dynamically-add-targets-to-a-prometheus-configuration
https://prometheus.io/docs/prometheus/latest/configuration/configuration/#file_sd_config
