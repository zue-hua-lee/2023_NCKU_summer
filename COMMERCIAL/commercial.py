import csv
import math
from estimate_price import estimate_price
import gurobipy as gp
from gurobipy import GRB

class ev:
    def __init__(self, name: int, time_in: int, time_out: int, soc_in: float, 
                 soc_out: float, soc_now: float, capacity: int, num_se: int):
        
        self.name = name            #電動車編號
        self.time_in = time_in      #進場時間
        self.time_out = time_out    #離場時間
        self.soc_in = soc_in/100    #進場SOC
        self.soc_out = soc_out/100  #離場SOC
        self.soc_now = soc_now/100  #現在SOC
        self.capacity = capacity    #電池容量
        self.num_se = num_se        #充電樁編號
        
        self.charge = 0.0           #所需充電量
        self.cost = 0.0             #車主花費
        
class se:
    def __init__(self, name: int, index_in_evlist: int, time_in: int, 
                 time_out: int):
        
        self.name = name                            #充電樁編號
        self.index_in_evlist = index_in_evlist      #所在車列中的index
        self.time_in = time_in                      #充電之電動車進場時間
        self.time_out = time_out                    #充電之電動車離場時間   

class power:
    def __init__(self, StationID: int, ChargeID: int, Power: int, TimeStamp: int):
        
        self.StationID = StationID
        self.ChargeID = ChargeID
        self.Power = Power
        self.State = 0
        self.TimeStamp = TimeStamp
        self.ev_ID = -1
        self.ev_soc = 0
        
class fcs:  #沒有新車加入
    def __init__(self, now_time: int, new_ev: int):
        
        self.now_time = now_time        #現在時間
        
        self.num_time = 288             #時間區間數
        self.earnings = 0.0             #充電收益
        self.cost = 0.0                 #購電電費
       
        self.ev_list = []                       #場內電動車
        self.ess = [0]*(self.num_time)          #每時段儲能電量
        self.ess_init = 0.5                       #儲能初始值定為0.5
        self.se_list = []                       #場內充電樁
        self.pnet = [0]*(self.num_time)         #淨負載
        self.Pbuy = [0]*(self.num_time)    #正淨負載
        self.get_FCS_info()
        if(new_ev == 1):
            self.update_ev_list()
    
    def read_file(self, file_name): #從本地端讀取資料
        try:
            with open(file_name, 'r', encoding='utf-8', errors='ignore', newline='') as file:
                csv_reader = csv.reader(file)
                header = next(csv_reader)
                info = []
                for row in csv_reader:
                    info.append(float(row[1]))
        except FileNotFoundError:
            print('文件未被找到')
        except Exception as e:       
            print('發生錯誤', e)
        return info
    
    def read_parameter(self): #從本地端讀取資料
        try:
            with open('cpos_parameter.csv', 'r', encoding='utf-8', errors='ignore',
                      newline='') as file:
                csv_reader = csv.reader(file)
                info = [0]
                for row in csv_reader:
                    if(row[1] == '1'):
                        self.efficiency = float(row[3])
                        self.ess_cap = int(row[4])
                        self.Pess = int(row[5])
                        self.Ptr = int(row[6])
                        self.ac_num_charge = int(row[7])
                        self.dc_num_charge = int(row[8])
                        self.ac_charge_price = int(row[9])
                        self.dc_charge_price = int(row[10])
                        self.ac_Pchar = int(row[11])
                        self.dc_Pchar = int(row[12])
        except FileNotFoundError:
            print('文件未被找到')
        except Exception as e:
            print('發生錯誤', e)
    
    def read_se_list(self):
        try:
            with open('se_list.csv', 'r', newline='') as file:
                csv_reader = csv.reader(file)
                header = next(csv_reader)
                for row in csv_reader:
                    temp_se = se(int(row[0]), -1, int(row[1]), int(row[2]))
                    self.se_list.append(temp_se)
        except FileNotFoundError:
            print('文件未被找到')
        except Exception as e:
            print('發生錯誤', e)
                    
    def read_ev_list(self):
        try:
            with open('ev_list.csv', 'r', newline='') as file:
                csv_reader = csv.reader(file)
                header = next(csv_reader)  #跳過第一行
                for row in csv_reader:
                    if(int(row[2]) > self.now_time): #離場時間大於現在時間才加進車列中
                        temp_ev = ev(int(row[0]), int(row[1]), int(row[2]),
                                     float(row[3]), float(row[4]), float(row[5]),
                                     int(row[6]), int(row[7]))
                        self.ev_list.append(temp_ev)
                        self.se_list[int(row[7])-1].index_in_evlist = len(self.ev_list)-1
                        self.se_list[int(row[7])-1].time_in = int(row[1])
                        self.se_list[int(row[7])-1].time_out = int(row[2])
                    elif(int(row[2]) <= self.now_time): #離場時間小於現在時間踢出車列
                        try:
                            with open('dep_ev.csv', 'a', newline='') as csvfile:
                                writer = csv.writer(csvfile)
                                writer.writerow(row)
                        except FileNotFoundError:
                            print('文件未被找到')
                        except Exception as e:
                            print('發生錯誤', e)
                
                for index in range(len(self.ev_list)):  #計算在場電動車所需充電量&取得預測電動車負載
                    self.ev_list[index].charge=(self.ev_list[index].soc_out - self.ev_list[index].soc_now) * self.ev_list[index].capacity
        except FileNotFoundError:
            print('文件未被找到')
        except Exception as e:
            print('發生錯誤', e)
        
    def get_FCS_info(self):
        self.load = self.read_file('load.csv')
        self.pv = self.read_file('pv.csv')
        self.ess = self.read_file('ess.csv')
        self.tou = self.read_file('tou.csv') 
        self.read_parameter()
        self.read_se_list()
        self.read_ev_list()
        
    def update_ev_list(self): #將新車加入
        try: 
            with open('new_ev.csv', 'r', newline='') as csvfile:
                reader = csv.reader(csvfile)
                for row in reader:
                    temp_ev = ev(int(row[0]), int(row[1]), int(row[2]),
                                 float(row[3]), float(row[4]), float(row[5]),
                                 int(row[6]), int(row[7]))
                self.ev_list.append(temp_ev)
                self.se_list[int(row[7])-1].index_in_evlist = len(self.ev_list)-1
                self.se_list[int(row[7])-1].time_in = int(row[1])
                self.se_list[int(row[7])-1].time_out = int(row[2])
        except FileNotFoundError:
            print('文件未被找到')
        except Exception as e:
            print('發生錯誤', e)
        
    def schedule(self):
        try:
            m = gp.Model("commercial_schedule")
            now_time = self.now_time
            ev_list = self.ev_list
            se_list = self.se_list
            load = self.load
            pv = self.pv
            tou = self.tou
            num_time = self.num_time
            efficiency = self.efficiency
            Pess = self.Pess
            ess_cap = self.ess_cap
            Ptr =self.Ptr
            ac_Pchar = self.ac_Pchar
            dc_Pchar = self.dc_Pchar

            ess_char = [0] * num_time           #充電量
            ess_dischar = [0] * num_time        #放電量
            ess_char_bool = [0] * num_time      #充電
            ess_dischar_bool = [0] * num_time   #放電

            total_cost = 0
            ess_penalty = 0
            ev_penalty = 0.0
            pc_penalty = 0.0
            Pnet = [0] * num_time
            Pbuy = [0] * num_time    #購入電量

            ess_cost = m.addVar(lb=0)              
            temp_charge=0.0      

            for t in range(now_time-1, num_time): #儲能初始值不等於最後值之懲罰
                ess_char[t] = m.addVar(lb=0, ub=Pess)
                ess_dischar[t] = m.addVar(lb=0, ub=Pess)
                ess_char_bool[t] = m.addVar(vtype = GRB.BINARY)
                ess_dischar_bool[t] = m.addVar(vtype = GRB.BINARY)
                m.addConstr(ess_char_bool[t] + ess_dischar_bool[t] == 1)
                m.addConstr(ess_char[t] - ess_char_bool[t] * Pess <= 0)
                m.addConstr(ess_dischar[t] - ess_dischar_bool[t] * Pess <= 0)

                temp_charge += ess_char[t] * efficiency - ess_dischar[t]

                m.addConstr(self.ess[t-1] * ess_cap + temp_charge >= ess_cap * 0.1)
                m.addConstr(self.ess[t-1] * ess_cap + temp_charge <= ess_cap * 0.9)
            ess_charge = (self.ess_init - self.ess[num_time-1]) * ess_cap #初始值要等於最後值
            m.addConstr(ess_charge - temp_charge <= ess_cost)
            ess_penalty = ess_cost * 50



            se_char = [[0]*len(se_list) for _ in range(num_time)] #每台充電樁在每個區間下的放電量
            for t in range(now_time-1, num_time):
                for index in range(len(se_list)):
                    if(se_list[index].index_in_evlist != -1):
                        if(ev_list[se_list[index].index_in_evlist].time_in < t+1 and ev_list[se_list[index].index_in_evlist].time_out > t+1): #只排入場下一時段和出場前一時段
                            if(ev_list[se_list[index].index_in_evlist].num_se <= self.ac_num_charge): #慢充
                                se_char[t][index] = m.addVar(lb=0, ub=ac_Pchar)
                            elif(ev_list[se_list[index].index_in_evlist].num_se > self.ac_num_charge): #快充
                                se_char[t][index] = m.addVar(lb=0, ub=dc_Pchar)

                
            ev_cost=[0.0] * len(ev_list)           
            for i in range(len(ev_list)): #未達要求充電量之罰金
                ev_cost[i]=m.addVar(lb=0)
                temp_charge = 0.0
                index = ev_list[i].num_se
                for t in range(now_time-1, ev_list[i].time_out):
                    temp_charge = temp_charge + se_char[t][index-1]*efficiency/12
                    m.addConstr(ev_list[i].soc_now*ev_list[i].capacity + temp_charge <= ev_list[i].capacity) 
                m.addConstr(float(ev_list[i].charge)-temp_charge <= ev_cost[i])
                ev_penalty = ev_penalty + ev_cost[i] * 5000  
                
            pc_cost=[0.0] * num_time                
            for t in range(now_time-1, num_time): #超過契約容量的罰金
                Pbuy[t] = m.addVar(lb=0)
                pc_cost[t]=m.addVar(lb=0)
                Pnet[t] = Pnet[t] + load[t] - pv[t]
                for index in range(len(se_list)):
                    Pnet[t] += se_char[t][index] * efficiency
                Pnet[t] += ess_char[t] - ess_dischar[t] * efficiency
                m.addConstr(Pnet[t] <= Pbuy[t])
                total_cost = total_cost + Pbuy[t] * tou[t] / 12 
                
                m.addConstr((Pnet[t]-Ptr) <= pc_cost[t])                #契約處罰
                pc_penalty = pc_penalty + pc_cost[t] * 10000
            
            m.setObjective(total_cost + ev_penalty + pc_penalty + ess_penalty, GRB.MINIMIZE)
            m.optimize()

            ### 取計算結果
            # 1.取充電樁充電功率        (印出來看+寫進csv)
            # 2.取場內車子現在充電電量  (寫進csv)
            # 3.取場內車子預估充電電量  (印出來看)
            # 4.取整個廠的總功率        (印出來看)
            

            # 1.取充電樁充電功率
            x_se_char = [[0.0]*(len(se_list)+1) for _ in range(num_time)]
            for t in range(now_time-1, num_time):
                for index in range(1,len(se_list)):
                    if(se_list[index].index_in_evlist != -1):
                        if(ev_list[se_list[index].index_in_evlist].time_in < t+1 and ev_list[se_list[index].index_in_evlist].time_out > t+1): #只排入場下一時段和出場前一時段
                            x_se_char[t][index] = se_char[t][index].x 
            for t in range(num_time):
                x_se_char[t][0] = t+1 
            with open('charger_power.csv', mode='w', newline='') as file:
                writer = csv.writer(file)
                header = ["time"]
                for index in range(len(se_list)):
                    temp = 'charger' + str(index+1)
                    header.append(temp)
                writer.writerow(header)
                writer.writerows(x_se_char)
 

            # 2.取場內車子現在充電電量
            ev_list_arr = [[0.0]*8 for _ in range(len(ev_list))]
            for i in range(len(ev_list)):
                index = ev_list[i].num_se
                add_soc = x_se_char[now_time-1][index] *efficiency /12 /ev_list[i].capacity
                ev_list[i].soc_now += add_soc
                ev_list_arr[i][0] = ev_list[i].name
                ev_list_arr[i][1] = ev_list[i].time_in
                ev_list_arr[i][2] = ev_list[i].time_out
                ev_list_arr[i][3] = ev_list[i].soc_in*100
                ev_list_arr[i][4] = ev_list[i].soc_out*100
                ev_list_arr[i][5] = round(ev_list[i].soc_now*100, 2)
                ev_list_arr[i][6] = ev_list[i].capacity
                ev_list_arr[i][7] = ev_list[i].num_se
            with open('ev_list.csv', 'w', newline='') as file:
                top_list = ['Number', 'Time_in', 'Time_out', 'Soc_in', 'Soc_out', 'Soc_now', 'EV_capacity', 'se_number']
                csv_writer = csv.writer(file)
                csv_writer.writerow(top_list)
                csv_writer.writerows(ev_list_arr)

        #     # # 3.取場內車子預估充電電量(看看用，之後會刪)(p.s. 上面的和這個不要一起開，不然soc_now會亂掉)
        #     # ev_soc = [[0.0]*(len(ev_list)) for _ in range(num_time)]
        #     # for i in range(len(ev_list)):
        #     #     index = ev_list[i].num_se
        #     #     soc = ev_list[i].soc_now
        #     #     for t in range(now_time-1, ev_list[i].time_out-1):
        #     #         add_soc = se_char[t][index].x *efficiency /12 /ev_list[i].capacity
        #     #         ev_soc[t][i] = soc + add_soc
        #     #         soc = ev_soc[t][i]
        #     # print("ev_soc")
        #     # for t in range(len(ev_soc)):
        #     #     print( t+1, ev_soc[t], sep="\t")
        #     # print("\n")

        #     # 4.取整個廠的總功率
        #     x_Pnet = [0.0] * num_time
        #     x_Pbuy = [0.0] * num_time
        #     for t in range(now_time-1, num_time):
        #         for index in range(1,len(se_list)):
        #                 x_Pnet[t] = x_Pnet[t] + x_se_char[t][index]
        #         x_Pbuy[t] = Pbuy[t].x
        #     # print("time","x_Pnet","x_Pbuy",sep="\t" )
        #     # for t in range(len(x_Pnet)):
        #     #     print(t+1,x_Pnet[t],x_Pbuy[t],sep="\t\t")
        #     # print("\n")

            #5.取得ess並寫入
            x_ess = [[0]*2 for _ in range(num_time)]
            for t in range(now_time-1, num_time):
                self.ess[t] = self.ess[t-1] + (ess_char[t].x*efficiency - ess_dischar[t].x)/12/self.ess_cap   
                #print(self.ess[t])
            for t in range(num_time):
                x_ess[t][0] = t+1
                x_ess[t][1] = self.ess[t]
            with open('ess.csv', mode='w', newline='') as file:
                writer = csv.writer(file)
                header = ["time","ess"]
                writer.writerow(header)
                writer.writerows(x_ess) 
            
            #更新樁列
            with open('se_list.csv', 'w', newline='') as csvfile:
                top_list = ['Number', 'Time_in', 'Time_out']
                csv_writer = csv.writer(csvfile)
                csv_writer.writerow(top_list)
            with open('se_list.csv', 'a', newline='') as csvfile:
                new_se_list = [[0]*3 for _ in range(len(self.se_list))]
                csv_writer = csv.writer(csvfile)
                for index in range(len(self.se_list)):
                    new_se_list[index][0] = str(self.se_list[index].name)
                    new_se_list[index][1] = str(self.se_list[index].time_in)
                    new_se_list[index][2] = str(self.se_list[index].time_out)                
                    csv_writer.writerow(new_se_list[index])
                    
            return_se_char = []
            for index in range(len(se_list)):
                temp_power = power(1, se_list[index].name, x_se_char[now_time][index], now_time)
                return_se_char.append(temp_power)
                
            for index in range(len(ev_list)):
                return_se_char[ev_list[index].num_se-1].state = 1
                return_se_char[ev_list[index].num_se-1].ev_soc = ev_list[index].soc_now*100
                

            return return_se_char
        
        except gp. GurobiError as e:
            print ('Error code ' + str(e. errno ) + ": " + str(e))
            

        # def renew_info(self):
        #     #更新車列
        #     with open('ev_list.csv', 'w', newline='') as csvfile:
        #         top_list = ['Number', 'Time_in', 'Time_out', 'Soc_in', 'Soc_out', 'Soc_now', 'EV_capacity', 'se_number']
        #         csv_writer = csv.writer(csvfile)
        #         csv_writer.writerow(top_list)
        #     with open('ev_list.csv', 'a', newline='') as csvfile:
        #         new_ev_list = [[0]*8 for _ in range(len(self.ev_list))]
        #         csv_writer = csv.writer(csvfile)
        #         for index in range(len(self.ev_list)):
        #             new_ev_list[index][0] = str(self.ev_list[index].name)
        #             new_ev_list[index][1] = str(self.ev_list[index].time_in)
        #             new_ev_list[index][2] = str(self.ev_list[index].time_out)                
        #             new_ev_list[index][3] = str(self.ev_list[index].soc_in)
        #             new_ev_list[index][4] = str(self.ev_list[index].soc_out)
        #             new_ev_list[index][5] = str(self.ev_list[index].soc_now + self.se_char[self.now_time-1][self.ev_list[index].num_se-1])
        #             new_ev_list[index][6] = str(self.ev_list[index].capacity)
        #             new_ev_list[index][7] = str(self.ev_list[index].num_se)
        #             csv_writer.writerow(new_ev_list[index])
        #     #更新樁列
        #     with open('se_list.csv', 'w', newline='') as csvfile:
        #         top_list = ['Number', 'Time_in', 'Time_out']
        #         csv_writer = csv.writer(csvfile)
        #         csv_writer.writerow(top_list)
        #     with open('se_list.csv', 'a', newline='') as csvfile:
        #         new_se_list = [[0]*3 for _ in range(len(self.se_list))]
        #         csv_writer = csv.writer(csvfile)
        #         for index in range(len(self.se_list)):
        #             new_se_list[index][0] = str(self.se_list[index].name)
        #             new_se_list[index][1] = str(self.se_list[index].time_in)
        #             new_se_list[index][2] = str(self.se_list[index].time_out)                
        #             csv_writer.writerow(new_se_list[index])
        #     #更新每台樁每個時段的放電狀況
        #     with open('se_char.csv', 'w', newline='') as csvfile:
        #         top_list = ['time']
        #         for index in range(len(self.se_list)):
        #             top_list.append(str(index+1))
        #         csv_writer = csv.writer(csvfile)
        #         csv_writer.writerow(top_list)
        #     with open('se_char.csv', 'a', newline='') as csvfile:
        #         new_se_char = [[0]*(len(self.se_list)+1) for _ in range(len(self.num_time))]
        #         csv_writer = csv.writer(csvfile)
        #         for t in range(len(self.num_time)):
        #             new_se_char[t][0] = str(t+1)
        #             for index in range(len(self.se_list)):
        #                 new_se_char[t][index+1] = str(self.se_char[t][index])               
        #             csv_writer.writerow(new_se_char[t])
        #     #更新儲能每個時段的電量
        #     with open('ess.csv', 'w', newline='') as csvfile:
        #         top_list = ['time','ess','charge','discharge']
        #         csv_writer = csv.writer(csvfile)
        #         csv_writer.writerow(top_list)
        #     with open('ess.csv', 'a', newline='') as csvfile:
        #         new_ess = [[0]*4 for _ in range(len(self.num_time))]
        #         csv_writer = csv.writer(csvfile)
        #         for t in range(len(self.num_time)):
        #             new_ess[t][0] = str(t+1)
        #             new_ess[t][1] = str(self.ess[t])
        #             new_ess[t][2] = str(self.ess_char[t])
        #             new_ess[t][3] = str(self.ess_dischar[t])      
        #             csv_writer.writerow(new_ess[t])
        
class fcs_new_ev: #有新車加入
    def __init__(self, now_time: int, name: int, time_in: int, time_out: int, soc_in: float, soc_out: float, capacity: int, char_type: int, location_x: float, 
               location_y: float):
        
        self.now_time = now_time        #現在時間
        
        self.num_time = 288             #時間區間數
        self.earnings = 0.0             #充電收益
        self.cost = 0.0                 #購電電費
       
        self.ev_list = []                       #場內電動車
        self.ess = [0]*(self.num_time)          #每時段儲能電量
        self.ess_init = 0.5                       #儲能初始值定為0.5
        self.se_list = []                       #場內充電樁
        self.pnet = [0]*(self.num_time)         #淨負載
        self.Pbuy = [0]*(self.num_time)    #正淨負載
        self.get_FCS_info()
        if(self.check(name, time_in, time_out, soc_in, soc_out, capacity, char_type, location_x, location_y) == 0):  #檢查新車是否到的了本場
            return    
        
    def check(self, name, time_in, time_out, soc_in, soc_out, capacity, char_type, location_x, location_y):
        distance = math.sqrt((self.location_x - location_x)**2 + (self.location_y - location_y)**2)
        remainder = soc_in * capacity #電動車剩餘電量
        if(remainder*0.01 < distance):
            print('電動車剩餘電量無法到達此場')
            return 0
        else:
            num_se = 0
            diff_time = 0
            if(char_type == 1): #慢充
                for index in range(0, self.ac_num_charge):
                    if(self.se_list[index].time_out - time_in < diff_time):
                        num_se = index+1
                        diff_time = self.se_list[index].time_out - time_in
            elif(char_type == 2): #快充
                for index in range(self.ac_num_charge, self.dc_num_charge):
                    if(self.se_list[index].time_out - time_in < diff_time):
                        num_se = index+1
                        diff_time = self.se_list[index].time_out - time_in
            if(num_se == 0):
                print('充電廠內沒有空位的充電樁')
                return 0
            else:
                add_ev = ev(name, time_in, time_out, soc_in, soc_out, soc_in, capacity, num_se)
                self.ev_list.append(add_ev)
                self.ev_list[len(self.ev_list)-1].charge = (self.ev_list[len(self.ev_list)-1].soc_out - self.ev_list[len(self.ev_list)-1].soc_in) * self.ev_list[len(self.ev_list)-1].capacity
                self.se_list[num_se-1].index_in_evlist = len(self.ev_list)-1
                return 1
    
    def read_file(self, file_name): #從本地端讀取資料
        try:
            with open(file_name, 'r', encoding='utf-8', errors='ignore', newline='') as file:
                csv_reader = csv.reader(file)
                header = next(csv_reader)
                info = []
                for row in csv_reader:
                    info.append(float(row[1]))
        except FileNotFoundError:
            print('文件未被找到')
        except Exception as e:       
            print('發生錯誤', e)
        return info
    
    def read_parameter(self): #從本地端讀取資料
        try:
            with open('cpos_parameter.csv', 'r', encoding='utf-8', errors='ignore',
                      newline='') as file:
                csv_reader = csv.reader(file)
                info = [0]
                for row in csv_reader:
                    if(row[1] == '1'):
                        self.efficiency = float(row[3])
                        self.ess_cap = int(row[4])
                        self.Pess = int(row[5])
                        self.Ptr = int(row[6])
                        self.ac_num_charge = int(row[7])
                        self.dc_num_charge = int(row[8])
                        self.ac_charge_price = int(row[9])
                        self.dc_charge_price = int(row[10])
                        self.ac_Pchar = int(row[11])
                        self.dc_Pchar = int(row[12])
                        self.location_x = float(row[13])
                        self.location_y = float(row[14])
        except FileNotFoundError:
            print('文件未被找到')
        except Exception as e:
            print('發生錯誤', e)
    
    def read_se_list(self):
        try:
            with open('se_list.csv', 'r', newline='') as file:
                csv_reader = csv.reader(file)
                header = next(csv_reader)
                for row in csv_reader:
                    temp_se = se(int(row[0]), -1, int(row[1]), int(row[2]))
                    self.se_list.append(temp_se)
        except FileNotFoundError:
            print('文件未被找到')
        except Exception as e:
            print('發生錯誤', e)
                    
    def read_ev_list(self):
        try:
            with open('ev_list.csv', 'r', newline='') as file:
                csv_reader = csv.reader(file)
                header = next(csv_reader)  #跳過第一行
                for row in csv_reader:
                    if(int(row[2]) > self.now_time): #離場時間大於現在時間才加進車列中
                        temp_ev = ev(int(row[0]), int(row[1]), int(row[2]),
                                     float(row[3]), float(row[4]), float(row[5]),
                                     int(row[6]), int(row[7]))
                        self.ev_list.append(temp_ev)
                        self.se_list[int(row[7])-1].index_in_evlist = len(self.ev_list)-1
                
                for index in range(len(self.ev_list)):  #計算在場電動車所需充電量&取得預測電動車負載
                    self.ev_list[index].charge=(self.ev_list[index].soc_out - self.ev_list[index].soc_now) * self.ev_list[index].capacity
        except FileNotFoundError:
            print('文件未被找到')
        except Exception as e:
            print('發生錯誤', e)
        
    def get_FCS_info(self):
        self.load = self.read_file('load.csv')
        self.pv = self.read_file('pv.csv')
        self.ess = self.read_file('ess.csv')
        self.tou = self.read_file('tou.csv') 
        self.read_parameter()
        self.read_se_list()
        self.read_ev_list()


    def schedule(self):
        try:
            m = gp.Model("commercial_schedule")
            now_time = self.now_time
            ev_list = self.ev_list
            se_list = self.se_list
            load = self.load
            pv = self.pv
            tou = self.tou
            num_time = self.num_time
            efficiency = self.efficiency
            Pess = self.Pess
            ess_cap = self.ess_cap
            Ptr =self.Ptr
            ac_Pchar = self.ac_Pchar
            dc_Pchar = self.dc_Pchar

            ess_char = [0] * num_time           #充電量
            ess_dischar = [0] * num_time        #放電量
            ess_char_bool = [0] * num_time      #充電
            ess_dischar_bool = [0] * num_time   #放電

            total_cost = 0
            ess_penalty = 0
            ev_penalty = 0.0
            pc_penalty = 0.0
            Pnet = [0] * num_time
            Pbuy = [0] * num_time    #購入電量

            ess_cost = m.addVar(lb=0)              
            temp_charge=0.0      

            for t in range(now_time-1, num_time): #儲能初始值不等於最後值之懲罰
                ess_char[t] = m.addVar(lb=0, ub=Pess)
                ess_dischar[t] = m.addVar(lb=0, ub=Pess)
                ess_char_bool[t] = m.addVar(vtype = GRB.BINARY)
                ess_dischar_bool[t] = m.addVar(vtype = GRB.BINARY)
                m.addConstr(ess_char_bool[t] + ess_dischar_bool[t] == 1)
                m.addConstr(ess_char[t] - ess_char_bool[t] * Pess <= 0)
                m.addConstr(ess_dischar[t] - ess_dischar_bool[t] * Pess <= 0)

                temp_charge += ess_char[t] * efficiency - ess_dischar[t]

                m.addConstr(self.ess[t-1] * ess_cap + temp_charge >= ess_cap * 0.1)
                m.addConstr(self.ess[t-1] * ess_cap + temp_charge <= ess_cap * 0.9)
            ess_charge = (self.ess_init - self.ess[num_time-1]) * ess_cap #初始值要等於最後值
            m.addConstr(ess_charge - temp_charge <= ess_cost)
            ess_penalty = ess_cost * 50



            se_char = [[0]*len(se_list) for _ in range(num_time)] #每台充電樁在每個區間下的放電量
            for t in range(now_time-1, num_time):
                for index in range(len(se_list)):
                    if(se_list[index].index_in_evlist != -1):
                        if(ev_list[se_list[index].index_in_evlist].time_in < t+1 and ev_list[se_list[index].index_in_evlist].time_out > t+1): #只排入場下一時段和出場前一時段
                            if(ev_list[se_list[index].index_in_evlist].num_se <= self.ac_num_charge): #慢充
                                se_char[t][index] = m.addVar(lb=0, ub=ac_Pchar)
                            elif(ev_list[se_list[index].index_in_evlist].num_se > self.ac_num_charge): #快充
                                se_char[t][index] = m.addVar(lb=0, ub=dc_Pchar)
                
            ev_cost=[0.0] * len(ev_list)           
            for i in range(len(ev_list)): #未達要求充電量之罰金
                ev_cost[i]=m.addVar(lb=0)
                temp_charge = 0.0
                index = ev_list[i].num_se
                for t in range(now_time-1, ev_list[i].time_out):
                    temp_charge = temp_charge + se_char[t][index-1]*efficiency/12
                    m.addConstr(ev_list[i].soc_now*ev_list[i].capacity + temp_charge <= ev_list[i].capacity) 
                m.addConstr(float(ev_list[i].charge)-temp_charge <= ev_cost[i])
                ev_penalty = ev_penalty + ev_cost[i] * 5000  
                
            pc_cost=[0.0] * num_time                
            for t in range(now_time-1, num_time): #超過契約容量的罰金
                Pbuy[t] = m.addVar(lb=0)
                pc_cost[t]=m.addVar(lb=0)
                Pnet[t] = Pnet[t] + load[t] - pv[t]
                for index in range(len(se_list)):
                    Pnet[t] += se_char[t][index] * efficiency
                Pnet[t] += ess_char[t] - ess_dischar[t] * efficiency
                m.addConstr(Pnet[t] <= Pbuy[t])
                total_cost = total_cost + Pbuy[t] * tou[t] / 12 
                
                m.addConstr((Pnet[t]-Ptr) <= pc_cost[t])                #契約處罰
                pc_penalty = pc_penalty + pc_cost[t] * 10000
            
            m.setObjective(total_cost + ev_penalty + pc_penalty + ess_penalty, GRB.MINIMIZE)
            m.optimize()
            
            x_Pnet = [0] * num_time
            
                    
            for t in range(now_time-1, num_time):
                for index in range(len(se_list)):
                    if(se_list[index].index_in_evlist != -1):
                        if(ev_list[se_list[index].index_in_evlist].time_in < t+1 and ev_list[se_list[index].index_in_evlist].time_out > t+1): #只排入場下一時段和出場前一時段   
                            x_Pnet[t] += se_char[t][index].x * efficiency
                            
            x_se_char = [0.0] * num_time
            total_charge = 0.0
            for t in range(now_time-1, num_time):
                if(ev_list[len(ev_list)-1].time_in < t+1 and ev_list[len(ev_list)-1].time_out > t+1):
                            total_charge += se_char[t][ev_list[len(ev_list)-1].num_se-1].x
                            x_se_char[t] = se_char[t][ev_list[len(ev_list)-1].num_se-1].x

            unit_price_of_ch, total_price_of_space = estimate_price(x_Pnet, x_se_char, tou, 1, ev_list[len(ev_list)-1].time_in, ev_list[len(ev_list)-1].time_out, Ptr)
            with open('new_ev.csv', 'w', newline='') as csvfile:
                new_ev = [0]*8
                csv_writer = csv.writer(csvfile)
                new_ev[0] = str(self.ev_list[len(ev_list)-1].name)
                new_ev[1] = str(self.ev_list[len(ev_list)-1].time_in)
                new_ev[2] = str(self.ev_list[len(ev_list)-1].time_out)                
                new_ev[3] = str(self.ev_list[len(ev_list)-1].soc_in * 100)
                new_ev[4] = str(self.ev_list[len(ev_list)-1].soc_out * 100)
                new_ev[5] = str(self.ev_list[len(ev_list)-1].soc_now * 100)
                new_ev[6] = str(self.ev_list[len(ev_list)-1].capacity)
                new_ev[7] = str(self.ev_list[len(ev_list)-1].num_se)
                csv_writer.writerow(new_ev)
                
            final_soc = int(self.ev_list[len(ev_list)-1].soc_in + total_charge/self.ev_list[len(ev_list)-1].capacity*100)
            return self.ev_list[len(ev_list)-1].name, final_soc,unit_price_of_ch, total_price_of_space
        
        except gp. GurobiError as e:
            print ('Error code ' + str(e. errno ) + ": " + str(e))
        
                
#myfcs = fcs_new_ev(12, 21, 12, 26, 37, 75, 100, 1, 4.2, 8.6)
#unit_price_of_ch, total_price_of_space, total_charge = myfcs.schedule()
#myfcs.update_ev_list()
                
myfcs_1 = fcs(13, -1)
se_char = myfcs_1.schedule()






    
        
        
        
