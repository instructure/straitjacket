#include <tunables/global>

profile straitjacket/compiled/c {
  #include <abstractions/base>
  #include <abstractions/straitjacket-base>

  set rlimit nproc <= 12,
}
