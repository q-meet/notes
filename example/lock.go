package main

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	close(ch1)
	close(ch2)
	for {
		select {
		case <-ch2:
			print("ch2 \n")
		case <-ch1:
			print("ch1 \n")
			//ch1 = nil
		}
		print("for \n")
		break
	}
}
