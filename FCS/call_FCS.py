import FCS
import sys

if __name__ == "__main__":

    args = sys.argv
    
    script_name = args[0]
    nowtime = int(args[1])
    ev_new = int(args[2])
    
    myroad = FCS.FCS(nowtime,ev_new)
    se_char= myroad.schedule()