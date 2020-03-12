// start up routines
.text
  .global _start
_start:
  mov 0(%rsp), %rdx # argc
  lea 8(%rsp), %rsi # argv
  jmp runtime.rt0_go

runtime.rt0_go:
  mov %rdx, %rax
  mov %rsi, %rbx
  mov %rbx, iruntime.argv+0(%rip)  # ptr
  mov %rax, iruntime.argv+8(%rip)  # len
  mov %rax, iruntime.argv+16(%rip) # cap
  mov $0, %rax
  mov $0, %rbx
  jmp _init_packages

iruntime.main:
  FUNCALL main.main
  mov $0,  %rdi
  FUNCALL iruntime.exit
