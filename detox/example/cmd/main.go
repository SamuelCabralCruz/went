package main

import (
	"errors"
	"fmt"
	"github.com/SamuelCabralCruz/went/detox"
	"github.com/SamuelCabralCruz/went/detox/example/pkg"
	"github.com/SamuelCabralCruz/went/detox/internal/common"
)

func fakedNamed(s string) (string, error) { return "ok", errors.New("coucou named function " + s) }

func main() {
	myMock := pkg.NewSomeMockClass()
	mockedHello := detox.When(myMock.Detox, myMock.Hello)
	//mockedHello2 := detox.When(myMock.Detox, myMock.Hello2)
	mockedPrepare := detox.When(myMock.Detox, myMock.Prepare)
	myMock2 := pkg.NewSomeMockClass()
	mocked2Hello := detox.When(myMock2.Detox, myMock2.Hello)
	myOtherMock := pkg.NewAnotherMockClass()
	mocked3Bye := detox.When(myOtherMock.Detox, myOtherMock.Bye)

	// real impl
	mockedHello.Call(pkg.Impl{}.Hello)
	fmt.Println(myMock.Hello("after named fake"))

	// react to args
	mockedHello.Call(func(s string) (string, error) {
		if s == "sam" {
			return "goat", nil
		} else {
			return "lcm", nil
		}
	})

	mockedHello.Call(fakedNamed)
	fakedInline := func(s string) (string, error) { return "ok", errors.New("coucou inline function " + s) }
	mocked2Hello.Call(fakedInline)
	fmt.Println(myMock.Hello("after named fake"))
	fmt.Println(myMock.Hello("sam"))
	fmt.Println(myMock2.Hello("after named fake"))

	mockedHello.Call(fakedInline)

	fmt.Println(myMock.Hello("before method reset"))
	mockedHello.Reset()
	//fmt.Println(myMock.Hello("after method reset")) // -> should fail
	myMock.Default(pkg.Impl{})
	fmt.Println(myMock.Hello("should call through real impl"))

	mockedHello.Call(func(s string) (string, error) { return "ok", errors.New("coucou anonymous " + s) })

	fmt.Println(myMock.Hello("before mock reset"))
	myMock.Reset()
	//fmt.Println(myMock.Hello("after mock reset")) // -> should fail

	mockedHello.CallOnce(fakedNamed)
	mockedHello.CallOnce(fakedInline)
	mockedHello.CallOnce(func(s string) (string, error) { return "ok", errors.New("coucou anonymous " + s) })
	mockedHello.Call(func(s string) (string, error) { return "ok", errors.New("coucou remainder " + s) })

	fmt.Println(myMock.Hello("1st"))
	fmt.Println(myMock.Hello("2nd"))
	fmt.Println(myMock.Hello("3rd"))
	fmt.Println(myMock.Hello("4th"))
	fmt.Println(myMock.Hello("5th"))
	fmt.Println(myMock.Hello("6th"))

	mockedPrepare.Call(func() pkg.Another { return myOtherMock })
	mocked3Bye.Call(func(s string) { fmt.Println("it works " + s) })
	myMock.Prepare().Bye("sam")

	//fmt.Println(mockedHello.Calls())
	//fmt.Println(mockedHello.CallsCount())
	//fmt.Println(mockedHello.NthCall(3))

	myMock.Reset()

	//fmt.Println(mockedHello.Calls())
	//fmt.Println(mockedHello.CallsCount())
	//fmt.Println(mockedHello.NthCall(3)) // should panic

	// fluent api
	mockedHello.CallOnce(func(s string) (string, error) {
		return "fluently faked 1", errors.New(s)
	})
	mockedHello.CallOnce(func(s string) (string, error) {
		return "fluently faked 2", errors.New(s)
	})
	mockedHello.Call(func(s string) (string, error) {
		return "fluently faked 3", errors.New(s)
	})
	fmt.Println(myMock.Hello("1st"))
	fmt.Println(myMock.Hello("2nd"))
	fmt.Println(myMock.Hello("3rd"))
	fmt.Println(myMock.Hello("4th"))

	mockedHello.Reset()

	mockedHello.Call(func(s string) (string, error) {
		return "fluently persistent unconditional", errors.New(s)
	})
	mockedHello.CallOnce(func(s string) (string, error) {
		return "fluently ephemeral unconditional", errors.New(s)
	})
	// since this one is not ephemeral, it can be called has many times as desired with 1st
	// should have priority over the previous one because it is specific to this use case
	mockedHello.WithArgs("1st").Call(func(s string) (string, error) {
		return "fluently conditional persistent", errors.New(s)
	})
	// will have priority over the previous one because it is ephemeral
	mockedHello.WithArgs("1st").CallOnce(func(s string) (string, error) {
		return "fluently conditional ephemeral", errors.New(s)
	})
	mockedHello.WithArgs("2nd").CallOnce(func(s string) (string, error) {
		return "fluently faked 2", errors.New(s)
	})
	fmt.Println(myMock.Hello("1st"))
	fmt.Println(myMock.Hello("2nd"))
	//fmt.Println(myMock.Hello("3rd"))
	//fmt.Println(myMock.Hello("4th"))
	fmt.Println(myMock.Hello("1st"))
	fmt.Println(myMock.Hello("1st"))
	fmt.Println(myMock.Hello("1st"))
	fmt.Println(myMock.Hello("1st"))
	fmt.Println(myMock.Hello("1st"))
	fmt.Println(myMock.Hello("2nd"))

	// TODO: CallThrough
	// (provide real) (could be done using fake implementation)
	// what if we could pass in an object to default to it...?
	// would have the lowest priority (fallback)

	// TODO: create custom matchers
	fmt.Println(mockedHello.Assert().HasBeenCalled())
	fmt.Println(mockedHello.Assert().HasBeenCalledTimes(8))
	fmt.Println(mockedHello.Assert().HasBeenCalledOnce())
	fmt.Println(mockedHello.Assert().HasBeenCalledWith("1st"))
	fmt.Println(mockedHello.Assert().HasBeenCalledWith("4th"))
	fmt.Println(mockedHello.Assert().HasNthCall(1, common.NewCall("2nd")))
	fmt.Println(mockedHello.Assert().HasBeenCalledTimesWith(2, "2nd"))
	fmt.Println(mockedHello.Assert().HasCalls(common.NewCall("2nd"), common.NewCall("2nd"), common.NewCall("2nd")))
	fmt.Println(mockedHello.Assert().HasCalls(common.NewCall("2nd"), common.NewCall("2nd")))
	fmt.Println(mockedHello.Assert().HasCalls(common.NewCall("2nd"), common.NewCall("1st"), common.NewCall("2nd")))
	fmt.Println(mockedHello.Assert().HasOrderedCalls(common.NewCall("2nd"), common.NewCall("1st"), common.NewCall("2nd")))
	fmt.Println(mockedHello.Assert().HasOrderedCalls(
		common.NewCall("1st"),
		common.NewCall("2nd"),
		common.NewCall("1st"),
		common.NewCall("1st"),
		common.NewCall("1st"),
		common.NewCall("1st"),
		common.NewCall("1st"),
		common.NewCall("2nd"),
	))
	fmt.Println(mockedHello.Assert().HasOrderedCalls(
		common.NewCall("1st"),
		common.NewCall("1st"),
		common.NewCall("2nd"),
		common.NewCall("1st"),
		common.NewCall("1st"),
		common.NewCall("1st"),
		common.NewCall("1st"),
		common.NewCall("2nd"),
	))
	fmt.Println(mockedHello.Assert().HasCalls(
		common.NewCall("1st"),
		common.NewCall("1st"),
		common.NewCall("2nd"),
		common.NewCall("1st"),
		common.NewCall("1st"),
		common.NewCall("1st"),
		common.NewCall("1st"),
		common.NewCall("2nd"),
	))

	// TODO: HaveBeenCalledNth(int) -> called n times
	// TODO: HaveCalls([][]any) -> any order
	// TODO: HaveBeenCalledWith([]any) -> contains a calls with provided args
	// TODO: HaveCallSequence([][]any) -> specific order

	// TODO: caveats
	// TODO: Mock return values (stub ramification)
}
