import FCS
import sys

if __name__ == "__main__":

    args = sys.argv
    
    script_name = args[0]
    nowtime = int(args[1])
    name = int(args[2])
    time_in = int(args[3])
    time_out = int(args[4])
    soc_in = float(args[5])
    soc_out = float(args[6])
    capacity = int(args[7])
    char_type = int(args[8])
    location_x =  float(args[9])
    location_y = float(args[10])
    
    myroad = FCS.FCS_new_ev(nowtime, name, time_in, time_out, soc_in, soc_out, capacity, char_type, location_x, location_y)
    se_char= myroad.schedule()