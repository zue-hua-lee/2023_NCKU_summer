import sys


### 二維陣列 
def your_python_function(num1:int,num2:int):
    matrix = [[0, 2, 3], [4, 0, 6], [7, 8, 9]]
    matrix[0][0] = num1
    matrix[1][1] = num2
    return matrix

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: example.py <num1> <num2>")
    num1 = int(sys.argv[1])
    num2 = int(sys.argv[2])
    
    matrix = your_python_function(num1,num2)
    for row in matrix:
        print(' '.join(map(str, row)))