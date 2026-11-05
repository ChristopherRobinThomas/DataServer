import streamlit as st
import pandas as pd
import numpy as np

all_dates = ["2026_11_1_22", "2026_11_1_24"]
date = st.selectbox("Dates", all_dates)

st.write("Hello World")
data = pd.DataFrame(np.random.randn(3, 1))
tab1, tab2 = st.tabs(["bar","plot"])
tab1.bar_chart(data, height=300)
tab2.line_chart(data, height=300)
