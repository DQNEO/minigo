// start up routines
.text
  .global _start
_start:
  movq %rsp,    %rbp # initial stack top addr
  movq 0(%rsp), %rdx # argc
  leaq 8(%rsp), %rsi # argv
  # get envp
  movq %rbp, %rax # stack top addr
  movq %rdx, %rbx # argc
  imulq $8,  %rbx # argc * 8
  addq %rbx, %rax # stack top addr + (argc * 8)
  addq $16,  %rax # + 16 (skip null and go to next) => envp
  movq %rax, iruntime.envp+0(%rip) # envp
  movq $0, %rax
  movq $0, %rbx
  jmp iruntime.rt0_go

iruntime.rt0_go:
  callq iruntime.args
  jmp _init_packages

iruntime.args:
  movq %rdx, %rax
  movq %rsi, %rbx
  movq %rbx, iruntime.argv+0(%rip)  # ptr
  movq %rax, iruntime.argv+8(%rip)  # len
  movq %rax, iruntime.argv+16(%rip) # cap
  ret

iruntime.main:
  FUNCALL main.main
  movq $0,  %rdi
  FUNCALL iruntime.exit
