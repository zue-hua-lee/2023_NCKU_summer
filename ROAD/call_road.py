import road
import sys

if __name__ == "__main__":

    args = sys.argv
    
    if (len(sys.argv) != 3):
        print('error')

    script_name = args[0]
    nowtime = int(args[1])
    ev_new = int(args[2])
    
    myroad = road.road(nowtime,ev_new)
    se_char= myroad.schedule()
    return_se_char = [[0]*5 for _ in range(len(se_char))]
    for index in range(len(se_char)):
        return_se_char[index][0] = se_char[index].StationID
        return_se_char[index][1] = se_char[index].ChargeID
        return_se_char[index][2] = se_char[index].Power
        return_se_char[index][3] = se_char[index].TimeStamp
        return_se_char[index][4] = se_char[index].ev_soc
    for row in return_se_char:
        print(' '.join(map(str, row)))
