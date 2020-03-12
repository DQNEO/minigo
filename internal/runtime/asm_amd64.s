// start up routines
.text
  .global	_start
_start:
  pop %rax # argc
  mov %rsp, %rbx # argv
  mov %rbx, iruntime.argv(%rip)    # ptr
  mov %rax, iruntime.argv+8(%rip)  # len
  mov %rax, iruntime.argv+16(%rip) # cap
  mov $0, %rax
  mov $0, %rbx
  jmp _init_packages

iruntime.main:
  FUNCALL main.main
  mov $0,  %rdi
  FUNCALL iruntime.exit
