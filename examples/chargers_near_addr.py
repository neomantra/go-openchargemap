# pip install osmnx pandas matplotlib 

import osmnx as ox
import pandas as pd
import subprocess
import matplotlib
from io import StringIO
from os import environ

ocm_key = environ.get("OCM_KEY")
if ocm_key == "":
    print("Please set the OCM_KEY environment variable to your OpenChargeMap API key")
    exit(1)

address = '362 Memorial Dr, Cambridge, MA'
dist = 1500   # meters

G = ox.graph_from_address(address, network_type='drive', dist=dist)
gdf_nodes, gdf_edges = ox.utils_graph.graph_to_gdfs(G)
bounds = gdf_nodes.total_bounds

# chargemeup wants "(lat1,lon1),(lat2,lon2)"
bounds_txt='(%f,%f),(%f,%f)' % (bounds[1], bounds[0], bounds[3], bounds[2])
json_data = subprocess.check_output(['chargemeup', '-k', ocm_key, '-b', bounds_txt])
print("found %d charging stations within %0.0f meters of '%s'" % (len(json_data), dist, address))
chargers_df = pd.read_json(StringIO(str(json_data, encoding='utf-8')))

# do a little extraction
totalConns, totalChargers = 0, 0
for i, ch in chargers_df.iterrows():
    addr = ch['AddressInfo']
    connections = ch['Connections']
    numConns = 0
    for conn in connections:
        if 'Quantity' in conn:
            numConns += conn['Quantity']
    if numConns != 0:
        totalChargers += 1
        totalConns += numConns
    print("id:%8d  geo: (%f,%f)  num_chargers: %d" % (ch['ID'], addr["Latitude"], addr["Longitude"], numConns))

print("total_chargers: %d  total_connections: %d" % (totalChargers, totalConns))