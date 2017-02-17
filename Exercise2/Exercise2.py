
from threading import Thread
import threading
import Queue




q1 = Queue.LifoQueue(1)

	
	
def thread_1(q1):

	for j in range(0,100):
		temp = q1.get()
		temp+=1
		print("A: ", temp)
		q1.put(temp)

def thread_2(q1):

	for n in range(0,99):
		temp = q1.get()
		temp-=1
		print("B: ", temp)
		q1.put(temp)
	


def main():
	i = 0
	q1.put(i)

	thread1 = Thread(target = thread_1, args=([q1]),)
	thread2 = Thread(target = thread_2, args=([q1]),)
	thread1.start()
	thread2.start()
	thread2.join()
	
	
	
	print q1.get()
	
	
	
	
main()
	
