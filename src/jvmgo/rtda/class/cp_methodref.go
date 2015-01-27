package class

import (
    //"fmt"
    cf "jvmgo/classfile"
)

type ConstantMethodref struct {
    className       string
    name            string
    descriptor      string
    cp              *ConstantPool
    method          *Method
}

type ConstantInterfaceMethodref struct {
    ConstantMethodref
}

// method resolution
func (self *ConstantMethodref) StaticMethod() (*Method) {
    if self.method == nil {
        method := self.findMethod(self.className)
        if method != nil && method.IsStatic() {
            if method.IsNative() {
                method.nativeMethod = findNativeMethod(method)
            }
            self.method = method
        } else {
            // todo
            panic("method not found!" + self.className + "." + self.name + "|" + self.descriptor)
        }
    }
    return self.method
}
func (self *ConstantMethodref) SpecialMethod() (*Method) {
    if self.method == nil {
        method := self.findMethod(self.className)
        if method != nil && !method.IsStatic() {
            if method.IsNative() {
                method.nativeMethod = findNativeMethod(method)
            }
            self.method = method
        } else {
            // todo
            panic("method not found!" + self.className + "." + self.name + "|" + self.descriptor)
        }
    }
    return self.method
}
func (self *ConstantMethodref) findMethod(className string) (*Method) {
    class := self.cp.class.classLoader.LoadClass(className)
    return class.GetMethod(self.name, self.descriptor)
}

// todo
func (self *ConstantMethodref) VirtualMethodArgCount() (uint) {
    return self.SpecialMethod().ArgCount()
}
func (self *ConstantMethodref) VirtualMethod(ref *Obj) (*Method) {
    classLoader := self.cp.class.classLoader
    className := ref.class.name //!!
    for {
        if className != "" {
            class := classLoader.LoadClass(className)
            method := class.GetMethod(self.name, self.descriptor)
            if method != nil && !method.IsStatic() {
                if method.IsNative() {
                    method.nativeMethod = findNativeMethod(method)
                }
                self.method = method
                return method
            } else {
                className = class.superClassName
            }
        } else {
            break
        }
    }
    // todo
    panic("method not found!" + self.className + "." + self.name + "|" + self.descriptor)
}

func newConstantMethodref(cp *ConstantPool, methodrefInfo *cf.ConstantMethodrefInfo) (*ConstantMethodref) {
    methodref := &ConstantMethodref{}
    methodref.cp = cp
    methodref.className = methodrefInfo.ClassName()
    methodref.name = methodrefInfo.Name()
    methodref.descriptor = methodrefInfo.Descriptor()
    return methodref
}

func newConstantInterfaceMethodref(cp *ConstantPool, methodrefInfo *cf.ConstantInterfaceMethodrefInfo) (*ConstantInterfaceMethodref) {
    methodref := &ConstantInterfaceMethodref{}
    methodref.cp = cp
    methodref.className = methodrefInfo.ClassName()
    methodref.name = methodrefInfo.Name()
    methodref.descriptor = methodrefInfo.Descriptor()
    return methodref
}
