# grinnodes
grin node statistics by finding all the nodes in the network.

---

This program use sqlite to store node info.

The locations, company information of IP adddresses come from  <https://ipapi.co> and <https://api.ipdata.co>

Use <https://www.mapbox.com> to display the node on the map.

# Running
* make the database file node.db in data folder use the sqlite.sql file.
* configuration file config.toml in the data folder:
    * mapbox_key:  token of mapbox
    * ipdata_key:  token of https://api.ipdata.co 
    * peer_active_duration: how long time from last seen in network the node is active in the seconds
    * client_http_addr: the Grin Node (<https://github.com/mimblewimble/grin>) API address
