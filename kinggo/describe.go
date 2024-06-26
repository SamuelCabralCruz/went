package kinggo

import (
	"github.com/SamuelCabralCruz/going/phi"
	"github.com/onsi/ginkgo/v2"
)

func DescribeFunction(function any, args ...any) bool {
	return ginkgo.Describe(phi.FunctionName(function), args...)
}

func DescribeType[T any](args ...any) bool {
	return ginkgo.Describe(phi.BaseTypeName[T](), args...)
}
