
from threading import Thread

i = 0
	
	
def thread_1():
	global i
	for j in range(0,100):
		i+=1
		print("B: ", i)

def thread_2():
	global i
	for n in range(0,99):
		i-=1
		print("A: ", i)


def main():
	thread1 = Thread(target = thread_1, args=(),)
	thread2 = Thread(target = thread_2, args=(),)
	thread1.start()
	thread2.start()
	thread2.join()
	print i
main()
	
