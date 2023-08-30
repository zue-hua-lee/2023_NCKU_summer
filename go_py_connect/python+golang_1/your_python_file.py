import sys



### 一維陣列 
def your_python_function(num1:int,num2:int):
    add1 = num1+num2+1
    add2 = num1+num2+2
    add3 = num1+num2+3
    return add1, add2, add3

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: example.py <num1> <num2>")
    num1 = int(sys.argv[1])
    num2 = int(sys.argv[2])

    ans1, ans2, ans3 = your_python_function(num1,num2)
    print(ans1,ans2,ans3)
