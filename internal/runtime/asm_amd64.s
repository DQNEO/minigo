// start up routines
.text
  .global _start
_start:
  mov %rsp,    %rbp # initial stack top addr
  mov 0(%rsp), %rdx # argc
  lea 8(%rsp), %rsi # argv
  # get envp
  mov %rsp, %rax
  mov %rdx, %rbx
  imul $8, %rbx
  addq %rbx, %rax
  addq $16, %rax
  mov %rax, iruntime.envp+0(%rip) # envp
  jmp iruntime.rt0_go

iruntime.rt0_go:
  call iruntime.args
  mov $0, %rax
  mov $0, %rbx
  jmp _init_packages

iruntime.args:
  mov %rdx, %rax
  mov %rsi, %rbx
  mov %rbx, iruntime.argv+0(%rip)  # ptr
  mov %rax, iruntime.argv+8(%rip)  # len
  mov %rax, iruntime.argv+16(%rip) # cap
  ret

iruntime.main:
  FUNCALL main.main
  mov $0,  %rdi
  FUNCALL iruntime.exit
