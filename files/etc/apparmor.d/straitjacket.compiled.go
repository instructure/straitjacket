#include <tunables/global>

profile straitjacket/compiled/go {
  #include <abstractions/base>
  #include <abstractions/straitjacket-base>

  set rlimit nproc <= 32,

  /var/local/straitjacket/tmp/execute/?* r,
}
