import csv
import os
import pandas as pd
import numpy as np
import math

### 計算價格        (總淨負載、充電功率、時間電價、慢充快充、入場時間、出場時間、充電效率、變壓器功率、時間段大小)
def estimate_price(total_load, p_charge, tou, ch_type, time_in, time_out, efficiency, p_tr, time_scale=5):
    t_time = float(time_scale)/60

    # 計算充電單價(公式 = 電網緊繃百分比*0.4*tou + 0.9*tou , 也就是收費是0.9~1.3倍的時間電價)
    total_price_of_ch = 0.0
    total_charge = 0.0
    for i in range(len(total_load)):
        if p_charge[i] != 0 :
            p_charge[i] *= efficiency
            percent_of_grid = float(total_load[i])/p_tr
            total_price_of_ch += (math.pow(percent_of_grid,3)*0.4 + 0.9)*tou[i]*p_charge[i]*t_time
            total_charge += p_charge[i] * t_time
    if (total_charge != 0):
        unit_price_of_ch = total_price_of_ch/total_charge
    else:
        print("this ev total_charge = 0 (msg from estimate_price.py)")
    # 計算佔位總價
    AC_price = 10   #(單位:每小時/元)
    DC_price = 50
    if (ch_type == 1):
        total_price_of_space = (time_out-time_in-2)*AC_price*t_time
    elif (ch_type == 2):
        total_price_of_space = (time_out-time_in-2)*DC_price*t_time
    #計算總價(充電價+佔位價)
    total_price = total_price_of_ch + total_price_of_space

    return round(unit_price_of_ch), round(total_price_of_space), round(total_price)


# ## 主程式
# if __name__ == "__main__":
#     ### 讀取資料
#     dir = 'D:\\實驗室\\111-2暑假\\最佳化程式周邊'
#     os.chdir(dir)
#     print("Start read data...")
#     df_data = pd.read_csv('estimate_price_data.csv')
#     df_tou = pd.read_csv('tou.csv')  
#     total_load = np.array(df_data['total_load'])   
#     total_load = np.insert(total_load,0,0) 
#     p_charge = np.array(df_data['p_charge'])
#     p_charge = np.insert(p_charge,0,0)
#     tou = np.array(df_tou['tou'])
#     tou = np.insert(tou,0,0)
#     ch_type = 1        # 1:慢充 / 2:快充
#     time_in = 99
#     time_out = 115
#     p_tr = 600      #變壓器容量

#     ### 計算價格
#     unit_price_of_ch, total_price_of_space = estimate_price(total_load, p_charge, tou, ch_type, time_in, time_out, p_tr)
#     print("unit_price_of_ch = ",unit_price_of_ch,"(kWh/$NTD)")
#     print("total_price_of_space = ",total_price_of_space,"($NTD)")


