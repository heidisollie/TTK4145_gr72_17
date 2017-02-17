from threading import Thread
import Queue


def thread1func(qu):
	for j in range(100):
		i = qu.get()
		i += 1
		qu.put(i)
	return 0

def thread2func(qu):
	for k in range(100 - 1):
		i = qu.get()
		i -= 1
		qu.put(i)
	return 0

def main():
	q = Queue.LifoQueue(1)
	q.put(0)
	thread1 = Thread(target = thread1func, args = [q])
	thread2 = Thread(target = thread2func, args = [q])

	thread1.start()
	thread2.start()

	thread1.join()
	thread2.join()

	i = q.get()
	print(i)


main()
