import FCS
import sys

if __name__ == "__main__":

    args = sys.argv
    
    if (len(sys.argv) != 3):
        print('error')

    script_name = args[0]
    nowtime = int(args[1])
    ev_new = int(args[2])
    
    myroad = FCS.fcs(nowtime,ev_new)
    se_char= myroad.schedule()
