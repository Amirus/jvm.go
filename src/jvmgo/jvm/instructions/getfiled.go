package instructions

import (
	"jvmgo/jvm/rtda"
	rtc "jvmgo/jvm/rtda/class"
)

// Fetch field from object
type getfield struct{ Index16Instruction }

func (self *getfield) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	ref := stack.PopRef()
	if ref == nil {
		frame.Thread().ThrowNPE()
		return
	}

	cp := frame.Method().ConstantPool()
	kFieldRef := cp.GetConstant(self.index).(*rtc.ConstantFieldref)
	field := kFieldRef.InstanceField()
	val := field.GetValue(ref)

	stack.Push(val)
}

// Get static field from class
type getstatic struct{ Index16Instruction }

func (self *getstatic) Execute(frame *rtda.Frame) {
	cp := frame.Method().Class().ConstantPool()
	kFieldRef := cp.GetConstant(self.index).(*rtc.ConstantFieldref)
	field := kFieldRef.StaticField()

	class := field.Class()
	if class.InitializationNotStarted() {
		frame.RevertNextPC() // undo getstatic
		frame.Thread().InitClass(class)
		return
	}

	val := field.GetStaticValue()
	frame.OperandStack().Push(val)
}
