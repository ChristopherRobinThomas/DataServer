import requests
import streamlit as st
import pandas as pd
import numpy as np

URL="http://localhost:9000"
def getServers():
    resp = requests.get(f"{URL}/servers")
    return resp.json()

def getEntries(date):
    resp = requests.get(f"{URL}/entries/{date}")
    return resp.json()

def getRuns(date, serverID):
    print(f"{URL}/runs/{date}/{serverID}")
    resp = requests.get(f"{URL}/runs/{date}/{serverID}")
    return resp.json()

servers = getServers()
server_names = [] 
for elem in servers:
    server_names.append(elem["name"])

selected_servers = st.multiselect("Choose severs", server_names, default=["oscip05"])
entry_lists = []

for s in selected_servers:
    res = getEntries(s)
    entry_lists.append(res)

selected_dates_per_server = []
selected_serverID_per_server = []
for i in range(len(selected_servers)):
    date_list = []
    server_list = []
    for elem in entry_lists[i]:
        date_list.append(elem["date"])
        server_list.append(elem["server"])

    selected_dates = st.multiselect(f"Choose Dates for servers:{server_names[i]}", date_list)
    selected_dates_per_server.append(selected_dates)
    selected_serverID_per_server.append(server_list)

for i in range(len(selected_dates_per_server)):
    serverName = selected_servers[i]
    for j in range(len(selected_dates_per_server[i])):
        data = getRuns(selected_dates_per_server[i][j], selected_serverID_per_server[i][j])
        st.info(data)
