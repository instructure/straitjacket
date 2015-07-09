#include <tunables/global>

profile straitjacket/compiled/d {
  #include <abstractions/base>
  #include <abstractions/straitjacket-base>

  set rlimit nproc <= 20,
}