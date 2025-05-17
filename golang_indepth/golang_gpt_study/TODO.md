# IMPORTANT REVISION POINTS FOR GOLANG.

1. Revise the template literal string from practice challenge 1.

#
2. Revise closures from lec 4

#
3. From lec 5, remember if you are embedding other struct in some struct dont
   need to specify the type in front of it so that you can direclty access the
   nested struct's fields from the outer struct without chaining the dot operator.

#
4. From lec 9, practice file operations. Once you write to a file and attemp to
   read from it without closing it you wont be able to read from it as the cursor
   will be at the last position and there will be nothing to read after it.
   So in that case you need to close the file and open it again so that the data
   could be read from the file.

#
5. Revise buffered/unbuffered channel. Time.After is a channel used to perform
   blocking operations.

#
6. Revise mutexes and send and receive channels.

#
7. I can say that ctx.Done() will only receive a value whenever I
   will call cancel() in my go program else it will never recieve a signal. If
   the cancel function is never called then it ctx.Done() will not ever receive
   a value.(Assuming the ctx was initialised with a WithCancel Method)

    # # # # # # # # # # # # # # # # # # # #
    When does ctx.Done() close?
    🔔 When you call cancel() (if using context.WithCancel()).
    🕒 When timeout/deadline is reached (if using context.WithTimeout() or context.WithDeadline()).
    🔗 If it's derived from another context that was cancelled.
    # # # # # # # # # # # # # # # # # # # #


#
8. WithTimeout and WithDeadline does not need cancel to be called explicitly.
   But it is a good practice to call cancel() function right after you instantiate
   the parent context with the cancel variable.

   time.Now() -> Gets the current time.
   time.Now().Add() -> If you want to do something after 5 seconds FROM NOW.



#
9.  # Context.WithValue
    type ctxKey string
    func mimicGo(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

    # Extract the value from the context.
	id := ctx.Value(ctxKey("id"))
	fmt.Println("Processing started for the id", id)
    }

    func main() {
	var wg sync.WaitGroup

    # Pass the value in the context witht the key
	ctx := context.WithValue(context.Background(), ctxKey("id"), 42)

	wg.Add(1)
	go mimicGo(&wg, ctx)

	wg.Wait()
    }









